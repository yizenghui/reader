package reader

import "testing"
import "fmt"

func Test_Demo(t *testing.T) {

	// a1, _ := Read("http://www.76wx.com/book/1563/3212972.html")
	// fmt.Println(a1)

	a2, _ := Read("http://www.aszw.org/book/38/38922/45061261.html")
	fmt.Println(a2)

}
