package reader

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

func Test_Demo(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := Read("http://book.zongheng.com/showchapter/523438.html")
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
