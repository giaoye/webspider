package model

import (
  "context"
  "go.mongodb.org/mongo-driver/bson"
  "go.mongodb.org/mongo-driver/mongo/options"
  "time"
  "webspider/pkg/store/mongo"
)

const (
  DATABASE string = "movie_store"
  MovieDetailsCol string = "movie_details"
  SummaryDetailsCol string = "summary_details"
  LinkDetailsCol string = "link_details"
  CategoryDetailsCol string = "category_details"
)

type Movie struct {
  ID         int64   `bson:"id"`          // ID
  Name       string  `bson:"name"`        // 电影名称
  SubTitle   string  `bson:"sub_title"`   // 副标题
  CategoryID int64   `bson:"category_id"` // 类型ID
  SummaryID  int64   `bson:"summary_id"`  // 简介ID
  LinkIds    []int64 `bson:"link_ids"`    // 链接ID列表
  CoverLink  string  `bson:"cover_link"`  // 封面
}

func (m Movie)Insert(ctx context.Context) (err error) {
  col, err := mongo.GetCollection(DATABASE, MovieDetailsCol)
  if err != nil {
    return err
  }
  _, err = col.InsertOne(ctx, &m)
  return err
}

type Link struct {
  ID   int64  `bson:"id"`   // ID
  Name string `bson:"name"` // 链接中文名称
  Type string `bson:"type"` // 链接类型：m3u8/xunlei/online
  Addr string `bson:"addr"` // 链接地址
}

func (l Link)Insert(ctx context.Context) (err error) {
  col, err := mongo.GetCollection(DATABASE, LinkDetailsCol)
  if err != nil {
    return err
  }
  _, err = col.InsertOne(ctx, &l)
  return err
}

type Summary struct {
  ID            int64     `bson:"id"`             // ID
  Directors     []string  `bson:"directors"`      // 导演
  Actors        []string  `bson:"actors"`         // 演员
  Region        string    `bson:"region"`         // 地区
  Lang          string    `bson:"language"`       // 语言
  UpdateTime    time.Time `bson:"update_time"`    // 更新时间
  Categories    string    `bson:"sub_categories"` // 类型
  ReleaseYear   string    `bson:"release_year"`   // 首映时间
  Content       string    `bson:"content"`        // 内容
  Duration      int       `bson:"duration"`       // 时长
}

func (s Summary)Insert(ctx context.Context) (err error) {
  col, err := mongo.GetCollection(DATABASE, SummaryDetailsCol)
  if err != nil {
    return err
  }
  _, err = col.InsertOne(ctx, &s)
  return err
}

type Category struct {
  ID   int64  `bson:"id"`   // ID
  Name string `bson:"name"` // 类型名
}

func (c Category)Insert(ctx context.Context) (err error) {
  col, err := mongo.GetCollection(DATABASE, CategoryDetailsCol)
  if err != nil {
    return err
  }
  _, err = col.InsertOne(ctx, &c)
  return err
}

func (c Category)Upsert(ctx context.Context) (err error) {
  col, err := mongo.GetCollection(DATABASE, CategoryDetailsCol)
  if err != nil {
    return err
  }
  opts := options.Update().SetUpsert(true)
  _, err = col.UpdateOne(ctx, bson.M{"name": c.Name}, bson.M{
    "$set": bson.M{
      "id": c.ID,
      "name": c.Name,
    },
  }, opts)
  return err
}

func (c Category)Find(ctx context.Context) (err error){
  client, _ := mongo.GetClient()
  err = client.Ping(ctx, nil)
  col, err := mongo.GetCollection(DATABASE, CategoryDetailsCol)
  if err != nil {
    return err
  }
  res := col.FindOne(ctx, bson.M{"name": c.Name})
  if res.Err() != nil {
    return res.Err()
  }
  return res.Decode(&c)
}