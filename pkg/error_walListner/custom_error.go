package erorwallistner

import "errors"

// Variable with connection errors.
var (
	ErrReplConnectionIsLost     = errors.New("replication connection to postgres is lost")
	ErrConnectionIsLost         = errors.New("db connection to postgres is lost")
	ErrMessageLost              = errors.New("messages are lost")
	ErrEmptyWALMessage          = errors.New("empty WAL message")
	ErrUnknownMessageType       = errors.New("unknown message type")
	ErrRelationNotFound         = errors.New("relation not found")
	ErrNotConnectedKafaProducer = errors.New("not connected kafka producer")
	ErrNotConnectedKafaConsumer = errors.New("not connected kafka consumer")
)

type ServiceErr struct {
	Caller string
	Err    error
}

func NewListenerError(caller string, err error) *ServiceErr {
	return &ServiceErr{Caller: caller, Err: err}
}

func (e *ServiceErr) Error() string {
	return e.Caller + ": " + e.Err.Error()
}
