package main

import (
	"fmt"
	"regexp"
)

func main() {

	re, _ := regexp.Compile(`(Volume|Tome|volume|tome|T[0-9]+|V[0-9]+|Vol[0-9]+|vol[0-9]+)`)
	r := re.FindStringSubmatch("[MFT] One Piece Scan T01-65 fr")
	fmt.Printf("%s\n", r)

	re, _ = regexp.Compile(`([0-9]+)`) // want to know what is in front of 'at'
	res := re.FindAllStringSubmatch("[MFT] One Piece Scan T01-T65 fr", -1)
	fmt.Printf("%v\n", res)

	matchVolume, _ := regexp.MatchString(`(Volume|Tome|volume|tome|T[0-9]+|V[0-9]+)`, "[MFT] One Piece Tome 01-65 fr")
	fmt.Println(matchVolume)
	// [MFT] One Piece Scan 740 fr

}
