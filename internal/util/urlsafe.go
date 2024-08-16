package util

import (
	"regexp"
	"strings"
)

func UrlSafe(val string) string {
	// Convert to lowercase
	slug := strings.ToLower(val)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove all non-url friendly characters (only keep alphanumeric, hyphens, and underscores)
	reg := regexp.MustCompile("[^a-z0-9-_]+")
	slug = reg.ReplaceAllString(slug, "")

	return slug
}
