package downloader

import (
	"regexp"
)

func checkForEmbeddedImages(bodyString string) [][]string {
	pattern := `xlink:href="([^"]+\.(jpg|png))"`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(bodyString, -1)
	return matches
}

func getDirName(url string) string {
	pattern := `\d+/(img|shade)/`
	re := regexp.MustCompile(pattern)
	match := re.FindString(url)
	return match
}
