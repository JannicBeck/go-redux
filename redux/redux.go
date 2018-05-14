package redux

type Action struct {
	Type string
}

type subscriber func(interface{})

type Store struct {
	subscribers []subscriber
	State       interface{}
	reducer     func(state interface{}, action Action) interface{}
}

func (s *Store) ReplaceReducer(reducer func(state interface{}, action Action) interface{}) {
	s.reducer = reducer
}

func (s *Store) GetState() interface{} {
	return s.State
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

	s.State = s.reducer(s.State.(int), action)

	for _, l := range s.subscribers {
		l(s.State)
	}

}
