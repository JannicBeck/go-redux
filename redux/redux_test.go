package redux

import (
	"testing"
)

const IncrementType = "Increment"
const DecrementType = "Decrement"

func counter(state State, action Action) State {
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

func Increment() Action {
	return Action{Type: IncrementType}
}

func Decrement() Action {
	return Action{Type: DecrementType}
}

func TestCreateStore(t *testing.T) {

	store := CreateStore(counter)

	if store.reducer == nil {
		t.Error("Reducer must not be nil")
	}

	if store.GetState() != 0 {
		t.Errorf("State expected: %v got %v", 0, store.GetState())
	}

}

func TestStore(t *testing.T) {

	store := CreateStore(counter)

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

func TestSubscription(t *testing.T) {

	store := CreateStore(counter)
	callbackCount := 0
	unsubscribe := store.Subscribe(func(state State) {
		callbackCount = callbackCount + 1
	})

	if len(store.subscribers) != 1 {
		t.Errorf("Expected subscriber count to be %v got %v", 1, len(store.subscribers))
	}

	store.Dispatch(Increment())
	if callbackCount != 1 {
		t.Errorf("Subscribe state: %v, want %v", callbackCount, 1)
	}

	store.Dispatch(Increment())

	if callbackCount != 2 {
		t.Errorf("Subscribe state: %v, want %v", callbackCount, 2)
	}

	unsubscribe()

	store.Dispatch(Increment())

	if callbackCount != 2 {
		t.Errorf("Unsubscribed state: %v, want %v", callbackCount, 2)
	}

}

func TestReplaceReducer(t *testing.T) {

	store := CreateStore(counter)
	store.ReplaceReducer(func(state State, action Action) State {
		switch action.Type {
		case IncrementType:
			return state.(int) + 10
		case DecrementType:
			return state.(int) - 10
		default:
			return state
		}
	})

	store.Dispatch(Increment())
	if store.GetState() != 10 {
		t.Errorf("%v got %v", 10, store.GetState())
	}

}
