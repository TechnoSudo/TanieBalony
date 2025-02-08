package custom

import (
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

	Towers []Tower `json:"towers"`

	AvailableEnemies []Enemy `json:"available_enemies"`

	AvailableTowers []Tower `json:"available_towers"`
}
type Tower struct {
	ID       int    `json:"id"`
	Type     string `json:"type"`
	Position Position `json:"position"`
	Attack    int    `json:"attack"`
	Range    int    `json:"range"`
	Cost     int   `json:"cost"`
}
type Enemy struct {
	ID       int      `json:"id"`
	Type     string   `json:"type"`
	PositionIndex int `json:"position_index"`
	Health   int      `json:"health"`
	Speed    int      `json:"speed"`
	Cost     int   `json:"cost"`
}
type Game struct {
	PlayerCount    int8
	AttackersCount int8
	DefendersCount int8
	GameState      GameState
	Players        []int
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
				{ X:50, Y:50 },
				{ X:100, Y:50 },
				{ X:150, Y:50 },
				{ X:200, Y:50 },
				{ X:250, Y:50 },
				{ X:300, Y:50 },
				{ X:350, Y:50 },
				{ X:400, Y:50 },
				{ X:450, Y:50 },
				{ X:500, Y:50 },
				{ X:500, Y:100 },
				{ X:500, Y:150 },
				{ X:500, Y:200 },
				{ X:500, Y:250 },
				{ X:500, Y:300 },
				{ X:500, Y:350 },
				{ X:450, Y:350 },
				{ X:400, Y:350 },
				{ X:350, Y:350 },
				{ X:300, Y:350 },
				{ X:250, Y:350 },
				{ X:200, Y:350 },
				{ X:150, Y:350 },
				{ X:150, Y:400 },
				{ X:150, Y:450 },
				{ X:150, Y:500 },
				{ X:150, Y:550 },
				{ X:200, Y:550 },
				{ X:250, Y:550 },
				{ X:300, Y:550 },
				{ X:350, Y:550 },
				{ X:400, Y:550 },
				{ X:450, Y:550 },
				{ X:500, Y:550 },
				{ X:550, Y:550 },
				{ X:600, Y:550 },
				{ X:650, Y:550 },
				{ X:700, Y:550 },
				{ X:700, Y:500 },
				{ X:700, Y:450 },
				{ X:700, Y:400 },
				{ X:700, Y:350 },
				{ X:700, Y:300 },
				{ X:700, Y:250 },
				{ X:700, Y:200 },
				{ X:700, Y:150 },
				{ X:700, Y:100 },
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
			// {
			// 	ID:   1,
			// 	Name: "Attacker One",
			// 	Gold: 100,
			// 	SpawnedEnemies: []struct {
			// 		Type  string `json:"type"`
			// 		Count int    `json:"count"`
			// 	}{
			// 		{Type: "Goblin", Count: 5},
			// 		{Type: "Orc", Count: 3},
			// 	},
			// },
		},
		Defenders: []Defender{
			// {
			// 	ID:   1,
			// 	Name: "Defender One",
			// 	Gold: 150,
			// },
		},
		Enemies: []Enemy{
			// {
			// 	ID:       1,
			// 	Type:     "Goblin",
			// 	Position: Position{X: 1, Y: 1},
			// 	Health:   30,
			// 	Speed:    10,
			// },
		},
		Towers: []Tower{
			// {
			// 	ID:       1,
			// 	Type:     "Archer Tower",
			// 	Position: Position{X: 2, Y: 2},
			// 	Attack:   15,
			// 	Range:    5,
			// },
		},
		AvailableEnemies: []Enemy{
			{
				ID:       1,
				Type:     "Red blun",
				PositionIndex: 0,
				Health:   30,
				Speed:    10,
				Cost:     10,
			},
			{
				ID:       2,
				Type:     "Green blun",
				PositionIndex: 0,
				Health:   80,
				Speed:    10,
				Cost:     40,
			},
			{
				ID:       3,
				Type:     "Blue blun",
				PositionIndex: 0,
				Health:   200,
				Speed:    10,
				Cost:     60,
			},
		},
		AvailableTowers: []Tower{
			{
				ID:       1,
				Type:     "Archer Tower",
				Position: Position{X: 2, Y: 2},
				Attack:   15,
				Range:    1,
				Cost:     100,
			},
			{
				ID:       1,
				Type:     "Archerer Tower",
				Position: Position{X: 2, Y: 2},
				Attack:   30,
				Range:    2,
				Cost:     200,
			},
			{
				ID:       1,
				Type:     "Archerest Tower",
				Position: Position{X: 2, Y: 2},
				Attack:   45,
				Range:    3,
				Cost:     400,
			},
		},
	}
}


type SafeMap struct {
	mu    sync.Mutex
	Store map[int]*Game
	logger  *slog.Logger
}

func (sm *SafeMap) GetFreeGame(userId int) (int, bool){
	slog.Debug("waiting for lock")
	sm.mu.Lock()
	slog.Debug("got the lock")
	defer sm.mu.Unlock()
	
	for gameId, game := range sm.Store {
		if game.PlayerCount < 4 {
			sm.Store[gameId].Players = append(sm.Store[gameId].Players, userId)
			if game.AttackersCount == 2 {
				sm.Store[gameId].PlayerCount +=1
			} else {
				sm.Store[gameId].DefendersCount +=1
			}
			sm.Store[gameId].PlayerCount +=1

			
			return gameId, false
		}
	}
	slog.Debug("unlocking")
	return -1, true

}



func (sm *SafeMap) GetPlayers(gameId int) []int {
	slog.Debug("waiting for lock")
	sm.mu.Lock()
	slog.Debug("got the lock")
	defer sm.mu.Unlock()
	return sm.Store[gameId].Players
}

func (sm *SafeMap) StartNewGame(userId int) int {
	slog.Debug("waiting for lock")
	sm.mu.Lock()
	slog.Debug("got the lock")
	defer sm.mu.Unlock()

	var newGameId int = 0 // TODO more games
	slog.Info("trying to create newGame")

	sm.Store[newGameId] = newGame()
	slog.Info("playerCOunt:" ,)
	sm.Store[newGameId].PlayerCount +=1
	sm.Store[newGameId].DefendersCount +=1
	sm.Store[newGameId].Players = append(sm.Store[newGameId].Players, userId)
	slog.Info("playerCOunt:" ,sm.Store[0].PlayerCount)
	slog.Info("unlocking")
	return newGameId
}


func NewSafeMap() *SafeMap {
	return &SafeMap{
		Store: make(map[int]*Game),
	}
}
func (sm *SafeMap) AddTower(gameId int, tower Tower) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	var gameState = sm.Store[gameId].GameState
	if gameState.Towers == nil {
		gameState.Towers = make([]Tower, 0)
	}
	gameState.Towers = append(gameState.Towers, tower)

	sm.Store[gameId].GameState = gameState

}

func (sm *SafeMap) AddEnemy(gameId int, enemy Enemy) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	println(&sm.Store[gameId].GameState.Enemies)
	sm.Store[gameId].GameState.Enemies = append(sm.Store[gameId].GameState.Enemies, enemy)
	println("done did enemy")
}

func (sm *SafeMap) MoveTime(gameId int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	var gameState = sm.Store[gameId].GameState
	// move enemies along path
	for enemyId, _ := range gameState.Enemies {
		sm.Store[gameId].GameState.Enemies[enemyId].PositionIndex += 1
		if sm.Store[gameId].GameState.Enemies[enemyId].PositionIndex >= len(sm.Store[gameId].GameState.Map.Path) {
			sm.Store[gameId].GameState.Enemies = append(sm.Store[gameId].GameState.Enemies[:enemyId], sm.Store[gameId].GameState.Enemies[enemyId+1:]...)
		}
	}


	// //deal dmg
	// for towerId, _ := range sm.Store[gameId].GameState.Towers {
	// 	sm.Store[gameId].GameState.Enemies[enemyId].PositionIndex += 1
	// }

}

func (sm *SafeMap) Put(key int, value *Game) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Store[key] = value
}

func (sm *SafeMap) Get(key int) (*Game, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	value, exists := sm.Store[key]
	return value, !exists
}

func (sm *SafeMap) GetGameState(key int) GameState {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	value := sm.Store[key].GameState
	return value
}

func (sm *SafeMap) GetGames() []int {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	var list []int 
	for id := range sm.Store {
		list = append(list, id)
	}
	return list
}


// func main() {
// 	safeMap := NewSafeMap()

// 	safeMap.Put(1, newGame())

// 	value, exists := safeMap.Get(1)
// 	if exists {
// 		println("key1:", value)
// 	} else {
// 		println("key1 not found")
// 	}
// }
