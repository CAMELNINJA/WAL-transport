package parser

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"time"

	"github.com/CAMELNINGA/cdc-postgres.git/internal/models"
	error_walListner "github.com/CAMELNINGA/cdc-postgres.git/pkg/error_walListner"
	"github.com/sirupsen/logrus"
)

// BinaryParser represent binary protocol parser.
type BinaryParser struct {
	byteOrder binary.ByteOrder
	msgType   byte
	buffer    *bytes.Buffer
}

// NewBinaryParser create instance of binary parser.
func NewBinaryParser(byteOrder binary.ByteOrder) *BinaryParser {
	return &BinaryParser{
		byteOrder: byteOrder,
	}
}

// ParseWalMessage parse postgres WAL message.
func (p *BinaryParser) ParseWalMessage(msg []byte, tx *models.WalTransaction) error {
	if len(msg) == 0 {
		return error_walListner.ErrEmptyWALMessage
	}

	p.msgType = msg[0]
	p.buffer = bytes.NewBuffer(msg[1:])

	switch p.msgType {
	case models.BeginMsgType:
		begin := p.getBeginMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"lsn": begin.LSN,
					"xid": begin.XID,
				}).
			Debugln("begin type message was received")

		tx.LSN = begin.LSN
		tx.BeginTime = &begin.Timestamp
	case models.CommitMsgType:
		commit := p.getCommitMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"lsn":             commit.LSN,
					"transaction_lsn": commit.TransactionLSN,
				}).
			Debugln("commit message was received")

		if tx.LSN > 0 && tx.LSN != commit.LSN {
			return fmt.Errorf("commit: %w", error_walListner.ErrMessageLost)
		}

		tx.CommitTime = &commit.Timestamp
	case models.OriginMsgType:
		logrus.Debugln("origin type message was received")
	case models.RelationMsgType:
		relation := p.getRelationMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"relation_id": relation.ID,
					"replica":     relation.Replica,
				}).
			Debugln("relation type message was received")

		if tx.LSN == 0 {
			return fmt.Errorf("commit: %w", error_walListner.ErrMessageLost)
		}

		rd := models.RelationData{
			Schema: relation.Namespace,
			Table:  relation.Name,
		}

		for _, rf := range relation.Columns {
			c := models.Column{
				Name:      rf.Name,
				ValueType: int(rf.TypeID),
				IsKey:     rf.Key,
			}
			rd.Columns = append(rd.Columns, c)
		}

		tx.RelationStore[relation.ID] = rd

	case models.TypeMsgType:
		logrus.Debugln("type message was received")
	case models.InsertMsgType:
		insert := p.getInsertMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"relation_id": insert.RelationID,
				}).
			Debugln("insert type message was received")

		action, err := tx.CreateActionData(
			insert.RelationID,
			nil,
			insert.NewRow,
			models.ActionKindInsert,
		)

		if err != nil {
			return fmt.Errorf("create action data: %w", err)
		}

		tx.Actions = append(tx.Actions, &action)
	case models.UpdateMsgType:
		upd := p.getUpdateMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"relation_id": upd.RelationID,
				}).
			Debugln("update type message was received")

		action, err := tx.CreateActionData(
			upd.RelationID,
			upd.OldRow,
			upd.NewRow,
			models.ActionKindUpdate,
		)
		if err != nil {
			return fmt.Errorf("create action data: %w", err)
		}

		tx.Actions = append(tx.Actions, &action)
	case models.DeleteMsgType:
		del := p.getDeleteMsg()

		logrus.
			WithFields(
				logrus.Fields{
					"relation_id": del.RelationID,
				}).
			Debugln("delete type message was received")

		action, err := tx.CreateActionData(
			del.RelationID,
			del.OldRow,
			nil,
			models.ActionKindDelete,
		)
		if err != nil {
			return fmt.Errorf("create action data: %w", err)
		}

		tx.Actions = append(tx.Actions, &action)
	default:
		return fmt.Errorf("%w : %s", error_walListner.ErrUnknownMessageType, []byte{p.msgType})
	}
	return nil
}

func (p *BinaryParser) getBeginMsg() models.Begin {
	return models.Begin{
		LSN:       p.readInt64(),
		Timestamp: p.readTimestamp(),
		XID:       p.readInt32(),
	}
}

func (p *BinaryParser) getCommitMsg() models.Commit {
	return models.Commit{
		Flags:          p.readInt8(),
		LSN:            p.readInt64(),
		TransactionLSN: p.readInt64(),
		Timestamp:      p.readTimestamp(),
	}
}

func (p *BinaryParser) getInsertMsg() models.Insert {
	return models.Insert{
		RelationID: p.readInt32(),
		NewTuple:   p.buffer.Next(1)[0] == models.NewTupleDataType,
		NewRow:     p.readTupleData(),
	}
}

func (p *BinaryParser) getDeleteMsg() models.Delete {
	return models.Delete{
		RelationID: p.readInt32(),
		KeyTuple:   p.charIsExists('K'),
		OldTuple:   p.charIsExists('O'),
		OldRow:     p.readTupleData(),
	}
}

func (p *BinaryParser) getUpdateMsg() models.Update {
	u := models.Update{}
	u.RelationID = p.readInt32()
	u.KeyTuple = p.charIsExists('K')
	u.OldTuple = p.charIsExists('O')
	if u.KeyTuple || u.OldTuple {
		u.OldRow = p.readTupleData()
	}

	u.OldTuple = p.charIsExists('N')
	u.NewRow = p.readTupleData()

	return u
}

func (p *BinaryParser) getRelationMsg() models.Relation {
	return models.Relation{
		ID:        p.readInt32(),
		Namespace: p.readString(),
		Name:      p.readString(),
		Replica:   p.readInt8(),
		Columns:   p.readColumns(),
	}
}

func (p *BinaryParser) readInt32() (val int32) {
	r := bytes.NewReader(p.buffer.Next(4))
	_ = binary.Read(r, p.byteOrder, &val)

	return
}

func (p *BinaryParser) readInt64() (val int64) {
	r := bytes.NewReader(p.buffer.Next(8))
	_ = binary.Read(r, p.byteOrder, &val)

	return
}

func (p *BinaryParser) readInt8() (val int8) {
	r := bytes.NewReader(p.buffer.Next(1))
	_ = binary.Read(r, p.byteOrder, &val)

	return
}

func (p *BinaryParser) readInt16() (val int16) {
	r := bytes.NewReader(p.buffer.Next(2))
	_ = binary.Read(r, p.byteOrder, &val)

	return
}

func (p *BinaryParser) readTimestamp() time.Time {
	ns := p.readInt64()

	return models.PostgresEpoch.Add(time.Duration(ns) * time.Microsecond)
}

func (p *BinaryParser) readString() (str string) {
	stringBytes, _ := p.buffer.ReadBytes(0)

	return string(bytes.Trim(stringBytes, "\x00"))
}

func (p *BinaryParser) readBool() bool {
	x := p.buffer.Next(1)[0]

	return x != 0
}

func (p *BinaryParser) charIsExists(char byte) bool {
	if p.buffer.Next(1)[0] == char {
		return true
	}
	_ = p.buffer.UnreadByte()

	return false
}

func (p *BinaryParser) readColumns() []models.RelationColumn {
	size := int(p.readInt16())
	data := make([]models.RelationColumn, size)

	for i := 0; i < size; i++ {
		data[i] = models.RelationColumn{
			Key:          p.readBool(),
			Name:         p.readString(),
			TypeID:       p.readInt32(),
			ModifierType: p.readInt32(),
		}
	}

	return data
}

func (p *BinaryParser) readTupleData() []models.TupleData {
	size := int(p.readInt16())
	data := make([]models.TupleData, size)

	for i := 0; i < size; i++ {
		sl := p.buffer.Next(1)

		switch sl[0] {
		case models.NullDataType:
			logrus.Debugln("tupleData: null data type")
		case models.ToastDataType:
			logrus.Debugln("tupleData: toast data type")
		case models.TextDataType:
			vSize := int(p.readInt32())
			data[i] = models.TupleData{Value: p.buffer.Next(vSize)}
		}
	}

	return data
}
