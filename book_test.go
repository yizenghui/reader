package reader

import "testing"

func Test_BookSplitSection(t *testing.T) {

	urlStr := `http://www.longfu8.com/263.html`
	info, err := GetBookContent(urlStr)
	if err != nil {

	}
	t.Fatal(info)

}
