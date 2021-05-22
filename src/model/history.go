package model

import (
	"fmt"
	"onair/src/shared"
)

type GameHistory struct {
	 shared.Key
	 Winner string `json:"winner"`
	 Host string `json:"host"`
	 CreatedAt int64 `json:"createdAt"`
}

func GetGameHistoryPK () string {
	return fmt.Sprintf("GH")
}

func GetGameHistorySK (roomId string) string {
	return fmt.Sprintf("%s", roomId)
}

type HistoryDetail struct {
	shared.Key
	UserType string `json:"userType"`
	Commands []string `json:"commands"`
}

func GetHistoryDetailPK (roomId string) string {
	return fmt.Sprintf("HD#%s", roomId)
}

func GetHistoryDetailSK (user string) string {
	return fmt.Sprintf("%s", user)
}