package helper

import (
	"fmt"
	"regexp"
	"strconv"
)

func checkForEmbeddedImages(bodyString string) [][]string {
	pattern := `xlink:href="([^"]+\.(jpg|png))"`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(bodyString, -1)
	return matches
}

func getDirName(url string) string {
	re := regexp.MustCompile(`(\d+)/(img|shade)/`)
	m := re.FindStringSubmatch(url)
	if len(m) != 3 {
		return ""
	}
	page, err := strconv.Atoi(m[1])
	if err != nil || page < 0 {
		return ""
	}
	return fmt.Sprintf("%04d/%s/", page, m[2])
}
