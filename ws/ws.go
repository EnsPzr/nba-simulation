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

var people []*websocket.Conn
var lock *sync.Mutex
var broadcastChan = make(chan Message)
var periodMap map[int][]ResultVM
var finalModel Message

func init() {
	lock = &sync.Mutex{}
	go broadcast()
	finalModel = Message{
		WsEvent: WsEventStart,
		Step:    0,
	}
}

func WsHandler(c *websocket.Conn) {
	wsConnect(c)
	broadcastChan <- finalModel
	for {
		// It does not crash when the first connection is closed. but falls in others.
		if _, msg, err := c.ReadMessage(); err != nil {
			fmt.Println("error => ", err)
			if websocket.IsUnexpectedCloseError(err) {
				wsDisconnect(c)
			}
			break
		} else {
			var event WsEvent
			err = json.Unmarshal(msg, &event)
			if err != nil {
				fmt.Println("Can not read message =>", err.Error())
				continue
			}
			switch event {
			case WsEventStart:
				if finalModel.WsEvent != WsEventDataReading {
					finalModel.WsEvent = WsEventDataReading
					broadcastChan <- finalModel
					initData()
					time.Sleep(time.Second * 5)
					finalModel.WsEvent = WsEventDataReady
					broadcastChan <- finalModel
				}
			case WsEventContinue:
				finalModel.WsEvent = WsEventDataSent
				broadcastChan <- finalModel
				time.Sleep(time.Second * 5)
				startCountingTime()
				finalModel.WsEvent = WsEventEnd
				broadcastChan <- finalModel
			}
		}

	}
}

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

func wsConnect(c *websocket.Conn) {
	lock.Lock()
	people = append(people, c)
	lock.Unlock()

}

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

func broadcast() {
	for {
		msg := <-broadcastChan
		lock.Lock()
		deleted := make([]*websocket.Conn, 0)
		for _, v := range people {
			err := v.WriteJSON(msg)
			if err != nil {
				if strings.Contains(err.Error(), "connection reset by peer") || strings.Contains(err.Error(), "broken pipe") {
					deleted = append(deleted, v)
				}
				fmt.Println("Can not write message =>", err.Error())
			}
		}
		lock.Unlock()
		for _, v := range deleted {
			wsDisconnect(v)
		}
	}
}
