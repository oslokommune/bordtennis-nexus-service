package status

import (
	"fmt"
)

type Data struct {
	TeamOne      int `json:"teamOne"`
	TeamTwo      int `json:"teamTwo"`
	InitialServe int `json:"initialServe"`
}

func Serialize(data Data) string {
	return fmt.Sprintf("%d;%d;%d", data.TeamOne, data.TeamTwo, data.InitialServe)
}

func Reset(status *Data) {
	status.TeamOne = 0
	status.TeamTwo = 0
	status.InitialServe = 1
}
