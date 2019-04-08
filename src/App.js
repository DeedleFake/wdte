// @format

import React, { useState, useReducer, useCallback, useMemo } from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group'

import ReactMarkdown from 'react-markdown'

import AceEditor from 'react-ace'
import './brace'

import { makeStyles } from '@material-ui/styles'

import { Menu, Dropdown, Message } from 'semantic-ui-react'

import initial from './initial'
import * as examples from './examples'

import * as wdte from './wdte'
import * as clipboard from './clipboard'

const useStyles = makeStyles((theme) => ({
	'@font-face': {
		fontFamily: 'Go Mono',
		src: 'url(assets/Go-Mono.ttf)',
	},

	main: {
		display: 'flex',
		flexDirection: 'row',
		//flexWrap: 'wrap-reverse',

		backgroundColor: '#EEEEEE',
		boxSizing: 'border-box',
		padding: 8,
		width: '100%',
		height: '100%',
		position: 'absolute',
	},

	column: {
		display: 'flex',
		flexDirection: 'column',

		flex: 1,
		margin: 8,
		overflowY: 'auto',
		minWidth: 300,
	},

	message: {
		marginBottom: '8px !important',
	},

	output: {
		minHeight: 300,
		fontFamily: 'Go-Mono',
		fontSize: 12,
		flex: 1,
		overflow: 'auto',
		padding: 8,
		boxShadow: 'inset 4px 4px 4px #AAAAAA',
		borderRadius: 8,
		backgroundColor: '#CCCCCC',
	},

	slide: {
		'&.enter': {
			transition: 'all 300ms',
			overflow: 'hidden',

			maxHeight: 0,

			'&.active': {
				maxHeight: '500px',
			},
		},

		'&.exit': {
			transition: 'all 300ms',
			overflow: 'hidden',

			maxHeight: '500px',

			'&.active': {
				maxHeight: 0,
			},
		},
	},
}))

const App = (props) => {
	const classes = useStyles()

	const [description, setDescription] = useState(initial.desc)
	const [input, setInput] = useState(() =>
		window.location.hash !== ''
			? decodeURIComponent(window.location.hash.substr(1))
			: initial.input,
	)
	const [output, setOutput] = useState('')

	const encodedInput = useMemo(() => encodeURIComponent(input), [input])

	const [messages, dispatchMessages] = useReducer(
		(state, action) => {
			switch (action.$) {
				case 'add':
					if (action.timeout != null) {
						setTimeout(() => {
							dispatchMessages({
								$: 'remove',
								id: state.id,
							})
						}, action.timeout)
					}

					return {
						...state,
						id: state.id + 1,
						[state.id]: {
							type: action.type,
							msg: action.msg,
						},
					}

				case 'remove':
					return Object.entries(state)
						.filter(([k, v]) => k !== action.id.toString())
						.reduce((acc, [k, v]) => ({ ...acc, [k]: v }), {})

				default:
					return state
			}
		},
		{ id: 0 },
	)

	const addMessage = useCallback(
		(type, msg, timeout = 3000) => {
			dispatchMessages({ $: 'add', type, msg, timeout })
		},
		[messages],
	)

	const runCode = useCallback(async () => {
		try {
			setOutput(await wdte.run(input))
		} catch (err) {
			setOutput(err.toString())
		}
	}, [input])

	const share = useCallback(() => {
		try {
			clipboard.copy(
				`${window.location.origin}${window.location.pathname}#${encodedInput}`,
			)
			addMessage('success', 'Link successfully copied to clipboard.')
		} catch (err) {
			addMessage('error', `Failed to copy to clipboard: ${err.toString()}`)
		}
	}, [encodedInput])

	return (
		<div className={classes.main}>
			<TransitionGroup component="div" className={classes.column}>
				{Object.entries(messages)
					.filter(([id, msg]) => !isNaN(id))
					.map(([id, msg]) => (
						<CSSTransition
							key={id}
							classNames={{
								enter: 'enter',
								enterActive: 'active',
								exit: 'exit',
								exitActive: 'active',
							}}
							timeout={300}
						>
							<Message
								className={[classes.message, classes.slide].join(' ')}
								{...{ [msg.type]: true }}
							>
								<p>{msg.msg}</p>
							</Message>
						</CSSTransition>
					))}

				<ReactMarkdown source={description} />
			</TransitionGroup>

			<div className={classes.column}>
				<Menu inverted>
					<Menu.Item onClick={runCode}>Run</Menu.Item>

					<Dropdown item text="Examples">
						<Dropdown.Menu>
							{Object.entries(examples).map(([id, example]) => (
								<Dropdown.Item
									key={id}
									value={id}
									onClick={(ev, data) => {
										setDescription(examples[data.value].desc)
										setInput(examples[data.value].input)
									}}
								>
									{example.name}
								</Dropdown.Item>
							))}
							{/*<Dropdown.Item onClick={onClickExample}>Canvas</Dropdown.Item>*/}
						</Dropdown.Menu>
					</Dropdown>

					<Menu.Item position="right" onClick={share}>
						Share
					</Menu.Item>
				</Menu>

				{/* TODO: Find a way to use classes instead. */}
				<AceEditor
					style={{
						width: null,
						height: null,
						minHeight: 300,
						flex: 1,
						borderRadius: 8,
					}}
					mode="wdte"
					theme="vibrant_ink"
					value={input}
					onChange={(val) => setInput(val)}
				/>

				<pre className={classes.output}>{output}</pre>
			</div>
		</div>
	)
}

export default App
