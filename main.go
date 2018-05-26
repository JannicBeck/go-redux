package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

func combineReducers(reducers map[string]redux.Reducer) func(redux.State, redux.Action) (redux.State, error) {
	return func(state redux.State, action redux.Action) (redux.State, error) {
		if state == nil {
			state = make(map[string]redux.State)
		}
		hasChanged := false
		var err error
		var nextState redux.State
		nextState = make(map[string]redux.State)
		for key, reducer := range reducers {
			previousStateForKey := state.(map[string]redux.State)[key]
			var nextStateForKey redux.State
			nextStateForKey, err = reducer(previousStateForKey, action)
			nextState.(map[string]redux.State)[key] = nextStateForKey
			hasChanged = hasChanged || nextStateForKey != previousStateForKey
		}
		if hasChanged {
			return nextState, err
		}
		return state, err
	}
}

// func rootReducer(state redux.State, action redux.Action) (redux.State, error) {
// 	var err error
// 	if state == nil {
// 		var initialState redux.State
// 		initialState = make(map[string]redux.State)
// 		s1 := initialState.(map[string]redux.State)["counter1"]
// 		s2 := initialState.(map[string]redux.State)["counter2"]
// 		initialState.(map[string]redux.State)["counter1"], err = counter.Counter(s1, action)
// 		initialState.(map[string]redux.State)["counter2"], err = counter.Counter(s2, action)
// 		return initialState, err
// 	}
// 	if action.Payload == "counter1" {
// 		var newState redux.State
// 		newState = make(map[string]redux.State)
// 		s := state.(map[string]redux.State)["counter1"]
// 		newCount, err := counter.Counter(s, action)
// 		newState.(map[string]redux.State)["counter1"] = newCount
// 		newState.(map[string]redux.State)["counter2"] = state.(map[string]redux.State)["counter2"]
// 		return newState, err
// 	} else if action.Payload == "counter2" {
// 		var newState redux.State
// 		newState = make(map[string]redux.State)
// 		s := state.(map[string]redux.State)["counter1"]
// 		newCount, err := counter.Counter(s, action)
// 		newState.(map[string]redux.State)["counter1"] = newCount
// 		newState.(map[string]redux.State)["counter2"] = state.(map[string]redux.State)["counter2"]
// 		return newState, err
// 	} else {
// 		return state, err
// 	}
// }

func main() {

	// reducerMap := make(map[string]redux.Reducer)
	// reducerMap["counter"] = counter.Counter
	// reducerMap["todos"] = todos.Todos
	// root := combineReducers(reducerMap)

	store := redux.CreateStore(counter.Counter)
	var printState redux.Subscriber
	printState = func(state redux.State) {
		fmt.Println(state)
	}
	unsubscribe := store.Subscribe(&printState)

	store.Dispatch(counter.Increment{})
	store.Dispatch(counter.Increment{})
	unsubscribe()
	store.Dispatch(counter.Increment{})
}
