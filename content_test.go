package reader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/sundy-li/html2article"
)

func Test_GetContentBody(t *testing.T) {

	url := "http://www.longfu8.com/246.html"
	type Article struct {
		// Basic
		Title       string `json:"title"`
		Content     string `json:"content"`
		Publishtime int64  `json:"publish_time"`
	}

	ext, err := html2article.NewFromUrl(url)
	if err != nil {
	}
	article, err := ext.ToArticle()
	if err != nil {
	}
	// fmt.Println(article)

	//parse the article to be readability
	article.Readable(url)

	// fmt.Println(article.Title, article.Publishtime)
	// md = html2md.Convert(article.ReadContent)

	t.Fatal(article.ReadContent)

}
func Test_Demo(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := Read("http://book.zongheng.com/showchapter/523438.html")
	t.Fatal(a2)

}
func Test_Content(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := Read("http://www.166xs.com/xiaoshuo/84/84625/18347752.html")
	t.Fatal(a2)

}

func Test_SplitSection(t *testing.T) {

	urlStr := `http://www.longfu8.com/263.html`
	info, err := GetContent(urlStr)
	if err != nil {

	}

	input := []byte(info.Content)
	unsafe := blackfriday.MarkdownCommon(input)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	bh := fmt.Sprintf(`
			<html>
			<head>
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
			<title>%v</title>
			<body>%v</body>
			</html>
			`, info.Title, string(content[:]))

	g, e := goquery.NewDocumentFromReader(strings.NewReader(bh))
	if e != nil {

	}
	// html := fmt.Sprintf(`<meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
	// 						<link rel="preload" href="https://yize.gitlab.io/css/main.css" as="style" />
	// 						%v`, string(content[:]))
	// return c.HTML(http.StatusOK, html)
	info.Content = g.Find("body").Text()
	// info.Content = string(content[:])

	c := strings.TrimSpace(info.Content)

	arr := strings.Split(c, " ")
	t.Fatal(len(arr))

}
