package api

import (
	"io"
	"net/http"
	"net/url"
)

type Requester struct{}

func (r *Requester) Get(path string, query map[string]string) (string, error) {
	baseURL := "https://api.modrinth.com/v2"
	u, err := url.Parse(baseURL + path)
	if err != nil {
		return "", err
	}
	q := u.Query()
	for k, v := range query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	response, err := http.Get(u.String())
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
