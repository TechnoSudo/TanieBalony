package main

import (
	// "crypto/rand"
	// "math/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	// "logger"
	// "os"
	"sync"
	"time"

	// "strings"
	"Goffer/custom"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	// "gorm.io/gorm"
)


const (
	websocketLifetime = 2 * time.Hour // Time after which we check if the user is still connected, if not then drop connection
)

type MetaMessage struct {
	Operation    string           `json:"operation"`
	Tower        custom.Tower     `json:"tower"`
	Enemy        custom.Enemy     `json:"enemy"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type GameState struct {
	Map struct {
		Name string `json:"name"`
		Path []struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"path"`
		Dimensions struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"dimensions"`
	} `json:"map"`
	Players struct {
		Attackers []struct {
			ID             int    `json:"id"`
			Name           string `json:"name"`
			Gold           int    `json:"gold"`
			SpawnedEnemies []struct {
				Type  string `json:"type"`
				Count int    `json:"count"`
			} `json:"spawnedEnemies"`
		} `json:"attackers"`
		Defenders []struct {
			ID           int    `json:"id"`
			Name         string `json:"name"`
			Gold         int    `json:"gold"`
			PlacedTowers []struct {
				ID       int    `json:"id"`
				Type     string `json:"type"`
				Position struct {
					X int `json:"x"`
					Y int `json:"y"`
				} `json:"position"`
				Level int `json:"level"`
			} `json:"placedTowers"`
		} `json:"defenders"`
	} `json:"players"`
	Enemies []struct {
		ID       int    `json:"id"`
		Type     string `json:"type"`
		Position struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"position"`
		Health int `json:"health"`
		Speed  int `json:"speed"`
	} `json:"enemies"`
}

type Game struct {
	PlayerCount int8
	AttackersCount int8
	DefendersCount int8
	GameState GameState
}

type ConnectionManager struct {
	Connections    sync.Map //userId to websocketConn
	Games		   custom.SafeMap //gameId to jsonState
	Players 	   sync.Map //userId to gameId
	isUpdated 	   bool
	updateMutex    sync.Mutex
}
func (connectionManager* ConnectionManager) BackgroundLogic() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout,  &slog.HandlerOptions{
		Level: lvl,
	}))
	//update game state, send updated state
	fpsTicker := time.NewTicker(100 * time.Millisecond)

	func() {
		for {
			select {
				case <- fpsTicker.C:
					var games = connectionManager.Games.GetGames()
					println("games:", games)
					for gameId :=  range games {
						println("moving time for ", gameId)
						connectionManager.Games.MoveTime(gameId)

						println("getGamestate")
						data:= connectionManager.Games.GetGameState(gameId)
						jsonData, _ := json.Marshal(data)

						println("get players")
						var players = connectionManager.Games.GetPlayers(gameId)
						println("got players", players)
						for userId := range players {
							println("userId: ",userId)
							retrieved_conn, ok := connectionManager.Connections.Load(userId)
							if !ok {
								logger.Error("No active connection found for userId:", userId)
								break
							}
							conn, ok := retrieved_conn.(*websocket.Conn) 
							if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
								logger.Error("Error sending game status:	", err)
							}
		
						}
					}
					logger.Debug("mainloop")
				
	
		}
		}
	}()
	fpsTicker.Stop()
}

func (connectionManager* ConnectionManager) HandleWebSocketConn(w http.ResponseWriter, r *http.Request) {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)
	logger := slog.New(slog.NewJSONHandler(os.Stdout,  &slog.HandlerOptions{
		Level: lvl,
	}))

	// var userId string = r.URL.Query().Get("userId");
	var userId int = 0
	println(userId)


	//verify TODO



	//add to game;
	var gameId int
	
	gameId, err2 := connectionManager.Games.GetFreeGame(userId)
	
	if err2 {
		gameId = connectionManager.Games.StartNewGame(userId)
	}
	connectionManager.Players.Store(userId, gameId)

	conn, err := upgrader.Upgrade(w, r, nil)
	
	connectionManager.Connections.Store(userId, conn)
	logger.Debug("Storing connection for userId: ", userId)
	// defer connectionManager.Connections.Delete(userId)
	if err != nil {
		logger.Error("Error upgrading to WebSocket:", err)
		return
	}


    ticker := time.NewTicker(websocketLifetime)
    defer ticker.Stop()

	go func(conn *websocket.Conn) {
		defer conn.Close()

		ticker := time.NewTicker(websocketLifetime)

		defer ticker.Stop()

		for {
			// channels for websocket signals
			messageChan := make(chan []byte)
			errorChan := make(chan error)

			// start a goroutine for websocket read
			go func() {
				_, msg, err := conn.ReadMessage()
				if err != nil {
					errorChan <- err
					return
				}
				messageChan <- msg
			}()

			select {
				case <-ticker.C:
					logger.Error("Websocket timeout reached. Closing the WebSocket connection.")
					conn.Close()
					return
				case msg := <-messageChan:
					logger.Info("Got message from client on websocket")
					var newReceivedMessage MetaMessage
					err = json.Unmarshal(msg, &newReceivedMessage)
					if err != nil {
						// If there's an error parsing JSON, respond with an error message
						slog.Error("Error parsing JSON:", err)
						slog.Error(newReceivedMessage.Operation, newReceivedMessage.Enemy, newReceivedMessage.Tower)
						errorMessage := fmt.Sprintf("Error: Invalid JSON format: %v", err)
						if err := conn.WriteMessage(websocket.TextMessage, []byte(errorMessage)); err != nil {
							slog.Error("Error sending error message:", err)
						}
						continue
					} else {
						// process correct message
						slog.Info("Processing correct message")
						slog.Info(newReceivedMessage.Operation)
						// TODO add error checking below
						if newReceivedMessage.Operation == "addTower" {
							connectionManager.Games.AddTower(gameId, newReceivedMessage.Tower)
						} else if newReceivedMessage.Operation == "addEnemy" {
							connectionManager.Games.AddEnemy(gameId, newReceivedMessage.Enemy)
						}
					}
				case err := <-errorChan:
					logger.Error("Error reading message:", err)
					conn.Close()
					return
			}
		}
	}(conn)
}

func main() {
	var port string = "8088"

	// application := services.NewMessageService(db)

	r := mux.NewRouter()
	connectionManager := ConnectionManager{}
	connectionManager.Games = *custom.NewSafeMap()
	
	mainRouter := r.PathPrefix("").Subrouter()

	mainRouter.HandleFunc("/game", connectionManager.HandleWebSocketConn)
	http.Handle("/game", r)
	go  connectionManager.BackgroundLogic()
	slog.Info("WebSocket server started on :" + port + " on /game endpoint")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		slog.Error("ListenAndServe:", err)
	}
}
