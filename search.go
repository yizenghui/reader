package reader

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)




// Find 返回小说章节目录地址及所收录最新章节
func Find(bookName string) {

	url := fmt.Sprintf("http://www.sodu.cc/result.html?searchstr=%v", bookName)

	g, err := Get(url)
	if err != nil {
		// fmt.Println(e)
		// return job, err
	}
	//

	var menu string

	g.Find("a").Each(func(i int, content *goquery.Selection) {
		// 书名

		n := strings.TrimSpace(content.Text())
		u, _ := content.Attr("href")

		if n == bookName {
			menu = u
		}
	})

	d, err := Get(menu)

	list := map[string]string{}

	d.Find(".main-html").Each(func(i int, content *goquery.Selection) {
		// 书名
		f := strings.TrimSpace(content.Find("a").Eq(1).Text())
		if _, ok := list[f]; !ok {
			m, _ := content.Find("a").Eq(1).Attr("href")
			n := strings.TrimSpace(content.Find("a").Eq(0).Text())
			u, _ := content.Find("a").Eq(0).Attr("href")
			list[f] = f
			fmt.Println(f, m, n, u)
		}
	})

	fmt.Println(menu)
}
