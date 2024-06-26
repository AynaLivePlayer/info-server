package storage

import (
	"errors"
	"github.com/rhine-tech/scene/composition/orm"
	"gorm.io/gorm"
	"infoserver/streamerstat"
)

type streamStatusImpl struct {
	db orm.Gorm `aperture:""`
}

func GormRepository(db orm.Gorm) streamerstat.StreamerStatusRepository {
	return &streamStatusImpl{db: db}
}

func (v *streamStatusImpl) Setup() error {
	err := v.db.RegisterModel(&streamerstat.StreamerStatus{})
	return err
}

func (v *streamStatusImpl) UpdateStatus(status *streamerstat.StreamerStatus) error {
	var s streamerstat.StreamerStatus
	err := v.db.DB().Where("platform = ? AND room_id = ?", status.Platform, status.RoomId).First(&s).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return v.db.DB().Create(status).Error
		}
		return err
	}
	return v.db.DB().Where("platform = ? AND room_id = ?", status.Platform, status.RoomId).Save(status).Error
}

func (v *streamStatusImpl) GetStatus(platform string, roomId string) (*streamerstat.StreamerStatus, error) {
	var status streamerstat.StreamerStatus
	err := v.db.DB().Where("platform = ? AND room_id = ?", platform, roomId).First(&status).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, streamerstat.ErrStatusNotFound
	}
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func (v *streamStatusImpl) GetStatusBatch(platform string, roomIds []string) ([]*streamerstat.StreamerStatus, error) {
	var statuses []*streamerstat.StreamerStatus
	err := v.db.DB().Where("platform = ? AND room_id IN ?", platform, roomIds).Find(&statuses).Error
	if err != nil {
		return nil, err
	}
	return statuses, nil
}
