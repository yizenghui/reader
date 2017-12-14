package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/lunny/html2md"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/yizenghui/reader"
)

type (
	//Stats 结构
	Stats struct {
		Uptime       time.Time      `json:"uptime"`
		RequestCount uint64         `json:"requestCount"`
		Statuses     map[string]int `json:"statuses"`
		mutex        sync.RWMutex
	}
)

//NewStats New Stats
func NewStats() *Stats {
	return &Stats{
		Uptime:   time.Now(),
		Statuses: make(map[string]int),
	}
}

// Process is the middleware function.
func (s *Stats) Process(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			c.Error(err)
		}
		s.mutex.Lock()
		defer s.mutex.Unlock()
		s.RequestCount++
		status := strconv.Itoa(c.Response().Status)
		s.Statuses[status]++
		return nil
	}
}

// Handle is the endpoint to get stats.
func (s *Stats) Handle(c echo.Context) error {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return c.JSON(http.StatusOK, s)
}

//GetContent 获取正文
func GetContent(c echo.Context) error {

	urlStr := c.QueryParam("url")

	info, err := reader.GetContent(urlStr)
	if err != nil {
		return c.String(http.StatusOK, "0")
	}

	input := []byte(info.Content)
	unsafe := blackfriday.MarkdownCommon(input)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	// html := fmt.Sprintf(`<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	// 						<link rel="preload" href="https://yize.gitlab.io/css/main.css" as="style" />
	// 						%v`, string(content[:]))
	// return c.HTML(http.StatusOK, html)
	info.Content = fmt.Sprintf(`%v`, string(content[:]))

	type Info struct {
		Title   string        `json:"title"`
		Content template.HTML `json:"content"`
		PubAt   string        `json:"pub_at"`
	}

	return c.JSON(http.StatusOK, Info{
		info.Title,
		template.HTML(info.Content),
		info.PubAt,
	})
}

//GetContent2 获取正文
func GetContent2(c echo.Context) error {
	urlStr := c.QueryParam("url")

	info, err := reader.GetContent(urlStr)
	if err != nil {
		return c.JSON(http.StatusOK, "0")
	}
	info.Content = html2md.Convert(info.Content)

	type Info struct {
		Title   string        `json:"title"`
		Content template.HTML `json:"content"`
		PubAt   string        `json:"pub_at"`
	}

	return c.JSON(http.StatusOK, info)
}

//GetList 获取列表
func GetList(c echo.Context) error {
	urlStr := c.QueryParam("url")
	if urlStr == "" {
		return c.JSON(http.StatusOK, "0")
	}
	links, _ := reader.GetList(urlStr)
	return c.JSON(http.StatusOK, links)
}

func main() {

	e := echo.New()
	e.Use(middleware.CORS())
	//-------------------
	// Custom middleware
	//-------------------
	// Stats
	s := NewStats()
	e.Use(s.Process)
	// 展示统计
	e.GET("/stats", s.Handle) // Endpoint to get stats

	e.GET("/list", GetList)
	e.GET("/show", GetContent)

	e.File("favicon.ico", "favicon.ico")
	// e.Logger.Fatal(e.Start(":8007"))
	e.Logger.Fatal(e.StartTLS(":1323", "cert.pem", "key.pem"))

}
