package main

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"onair/src/config"
	"onair/src/database"
	"onair/src/module/history"
	"onair/src/utils"
	"os"

	"github.com/gofiber/fiber/v2/middleware/cors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/google/uuid"
)

func GetRandomCode() string {
	code := uuid.NewString()
	return code
}

var animals = []string{"고양이", "독수리", "고슴도치", "마모트"}
var lastId = 0

type Client struct {
	Id         int
	Name       string
	IsPlayed   bool
	AnimalType int
	SkillUsed  bool
} // Add more data to this type if needed

var roomDB = map[string]Room{}

type Chat struct {
	Message   string
	Client    Client
	CreatedAt string
}

type Movement struct {
	IsLeft bool
}

// 채팅 개발 (한결)

// 유저이동
// 	이동 방향 (홍두)

// 탄막스킬
// 탄막 유저 번호 (홍두)

// 탄막생성
// 좌표, 탄막 종류, 탄막 유저 번호 (한결)

// 탄막이동
// 탄막 이동 방향, 탄막 유저 번호 (홍두)

// 방에 더 못들어오게하는거 (홍두) V

// 유저 정보 받기 (한결)

// 방 정보 받기 (한결)

// 더미 실행코드 (홍두)

// 배포 (한결)

type Room struct {
	Code       string
	Host       *Client
	Clients    map[*websocket.Conn]Client
	Register   chan *websocket.Conn
	Broadcast  chan string
	Unregister chan *websocket.Conn
	MaxClients int
	Difficulty int
	IsStart    bool
}

func NewRoom(room Room) {
	room.Clients = make(map[*websocket.Conn]Client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	room.Register = make(chan *websocket.Conn)
	room.Broadcast = make(chan string)
	room.Unregister = make(chan *websocket.Conn)
	roomDB[room.Code] = room
}

func runRoom(room Room) {
	for {
		select {
		case connection := <-room.Register:
			lastId += 1
			client := Client{
				lastId,
				fmt.Sprintf("%v", connection.Params("name")),
				false,
				rand.Intn(len(animals)),
				false,
			}
			if len(room.Clients) == 0 {
				room.Host = &client
			}
			room.Clients[connection] = client
			log.Println("connection registered")
			clientByte, _ := json.Marshal(client)
			if err := connection.WriteMessage(websocket.TextMessage, clientByte); err != nil {
				log.Println("write error:", err)
				connection.WriteMessage(websocket.CloseMessage, []byte{})
				connection.Close()
				delete(room.Clients, connection)
			}

		case message := <-room.Broadcast:
			log.Println("message received:", message)

			// Send the message to all clients
			for connection := range room.Clients {
				if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
					log.Println("write error:", err)

					connection.WriteMessage(websocket.CloseMessage, []byte{})
					connection.Close()
					delete(room.Clients, connection)
				}
			}

		case connection := <-room.Unregister:
			// Remove the Client from the hub
			delete(room.Clients, connection)
			log.Println("connection unregistered")
		}
	}
}

func main() {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))
	app.Use("ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the Client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	app.Post("/api/create/room", func(ctx *fiber.Ctx) error {
		code := GetRandomCode()
		type roomBody struct {
			MaxClients int `json:"maxClients" validate:"min=50,max=100"`
			Difficulty int `json:"difficulty" validate:"min=1,max=3"`
		}
		roomInput := new(roomBody)
		if err := ctx.BodyParser(roomInput); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		errs := utils.Validate(roomInput)
		if len(errs) > 0 {
			return fiber.NewError(fiber.StatusBadRequest, "validation error")
		}
		room := Room{
			Code:       code,
			MaxClients: roomInput.MaxClients,
			Difficulty: roomInput.Difficulty,
		}
		NewRoom(room)
		go runRoom(roomDB[code])
		return ctx.JSON(fiber.Map{
			"code":       code,
			"maxClients": roomInput.MaxClients,
			"difficulty": roomInput.Difficulty,
		})
	})

	app.Get("/ws/:code/join/:name", websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the Client and close the connection
		room := roomDB[c.Params("code")]
		defer func() {
			room.Unregister <- c
			c.Close()
		}()

		// Register the Client
		room.Register <- c

		if len(room.Clients) >= room.MaxClients {
			return
		}

		for {
			messageType, message, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Println("read error:", err)
				}

				return // Calls the deferred function, i.e. closes the connection on error
			}

			if messageType == websocket.TextMessage {
				// Broadcast the received message
				//event := Event{}
				//err := json.Unmarshal(message, &event)
				//if err != nil {
				//	log.Fatal(err.Error())
				//}
				room.Broadcast <- string(message)
			} else {
				log.Println("websocket message received of type", messageType)
			}
		}
	}))

	log.Printf("[%s] START SERVER ON %s", os.Getenv("MODE"), config.GetEnv("PORT"))

	db := database.GetDB()
	historyGroup := app.Group("api/history")
	history.InitModule(historyGroup, db)

	//db := database.GetNewConnection(config.DSN, &gorm.Config{})
	//bookRepository := repository.NewBookRepository(db)
	//bookUseCase := usecase.NewBookUseCase(bookRepository)
	//bookHandler := handler.NewBookHandler(bookUseCase)
	//
	//bookRouter := app.Group("books")
	//{
	//	bookRouter.Get("", bookHandler.GetAllBooks)
	//	bookRouter.Get("/:id", bookHandler.GetBook)
	//	bookRouter.Post("", bookHandler.CreateBook)
	//	bookRouter.Post("/:id", bookHandler.UpdateBook)
	//	bookRouter.Delete("/:id", bookHandler.DeleteBook)
	//}

	log.Fatal(app.Listen(config.GetEnv("PORT")))
}
