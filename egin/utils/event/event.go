package event

import (
	"sync"
	"time"

	"github.com/daodao97/egin/utils/goroutine"
)

type EventData struct {
	Name string
}

type EventBusStruct struct {
	events    chan EventData
	listeners sync.Map
}

func (bus *EventBusStruct) Trigger(event EventData) {
	bus.events <- event
}

func (bus *EventBusStruct) On(eventName string, handler func(data EventData)) {
	bus.listeners.Store(eventName, handler)
}

var EventBus EventBusStruct

func init() {
	EventBus = EventBusStruct{events: make(chan EventData, 20)}
	go func() {
		for range time.Tick(time.Second) {
			select {
			case data, ok := <-EventBus.events:
				if ok {
					eventName := data.Name
					handler, ok := EventBus.listeners.Load(eventName)
					if ok {
						callback := handler.(func(data EventData))
						goroutine.Go(func() {
							callback(data)
						})
					}
				}
			}
		}
	}()
}
