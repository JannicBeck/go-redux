package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

func main() {
	store := redux.CreateStore(counter.Counter)
	unsubscribe := store.Subscribe(func(state redux.State) {
		fmt.Println(state)
	})

	store.Dispatch(counter.Increment())
	store.Dispatch(counter.Increment())
	unsubscribe()
	store.Dispatch(counter.Increment())
}
