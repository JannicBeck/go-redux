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

type StoreBase interface {
	GetState() State
	Dispatch(Action) Action
}

type DynamicStoreBase struct {
	isDispatching bool
	reducer       Reducer
	state         State
	onChange      func(State, Action)
}

type Store struct {
	GetState       func() State
	ReplaceReducer func(Reducer)
	Subscribe      func(subscriber *Subscriber) Unsubscribe
	Dispatch       func(Action) Action
}

func (store DynamicStoreBase) Dispatch(action Action) Action {

	if store.isDispatching {
		log.Fatal("Reducers may not dispatch actions.")
	}

	store.isDispatching = true
	newState, err := store.reducer(store.state, action)

	if err != nil {
		log.Fatal(err)
	} else {
		store.isDispatching = false
		store.state = newState
		store.onChange(newState, action)
	}

	return action
}

func (store DynamicStoreBase) GetState() State {
	return store.state
}

func CreateStoreBase(reducer Reducer, initialState State, onChange OnChange) StoreBase {

	return DynamicStoreBase{
		false,
		reducer,
		initialState,
		onChange,
	}

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

	return Store{
		GetState,
		ReplaceReducer,
		Subscribe,
		Dispatch,
	}

}
