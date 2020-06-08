package repo

import (
  "context"
  "crypto/md5"
  "encoding/hex"
  "strconv"
  "webspider/pkg/model"
  "webspider/pkg/snowflake"
  "webspider/pkg/store/redis"
)

func InsertMovie(movie *model.Movie, summary *model.Summary, category *model.Category,
  m3u8Link []*model.Link, onlineLink []*model.Link, downLink []*model.Link) error {
  node, err := snowflake.GetNode()
  if err != nil {
    return err
  }
  cache := redis.GetClient()
  movieId := node.Generate()
  summaryId := node.Generate()
  movie.ID = int64(movieId)
  movie.SummaryID = int64(summaryId)
  summary.ID = int64(summaryId)
  
  categoryKey := getCategoryKey(category.Name)
  exist, err := cache.Exists(categoryKey).Result()
  if err != nil && err != redis.Nil {
    return err
  }
  var categoryId int64
  if exist == 0 {
    cntKey := getCategoriesCntKey()
    _, err = cache.Get(cntKey).Result()
    if err != redis.Nil && err != nil {
      return err
    }
    value, err := cache.Incr(cntKey).Result()
    if err != nil {
      return err
    }
    categoryId = value
    cache.Set(getCategoryKey(category.Name), value, 0)
  } else {
    cnt, err := cache.Get(categoryKey).Result()
    if err != nil {
      return err
    }
    id, err := strconv.Atoi(cnt)
    if err != nil {
      return err
    }
    categoryId = int64(id)
  }
  movie.CategoryID = int64(categoryId)
  category.ID = int64(categoryId)
  for _, link := range m3u8Link {
    link.ID = int64(node.Generate())
    movie.LinkIds = append(movie.LinkIds, link.ID)
    link.Insert(context.Background())
  }
  for _, link := range onlineLink {
    link.ID = int64(node.Generate())
    movie.LinkIds = append(movie.LinkIds, link.ID)
    link.Insert(context.Background())
  }
  for _, link := range downLink {
    link.ID = int64(node.Generate())
    movie.LinkIds = append(movie.LinkIds, link.ID)
    link.Insert(context.Background())
  }
  
  err = category.Upsert(context.Background())
  if err != nil {
    return err
  }
  err = summary.Insert(context.Background())
  if err != nil {
    return err
  }
  err = movie.Insert(context.Background())
  if err != nil {
    return err
  }
  return nil
}

func getCategoryKey(name string) string {
  md5New := md5.New()
  sum := md5New.Sum([]byte(name))
  return "categories:"+string(hex.EncodeToString(sum))
}

func getCategoriesCntKey() string {
  return "categoriescnt"
}