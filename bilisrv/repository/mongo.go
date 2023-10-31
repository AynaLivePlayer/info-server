package repository

import (
	"context"
	"github.com/rhine-tech/scene/drivers/repos"
	"github.com/rhine-tech/scene/model"
	"github.com/spf13/cast"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"infoserver/bilisrv"
)

type mongoImpl struct {
	*repos.MongoCollectionRepo[bilisrv.BiliCredential]
}

func (m *mongoImpl) RepoImplName() string {
	return "bilisrv.repository.mongo"
}

func NewBiliCredentialMongoRepo(cfg model.DatabaseConfig) bilisrv.BiliCredentialRepo {
	repo := &mongoImpl{
		MongoCollectionRepo: repos.NewMongoCollectionRepo[bilisrv.BiliCredential](cfg, "bilibili_credentials"),
	}
	return repo
}

func (m *mongoImpl) Size() int {
	count, _ := m.Collection().CountDocuments(context.Background(), bson.M{})
	return cast.ToInt(count)
}

func (m *mongoImpl) Get(uid int) (bilisrv.BiliCredential, bool) {
	credential, err := m.FindOne(bson.M{"uid": uid})
	return credential, err == nil
}

func (m *mongoImpl) GetRandom() (bilisrv.BiliCredential, bool) {
	pipeline := bson.A{bson.M{"$sample": bson.M{"size": 1}}}

	cursor, err := m.Collection().Aggregate(context.Background(), pipeline)
	if err != nil {
		return bilisrv.BiliCredential{}, false
	}
	defer cursor.Close(context.Background())

	var credentials []bilisrv.BiliCredential
	if err := cursor.All(context.Background(), &credentials); err != nil {
		return bilisrv.BiliCredential{}, false
	}

	if len(credentials) == 0 {
		return bilisrv.BiliCredential{}, false
	}

	return credentials[0], true
}

func (m *mongoImpl) Upsert(credential bilisrv.BiliCredential) error {
	_, err := m.Collection().UpdateOne(
		context.Background(),
		bson.M{"uid": credential.UID},
		bson.M{"$set": credential},
		options.Update().SetUpsert(true),
	)
	return err
}

func (m *mongoImpl) Delete(uid int) error {
	val, err := m.Collection().DeleteOne(context.Background(), bson.M{"uid": uid})
	if err != nil {
		return err
	}
	if val.DeletedCount == 0 {
		return bilisrv.ErrNoSuchCredential.WithDetailStr(cast.ToString(uid))
	}
	return nil
}

func (m *mongoImpl) ListCredentials(offset int64, limit int64) ([]bilisrv.BiliCredential, int) {
	results, i, err := m.FindPagination(bson.M{}, bson.M{}, offset, limit)
	if err != nil || i == 0 {
		return []bilisrv.BiliCredential{}, 0
	}
	return results, i
}
