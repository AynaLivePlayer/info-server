package statusstore

import (
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/model/query"
	"scene-service/meowsic"
)

func NewGorm() orm.GenericRepository[meowsic.SourceStatus] {
	return orm.NewGormRepository[meowsic.SourceStatus](nil, map[query.Field]string{})
}
