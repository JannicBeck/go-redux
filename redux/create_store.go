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

type storeBase struct {
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

func (store *storeBase) Dispatch(action Action) Action {

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

func (store *storeBase) GetState() State {
	return store.state
}

func createStoreBase(reducer Reducer, initialState State, onChange OnChange) StoreBase {

	return &storeBase{
		false,
		reducer,
		initialState,
		onChange,
	}

}

type CreateStoreBase func(Reducer, State, OnChange) StoreBase

func CreateStore(reducer Reducer, initialState State, enhancer func(CreateStoreBase) CreateStoreBase) Store {

	if reducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	if initialState == nil {
		initialReducerState, err := reducer(nil, InitAction{})

		if err != nil {
			log.Fatal("Error when producing initial state with reducer")
		}

		if initialReducerState == nil {
			log.Fatal(noInitialStateProducedErrMsg)
		}

		initialState = initialReducerState
	}

	var subscribers Subscribers

	onChange := func(state State, action Action) {
		for _, sub := range subscribers {
			(*sub)(state, action)
		}
	}

	if enhancer == nil {
		enhancer = func(x CreateStoreBase) CreateStoreBase {
			return x
		}
	}

	createFinalStoreBase := enhancer(createStoreBase)

	storeBase := createFinalStoreBase(reducer, initialState, onChange)

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

	ReplaceReducer := func(nextReducer Reducer) {

		if nextReducer == nil {
			log.Fatal("Expected the nextReducer to be a function.")
		}

		var nextInitialState State
		if storeBase != nil {
			nextInitialState = GetState()
		} else {
			nextInitialState = initialState
		}

		storeBase = createStoreBase(nextReducer, nextInitialState, onChange)
		Dispatch(InitAction{})

		reducer = nextReducer
	}

	return Store{
		GetState,
		ReplaceReducer,
		Subscribe,
		Dispatch,
	}

}
