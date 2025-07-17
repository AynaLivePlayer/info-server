package repository

import "infoserver/version"

type tableVersion struct {
	Version     uint32 `gorm:"primary_key"`
	Note        string `gorm:"column:note"`
	ReleaseDate int64  `gorm:"column:release_date"`
	FileID      string `gorm:"column:file_id"`
}

func (tableVersion) TableName() string {
	return version.Lens.TableName("version")
}
