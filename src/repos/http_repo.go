package repos

import (
	"bytes"
	"encoding/json"
	"image"
	"net/http"
)

type HttpRepo struct {
	http.Client
}

func (hr *HttpRepo) JsonRequest(method string, url string, body interface{}) (*http.Response, error) {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}

	//req.Header.Set("Content-Type", "application/json")
	return hr.Do(req)
}

func (hr *HttpRepo) DecodedJsonRequest(
	method string, url string, body interface{},
	decodeInto interface{},
) (*http.Response, error) {
	res, err := hr.JsonRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	err = json.NewDecoder(res.Body).Decode(decodeInto)
	return res, err

}

func (hr *HttpRepo) ImageDecodedJsonRequest(
	method string, url string, body interface{},
) (*image.Image, *string, error) {
	res, err := hr.JsonRequest(method, url, body)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()
	img, ext, err := image.Decode(res.Body)
	return &img, &ext, err
}
