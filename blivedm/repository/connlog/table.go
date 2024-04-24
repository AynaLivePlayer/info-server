package connlog

import "infoserver/blivedm"

type tableLog struct {
	ID     int    `gorm:"column:id;primaryKey;autoIncrement"`
	RoomID int    `gorm:"column:room_id"`
	Source string `gorm:"column:source"`
	Time   int64  `gorm:"column:time"`
}

func (tableLog) TableName() string {
	return blivedm.Lens.TableName("connection_log")
}
