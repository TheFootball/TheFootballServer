package main

import "C"
import (
	"log"
	"onair/src/config"
	"onair/src/database"
	"onair/src/module/history"

	"os"

	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func getRandomCode() string {
	// TODO: uuid4
	return "ABC"
}

func (room *Room) publishRoomMessage(message []byte, code string, cache *redis.Client) {
	err := cache.Publish(code, message).Err()

	if err != nil {
		log.Println(err)
	}
}

func (room *Room) subscribeToRoomMessages(code string, cache *redis.Client) {
	pubsub := cache.Subscribe(code)

	ch := pubsub.Channel()

	for msg := range ch {
		room.broadcastToClientsInRoom([]byte(msg.Payload))
	}
}

type client struct{} // Add more data to this type if needed

var roomDB = map[string]Room{}

// TODO: redis달아야함

type Room struct {
	Clients    map[*websocket.Conn]client
	Register   chan *websocket.Conn
	Broadcast  chan string
	Unregister chan *websocket.Conn
}

func NewRoom(code string, cache *redis.Client) Room {
	room := Room{}
	room.Clients = make(map[*websocket.Conn]client) // Note: although large maps with pointer-like types (e.g. strings) as keys are slow, using pointers themselves as keys is acceptable and fast
	room.Register = make(chan *websocket.Conn)
	room.Broadcast = make(chan string)
	room.Unregister = make(chan *websocket.Conn)
	// roomDB[code] = room
	cache.Set(code, room, 99999)
	return room
}

func runRoom(room Room) {
	for {
		select {
		case connection := <-room.Register:
			room.Clients[connection] = client{}
			log.Println("connection registered")

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
			// Remove the client from the hub
			delete(room.Clients, connection)
			log.Println("connection unregistered")
		}
	}
}

var cache *redis.Client

func main() {
	cache = database.GetCacheClient()
	app := fiber.New()
	app.Use("ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) { // Returns true if the client requested upgrade to the WebSocket protocol
			return c.Next()
		}
		return c.SendStatus(fiber.StatusUpgradeRequired)
	})

	app.Get("/create/", func(ctx *fiber.Ctx) error {
		code := getRandomCode()
		room := NewRoom(code, cache)
		go runRoom(room)
		return ctx.JSON(fiber.Map{
			"code": code,
		})
	})

	app.Get("/ws/:code/join", websocket.New(func(c *websocket.Conn) {
		// When the function returns, unregister the client and close the connection
		// room := roomDB[c.Params("code")]
		code := c.Params("code")
		room := cache.Get(code)
		defer func() {
			room.Unregister <- c
			c.Close()
		}()

		// Register the client
		room.Register <- c

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
