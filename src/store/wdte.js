import * as actions from './actions'

import * as wdte from '../wdte.go'

export const runWDTE = async (input) => {
	return {
		type: actions.RUN_WDTE,
		output: wdte.eval(input),
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
