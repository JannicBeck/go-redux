package redux

func CombineReducers(reducers map[string]Reducer) func(State, Action) (State, error) {
	return func(state State, action Action) (State, error) {
		if state == nil {
			state = make(map[string]State)
		}
		var err error
		var nextState State
		nextState = make(map[string]State)
		for key, reducer := range reducers {
			previousStateForKey := state.(map[string]State)[key]
			var nextStateForKey State
			nextStateForKey, err = reducer(previousStateForKey, action)
			nextState.(map[string]State)[key] = nextStateForKey
		}

		return nextState, err
	}
}
