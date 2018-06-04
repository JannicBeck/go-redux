package redux

type MiddlewareStore struct {
	StoreBase
	middlewareChain []MiddlewareC
}

type Dispatch func(Action) Action

type Middleware func(StoreBase) func(Dispatch) Dispatch
type MiddlewareC func(Dispatch) Dispatch

func (store *MiddlewareStore) Dispatch(action Action) Action {

	// TODO figure out why store.Dispatch loops infinite here
	dispatchWithMiddleware := store.StoreBase.Dispatch

	for _, middleware := range store.middlewareChain {
		dispatchWithMiddleware = middleware(dispatchWithMiddleware)
	}
	return dispatchWithMiddleware(action)
}

func ApplyMiddleware(middleware []Middleware) func(CreateStoreBase) CreateStoreBase {
	return func(createStoreBase CreateStoreBase) CreateStoreBase {

		return func(reducer Reducer, initialState State, onChange OnChange) StoreBase {

			storeBase := createStoreBase(reducer, initialState, onChange)

			var middlewareChain []MiddlewareC

			for _, middleware := range middleware {
				middlewareChain = append(middlewareChain, middleware(storeBase))
			}

			return &MiddlewareStore{
				storeBase,
				middlewareChain,
			}
		}
	}

}
