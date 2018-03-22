import * as actions from './actions'

export const runWDTE = async (input) => {
	return {
		type: actions.RUN_WDTE,
		output: 'Not implemented.',
	}
}

const inital = {
	output: '',
}

export default (state = inital, action) => {
	switch (action.type) {
		case actions.RUN_WDTE:
			return {
				...state,
				output: action.output,
			}

		default:
			return state
	}
}
