package validate

import "net/url"

func URL(longUrl string) bool {
	u, err := url.ParseRequestURI(longUrl)
	if err != nil {
		return false
	}

if u.Scheme == "" || u.Host == "" || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}

	return true
}
