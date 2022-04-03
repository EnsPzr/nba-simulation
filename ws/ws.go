package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"sort"
	"strings"
	"sync"
	"time"
)

// This variable is used to store the list of all the connected clients.
var people []*websocket.Conn

// This variable was used to avoid conflicts.
var lock *sync.Mutex

// This channel was used to broadcast.
var broadcastChan = make(chan Message)

// This map store all period state.
var periodMap map[int][]ResultVM

// This structure is used return result.
var finalModel Message

func init() {
	// İnitialize lock
	lock = &sync.Mutex{}
	// initialize broadcast loop
	go broadcast()
	// initialize final model
	finalModel = Message{
		WsEvent: WsEventStart,
		Step:    0,
	}
}

// WsHandler This function call when client connect.
func WsHandler(c *websocket.Conn) {
	// Client connected then add to people list.
	wsConnect(c)
	// Client connected then send message to clients.
	broadcastChan <- finalModel

	// Read message from client.
	for {
		// It does not crash when the first connection is closed. but falls in others.
		// If any error, client connection close.
		if _, msg, err := c.ReadMessage(); err != nil {
			fmt.Println("error => ", err)
			if websocket.IsUnexpectedCloseError(err) {
				// Client connection close. Remove from people list.
				wsDisconnect(c)
			}
			break
		} else {
			var event WsEvent
			// Read message body.
			err = json.Unmarshal(msg, &event)
			if err != nil {
				fmt.Println("Can not read message =>", err.Error())
				continue
			}
			switch event {
			// If client send start message, then start simulation
			case WsEventStart:
				if finalModel.WsEvent != WsEventDataReading {
					// Change status and send message to clients.
					finalModel.WsEvent = WsEventDataReading
					broadcastChan <- finalModel
					// initialize period map
					initData()
					time.Sleep(time.Second * 5)
					// Change status and send data ready message to clients.
					finalModel.WsEvent = WsEventDataReady
					broadcastChan <- finalModel
				}
				// If client send continue message, simulation is started.
			case WsEventContinue:
				// Change status and send message to clients.
				finalModel.WsEvent = WsEventDataSent
				broadcastChan <- finalModel
				time.Sleep(time.Second * 5)
				// start 240 second simulation. and every 5 second send data to clients.
				startCountingTime()
				// Change status and send message to clients. (Simulation is finished)
				finalModel.WsEvent = WsEventEnd
				broadcastChan <- finalModel
				// If client send restart message, then restart simulation.
			case WsEventRestart:
				finalModel = Message{
					WsEvent: WsEventDataReady,
				}
				broadcastChan <- finalModel
			}
		}

	}
}

// This function send every 5 second data to clients.
// Function read data from periodMap and send to clients.
// Function send message then 5 second delay.
func startCountingTime() {
	for i := 0; i < 48; i++ {
		finalModel.Step = i + 1
		finalModel.RealTime = i + 1
		finalModel.VirtualTime = (i + 1) * 5
		periods := periodMap[i]
		sort.Slice(periods, func(i, j int) bool {
			return periods[i].GameID < periods[j].GameID
		})
		finalModel.ResultVM = periods
		broadcastChan <- finalModel
		if i != 47 {
			time.Sleep(time.Second * 5)
		}
	}
}

// This function append new client to people list.
// Use lock to avoid conflicts.
func wsConnect(c *websocket.Conn) {
	lock.Lock()
	people = append(people, c)
	lock.Unlock()

}

// This function remove client to people list.
// Use lock to avoid conflicts.
func wsDisconnect(c *websocket.Conn) {
	lock.Lock()
	for i, v := range people {
		if v == c {
			people = append(people[:i], people[i+1:]...)
			break
		}
	}
	lock.Unlock()
}

// This function send message to all clients.
// Function read broadcastChan and send message to all clients.
// Use lock. Beacuse to prevent array size from changing while looping.
// If array size change, it will cause panic.
func broadcast() {
	for {
		msg := <-broadcastChan
		lock.Lock()
		deleted := make([]*websocket.Conn, 0)
		for _, v := range people {
			err := v.WriteJSON(msg)
			if err != nil {
				// The first person to leave cannot be caught.
				// This block of code ensures that people who leave are caught.
				if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "broken pipe") {
					deleted = append(deleted, v)
				}
				fmt.Println("Can not write message =>", err.Error())
			}
		}
		lock.Unlock()
		// If any leave, remove from people list.
		for _, v := range deleted {
			wsDisconnect(v)
		}
	}
}
