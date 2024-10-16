package posts

import (
	"regexp"
)

// ExtractHashtags extracts valid hashtags from a given content string.
func ExtractHashtags(content string) []string {

	re := regexp.MustCompile(`(?i)#([a-z0-9_]+)`)
	matches := re.FindAllStringSubmatch(content, -1)

	var hashtags []string
	for _, match := range matches {
		hashtags = append(hashtags, match[1]) // Append only the hashtag content (without the # symbol)
	}
	return hashtags
}
