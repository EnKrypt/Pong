package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

var state GameState
var clients []*websocket.Conn

/* Some of these objects have odd numbered dimensions. This is so that
 * those objects have a well defined center pixel that can denote the
 * position of the object without ambiguity.
 */
const (
	arenaWidth         = 800
	arenaHeight        = 600
	ballSize           = 21
	paddleWidth        = 21
	paddleHeight       = 121
	scoreLimit         = 11
	velocityMultiplier = 5
)

/* The position of the ball and paddles represent the object's center,
 * but we have to detect collision by checking the outer bounds of the
 * objects, so we calculate the permissible bounds for those objects
 * in advance relative to their center.
 */
var ballCollisionBounds = [2][2]int{
	{paddleWidth + ((ballSize - 1) / 2), (ballSize - 1) / 2},
	{arenaWidth - paddleWidth - ((ballSize - 1) / 2), arenaHeight - ((ballSize - 1) / 2)},
}
var paddleMovementBounds = [2]int{(paddleHeight + 1) / 2, arenaHeight - ((paddleHeight + 1) / 2)}

func main() {
	http.HandleFunc("/", pongServer)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pongServer(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	/* We maintain a list of actively connected websocket clients so
	 * that we can broadcast the game state to all connected participants.
	 * When a client gets dropped, we remove it from this list.
	 */
	clients = append(clients, ws)
	defer func() {
		for i, client := range clients {
			if client == ws {
				clients = append(clients[:i], clients[i+1:]...)
				break
			}
		}
	}()
	if !state.inProgress { // If there is no game, signal that we are ready to start one
		if err := ws.WriteMessage(websocket.TextMessage, []byte("READY")); err != nil {
			log.Println(err)
		}
	}
	for {
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		params := strings.Split(string(data), " ")
		if params[0] == "START" && !state.inProgress {
			go start() // Run the goroutine to start a game
		} else if len(params) > 1 {
			/* Refer to keyboardEvents.js file in the frontend folder.
			 * When indicated by the client, set only the motion state.
			 * The position calculation that happens per tick will take
			 * care of the rest.
			 */
			if state.inProgress && params[1] == "0" {
				if params[0] == "W" {
					state.player1.moving = "up"
				} else if params[0] == "S" {
					state.player1.moving = "down"
				} else if params[0] == "U" {
					state.player2.moving = "up"
				} else if params[0] == "D" {
					state.player2.moving = "down"
				}
				/* Stop moving only if the user lifts the key and the paddle
				 * was moving in the direction that they have lifted away from.
				 */
			} else if params[1] == "1" && state.inProgress {
				if params[0] == "W" && state.player1.moving == "up" {
					state.player1.moving = "no"
				} else if params[0] == "S" && state.player1.moving == "down" {
					state.player1.moving = "no"
				} else if params[0] == "U" && state.player2.moving == "up" {
					state.player2.moving = "no"
				} else if params[0] == "D" && state.player2.moving == "down" {
					state.player2.moving = "no"
				}
			}
		}
	}
}

/* This is the heart of the logic.
 * Every time a new game starts, we reset the game state.
 * The server then maintains a tick rate as per a given interval.
 * Every tick, we perform calculations that progress the game by one frame
 * and broadcast the newly calculated game state to all connected clients.
 * User interaction happens concurrently and affects the calculations
 * in the next tick.
 */
func start() {
	state = GameState{inProgress: true, player1: PlayerState{position: arenaHeight / 2, moving: "no"}, player2: PlayerState{position: arenaHeight / 2, moving: "no"}} // The paddles must start from the center of the arena
	state.GenerateNewBallState()
	ticker := time.NewTicker(15 * time.Millisecond) // 15ms frame time means 65+ frames per second
	for state.inProgress {
		state.GameTick()
		<-ticker.C
	}
	ticker.Stop()
}

// Broadcast to all clients that the game has been won by a player and a new game can begin
func endGame(player string) {
	state.inProgress = false
	for _, client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("READY %v", player))); err != nil {
			log.Println(err)
		}
	}
}

// Broadcast the current game state to all connected clients for the frontend to render
func announceState() {
	for _, client := range clients {
		if err := client.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%v %v %v %v %v %v", state.ball.position.x, state.ball.position.y, state.player1.position, state.player2.position, state.player1.score, state.player2.score))); err != nil {
			log.Println(err)
		}
	}
}
