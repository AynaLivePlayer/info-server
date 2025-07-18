package service

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	"github.com/rhine-tech/scene/model/query"
	"github.com/spf13/cast"
	"infoserver/blivedm"
	"infoserver/blivedm/gen/fields"
	"infoserver/streamerstat"
)

type connlogImpl struct {
	repo       blivedm.ConnectionLogRepository    `aperture:""`
	streamStat streamerstat.IStreamerStatsService `aperture:""`
	log        logger.ILogger                     `aperture:""`
}

func (c *connlogImpl) Setup() error {
	return nil
}

func ConnectionLogService(repo blivedm.ConnectionLogRepository) blivedm.ConnectionLogService {
	return &connlogImpl{repo: repo}
}

func (c *connlogImpl) SrvImplName() scene.ImplName {
	return blivedm.Lens.ImplNameNoVer("ConnectionLogService")
}

func (c *connlogImpl) AddEntry(roomId int, source string, time int64) error {
	err := c.repo.AddEntry(roomId, source, time)
	if err != nil {
		c.log.ErrorW("Failed to add connection log entry", "error", err)
		return blivedm.ErrFailToAddEntry
	}
	return err
}

func (c *connlogImpl) ListEntries(offset, limit int64, filters ...query.Option) (model.PaginationResult[blivedm.ConnectionLog], error) {
	result, err := c.repo.GetEntries(offset, limit, filters...)
	if err != nil {
		c.log.ErrorW("Failed to get connection log entries", "error", err)
		return result, blivedm.ErrFailToListLogEntry
	}
	return result, nil
}

func (c *connlogImpl) ListEntryBySource(source string, offset, limit int64) (model.PaginationResult[blivedm.ConnectionLog], error) {
	result, err := c.repo.GetEntries(offset, limit, fields.ConnectionLog.Source.Equal(source))
	if err != nil {
		c.log.ErrorW("Failed to get connection log entry by source", "error", err)
		return result, blivedm.ErrFailToListLogEntry
	}
	return result, nil
}

func (c *connlogImpl) ListEntryByRoomID(roomID int, offset, limit int64) (model.PaginationResult[blivedm.ConnectionLog], error) {
	result, err := c.repo.GetEntries(offset, limit, fields.ConnectionLog.RoomID.Equal(roomID))
	if err != nil {
		c.log.ErrorW("Failed to get connection log entry by roomID", "error", err)
		return result, blivedm.ErrFailToListLogEntry
	}
	return result, nil
}

func (c *connlogImpl) ListEntriesByTimeRange(startTime, endTime int64, offset, limit int64) (model.PaginationResult[blivedm.ConnectionLog], error) {
	result, err := c.repo.GetEntries(offset, limit, fields.ConnectionLog.Time.GreaterOrEqual(startTime), fields.ConnectionLog.Time.LessOrEqual(endTime))
	if err != nil {
		c.log.ErrorW("Failed to get connection log entry by time range", "error", err)
		return result, blivedm.ErrFailToListLogEntry
	}
	return result, nil
}

func (c *connlogImpl) GetRoomLog(offset int64, limit int64) (model.PaginationResult[model.JsonResponse], error) {
	r, err := c.repo.GetRoomLog(offset, limit)
	if err != nil {
		return r, blivedm.ErrFailToListLogEntry
	}
	roomIds := make([]string, 0)
	for _, val := range r.Results {
		roomIds = append(roomIds, cast.ToString(val["room_id"]))
	}
	batch, err := c.streamStat.GetStatusBatch("bilibili", roomIds)
	if err == nil {
		for i, _ := range r.Results {
			r.Results[i]["status"] = batch[i]
		}
	}
	return r, nil
}
