package providers

import "fmt"

// Providers available
var Providers = map[string]Provider{
	"T411": {"http://www.t411.io", "html"},
}

// Provider contains informations and options about providers
type Provider struct {
	URL  string
	Type string
}

// Result of provider search
type Result struct {
	Name         string
	DetailURL    string
	TorrentURL   string
	IsVolume     bool
	IsChapter    bool
	RefNumbers   []string
	IsMultifiles bool
	Format       string
}

// Display display the result
func (r *Result) Display() {
	fmt.Println(r.Name)
	if r.IsVolume {
		fmt.Printf(" - Volume : %v\n", r.RefNumbers)
	} else if r.IsChapter {
		fmt.Printf(" - Chapter : %v\n", r.RefNumbers)
	}
	fmt.Println(" - URL D : ", r.DetailURL)
	fmt.Println(" - URL T : ", r.TorrentURL)
}

// IsAvailableProvider provider
func IsAvailableProvider(provider string) bool {
	_, exists := Providers[provider]
	return exists
}

// DisplayProvidersList liste all available providers in the output
func DisplayProvidersList() {
	fmt.Println("Available providers are :")
	for k := range Providers {
		fmt.Printf("- %s \n", k)
	}
}
