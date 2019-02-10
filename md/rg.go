package md

import "regexp"

var (
	// https://gist.github.com/jbroadway/2836900

	// ItalicRegexes $1
	ItalicRegexes = []*regexp.Regexp{
		regexp.MustCompile(`_(.*?)_`),
		regexp.MustCompile(`\*(.*?)\*`),
	}

	// BoldRegex $1
	BoldRegex = regexp.MustCompile(`\*\*(.*?)\*\*`)

	// StrikethroughRegex $1
	StrikethroughRegex = regexp.MustCompile(`\~\~(.*?)\~\~`)

	// SpoilerRegex $1
	SpoilerRegex = regexp.MustCompile(`\|\|(.*?)\|\|`)

	// UnderlineRegex $1
	UnderlineRegex = regexp.MustCompile(`__(.*?)__`)
)
