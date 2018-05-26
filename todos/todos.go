package todos

import "github.com/jannicbeck/redux/redux"

type AddTodo struct {
	Id   string
	Text string
}

func (a AddTodo) Type() string {
	return "ADD_TODO"
}

type ToggleTodo struct {
	Id string
}

func (a ToggleTodo) Type() string {
	return "TOGGLE_TODO"
}

type Todo struct {
	Id        string
	Text      string
	Completed bool
}

func Todos(state redux.State, action redux.Action) (redux.State, error) {
	if state == nil {
		state = make([]Todo, 0)
	}
	switch action.Type() {
	case "ADD_TODO":
		todo := Todo{Id: action.(AddTodo).Id, Text: action.(AddTodo).Text, Completed: false}
		return append(state.([]Todo), todo), nil
	case "TOGGLE_TODO":
		var todoIdx int
		var toggleThisTodo Todo
		for i, todo := range state.([]Todo) {
			if todo.Id == action.(ToggleTodo).Id {
				todoIdx = i
				toggleThisTodo = todo
			}
		}
		toggledTodo := Todo{Id: toggleThisTodo.Id, Text: toggleThisTodo.Text, Completed: !toggleThisTodo.Completed}
		newState := make([]Todo, len(state.([]Todo)))
		copy(newState, state.([]Todo))
		newState[todoIdx] = toggledTodo
		return newState, nil
	default:
		return state, nil
	}
}

type SetVisibilityFilter struct {
	filter string
}

func (filter SetVisibilityFilter) Type() string {
	return "SET_VISIBILITY_FILTER"
}

func VisibilityFilter(state redux.State, action redux.Action) (redux.State, error) {
	if state == nil {
		state = "SHOW_ALL"
	}
	switch action.Type() {
	case "SET_VISIBILITY_FILTER":
		return action.(SetVisibilityFilter).filter, nil
	default:
		return state, nil
	}
}
