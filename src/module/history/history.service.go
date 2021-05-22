package history

import (
	"onair/src/database"
	"onair/src/model"
	"onair/src/shared"
	"time"
)

type service struct {
	DB *database.DB
}

func (s *service) getAllHistory() *[]model.GameHistory {
	var gameHistories []model.GameHistory
	err := s.DB.CoreTable.Get("PK", model.GetGameHistoryPK()).All(&gameHistories)
	if err != nil {
		panic(err)
	}

	return &gameHistories
}

func (s *service) createGameHistory(dto *createGameHistoryDTO) error {
	gameHistory := model.GameHistory{
		Key: shared.Key{
			PK: model.GetGameHistoryPK(),
			SK: model.GetGameHistorySK(dto.RoomId),
		},
		Host: dto.Host,
		Winner: dto.Winner,
		CreatedAt: time.Now().Unix(),
	}
	return s.DB.CoreTable.Put(&gameHistory).Run()
}

func newService(DB *database.DB) *service {
	return &service{DB: DB}
}