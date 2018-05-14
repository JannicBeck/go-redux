package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

func main() {
	state := counter.Counter(1, counter.Increment())
	state = counter.Counter(state, counter.Decrement())
	store := redux.Store{State: 0}
	store.ReplaceReducer(counter.Counter)
	unsubscribe := store.Subscribe(func(state interface{}) {
		fmt.Println(state)
	})
	store.Dispatch(counter.Increment())
	store.Dispatch(counter.Increment())
	unsubscribe()
	store.Dispatch(counter.Increment())
}
