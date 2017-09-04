package reader

import (
	"errors"

	"github.com/lunny/html2md"
	"github.com/sundy-li/html2article"
)

// Read 读url中的正文 解释返回 markdown 格式正文
func Read(url string) (md string, err error) {
	if url == "" {
		return "", errors.New("url不能为空")
	}

	ext, err := html2article.NewFromUrl(url)
	if err != nil {
		return "", err
	}
	article, err := ext.ToArticle()
	if err != nil {
		return "", err
	}
	//parse the article to be readability
	article.Readable(url)

	md = html2md.Convert(article.ReadContent)

	return md, nil

}
