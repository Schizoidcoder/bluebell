package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
)

var ctx = context.Background()

func InsertOne(TableName string, Doc interface{}) error {
	coll := mdb.Collection(TableName)
	_, err := coll.InsertOne(ctx, Doc)
	return err
}

func InsertMany(TableName string, Docs []interface{}) error {
	coll := mdb.Collection(TableName)
	_, err := coll.InsertMany(ctx, Docs)
	return err
}

//Todo:update,delete
//func UpdateOne(TableName string, filter interface{}, ) error {
//	coll := mdb.Collection(TableName)
//	bson.D{}
//	coll.UpdateOne()
//}

func FindOne(TableName string, filter interface{}) (result bson.M, err error) {
	coll := mdb.Collection(TableName)
	res := coll.FindOne(context.TODO(), filter)
	//var result bson.M
	err = res.Decode(&result)
	if err != nil {
		return nil, err
	}
	//encode, err = encoder.Encode(result, encoder.SortMapKeys)
	return
}

func FindManyByOneCon(TableName string, filter interface{}) (results []bson.M, err error) {
	coll := mdb.Collection(TableName)
	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	//序列化
	//encodes = make([][]byte, len(results))
	//for _, doc := range results {
	//	encode, err := encoder.Encode(doc, encoder.SortMapKeys)
	//	if err != nil {
	//		return nil, err
	//	}
	//	encodes = append(encodes, encode)
	//}
	return
}
