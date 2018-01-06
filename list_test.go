package reader

import (
	"fmt"
	"net/url"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
)

func Test_GetList(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := GetList("http://book.zongheng.com/showchapter/523438.html")
	t.Fatal(a2)

}

func Test_GetList2(t *testing.T) {

	a2, _ := GetContent("http://book.zongheng.com/showchapter/523438.html")

	bh := fmt.Sprintf(`
			<html>
			<head>
			<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
			<title>龙符无弹窗,龙符最新章节全文阅读,梦入神机的小说-纵横中文网</title>
			<body>
			%v
			</body>
			</html>
			`, a2)
	c := strings.NewReader(bh)
	t.Fatal(c)

}

func Test_GetTag(t *testing.T) {

	s := GetTag("http://book.zongheng.com/showchapter/523438.html")

	var edu = map[string]int{}
	for _, k := range strings.Split(s, ",") {
		if v, ok := edu[k]; ok && k != "" {
			v++
			edu[k] = v
		} else {
			edu[k] = 1
		}
	}
	t.Fatal(edu)

}

func Test_Cleaning(t *testing.T) {

	urlStr := `http://www.longfuxiaoshuo.com/`
	// urlStr = `http://www.17k.com/list/2657648.html`
	// urlStr = `http://www.3dllc.cc/html/86/86244/`
	g, e := goquery.NewDocument(urlStr)
	if e != nil {

	}
	link, _ := url.Parse(urlStr)

	var links []Link
	// fmt.Println(g.Text())
	g.Find("a").Each(func(i int, content *goquery.Selection) {
		// 书名
		n := strings.TrimSpace(content.Text())
		u, _ := content.Attr("href")
		if strings.Index(u, "/") == 0 && strings.Index(u, "//") != 0 {
			u = fmt.Sprintf(`%v://%v%v`, link.Scheme, link.Host, u)
		}
		links = append(links, Link{
			n,
			u,
		})
	})
	edu := Cleaning(links)
	t.Fatal(edu)
}
