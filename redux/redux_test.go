package redux

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

func TestCreateStoreWithoutInitialStateFatal(t *testing.T) {

	fn := func() {
		CreateStore(func(state State, action Action) State {
			return state
		})
	}
	command := "-test.run=TestCreateStoreWithoutInitialStateFatal"

	crashTest(t, fn, command, noInitialStateProducedErrMsg)
}

func TestCreateStoreWithoutReducerFatal(t *testing.T) {

	fn := func() {
		CreateStore(nil)
	}
	command := "-test.run=TestCreateStoreWithoutReducerFatal"

	crashTest(t, fn, command, noReducerProvidedErrMsg)
}

func crashTest(t *testing.T, fn func(), command string, errMsg string) {
	// Only run the failing part when a specific env variable is set
	if os.Getenv("BE_CRASHER") == "1" {
		fn()
		return
	}

	// Start the actual test in a different subprocess
	cmd := exec.Command(os.Args[0], command)
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	stdout, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	// Check that the log fatal message is what we expected
	gotBytes, _ := ioutil.ReadAll(stdout)
	got := string(gotBytes)

	if !strings.HasSuffix(got[:len(got)-1], errMsg) {
		t.Fatalf("Unexpected log message. Got %s but should contain %s", got[:len(got)-1], errMsg)
	}

	// Check that the program exited
	err := cmd.Wait()
	if e, ok := err.(*exec.ExitError); !ok || e.Success() {
		t.Fatalf("Process ran with err %v, want exit status 1", err)
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
		if state == nil {
			state = 0
		}
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
