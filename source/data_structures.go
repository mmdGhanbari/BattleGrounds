package battle_grounds

import (
    "time"
)

type Player struct {
    Id int
    Username string
    Password string
    Level int
    XP int
}

type PlayerSign struct {
    Username string
    Level int
}

type Match struct {
    Player1 *Connection
    Player2 *Connection
    ReadyState map[int]bool
    State string   // lobby, playing, finished
    CreationTime time.Time
}

type PlayerInGameData struct {
	Cards map[int]int	// [card id] -> number of that card
	Fighters []Fighter	//  list of live fighters in the ground
}

type FighterCard struct {
	Id int
	Name string
	Attack int
	Defence int
	Health int
	MaxUse int
	Type string
}

type Fighter struct {
	Id int
	Health int
	Position int
}
