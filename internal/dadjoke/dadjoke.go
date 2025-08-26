package dadjoke

import (
	"io"
	"net/http"
)

// https://icanhazdadjoke.com/api
func GetJoke(jokeID string) (string, error) {
	path := ""
	if jokeID != "" {
		path = "/j/" + jokeID
	}

	req, err := http.NewRequest("GET", "https://icanhazdadjoke.com"+path, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "https://github.com/mvndaai/go_wasm_tools")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
