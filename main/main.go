package main

import (
	"flag"
	"fmt"
	"net/url"

	p "github.com/bpetetot/bookie/providers"
)

func main() {

	// Parsing command line arguments
	var provider = flag.String("provider", "t411", "Provider name")
	var search = flag.String("search", "", "Searched book or manga")
	var typeSearch = flag.String("type", "volume", "Set volume or chapter")
	var number = flag.Int("number", 1, "Volume or Chapter number")
	flag.Parse()

	if *provider == "" || *search == "" {
		flag.PrintDefaults()
	}

	// Search values into providers
	if p.IsAvailableProvider(*provider) {
		fmt.Printf("Search '%s' with '%s' provider. \n", *search, *provider)
		var urlp = p.Init(*search)
		var body = p.Search(urlp)
		var results = p.ParseSearch(body)
		var uResults = p.AnalyzeSearch(*results, *typeSearch, *number)

		for _, r := range uResults {
			// Get details
			urlp, err := url.Parse("http:" + r.DetailURL)
			if err != nil {
				panic("Unable to parse 'T411' provider URL")
			}
			body = p.GetDetails(urlp)
			r.TorrentURL = p.ParseDetails(body)
			r.Display()

			// Download torrent file
			urlp, err = url.Parse("http://www.t411.in" + r.TorrentURL)
			if err != nil {
				panic("Unable to parse 'T411' provider URL")
			}
			err = p.DownloadTorrentFile(urlp, "D:\\temp\\")
			if err != nil {
				fmt.Printf("Error downloading torrent %s\n", err)
			}
		}

	} else {
		fmt.Printf("Provider '%s' not found. \n", *provider)
		p.DisplayProvidersList()
	}

}
