package utils

import (
	"net/url"
)

func SetQueries(path string, queries map[string]string) (string, error) {
	u, err := url.Parse(path)
	if err != nil {
		return "", err
	}

	queryString := u.Query()
	for key, value := range queries {
		queryString.Set(key, value)
	}
	u.RawQuery = queryString.Encode()

	return u.String(), nil
}
