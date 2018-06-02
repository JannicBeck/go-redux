package redux

import (
	"log"
)

type OnChange func(Subscribers, State, Action)

type StoreBase struct {
	isDispatching bool
	reducer       Reducer
	state         State
	onChange      OnChange
}

func (store *StoreBase) GetState() State {
	return store.state
}

func (store *StoreBase) Dispatch(action Action, subscribers Subscribers) Action {

	if action.Type == nil {
		log.Fatal(`Actions may not have a nil Type property.`)
	}

	if store.isDispatching {
		log.Fatal("Reducers may not dispatch actions.")
	}

	store.isDispatching = true
	newState, err := store.reducer(store.state, action)

	if err != nil {
		log.Fatal(err)
	} else {
		store.isDispatching = false
		store.setState(newState)
		// TODO find a way to abstract from subscribers here (e.g. don't pass them to Dispatch)
		onChange(subscribers, newState, action)
	}

	return action
}

func (store *StoreBase) setState(state State) {
	store.state = state
}

func CreateStoreBase(reducer Reducer, initialState State, onChange OnChange) StoreBase {
	store := StoreBase{}
	store.setState(initialState)
	store.reducer = reducer
	store.onChange = onChange
	return store
}
