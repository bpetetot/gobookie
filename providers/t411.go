package providers

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"golang.org/x/net/html"
)

var header = map[string]string{
	"Content-Type": "application/x-www-form-urlencoded",
	"User-Agent":   "Mozilla/5.0 (Windows; U; Windows NT 5.1; en-US; rv:1.9.1.7) Gecko/20091221 Firefox/3.5.7 (.NET CLR 3.5.30729)",
}

var params = map[string]string{
	"name":        "",
	"search":      "",
	"description": "",
	"file":        "",
	"user":        "",
	"cat":         "404",
	//"subcat":      "409", // mangas
	"subcat": "406", // BD
	"submit": "Recherche",
	"order":  "added",
	"type":   "desc",
}

var cookies = []http.Cookie{
	{Name: "uid", Value: "7531997"},
	{Name: "pass", Value: "93b2d1460f1454ee349f52ac5bd611dcb9aa741d"},
	{Name: "authKey", Value: "2b3768e03113071bff127b26fb3c22ef"},
}

// Init the provider execution
func Init(search string) *url.URL {

	params["name"] = search
	params["search"] = "@name " + search

	// Build URL
	var urlp *url.URL
	urlp, err := url.Parse("http://www.t411.in")
	if err != nil {
		panic("Unable to parse 'T411' provider URL")
	}
	urlp.Path += "/torrents/search/"

	// Add parameters
	parameters := url.Values{}
	for k, p := range params {
		parameters.Add(k, p)
	}
	urlp.RawQuery = parameters.Encode()

	fmt.Printf("Encoded URL is %q\n", urlp.String())

	return urlp
}

// Search execute search to the provider
func Search(urlp *url.URL) io.Reader {
	fmt.Printf("Execute %q\n", urlp.String())

	// Create HTTP request
	client := http.Client{}
	request, err := http.NewRequest("GET", urlp.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Add headers
	for k, h := range header {
		request.Header.Add(k, h)
	}
	// Add Cookies
	for _, c := range cookies {
		request.AddCookie(&c)
	}
	// Execute request
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	//defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	return resp.Body
}

// ParseSearch body of the provider search request
func ParseSearch(body io.Reader) *[]Result {
	var results = make([]Result, 0)

	doc, err := html.Parse(body)
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			var found = false
			var result = Result{}
			for _, a := range n.Attr {

				if a.Key == "href" && strings.Contains(a.Val, "www.t411.in/torrents") {
					result.DetailURL = a.Val
					found = true
				}
				if a.Key == "title" && found {
					result.Name = strings.TrimSpace(a.Val)
					results = append(results, result)
				}

			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return &results
}

// AnalyzeSearch make an string analysis to identify results
func AnalyzeSearch(results []Result, typeSearch string, number int) []Result {
	var newResults = make([]Result, 0)
	for i := range results {
		result := &results[i]
		for _, analyzer := range Analyzers {
			analyzer.analyze(result)
			if (typeSearch == "volume" && result.IsVolume) || (typeSearch == "chapter" && result.IsChapter) {
				newResults = append(newResults, *result)
			}
		}
	}
	return newResults
}

// GetDetails execute search to get details
func GetDetails(urlp *url.URL) io.Reader {
	fmt.Printf("Execute %q\n", urlp.String())

	// Create HTTP request
	client := http.Client{}
	request, err := http.NewRequest("GET", urlp.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	// Add headers
	for k, h := range header {
		request.Header.Add(k, h)
	}
	// Add Cookies
	for _, c := range cookies {
		request.AddCookie(&c)
	}
	// Execute request
	resp, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}

	//defer resp.Body.Close()
	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	return resp.Body
}

// ParseDetails body of the provider search request
func ParseDetails(body io.Reader) string {
	doc, err := html.Parse(body)

	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node) string
	f = func(n *html.Node) string {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key == "href" && strings.Contains(a.Val, "/torrents/download/") {
					return a.Val
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			var val = f(c)
			if val != "" {
				return val
			}
		}
		return ""
	}
	return f(doc)
}

// DownloadTorrentFile execute search to get details
func DownloadTorrentFile(urlp *url.URL, dest string) error {
	fmt.Printf("Download Torrent File %q\n", urlp.String())

	// Create HTTP request
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	request, err := http.NewRequest("GET", urlp.String(), nil)
	if err != nil {
		return err
	}
	// Add headers
	for k, h := range header {
		request.Header.Add(k, h)
	}
	// Add Cookies
	for _, c := range cookies {
		request.AddCookie(&c)
	}
	// Execute request
	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	//fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

	var filename = resp.Header.Get("Content-Disposition")
	filename = strings.Replace(filename, "attachment; filename=", "", 1)
	filename = strings.Replace(filename, "\"", "", 2)

	defer resp.Body.Close()
	d, err := ioutil.ReadAll(resp.Body)

	if writeFile(dest+filename, d) == nil {
		fmt.Printf("saved %s as %s\n", urlp.String(), dest+filename)
	}

	return nil
}

func writeFile(dst string, d []byte) error {
	err := ioutil.WriteFile(dst, d, 0444)
	if err != nil {
		log.Fatal(err)
	}
	return err
}
