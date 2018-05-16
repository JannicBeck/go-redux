package redux

import (
	"testing"
)

const IncrementType = "Increment"
const DecrementType = "Decrement"

func Counter(state interface{}, action Action) interface{} {
	switch action.Type {
	case IncrementType:
		return state.(int) + 1
	case DecrementType:
		return state.(int) - 1
	default:
		return state
	}
}

func Increment() Action {
	return Action{Type: IncrementType}
}

func Decrement() Action {
	return Action{Type: DecrementType}
}

func Test(t *testing.T) {

	store := Store{State: 0}
	store.ReplaceReducer(Counter)

	var tests = []struct {
		a Action
		s int
	}{
		{Increment(), 1},
		{Increment(), 2},
		{Decrement(), 1},
	}
	for _, c := range tests {
		store.Dispatch(c.a)
		got := store.GetState()
		if got != c.s {
			t.Errorf("Dispatch(%q) == %q, want %q", c.a, got, c.s)
		}
	}

}
