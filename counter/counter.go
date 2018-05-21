package counter

import "github.com/jannicbeck/redux/redux"

const IncrementType = "Increment"
const DecrementType = "Decrement"

func Counter(state redux.State, action redux.Action) redux.State {
	if state == nil {
		state = 0
	}
	switch action.Type {
	case IncrementType:
		return state.(int) + 1
	case DecrementType:
		return state.(int) - 1
	default:
		return state
	}
}

func Increment() redux.Action {
	return redux.Action{Type: IncrementType}
}

func Decrement() redux.Action {
	return redux.Action{Type: DecrementType}
}
