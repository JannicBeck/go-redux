package redux

import (
	"log"
)

const noInitialStateProducedErrMsg = `Error: No initialState produced by the supplied reducer.
	Please make sure to check state == nil and assign to it an initial value inside your reducer.
	If you don't know the initial state inside your reducer, you might want to use CreateStoreWithState.`

const noReducerProvidedErrMsg = "Reducer must not be nil"

type Action interface {
	Type() string
}

type State interface{}

type Reducer func(State, Action) (State, error)

type Subscriber func(State, Action)
type Subscribers []*Subscriber

type StoreBase struct {
	isDispatching bool
	reducer       Reducer
	state         State
}

func (store *StoreBase) GetState() State {
	return store.state
}

func (store *StoreBase) Dispatch(action Action, subscribers Subscribers) {

	if store.isDispatching {
		log.Fatal("Reducers may not dispatch actions.")
	}

	state, err := store.reducer(store.state, action)

	if err != nil {
		log.Fatal(err)
	} else {
		store.setState(state)
		onChange(subscribers, store.GetState(), action)
	}
}

func onChange(subscribers Subscribers, state State, action Action) {
	for _, sub := range subscribers {
		(*sub)(state, action)
	}
}

func (store *StoreBase) setState(state State) {
	store.state = state
}

func CreateStoreBase(reducer Reducer, initialState State) StoreBase {
	store := StoreBase{}
	store.reducer = reducer
	store.setState(initialState)
	return store
}

type Store struct {
	storeBase     StoreBase
	subscribers   Subscribers
	reducer       Reducer
	isDispatching bool
}

func (store *Store) ReplaceReducer(nextReducer Reducer) {

	if nextReducer == nil {
		log.Fatal("Expected the nextReducer to be a function.")
	}

	if nextReducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	initialState, err := nextReducer(nil, InitAction{})
	if err != nil {
		log.Fatal("Error when producing initial state")
	}
	if initialState == nil {
		log.Fatal(noInitialStateProducedErrMsg)
	}

	store.storeBase = CreateStoreBase(nextReducer, initialState)

	store.reducer = nextReducer
}

type InitAction struct {
}

func (i InitAction) Type() string {
	return "@@Gorux/INIT"
}

func CreateStore(reducer Reducer) Store {

	if reducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	store := Store{}
	store.ReplaceReducer(reducer)
	return store
}

func (store *Store) GetState() State {
	return store.storeBase.GetState()
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

func (store *Store) Dispatch(action Action) {
	store.storeBase.Dispatch(action, store.subscribers)
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
