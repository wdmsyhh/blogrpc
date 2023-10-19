package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var MongoDatabase *mongo.Database

func main() {
	// 使用 mongo-driver 访问 GridFS
	mongoClient, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@127.0.0.1:27012"))
	if err != nil {
		log.Fatal(err)
	}
	defer mongoClient.Disconnect(context.Background())

	MongoDatabase = mongoClient.Database("testdb")

	err = GridFsUploadWithID("", "aaaaa", "name1", []byte("hhhhh"))
	if err != nil {
		fmt.Println(err.Error())
	}
}

func getGridFsBucket(collName string) *gridfs.Bucket {
	var bucket *gridfs.Bucket
	// 使用默认文件集合名称
	if collName == "" || collName == options.DefaultName {
		bucket, _ = gridfs.NewBucket(MongoDatabase)
	} else {
		// 使用传入的文件集合名称
		bucketOptions := options.GridFSBucket().SetName(collName)
		bucket, _ = gridfs.NewBucket(MongoDatabase, bucketOptions)
	}
	return bucket
}

// 上传文件
// collName:文件集合名称 fileID:文件ID，必须唯一，否则会覆盖
// fileName:文件名称 fileContent:文件内容
func GridFsUploadWithID(collName, fileID, fileName string, fileContent []byte) error {
	bucket := getGridFsBucket(collName)
	err := bucket.UploadFromStreamWithID(fileID, fileName, bytes.NewBuffer(fileContent))
	if err != nil {
		return err
	}
	return nil
}

// 下载文件
// 返回文件内容
func GridFsDownload(collName, fileID string) (fileContent []byte, err error) {
	bucket := getGridFsBucket(collName)
	fileBuffer := bytes.NewBuffer(nil)
	if _, err = bucket.DownloadToStream(fileID, fileBuffer); err != nil {
		return nil, err
	}
	return fileBuffer.Bytes(), nil
}

// 删除文件
func GridFsDelete(collName, fileID string) error {
	bucket := getGridFsBucket(collName)
	if err := bucket.Delete(fileID); err != nil && !errors.Is(err, gridfs.ErrFileNotFound) {
		return err
	}
	return nil
}
