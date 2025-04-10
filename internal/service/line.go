package service

import (
	"encoding/json"
	"io"
	"net/http"
)

type LineResponseWrapper struct {
	Lines any `json:"lines"`
}

type LineBaseballResponse struct {
	Baseball float32 `json:"BASEBALL,string,omitempty"`
}

type LineService struct {
}

func (s *LineService) Fetch() (*LineBaseballResponse, error) {
	resp, err := http.Get("http://localhost:8000/api/v1/lines/baseball")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	payload := &LineBaseballResponse{}

	err = json.Unmarshal(body, &LineResponseWrapper{payload})
	
	return payload, err
}
