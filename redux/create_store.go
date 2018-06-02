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

type InitAction struct {
}

func (i InitAction) Type() string {
	return "@@Gorux/INIT"
}

type State interface{}

type Reducer func(State, Action) (State, error)

type Subscriber func(State, Action)
type Subscribers []*Subscriber
type Unsubscribe func()

type OnChange func(State, Action)

// TODO: to interface?
type StoreBase struct {
	GetState func() State
	Dispatch func(Action) Action
}

// TODO: to interface?
type Store struct {
	GetState       func() State
	ReplaceReducer func(Reducer)
	Subscribe      func(subscriber *Subscriber) Unsubscribe
	Dispatch       func(Action) Action
}

func CreateStoreBase(reducer Reducer, initialState State, onChange OnChange) StoreBase {

	store := StoreBase{}
	var state State

	GetState := func() State {
		return state
	}

	var isDispatching bool

	Dispatch := func(action Action) Action {

		if isDispatching {
			log.Fatal("Reducers may not dispatch actions.")
		}

		isDispatching = true
		newState, err := reducer(state, action)

		if err != nil {
			log.Fatal(err)
		} else {
			isDispatching = false
			state = newState
			onChange(newState, action)
		}

		return action
	}

	state = initialState
	store.GetState = GetState
	store.Dispatch = Dispatch
	return store
}

func CreateStore(reducer Reducer) Store {

	if reducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	initialState, err := reducer(nil, InitAction{})

	if err != nil {
		log.Fatal("Error when producing initial state")
	}
	if initialState == nil {
		log.Fatal(noInitialStateProducedErrMsg)
	}

	store := Store{}

	var subscribers Subscribers

	onChange := func(state State, action Action) {
		for _, sub := range subscribers {
			(*sub)(state, action)
		}
	}

	storeBase := CreateStoreBase(reducer, initialState, onChange)

	ReplaceReducer := func(nextReducer Reducer) {

		if nextReducer == nil {
			log.Fatal("Expected the nextReducer to be a function.")
		}

		reducer = nextReducer
	}

	GetState := func() State {
		return storeBase.GetState()
	}

	Subscribe := func(subscriber *Subscriber) Unsubscribe {

		if *subscriber == nil {
			log.Fatal(`Subscriber must not be nil`)
		}

		subscribers = append(subscribers, subscriber)
		isSubscribed := true

		return func() {
			if !isSubscribed {
				return
			}

			for i := len(subscribers) - 1; i >= 0; i-- {
				sub := subscribers[i]

				if sub == subscriber {
					subscribers[i] = subscribers[len(subscribers)-1]
					subscribers[len(subscribers)-1] = nil
					subscribers = subscribers[:len(subscribers)-1]
					break
				}

			}

			isSubscribed = false

		}

	}

	Dispatch := func(action Action) Action {
		return storeBase.Dispatch(action)
	}

	store.ReplaceReducer = ReplaceReducer
	store.GetState = GetState
	store.Subscribe = Subscribe
	store.Dispatch = Dispatch

	return store
}
