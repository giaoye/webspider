package mongo

import (
  "context"
  "go.mongodb.org/mongo-driver/mongo"
  "go.mongodb.org/mongo-driver/mongo/options"
  "sync"
  "time"
)

var (
  once sync.Once
  client *mongo.Client
  
  ErrNotDocs error = mongo.ErrNoDocuments
)

// TODO: 改为可配置
func NewClient() (*mongo.Client, error) {
  ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
  return mongo.Connect(ctx, options.Client().ApplyURI("mongodb://114.67.169.177:27017"))
}

func GetClient() (*mongo.Client, error) {
  var err error
  once.Do(func() {
    client, err = NewClient()
  })
  return client, err
}


func GetCollection(dbName, colName string) (*mongo.Collection, error) {
  client, err := GetClient()
  if err != nil {
    return nil, err
  }
  return client.Database(dbName).Collection(colName), nil
}

