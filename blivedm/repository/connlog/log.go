package connlog

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/infrastructure/logger"
	"github.com/rhine-tech/scene/model"
	"github.com/rhine-tech/scene/model/filter"
	"infoserver/blivedm"
)

type logRepo struct {
	db  orm.Gorm       `aperture:""`
	log logger.ILogger `aperture:""`
}

func (l *logRepo) Setup() error {
	l.log = l.log.WithPrefix(l.RepoImplName().Identifier())
	return l.db.DB().AutoMigrate(&tableLog{})
}

func (l *logRepo) RepoImplName() scene.ImplName {
	return blivedm.Lens.ImplName("ConnectionLogRepository", "gorm")
}

func GormRepo(db orm.Gorm) blivedm.ConnectionLogRepository {
	return &logRepo{db: db}
}

func (l *logRepo) AddEntry(roomId int, source string, tim int64) error {
	record := tableLog{
		RoomID: roomId,
		Source: source,
		Time:   tim,
	}
	err := l.db.DB().Create(&record).Error
	if err != nil {
		l.log.ErrorW("Failed to create room connection log", "error", err)
	}
	return l.db.DB().Create(&record).Error
}

func (l *logRepo) GetEntries(offset int64, limit int64, filters ...filter.Filter) (result model.PaginationResult[blivedm.ConnectionLog], err error) {
	var records []tableLog
	err = l.db.DB().
		Where(l.db.WithFieldMapper(fieldMapper).BuildFilter(filters...)).
		Offset(int(offset)).Limit(int(limit)).
		Order("time desc").
		Find(&records).Error
	if err != nil {
		l.log.ErrorW("Failed to get room connection logs", "error", err)
		return model.PaginationResult[blivedm.ConnectionLog]{}, err
	}
	result.Results = make([]blivedm.ConnectionLog, len(records))
	for i, record := range records {
		result.Results[i] = blivedm.ConnectionLog{
			RoomID: record.RoomID,
			Source: record.Source,
			Time:   record.Time,
		}
	}
	l.db.DB().Model(&tableLog{}).Where(l.db.WithFieldMapper(fieldMapper).BuildFilter(filters...)).Count(&result.Total)
	result.Offset = offset
	result.Count = int64(len(records))
	return result, nil
}

func (l *logRepo) GetRoomLog(offset int64, limit int64) (result model.PaginationResult[model.JsonResponse], err error) {
	var records []map[string]interface{}
	//select room_id, count(*) as count, max(time) as last_time from blivedm_connection_log
	//	group by room_id
	//	order by last_time DESC
	query := l.db.DB().
		Select("room_id, count(*) as count, max(time) as last_time").
		Table("blivedm_connection_log").
		Group("room_id").
		Order("last_time DESC").
		Offset(int(offset)).Limit(int(limit))
	err = query.Find(&records).Error
	if err != nil {
		l.log.ErrorW("fail to get room connection log", "error", err)
		return model.PaginationResult[model.JsonResponse]{}, nil
	}
	query.Count(&result.Total)
	result.Offset = offset
	for _, r := range records {
		result.Results = append(result.Results, r)
	}
	result.Count = int64(len(records))
	return result, nil
}
