package redux

import (
	"log"
)

const noInitialStateProducedErrMsg = `Error: No initialState produced by the supplied reducer.
	Please make sure to check state == nil and assign to it an initial value inside your reducer.
	If you don't know the initial state inside your reducer, you might want to use CreateStoreWithState.`

const noReducerProvidedErrMsg = "Reducer must not be nil"

type Action struct {
	Type    string
	Payload interface{}
}

type State interface{}

type Reducer func(State, Action) (State, error)

type Subscriber func(State)

type Store struct {
	subscribers   []*Subscriber
	state         State
	reducer       Reducer
	isDispatching bool
}

func CreateStore(reducer Reducer) Store {

	if reducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	initialState, err := reducer(nil, Action{})
	if err != nil {
		log.Fatal("Error when producing initial state")
	}
	if initialState == nil {
		log.Fatal(noInitialStateProducedErrMsg)
	}
	store := Store{}
	store.setState(initialState)
	store.ReplaceReducer(reducer)
	return store
}

func CreateStoreWithState(reducer Reducer, preloadedState State) Store {
	store := CreateStore(reducer)
	store.setState(preloadedState)
	return store
}

func (store *Store) ReplaceReducer(nextReducer Reducer) {

	if nextReducer == nil {
		log.Fatal("Expected the nextReducer to be a function.")
	}

	store.reducer = nextReducer
}

func (store *Store) GetState() State {
	if store.isDispatching {
		log.Fatal(`You may not call store.getState() while the reducer is executing.
			The reducer has already received the state as an argument.
			Pass it down from the top reducer instead of reading it from the store.`)
	}
	return store.state
}

func (store *Store) setState(state State) {
	store.state = state
}

func (store *Store) Subscribe(subscriber *Subscriber) func() {

	if *subscriber == nil {
		log.Fatal(`Subscriber must not be nil`)
	}

	if store.isDispatching {
		log.Fatal(`You may not call store.subscribe() while the reducer is executing.
			If you would like to be notified after the store has been updated, subscribe from a
			component and invoke store.getState() in the callback to access the latest state.
			See https://redux.js.org/api-reference/store#subscribe(listener) for more details.`)
	}

	addSubscriber(store, subscriber)
	isSubscribed := true

	return func() {
		if !isSubscribed {
			return
		}

		if store.isDispatching {
			log.Fatal(`You may not unsubscribe from a store listener while the reducer is executing.
				See https://redux.js.org/api-reference/store#subscribe(listener) for more details.`)
		}

		removeSubscriber(store, subscriber)
		isSubscribed = false

	}

}

func addSubscriber(store *Store, subscriber *Subscriber) {
	store.subscribers = append(store.subscribers, subscriber)
}

func removeSubscriber(store *Store, subscriber *Subscriber) {

	for i := len(store.subscribers) - 1; i >= 0; i-- {
		sub := store.subscribers[i]

		if sub == subscriber {
			store.subscribers[i] = store.subscribers[len(store.subscribers)-1]
			store.subscribers[len(store.subscribers)-1] = nil
			store.subscribers = store.subscribers[:len(store.subscribers)-1]
			break
		}

	}

}

func (store *Store) Dispatch(action Action) {

	if store.isDispatching {
		log.Fatal("Reducers may not dispatch actions.")
	}

	state, err := store.reducer(store.state.(int), action)

	// or log.Fatal?
	if err != nil {
		log.Print(err)
	} else {
		store.setState(state)
		notifySubscribers(store)
	}

}

func notifySubscribers(store *Store) {
	for _, sub := range store.subscribers {
		(*sub)(store.state)
	}
}
