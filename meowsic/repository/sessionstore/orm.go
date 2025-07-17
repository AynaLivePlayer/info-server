package sessionstore

import (
	"github.com/rhine-tech/scene/composition/orm"
	"github.com/rhine-tech/scene/model/query"
	"scene-service/meowsic"
)

type gormImpl struct {
	db       orm.Gorm `aperture:""`
	internal *orm.GormRepository[meowsic.SourceSession]
}

func (g *gormImpl) Setup() error {
	g.internal = orm.NewGormRepository[meowsic.SourceSession](g.db, map[query.Field]string{})
	return g.internal.Setup()
}

func (g *gormImpl) GetSession(provider string) (string, bool, error) {
	val, ok, err := g.internal.FindFirst(query.Field("provider").Equal(provider))
	if err != nil {
		return "", ok, err
	}
	return val.Session, ok, err
}

func (g *gormImpl) RemoveSession(provider string) error {
	return g.internal.Delete(query.Field("provider").Equal(provider))
}

func (g *gormImpl) UpdateSession(provider string, session string) error {
	_, ok, _ := g.internal.FindFirst(query.Field("provider").Equal(provider))
	if ok {
		return g.internal.Update(map[string]interface{}{
			"session": session,
		}, query.Field("provider").Equal(provider))
	}
	return g.internal.Create(&meowsic.SourceSession{
		Session:  session,
		Provider: provider,
	})
}

func NewGorm() meowsic.ISessionStorage {
	return &gormImpl{}
}
