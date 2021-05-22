package history

type command struct {
	typing string
	executionTime int64;
}

type userCommand struct {
	User string `json:"user"`
	UserType string `json:"userType"`
	Commands []command `json:"commands"`
}

type createGameHistoryDTO struct {
	Winner string `json:"winner"`
	Host string `json:"host"`
	RoomId string `json:"roomId"`
	UserCommands []userCommand `json:"userCommands"`
}