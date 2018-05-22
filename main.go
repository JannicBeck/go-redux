package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

func rootReducer(state redux.State, action redux.Action) (redux.State, error) {
	var err error
	if state == nil {
		state = make(map[string]int)
		state.(map[string]redux.State)["counter1"], err = counter.Counter(state, action)
		state.(map[string]redux.State)["counter2"], err = counter.Counter(state, action)
	}
	if action.Payload == "counter1" {
		newState := make(map[string]int)
		s := state.(map[string]redux.State)["counter1"]
		newCount, err := counter.Counter(s, action)
		newState["counter1"] = newCount.(int)
		newState["counter2"] = state.(map[string]int)["counter1"]
		return newState, err
	} else if action.Payload == "counter2" {
		newState := make(map[string]int)
		s := state.(map[string]redux.State)["counter2"]
		newCount, err := counter.Counter(s, action)
		newState["counter1"] = state.(map[string]int)["counter1"]
		newState["counter2"] = newCount.(int)
		return newState, err
	} else {
		return state, err
	}
}

func main() {

	store := redux.CreateStore(counter.Counter)
	var printState redux.Subscriber
	printState = func(state redux.State) {
		fmt.Println(state)
	}
	unsubscribe := store.Subscribe(&printState)

	store.Dispatch(counter.Increment())
	store.Dispatch(counter.Increment())
	unsubscribe()
	store.Dispatch(counter.Increment())
}
