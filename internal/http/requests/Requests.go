package requests

import (
	"../../app"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

type DeployRequest struct {
	Token string `json:"token"`
	Repo  string `json:"repo"`
}
type FindResponse struct {
	Running  *[]app.AppJSON `json:"running"`
	Deployed *[]app.AppJSON `json:"deployed"`
}
type ErrorResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
type SuccessResponse struct {
	Message string  `json:"message"`
	App     app.AppJSON `json:"app"`
}
type SettingsRequest struct {
	Id       string            `json:"id"`
	Settings map[string]string `json:"settings"`
}

func (s *SettingsRequest) Read(body *io.ReadCloser) {
	jsonBytes, _ := ioutil.ReadAll(*body)
	err := json.Unmarshal(jsonBytes, s)
	if err != nil {
		fmt.Println(err)
	}
}
