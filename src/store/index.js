import {
	createStore,
	combineReducers,
	applyMiddleware,
} from 'redux'

import wdte from './wdte'

export * from './wdte'

const asyncMiddleware = (store) => (next) => (action) => {
	if (typeof(action) === 'function') {
		return action({
			state: store.getState(),
			dispatch: store.dispatch,
		})
	}

	if (typeof(action.then) === 'function') {
		return Promise.resolve(action).then(store.dispatch)
	}

	return next(action)
}

export default createStore(
	combineReducers({
		wdte,
	}),

	applyMiddleware(
		asyncMiddleware,
	),
)
