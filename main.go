package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

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
