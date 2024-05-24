package connlog

import (
	"github.com/rhine-tech/scene/model/query"
	"infoserver/blivedm"
	"infoserver/blivedm/gen/fields"
)

type tableLog struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID int    `gorm:"column:room_id"`
	Source string `gorm:"column:source"`
	Time   int64  `gorm:"column:time"`
}

func (tableLog) TableName() string {
	return blivedm.Lens.TableName("connection_log")
}

var fieldMapper = map[query.Field]string{
	fields.ConnectionLog.RoomID: (tableLog{}).TableName() + "." + "room_id",
	fields.ConnectionLog.Source: (tableLog{}).TableName() + "." + "source",
	fields.ConnectionLog.Time:   (tableLog{}).TableName() + "." + "time",
}
