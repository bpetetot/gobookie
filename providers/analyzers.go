package providers

import "regexp"

// Analyzers available analyzers
var Analyzers = []Analyzer{
	VolumeAnalyzer{},
}

// Analyzer Analyzer interface
type Analyzer interface {
	analyze(Analyzed *Result) *Result
}

// VolumeAnalyzer Analyze volumes
type VolumeAnalyzer struct {
}

var volumesWords = "(Volume|Volumes|volume|volumes|Tome|Tomes|tome|tomes|TOME|TOMES|T[0-9]+|V[0-9]+|Vol[0-9]+|vol[0-9]+)"

func (v VolumeAnalyzer) analyze(analyzed *Result) *Result {
	matchVolume, _ := regexp.MatchString(volumesWords, analyzed.Name)
	if matchVolume {
		analyzed.IsVolume = true
		re, _ := regexp.Compile(`([0-9]+)`) // want to know what is in front of 'at'
		res := re.FindAllStringSubmatch(analyzed.Name, -1)
		var elements = make([]string, 0)
		for _, r := range res {
			elements = append(elements, r[0])
		}
		analyzed.RefNumbers = elements
	} else {
		analyzed.IsChapter = true
		re, _ := regexp.Compile(`([0-9]+)`) // want to know what is in front of 'at'
		res := re.FindAllStringSubmatch(analyzed.Name, -1)
		var elements = make([]string, 0)
		for _, r := range res {
			elements = append(elements, r[0])
		}
		analyzed.RefNumbers = elements
	}
	return analyzed
}
