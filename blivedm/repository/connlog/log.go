package connlog

import (
	"github.com/rhine-tech/scene"
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/infrastructure/logger"
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
