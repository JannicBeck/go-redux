package counter

import "github.com/jannicbeck/redux/redux"

const IncrementType = "Increment"
const DecrementType = "Decrement"

func Counter(state redux.State, action redux.Action) (redux.State, error) {
	if state == nil {
		state = 0
	}
	switch action.Type() {
	case IncrementType:
		return state.(int) + 1, nil
	case DecrementType:
		return state.(int) - 1, nil
	default:
		return state, nil
	}
}

type Increment struct {
}

func (inc Increment) Type() string {
	return IncrementType
}

type Decrement struct {
}

func (dec Decrement) Type() string {
	return DecrementType
}
