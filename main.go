package main

import (
  "fmt"
  "go.uber.org/zap"
  "strconv"
  "strings"
  "time"
  "webspider/pkg/log"
  "webspider/pkg/model"
  "webspider/pkg/repo"
)
import "github.com/gocolly/colly"

var logger = log.GetLogger()
const URL string = "http://www.zuidazy5.com"
func main() {
  c := colly.NewCollector()
  GetCategoryURLs(c)
  GetMovieURLs(c)
  GetMovieDetail(c)
  c.OnRequest(func(r *colly.Request) {
    logger.Info("Visiting", zap.String("url", r.URL.String()))
  })
  
  c.Visit(URL)
}


func GetCategoryURLs(c *colly.Collector) {
  c.OnHTML(`ul[id=sddm]`, func(e *colly.HTMLElement) {
    requestURL := e.Request.URL.String()
    if requestURL == URL { // 主页面
      fmt.Println("request:", e.Request.URL)
      e.ForEach("a", func(i int, element *colly.HTMLElement) {
        path := element.Attr("href")
        if path != "/" {
          fmt.Println(element.Text, URL+path)
          e.Request.Visit(URL+path) // 访问各电影类型列表页
        }
      })
    }
  })
}

func GetMovieURLs(c *colly.Collector) {
  c.OnHTML(`div[class=xing_vb]`, func(e *colly.HTMLElement) {
    requestURL := e.Request.URL.String()
    if strings.Contains(requestURL, "vod-type-id") { // 电影类型列表页
      e.ForEach("li", func(i int, element *colly.HTMLElement) {
        vb4 := element.DOM.Find("span[class=xing_vb4]")
        var movieName string
        var movieUrl string
        if vb4.Text() != "" {
          movieName = vb4.Text()
          movieUrl, _ = vb4.Find("a").Attr("href")
          fmt.Println("vb4:", movieName)
          fmt.Println(movieUrl)
        }
        vb5 := element.DOM.Find("span[class=xing_vb5]")
        var category string
        if vb5.Text() != "" {
          category = vb5.Text()
          fmt.Println("vb5:", category)
        }
        vb6 := element.DOM.Find("span[class=xing_vb6]")
        var updateTime string
        if vb6.Text() != "" {
          updateTime = vb6.Text()
          fmt.Println("vb6:", updateTime)
          time.Sleep(25*time.Millisecond)
          element.Request.Visit(URL+movieUrl)
        }
        pageLink := element.DOM.Find("a[class=pagelink_a]").First()
        if pageLink.Text() == "下一页" {
          fmt.Println("pageLink:", pageLink.Text())
          path, exist := pageLink.Attr("href")
          if exist { // 访问下一页
            element.Request.Visit(URL+path)
          }
        }
      })
    }
  })
}

func GetMovieDetail(c *colly.Collector) {
  c.OnHTML(`div[class=warp]`, func(e *colly.HTMLElement) {
    requestURL := e.Request.URL.String()
    if strings.Contains(requestURL, "vod-detail-id") { // 电影详情页
      movie := model.Movie{}
      summary := model.Summary{}
      category := model.Category{}
      m3u8Link := make([]*model.Link, 0)
      onlineLink := make([]*model.Link, 0)
      downLink := make([]*model.Link, 0)
      e.ForEach(`div`, func(i int, element *colly.HTMLElement) {
        //fmt.Println("[vodBox]class:", element.Attr("class"))
        if element.Attr("class") == "vodImg" {
          movie.CoverLink = element.ChildAttr("img", "src")
          //fmt.Println("coverUrl:", movie.CoverLink)
        }
        if element.Attr("class") == "vodInfo" {
          element.ForEach("div", func(i int, div *colly.HTMLElement) {
            if div.Attr("class") == "vodh" {
              movie.Name = div.ChildText("h2")
              movie.SubTitle = div.ChildText("span")
            }
            if div.Attr("class") == "vodinfobox" {
              div.ForEach("ul", func(i int, ul *colly.HTMLElement) {
                ul.ForEach("li", func(i int, li *colly.HTMLElement) {
                  //fmt.Println("li:", li.Text)
                  text := li.ChildText("span")
                  if strings.HasPrefix(li.Text, "导演") {
                    summary.Directors = append(summary.Directors, text)
                  }
                  if strings.HasPrefix(li.Text, "主演") {
                    summary.Actors = append(summary.Actors, text)
                  }
                  if strings.HasPrefix(li.Text, "类型") {
                    categories := strings.Split(text, " ")
                    category.Name = categories[0]
                    summary.Categories = text
                  }
                  if strings.HasPrefix(li.Text, "语言") {
                    summary.Lang = text
                  }
                  if strings.HasPrefix(li.Text, "地区") {
                    summary.Region = text
                  }
                  if strings.HasPrefix(li.Text, "上映") {
                    summary.ReleaseYear = text
                  }
                  if strings.HasPrefix(li.Text, "片长") {
                    dur, err := strconv.Atoi(text)
                    if err != nil {
                      return
                    }
                    summary.Duration = dur
                  }
                  if strings.HasPrefix(li.Text, "更新") {
                    updateTime, err := time.Parse("2006-01-02 15:04:05", text)
                    if err != nil {
                      return
                    }
                    summary.UpdateTime = updateTime
                  }
                  
                  if li.Attr("class") == "cont" {
                    content := strings.TrimSpace(text)
                    summary.Content = content
                  }
                })
              })
            }
          })
        }
        if element.Attr("id") == "play_1" { // m3u8播放链接
          element.ForEach("li", func(i int, li *colly.HTMLElement) {
            kv := strings.Split(li.Text, "$")
            m3u8Link = append(m3u8Link, &model.Link{
              Type: "m3u8",
              Name: kv[0],
              Addr: kv[1],
            })
            //fmt.Println("m3u8:", li.Text)
          })
        }
        if element.Attr("id") == "play_2" { // 在线播放链接
          element.ForEach("li", func(i int, li *colly.HTMLElement) {
            kv := strings.Split(li.Text, "$")
            onlineLink = append(onlineLink, &model.Link{
              Type: "online",
              Name: kv[0],
              Addr: kv[1],
            })
            //fmt.Println("online:", li.Text)
          })
        }
        if element.Attr("id") == "play_3" { // 下载链接
          element.ForEach("li", func(i int, li *colly.HTMLElement) {
            kv := strings.Split(li.Text, "$")
            downLink = append(downLink, &model.Link{
              Type: "download",
              Name: kv[0],
              Addr: kv[1],
            })
            //fmt.Println("download:", li.Text)
          })
        }
      })
      err := repo.InsertMovie(&movie, &summary, &category, m3u8Link, onlineLink, downLink)
      if err != nil {
        logger.Error("InsertMovie error", zap.Error(err),
          zap.Reflect("category:", category),
          zap.Reflect("summary:", summary),
          zap.Reflect("movie:", movie),
          zap.Reflect("m3u8", m3u8Link),
          zap.Reflect("online", onlineLink),
          zap.Reflect("download", downLink))
      }
    }
  })
}