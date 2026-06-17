package validate

import "net/url"

func URL(longUrl string) bool {
	u, err := url.ParseRequestURI(longUrl)
	if err != nil {
		return false
	}

	if u.Scheme == "" || u.Host == "" {
		return false
	}

	return true
}
