package reader

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

//Get client  html
func Get(url string) (*goquery.Document, error) {

	client := &http.Client{}

	reqest, err := http.NewRequest("GET", url, nil)

	if err != nil {
		//	 handle error
		return nil, err

	}
	reqest.Header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 9_1 like Mac OS X) AppleWebKit/601.1.46 (KHTML, like Gecko) Version/9.0 Mobile/13B143 Safari/601.1")
	//	time.Sleep(10 * time.Second)
	// fmt.Println(mobile_info_url)

	resp, err := client.Do(reqest)

	if err != nil {
		//	 handle error
		return nil, err
		
	}
	g, err := goquery.NewDocumentFromReader(resp.Body)
	resp.Body.Close()
	return g, err
}
