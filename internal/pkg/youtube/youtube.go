package youtube

import (
	"net/url"
	"regexp"
	"strings"
)

var pathIDPattern = regexp.MustCompile(`/(shorts|embed|v)/([a-zA-Z0-9_-]{11})`)

// ParseVideoID extracts the 11-character YouTube video ID from common URL formats.
func ParseVideoID(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}

	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	host := strings.ToLower(u.Host)
	if strings.Contains(host, "youtube.com") {
		if v := u.Query().Get("v"); v != "" {
			return v
		}
		if m := pathIDPattern.FindStringSubmatch(u.Path); len(m) == 3 {
			return m[2]
		}
	}

	if host == "youtu.be" {
		id := strings.TrimPrefix(u.Path, "/")
		id = strings.Split(id, "?")[0]
		if len(id) == 11 {
			return id
		}
	}

	return ""
}
