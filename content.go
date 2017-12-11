package reader

import (
	"errors"
	"fmt"

	"github.com/sundy-li/html2article"
)

// Info 链接
type Info struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	PubAt   string `json:"pub_at"`
}

// GetContent 读url中的正文 解释返回 markdown 格式正文
func GetContent(url string) (info Info, err error) {

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
	fmt.Println(article)

	//parse the article to be readability
	article.Readable(url)

	// fmt.Println(article.Title, article.Publishtime)
	// md = html2md.Convert(article.ReadContent)

	info.Title = article.Title
	info.Content = article.ReadContent
	// info.PubAt = Publishtime
	return info, nil

}
