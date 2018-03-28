package reader

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sundy-li/html2article"
)

// Data 链接
type Data struct {
	// Basic
	Title string `json:"title"`
	Links []Link
}

// Link 链接
type Link struct {
	// Basic
	Title string `json:"title"`
	URL   string `json:"url"`
}

var exp = []string{
	`?`,
	`&`,
	`#`,
	`/`,
	`=`,
	// `.`, // 不能把这个点去掉
}

//CheckStrIsLink 检查字符串是否支持的链接
func CheckStrIsLink(urlStr string) error {

	link, err := url.Parse(urlStr)

	if err != nil {
		return err
	}

	if link.Scheme == "" {
		return errors.New("Scheme Fatal")
	}

	if link.Host == "" {
		return errors.New("Host Fatal")
	}
	return nil
}

// GetList 获取列表，过滤零散的链接 (适用小说类)
func GetList(urlStr string) (data Data, err error) {
	err = CheckStrIsLink(urlStr)
	if err != nil {
		return
	}

	resp, err := http.Get(urlStr)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bs, _ := ioutil.ReadAll(resp.Body)
	htmlStr := string(bs)
	htmlStr = html2article.DecodeHtml(resp.Header, htmlStr, htmlStr)

	g, e := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))

	// g, e := goquery.NewDocument(urlStr)
	if e != nil {
		return data, e
	}

	data.Title = g.Find("title").Text()

	link, _ := url.Parse(urlStr)

	var links []Link
	// fmt.Println(g.Text())
	g.Find("a").Each(func(i int, content *goquery.Selection) {
		// 书名
		n := strings.TrimSpace(content.Text())
		u, _ := content.Attr("href")

		if strings.Index(u, "java") != 0 {
			if strings.Index(u, "//") == 0 {
				u = fmt.Sprintf(`%v:%v`, link.Scheme, u)
			} else if strings.Index(u, "/") == 0 {
				u = fmt.Sprintf(`%v://%v%v`, link.Scheme, link.Host, u)
			} else if strings.Index(u, "#") != 0 && strings.Index(u, "http") != 0 {
				u = fmt.Sprintf(`%v://%v%v%v`, link.Scheme, link.Host, link.Path, u)
			}
			u = strings.Replace(u, " ", "", -1)
			u = strings.Replace(u, "　", "", -1)
			// 去除换行符
			u = strings.Replace(u, "\n", "", -1)
			u = strings.Replace(u, "\t", "", -1)

			links = append(links, Link{
				n,
				u,
			})
		}
	})
	// fmt.Println(data)
	// data.Links = links
	data.Links = Cleaning(links)

	if len(data.Links) < 20 { // 这里面是兼容处理，如果
		exp = []string{
			`?`,
			`&`,
			`#`,
			`/`,
			`=`,
			`.`, // 不能把这个点去掉
		}
		// data.Links = Cleaning(links)
	}
	// log.Fatal(data.Links)
	return data, nil

}

// GetListByContent 获取正文中的链接
func GetListByContent(urlStr string) (data Data, err error) {
	a2, _ := GetContent(urlStr)

	data.Title = a2.Title
	bh := fmt.Sprintf(`
		<html>
		<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
		<title>%v</title>
		<body>
		%v
		</body>
		</html>
		`, a2.Title, a2.Content)
	c := strings.NewReader(bh)

	g, e := goquery.NewDocumentFromReader(c)

	if e != nil {

	}
	link, _ := url.Parse(urlStr)

	// fmt.Println(g.Text())
	g.Find("a").Each(func(i int, content *goquery.Selection) {
		// 书名
		n := strings.TrimSpace(content.Text())
		u, _ := content.Attr("href")
		if strings.Index(u, "/") == 0 && strings.Index(u, "//") != 0 {
			u = fmt.Sprintf(`%v://%v%v`, link.Scheme, link.Host, u)
		}
		data.Links = append(data.Links, Link{
			n,
			u,
		})
	})
	return data, nil

}

//Cleaning 清洗数据
func Cleaning(links []Link) (newlinks []Link) {
	// 拆分链接字符占比重
	var edu = map[string]int{}
	for _, link := range links { //所有链接
		s := GetTag(link.URL)
		for _, k := range strings.Split(s, ",") { //链接分拆统计
			if v, ok := edu[k]; ok && k != "" && k != " " {
				v++
				edu[k] = v
			} else {
				edu[k] = 1
			}
		}
	}

	var mw = 0
	var maxWeight = 0.0

	// log.Fatal(edu)
	for _, v := range edu {
		if v > 10 {
			v = 10
		}
		mw += v
	}

	// 找出最大重量
	for _, link := range links {
		s := GetTag(link.URL)
		w := 0
		for _, k := range strings.Split(s, ",") {
			if v, ok := edu[k]; ok {
				w += v
			}
		}
		if (float64(w) / float64(mw)) > maxWeight {
			maxWeight = float64(w) / float64(mw)
		}
		// wg[link.URL] = w
	}
	var pro = maxWeight * 0.30
	// 这个链接的重量
	var wg = map[string]int{}
	for _, link := range links {
		s := GetTag(link.URL)
		w := 0
		for _, k := range strings.Split(s, ",") {
			if v, ok := edu[k]; ok {
				w += v
			}
		}
		if float64(w) > (float64(mw) * float64(pro)) {
			wg[link.URL] = w
		}
		// wg[link.URL] = w
	}

	// log.Fatal(links)
	var crp = map[string]int{}
	for _, link := range links {
		if _, ok := crp[link.URL]; !ok && link.Title != "" {
			crp[link.URL] = 1
			if _, ok := wg[link.URL]; ok && link.Title != "" {
				newlinks = append(newlinks, link)
			}
		}

	}
	return newlinks
}

//GetTag 获取特点
func GetTag(urlStr string) string {

	// todo .htm .html .shtml 转换成非点分割
	// urlStr = strings.Replace(urlStr, `.html`, `/html`, -1)
	// urlStr = strings.Replace(urlStr, `.shtml`, `/shtml`, -1)
	// urlStr = strings.Replace(urlStr, `.htm`, `/htm`, -1)

	link, _ := url.Parse(urlStr)

	// link.Path =
	for _, t := range exp {
		// u := fmt.Sprintf(`%v`, link.Path)
		link.Path = strings.Replace(link.Path, t, ",", -1)
	}
	return link.Path
}
