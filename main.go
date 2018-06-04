package main

import (
	"fmt"

	"github.com/jannicbeck/redux/counter"
	"github.com/jannicbeck/redux/redux"
)

type StoreBaseLog struct {
	redux.StoreBase
}

func (store *StoreBaseLog) Dispatch(action redux.Action) redux.Action {
	fmt.Println(action.Type())
	return store.StoreBase.Dispatch(action)
}

func logEnhancer(createStoreBase redux.CreateStoreBase) redux.CreateStoreBase {
	return func(reducer redux.Reducer, initialState redux.State, onChange redux.OnChange) redux.StoreBase {

		return &StoreBaseLog{
			createStoreBase(reducer, initialState, onChange),
		}

	}
}

func logger(store redux.StoreBase) func(redux.Dispatch) redux.Dispatch {
	return func(next redux.Dispatch) redux.Dispatch {
		return func(action redux.Action) redux.Action {
			fmt.Println("dispatching", action.Type())
			result := next(action)
			fmt.Println("next state", store.GetState())
			return result
		}
	}
}

func main() {

	var middleware []redux.Middleware
	middleware = append(middleware, logger)

	store := redux.CreateStore(counter.Counter, nil, redux.ApplyMiddleware(middleware))
	var printState redux.Subscriber
	printState = func(state redux.State, action redux.Action) {
		fmt.Println(state, action.Type())
	}
	store.Subscribe(&printState)
	store.Dispatch(counter.Increment{})
	store.Dispatch(counter.Increment{})

	// reducerMap := make(map[string]redux.Reducer)
	// reducerMap["counter"] = counter.Counter
	// todosMap := make(map[string]redux.Reducer)
	// todosMap["todos"] = todos.Todos
	// todosMap["visibilityFilter"] = todos.VisibilityFilter
	// todosReducer := redux.CombineReducers(todosMap)
	// reducerMap["todos"] = todosReducer
	// root := redux.CombineReducers(reducerMap)

	// store := redux.CreateStore(root, nil, nil)
	// var printState redux.Subscriber
	// printState = func(state redux.State, action redux.Action) {
	// 	fmt.Println(state)
	// }
	// unsubscribe := store.Subscribe(&printState)
	// store.Dispatch(todos.AddTodo{Id: "1", Text: "First"})
	// store.Dispatch(todos.AddTodo{Id: "2", Text: "Second"})
	// store.Dispatch(todos.AddTodo{Id: "3", Text: "Third"})
	// store.Dispatch(todos.ToggleTodo{Id: "2"})
	// unsubscribe()

}
