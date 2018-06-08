package battle_grounds

import (
    "time"
    "encoding/json"
)

var requestPool = make([]*Connection, 0, 10)
var matchPool = make([]*Match, 0, 10)

func MatchRequest(newReq *Connection) {
    opp := FindProperOpponent(newReq.player.Level)

    if opp == nil {
        requestPool = append(requestPool, newReq)
    } else {
		newMatch := Match{opp, newReq, map[int]bool{newReq.player.Id: false, opp.player.Id: false}, "lobby", time.Now()}
		newReq.match = &newMatch
		opp.match = &newMatch
		matchPool = append(matchPool, &newMatch)
		InformMatchPlayers(&newMatch)

		for i, p := range requestPool {
			if p == opp {
				requestPool = append(requestPool[:i], requestPool[i+1:]...)
			}
		}
	}
}

func FindProperOpponent(level int) *Connection {
    for _, req := range requestPool {
        if req.player.Level == level {
            return req
        }
    }
    return nil
}

func InformMatchPlayers(match *Match) {
    data := make([]PlayerSign, 2)
    data[0] = GetPlayerSignByPlayer(match.Player1.player)
    data[1] = GetPlayerSignByPlayer(match.Player2.player)

	msg := "matchFounded"
	outMessage := struct {
		Message string
		Players []PlayerSign
	} {msg, data}

	outJson, _ := json.Marshal(outMessage)
	writeSocket(string(outJson), match.Player1.conn)
	writeSocket(string(outJson), match.Player2.conn)
}

func SendReadyState(match *Match, username string) {
	outMessage := struct {
		Message string
		Username string
	} {"updateReadyState", username}

	outJson, _ := json.Marshal(outMessage)
	writeSocket(string(outJson), match.Player1.conn)
	writeSocket(string(outJson), match.Player2.conn)
}

//func ResetReadyStates(match *Match) {
//    for key, _ := range match.ReadyState {
//        delete(match.ReadyState, key)
//    }
//}

func ReadyMatch(conn *Connection) {
    conn.match.ReadyState[conn.player.Id] = true
    SendReadyState(conn.match, conn.player.Username)

    if conn.match.ReadyState[conn.match.Player1.player.Id] && conn.match.ReadyState[conn.match.Player2.player.Id] {
        SendInitialMatchData(conn.match)
        go SendStartMatchAfterSeconds(conn.match, 3)
    }
}

func SendInitialMatchData(match *Match) {

}

func SendStartMatchAfterSeconds(match *Match, t time.Duration) {
    time.Sleep(t * time.Second)

    match.State = "playing"
	sendJsonMessageToClient("startMatch", match.Player1)
	sendJsonMessageToClient("startMatch", match.Player2)
}

func LeaveMatch(conn *Connection) {
	outMessage := struct {
		Message string
		Username string
	} {"matchCanceled", conn.player.Username}
	outJson, _ := json.Marshal(outMessage)

	writeSocket(string(outJson), conn.match.Player1.conn)
	writeSocket(string(outJson), conn.match.Player2.conn)

	RemoveMatchFromPool(conn.match)
}

func RemoveMatchFromPool(match *Match) {
    for i, m := range matchPool {
        if m == match {
            matchPool = append(matchPool[:i], matchPool[i + 1:]...)
            m.Player1.match = nil
            m.Player2.match = nil
        }
    }
}
