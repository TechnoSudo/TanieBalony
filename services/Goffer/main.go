package main

import (
	"encoding/json"
	"log/slog"
	"log"
	"net/http"

	// "slog"
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
}


func (connectionManager* ConnectionManager) HandleWebSocketConn(w http.ResponseWriter, r *http.Request) {

	// var userId string = r.URL.Query().Get("userId");
	var userId = "test1"


	//verify TODO



	//add to game;
	var gameId int
	gameId, err := connectionManager.Games.GetFreeGame()
	
	if err {
		// add to game
		gameId = connectionManager.Games.StartNewGame()

	}
	connectionManager.Players.Store(userId, gameId)


	// start websocket
	conn, err2 := upgrader.Upgrade(w, r, nil)

	if err2 != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}


    ticker := time.NewTicker(websocketLifetime)
    defer ticker.Stop()

	go func(conn *websocket.Conn) {
		defer conn.Close()

		ticker := time.NewTicker(websocketLifetime)
		fpsTicker := time.NewTicker(500 * time.Millisecond)
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
					log.Println("Websocket timeout reached. Closing the WebSocket connection.")
					conn.Close()
					return
				case <-fpsTicker.C:
					slog.Info("mainloop")
					data:= connectionManager.Games.GetGameState(gameId)
					jsonData, _ := json.Marshal(data)
					if err := conn.WriteMessage(websocket.TextMessage, jsonData); err != nil {
						slog.Error("Error sending game status:	", err)
					}
				case _ = <-messageChan:
					slog.Info("Got message from client on websocket")
				case err := <-errorChan:
					log.Println("Error reading message:", err)
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

	slog.Info("WebSocket server started on :" + port + " on /game endpoint")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
