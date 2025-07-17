package blivedm

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/model"
	"github.com/rhine-tech/scene/model/query"
)

type ConnectionLog struct {
	RoomID int    `json:"room_id" bson:"room_id"`
	Source string `json:"source" bson:"source"`
	Time   int64  `json:"time" bson:"time"`
}

type ConnectionLogRepository interface {
	scene.Named
	AddEntry(roomId int, source string, time int64) error
	GetEntries(offset int64, limit int64, options ...query.Option) (result model.PaginationResult[ConnectionLog], err error)
	GetRoomLog(offset int64, limit int64) (model.PaginationResult[model.JsonResponse], error)
}

type ConnectionLogService interface {
	scene.Service
	AddEntry(roomId int, source string, time int64) error
	ListEntries(offset, limit int64, options ...query.Option) (model.PaginationResult[ConnectionLog], error)
	GetRoomLog(offset int64, limit int64) (model.PaginationResult[model.JsonResponse], error)
	ListEntryBySource(source string, offset, limit int64) (model.PaginationResult[ConnectionLog], error)
	ListEntryByRoomID(roomID int, offset, limit int64) (model.PaginationResult[ConnectionLog], error)
	ListEntriesByTimeRange(startTime, endTime int64, offset, limit int64) (model.PaginationResult[ConnectionLog], error)
}
