package reader

import "testing"

func Test_Demo(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := Read("http://book.zongheng.com/showchapter/523438.html")
	t.Fatal(a2)

}
