package mongo

import (
	"bluebell/settings"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var mdb *mongo.Database

func Init(cfg *settings.MongoConfig) (err error) {
	dsn := fmt.Sprintf("mongodb://%s:%d", cfg.Host, cfg.Port)
	var mdbClient *mongo.Client
	mdbClient, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(dsn))
	if err != nil {
		zap.L().Error("connect mongo error", zap.Error(err))
		return
	}
	mdb = mdbClient.Database(cfg.DbName)
	fmt.Println(mdb.Name())
	return
}
