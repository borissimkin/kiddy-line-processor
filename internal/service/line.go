package service

import (
	"io"
	"net/http"
)

type LineService struct {
}

func (s *LineService) Fetch() (any, error) {
	resp, err := http.Get("http://localhost:8000/api/v1/lines/baseball")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	return body, nil
}
