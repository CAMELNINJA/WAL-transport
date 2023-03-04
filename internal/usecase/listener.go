package usecase

import (
	"context"

	"github.com/CAMELNINGA/cdc-postgres/internal/models"
	"github.com/jackc/pgx"
)

const errorBufferSize = 100

// Logical decoding plugin.
const pgOutputPlugin = "pgoutput"

type publisher interface {
	Publish(string, models.Event) error
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
