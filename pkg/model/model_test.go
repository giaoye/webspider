package model

import (
  "context"
  "strconv"
  "testing"
  "webspider/pkg/snowflake"
  "webspider/pkg/store/mongo"
  "webspider/pkg/store/redis"
)

var redisClient = redis.GetClient()

func TestMovie_Insert(t *testing.T) {
  node, err := snowflake.GetNode()
  if err != nil {
    t.Fatal("GetNode error:", err)
    return
  }
  c := Category{
    Name: "动作片",
  }
  err = c.Find(context.Background())
  var categoryId int
  if err == mongo.ErrNotDocs {
    res := redisClient.Get("category:maxId")
    if res.Err() != nil && res.Err() != redis.Nil {
      t.Error("[redis]Get:", err)
      return
    }
    redisClient.Incr("category:maxId")
    categoryId, err = strconv.Atoi(res.Val())
    if err != nil {
      t.Error("Atoi:", err)
      return
    }
  } else if err != nil {
    t.Error("[mongo]Find:", err)
    return
  } else {
    categoryId = int(c.ID)
  }
  movie := Movie{
    ID: int64(node.Generate()),
    Name: "天下第一",
    CategoryID: int64(categoryId),
    SummaryID: int64(node.Generate()),
  }
  err = movie.Insert(context.Background())
  if err != nil {
    t.Error(err)
  }
}
