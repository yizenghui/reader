package reader

import "testing"

func Test_GetContent(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := GetContent("http://www.aszw.org/book/38/38922/45061261.html")
	t.Fatal(a2)

}
