// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    response, err := UnmarshalApodResponse(bytes)
//    bytes, err = response.Marshal()

package apodRequester

import "encoding/json"

func UnmarshalApodResponse(data []byte) (ApodResponse, error) {
	var r ApodResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ApodResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ApodResponse struct {
	Copyright      string `json:"copyright"`
	Date           string `json:"date"`
	Explanation    string `json:"explanation"`
	Hdurl          string `json:"hdurl"`
	MediaType      string `json:"media_type"`
	ServiceVersion string `json:"service_version"`
	Title          string `json:"title"`
	URL            string `json:"url"`
	Error          Error  `json:"error"`
}

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
