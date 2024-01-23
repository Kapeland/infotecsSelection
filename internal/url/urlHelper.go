package url

import (
	"net/url"
	"strings"
)

func ParseURL(varURL *url.URL) []string {
	return strings.Split(varURL.Path, "/")
}
