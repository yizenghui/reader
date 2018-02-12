package reader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/lunny/html2md"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/sundy-li/html2article"
)

// BookSection 小说段落 字数
type BookSection struct {
	Text string `json:"text"`
}

// BookInfo 链接
type BookInfo struct {
	Title   string        `json:"title"`
	Content []BookSection `json:"content"`
	PubAt   string        `json:"pub_at"`
}

// GetBookContent 获取小说正文，返回标题与段落内容
func GetBookContent(url string) (info BookInfo, err error) {

	type Article struct {
		// Basic
		Title       string `json:"title"`
		Content     string `json:"content"`
		Publishtime int64  `json:"publish_time"`
	}
	if url == "" {
		return info, errors.New("url不能为空")
	}

	ext, err := html2article.NewFromUrl(url)
	if err != nil {
		return info, err
	}
	article, err := ext.ToArticle()
	if err != nil {
		return info, err
	}
	// fmt.Println(article)

	//parse the article to be readability
	article.Readable(url)

	// fmt.Println(article.Title, article.Publishtime)
	// md = html2md.Convert(article.ReadContent)

	info.Title = article.Title

	md := html2md.Convert(article.ReadContent)
	input := []byte(md)
	unsafe := blackfriday.MarkdownCommon(input)
	content := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	c := strings.TrimSpace(fmt.Sprintf(`%v`, string(content[:])))

	c = strings.Replace(c, `<p>`, ``, -1)
	c = strings.Replace(c, `<code>`, ``, -1)
	c = strings.Replace(c, `</code>`, ``, -1)
	c = strings.Replace(c, `<pre>`, ``, -1)
	c = strings.Replace(c, `</pre>`, ``, -1)

	c = strings.Replace(c, `<br/>`, `</p>`, -1)
	c = strings.Replace(c, `<br />`, `</p>`, -1)
	c = strings.Replace(c, `<br>`, `</p>`, -1)
	c = strings.Replace(c, `\n`, `</p>`, -1)
	// arr := strings.Split(c, "<p>")
	arr := strings.Split(c, "</p>")

	if len(arr) > 2 {
		for _, v := range arr {
			text := strings.TrimSpace(v)
			if text != "" {
				info.Content = append(info.Content, BookSection{
					Text: text,
				})
			}
		}
	}
	// info.PubAt = Publishtime
	return info, nil

}

// GetBookMenu 获取小说目录
func GetBookMenu(urlStr string) (data Data, err error) {
	return GetList(urlStr)
}
