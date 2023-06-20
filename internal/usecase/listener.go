package usecase

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/CAMELNINGA/WAL-transport.git/config"
	"github.com/CAMELNINGA/WAL-transport.git/internal/models"
	"github.com/CAMELNINGA/WAL-transport.git/internal/sanitize"
	error_walListner "github.com/CAMELNINGA/WAL-transport.git/pkg/error_walListner"

	"github.com/jackc/pgx"
	"github.com/sirupsen/logrus"
)

const errorBufferSize = 100

// Logical decoding plugin.
const pgOutputPlugin = "pgoutput"

type publisher interface {
	Publish(string, models.Message) error
}

type parser interface {
	ParseWalMessage([]byte, *models.WalTransaction) error
}

type replication interface {
	CreateReplicationSlotEx(slotName, outputPlugin string) (consistentPoint string, snapshotName string, err error)
	DropReplicationSlot(slotName string) (err error)
	StartReplication(slotName string, startLsn uint64, timeline int64, pluginArguments ...string) (err error)
	WaitForReplicationMessage(ctx context.Context) (*pgx.ReplicationMessage, error)
	SendStandbyStatus(k *pgx.StandbyStatus) (err error)
	IsAlive() bool
	Close() error
}

type repository interface {
	CreatePublication(name string) error
	GetSlotLSN(slotName string) (string, error)
	IsAlive() bool
	Close() error
}

type saver interface {
	SaveData(ctx context.Context, tx *models.WalTransaction) error
}

// Listener main service struct.
type Listener struct {
	cfg        config.Config
	log        *logrus.Entry
	mu         sync.RWMutex
	slotName   string
	publisher  publisher
	replicator replication
	repository repository
	parser     parser
	sanitizer  sanitize.Handler
	lsn        uint64
	errChannel chan error
}

// NewWalListener create and initialize new service instance.
func NewWalListener(
	log *logrus.Entry,
	slotName string,
	repo repository,
	repl replication,
	parser parser,
	publisher publisher,
	sanitizer sanitize.Handler,
) *Listener {
	return &Listener{
		log:        log,
		slotName:   slotName,
		repository: repo,
		replicator: repl,
		parser:     parser,
		publisher:  publisher,
		sanitizer:  sanitizer,
		errChannel: make(chan error, errorBufferSize),
	}
}

const (
	protoVersion    = "proto_version '1'"
	publicationName = "wal-transport"
)

// Process is main service entry point.
func (l *Listener) Process(ctx context.Context) error {
	logger := l.log.WithField("slot_name", l.slotName)

	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	logger.Infoln("service was started")

	if err := l.repository.CreatePublication(publicationName); err != nil {
		logger.WithError(err).Warnln("skip create publication")
	}

	slotIsExists, err := l.slotIsExists(ctx)
	if err != nil {
		return fmt.Errorf("slot is exists: %w", err)
	}

	if !slotIsExists {
		consistentPoint, _, err := l.replicator.CreateReplicationSlotEx(l.slotName, pgOutputPlugin)
		if err != nil {
			return fmt.Errorf("create replication slot: %w", err)
		}

		lsn, err := pgx.ParseLSN(consistentPoint)
		if err != nil {
			return fmt.Errorf("parse lsn: %w", err)
		}

		l.setLSN(lsn)
		logger.Infoln("new slot was created")
	} else {
		logger.Infoln("slot already exists, LSN updated")
	}

	go l.Stream(ctx)

	refresh := time.NewTicker(time.Duration(30))
	defer refresh.Stop()

	var svcErr *error_walListner.ServiceErr

ProcessLoop:
	for {
		select {
		case <-refresh.C:
			if !l.replicator.IsAlive() {
				return fmt.Errorf("replicator: %w", error_walListner.ErrReplConnectionIsLost)
			}

			if !l.repository.IsAlive() {
				return fmt.Errorf("repository: %w", error_walListner.ErrConnectionIsLost)
			}
		case err := <-l.errChannel:
			if errors.As(err, &svcErr) {
				return err
			}

			l.log.WithError(err).Errorln("received error")
		case <-ctx.Done():
			logger.Debugln("context was canceled")

			if err := l.Stop(); err != nil {
				logger.WithError(err).Errorln("listener stop error")
			}

			break ProcessLoop
		}
	}

	return nil
}

// Stream receive event from PostgreSQL.
// Accept message, apply filter and  publish it in NATS server.
func (l *Listener) Stream(ctx context.Context) {
	if err := l.replicator.StartReplication(
		l.slotName,
		l.readLSN(),
		-1,
		protoVersion,
		publicationNames(publicationName),
	); err != nil {
		l.errChannel <- error_walListner.NewListenerError("StartReplication()", err)

		return
	}

	go l.SendPeriodicHeartbeats(ctx)

	tx := models.NewWalTransaction()

	for {
		if err := ctx.Err(); err != nil {
			l.errChannel <- error_walListner.NewListenerError("read msg", err)
			break
		}

		msg, err := l.replicator.WaitForReplicationMessage(ctx)
		if err != nil {
			l.errChannel <- error_walListner.NewListenerError("WaitForReplicationMessage()", err)
			continue
		}

		if msg != nil {
			if msg.WalMessage != nil {
				l.log.WithField("wal", msg.WalMessage.WalStart).Debugln("receive wal message")

				if err := l.parser.ParseWalMessage(msg.WalMessage.WalData, tx); err != nil {
					l.log.WithError(err).Errorln("msg parse failed")
					l.errChannel <- fmt.Errorf("unmarshal wal message: %w", err)

					continue
				}
				var newActions []*models.ActionData
				for _, actions := range tx.Actions {
					v := l.sanitizer.Handle(actions)
					if v != nil {
						newActions = append(newActions, v)
					}

				}
				tx.Actions = newActions
				l.log.Info(tx.RelationStore)
				l.log.Info(tx)
				//TODO: interfase work change wal logs to json file
				if tx.CommitTime != nil {
					l.log.WithField("commit_time", tx.CommitTime).Debugln("commit transaction")
					message := tx.CreateMessges()
					for _, event := range message {
						subjectName := event.SubjectName()
						if err = l.publisher.Publish(subjectName, event); err != nil {
							l.errChannel <- fmt.Errorf("publish message: %w", err)
							continue
						}

						l.log.WithFields(logrus.Fields{
							"subject": subjectName,
							"action":  event.Action,
							"table":   event.Table,
							"lsn":     l.readLSN(),
						}).Infoln("event was sent")
					}

					tx.Clear()
				}

				if msg.WalMessage.WalStart > l.readLSN() {
					if err = l.AckWalMessage(msg.WalMessage.WalStart); err != nil {
						l.errChannel <- fmt.Errorf("acknowledge wal message: %w", err)
						continue
					}

					l.log.WithField("lsn", l.readLSN()).Debugln("ack wal msg")
				}
			}

			if msg.ServerHeartbeat != nil {
				//FIXME panic if there have been no messages for a long time.
				l.log.WithFields(logrus.Fields{
					"server_wal_end": msg.ServerHeartbeat.ServerWalEnd,
					"server_time":    msg.ServerHeartbeat.ServerTime,
				}).Debugln("received server heartbeat")

				if msg.ServerHeartbeat.ReplyRequested == 1 {
					l.log.Debugln("status requested")

					if err = l.SendStandbyStatus(); err != nil {
						l.errChannel <- fmt.Errorf("send standby status: %w", err)
					}
				}
			}
		}
	}
}

// Stop is a finalizer function.
func (l *Listener) Stop() error {
	if err := l.repository.Close(); err != nil {
		return fmt.Errorf("repository close: %w", err)
	}

	if err := l.replicator.Close(); err != nil {
		return fmt.Errorf("replicator close: %w", err)
	}

	l.log.Infoln("service was stopped")

	return nil
}

// SendPeriodicHeartbeats send periodic keep alive heartbeats to the server.
func (l *Listener) SendPeriodicHeartbeats(ctx context.Context) {
	//todod : add HeartbeatInterval
	heart := time.NewTicker(time.Duration(100))
	defer heart.Stop()

	for {
		select {
		case <-ctx.Done():
			l.log.WithField("func", "SendPeriodicHeartbeats").
				Infoln("context was canceled, stop sending heartbeats")

			return
		case <-heart.C:
			{
				if err := l.SendStandbyStatus(); err != nil {
					l.log.WithError(err).Errorln("failed to send status heartbeat")
					continue
				}

				l.log.Debugln("sending periodic status heartbeat")
			}
		}
	}
}

// SendStandbyStatus sends a `StandbyStatus` object with the current RestartLSN value to the server.
func (l *Listener) SendStandbyStatus() error {
	standbyStatus, err := pgx.NewStandbyStatus(l.readLSN())
	if err != nil {
		return fmt.Errorf("unable to create StandbyStatus object: %w", err)
	}

	standbyStatus.ReplyRequested = 0

	if err := l.replicator.SendStandbyStatus(standbyStatus); err != nil {
		return fmt.Errorf("unable to send StandbyStatus object: %w", err)
	}

	return nil
}

// AckWalMessage acknowledge received wal message.
func (l *Listener) AckWalMessage(lsn uint64) error {
	l.setLSN(lsn)

	if err := l.SendStandbyStatus(); err != nil {
		return fmt.Errorf("send status: %w", err)
	}

	return nil
}

func (l *Listener) readLSN() uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	return l.lsn
}

// slotIsExists checks whether a slot has already been created and if it has been created uses it.
func (l *Listener) slotIsExists(ctx context.Context) (bool, error) {
	restartLSNStr, err := l.repository.GetSlotLSN(l.slotName)
	if err != nil {
		return false, err
	}

	if len(restartLSNStr) == 0 {
		l.log.WithField("slot_name", l.slotName).Warningln("restart LSN not found")
		return false, nil
	}

	lsn, err := pgx.ParseLSN(restartLSNStr)
	if err != nil {
		return false, fmt.Errorf("parse lsn: %w", err)
	}

	l.setLSN(lsn)

	return true, nil
}

func publicationNames(publication string) string {
	return fmt.Sprintf(`publication_names '%s'`, publication)
}

func (l *Listener) setLSN(lsn uint64) {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.lsn = lsn
}
