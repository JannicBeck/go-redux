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

type Reducer func(State, Action) State

type subscriber func(State)

type Store struct {
	subscribers []subscriber
	state       State
	reducer     func(state State, action Action) State
}

func CreateStore(reducer Reducer) Store {

	if reducer == nil {
		log.Fatal(noReducerProvidedErrMsg)
	}

	initialState := reducer(nil, Action{})
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

func (s *Store) ReplaceReducer(reducer Reducer) {
	s.reducer = reducer
}

func (s *Store) GetState() State {
	return s.state
}

func (s *Store) setState(state State) {
	s.state = state
}

func (s *Store) Subscribe(l subscriber) func() {

	s.subscribers = append(s.subscribers, l)
	i := len(s.subscribers) - 1

	return func() {
		s.subscribers[i] = s.subscribers[len(s.subscribers)-1]
		s.subscribers[len(s.subscribers)-1] = nil
		s.subscribers = s.subscribers[:len(s.subscribers)-1]
	}

}

func (s *Store) Dispatch(action Action) {

	s.state = s.reducer(s.state.(int), action)

	for _, l := range s.subscribers {
		l(s.state)
	}

}
