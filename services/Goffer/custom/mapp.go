package custom

import (
	"fmt"
	"sync"
	// "math/rand"
	"log/slog"
)

type Position struct {
	X int `json:X`
	Y int `json:Y`
}

type Defender struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Gold         int    `json:"gold"`
}

type Attacker struct {
	ID             int    `json:"id"`
	Name           string `json:"name"`
	Gold           int    `json:"gold"`
	SpawnedEnemies []struct {
		Type  string `json:"type"`
		Count int    `json:"count"`
	} `json:"spawnedEnemies"`
} 

type GameState struct {
	Map struct {
		Name       string    `json:"name"`
		Path       []Position `json:"path"`
		Dimensions struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"dimensions"`
	} `json:"map"`

	Attackers []Attacker `json:"attackers"`

	Defenders []Defender `json:"defenders"` 

	Enemies []Enemy `json:"enemies"`

	Towers []Tower `json:"Towers"`
}
type Tower struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Position Position `json:"position"`
	Attack    int    `json:"attack"`
	Range    int    `json:"range"`
}
type Enemy struct {
	ID       int      `json:"id"`
	Type     string   `json:"type"`
	Position Position `json:"position"`
	Health   int      `json:"health"`
	Speed    int      `json:"speed"`
}

type Game struct {
	PlayerCount    int8
	AttackersCount int8
	DefendersCount int8
	GameState      GameState
}

func newGame() *Game {
	return &Game{
		PlayerCount:    0,
		AttackersCount: 0,
		DefendersCount: 0,
		GameState: initGameState(),
	}
}

func initGameState() GameState {
	slog.Info("Creating new Game state")
	return GameState{
		Map: struct {
			Name       string     `json:"name"`
			Path       []Position `json:"path"`
			Dimensions struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			} `json:"dimensions"`
		}{
			Name: "Fantasy World",
			Path: []Position{
				{ X:0, Y:50 },
				{ X:500, Y:50 },
				{ X:500, Y:350 },
				{ X:150, Y:350 },
				{ X:150, Y:550 },
				{ X:700, Y:550 },
				{ X:700, Y:50 },
			},
			Dimensions: struct {
				Width  int `json:"width"`
				Height int `json:"height"`
			}{
				Width:  800,
				Height: 600,
			},
		},
		Attackers: []Attacker{
			{
				ID:   1,
				Name: "Attacker One",
				Gold: 100,
				SpawnedEnemies: []struct {
					Type  string `json:"type"`
					Count int    `json:"count"`
				}{
					{Type: "Goblin", Count: 5},
					{Type: "Orc", Count: 3},
				},
			},
		},
		Defenders: []Defender{
			{
				ID:   1,
				Name: "Defender One",
				Gold: 150,
			},
		},
		Enemies: []Enemy{
			{
				ID:       1,
				Type:     "Goblin",
				Position: Position{X: 1, Y: 1},
				Health:   30,
				Speed:    10,
			},
		},
		Towers: []Tower{
			{
				ID:       1,
				Type:     "Archer Tower",
				Position: Position{X: 2, Y: 2},
				Attack:   15,
				Range:    5,
			},
		},
	}
}


type SafeMap struct {
	mu    sync.Mutex
	store map[int]*Game
	logger  *slog.Logger
}

func (sm *SafeMap) GetFreeGame() (int, bool){
	slog.Info("waiting for lock")
	sm.mu.Lock()
	slog.Info("got the lock")
	defer sm.mu.Unlock()
	
	for gameId, game := range sm.store {
		if game.PlayerCount < 4 {
			return gameId, false
		}
	}
	slog.Info("unlocking")
	return -1, true

}

func (sm *SafeMap) StartNewGame() int {
	slog.Info("waiting for lock")
	sm.mu.Lock()
	slog.Info("got the lock")
	defer sm.mu.Unlock()

	var newGameId int
	newGameId = 1
	slog.Info("trying to create newGame")
	varnewGaem := newGame()
	if varnewGaem == nil {
		slog.Info("test")
	}

	sm.store[newGameId] = newGame()
	slog.Info("unlocking")
	return newGameId
}


func NewSafeMap() *SafeMap {
	return &SafeMap{
		store: make(map[int]*Game),
	}
}
func (sm *SafeMap) AddTower(gameId int, tower Tower) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var gameState = sm.store[gameId].GameState

	gameState.Towers = append(gameState.Towers, tower)

	sm.store[gameId].GameState = gameState

}

// AddEnemy adds an enemy to the game state
func (sm *SafeMap) AddEnemy(gameId int, enemy Enemy) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.store[gameId].GameState.Enemies = append(sm.store[gameId].GameState.Enemies, enemy)


	sm.logger.Info("Added Enemy", slog.Int("EnemyID", enemy.ID), slog.String("EnemyType", enemy.Type), slog.Int("X", enemy.Position.X), slog.Int("Y", enemy.Position.Y))
}

func (sm *SafeMap) MoveTime(key int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
}

func (sm *SafeMap) Put(key int, value *Game) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.store[key] = value
}

func (sm *SafeMap) Get(key int) (*Game, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	value, exists := sm.store[key]
	return value, !exists
}

func (sm *SafeMap) GetGameState(key int) GameState {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	value := sm.store[key].GameState
	return value
}

func main() {
	safeMap := NewSafeMap()

	safeMap.Put(1, newGame())

	value, exists := safeMap.Get(1)
	if exists {
		fmt.Println("key1:", value)
	} else {
		fmt.Println("key1 not found")
	}
}
