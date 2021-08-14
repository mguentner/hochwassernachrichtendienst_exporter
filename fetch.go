package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	RequestError = errors.New("Expected HTTP 200")
	NotFoundError = errors.New("Station probably invalid")
)

func FetchLevels(station string) (io.Reader, error) {
	url := fmt.Sprintf("https://m.hnd.bayern.de/pegel.php?pgnr=%s", station)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, RequestError
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, NotFoundError
	}
	reader := bytes.NewReader(body)
	return reader, nil
}
