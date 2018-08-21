// @format

import React, { Component } from 'react'
import { CSSTransition, TransitionGroup } from 'react-transition-group'

import ReactMarkdown from 'react-markdown'

import AceEditor from 'react-ace'
import './brace'

import injectSheet from 'react-jss'

import { Menu, Dropdown, Message } from 'semantic-ui-react'

import Clipboard from 'clipboard'

import initial from './initial'
import * as examples from './examples'

import * as wdte from './wdte'

const styles = {
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

	input: {
		width: null,
		height: null,
		minHeight: 300,
		flex: 1,
		borderRadius: 8,
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
}

class App extends Component {
	state = {
		description: initial.desc,
		input: initial.input,
		output: '',

		messages: {},
	}

	setVal = (k, f) => (val) =>
		this.setState({
			[k]: (f || ((v) => v))(val),
		})

	componentDidMount() {
		this.clipboard = new Clipboard('#share', {
			text: () =>
				`${window.location.origin}${
					window.location.pathname
				}#${encodeURIComponent(this.state.input)}`,
		})

		this.clipboard.on('success', (ev) => {
			this.addMessage('success', 'Link successfully copied to clipboard.')
			window.location.hash = `#${encodeURIComponent(this.state.input)}`
		})

		this.clipboard.on('error', (ev) => {
			this.addMessage('error', 'Failed to copy to clipboard.')
		})
	}

	componentWillUnmount() {
		this.clipboard.destroy()
	}

	addMessage = (type, msg, timeout = 3000) => {
		let id = new Date().getTime().toString()

		this.setState(
			{
				messages: {
					...this.state.messages,
					[id]: {
						type,
						msg,
					},
				},
			},
			() =>
				setTimeout(
					() =>
						this.setState({
							messages: Object.entries(this.state.messages).reduce(
								(acc, [k, v]) => (k !== id ? { ...acc, [k]: v } : acc),
								{},
							),
						}),
					timeout,
				),
		)
	}

	onRun = async () => {
		try {
			this.setState({
				output: await wdte.run(this.state.input),
			})
		} catch (err) {
			this.setState({
				output: err.toString(),
			})
		}
	}

	onClickExample = (ev, data) => {
		this.setState({
			description: examples[data.value].desc,
			input: examples[data.value].input,
		})
	}

	render() {
		return (
			<div className={this.props.classes.main}>
				<TransitionGroup component="div" className={this.props.classes.column}>
					{Object.entries(this.state.messages).map(([id, msg]) => (
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
								className={[
									this.props.classes.message,
									this.props.classes.slide,
								].join(' ')}
								{...{ [msg.type]: true }}
							>
								<p>{msg.msg}</p>
							</Message>
						</CSSTransition>
					))}

					<ReactMarkdown source={this.state.description} />
				</TransitionGroup>

				<div className={this.props.classes.column}>
					<Menu inverted>
						<Menu.Item onClick={this.onRun}>Run</Menu.Item>

						<Dropdown item text="Examples">
							<Dropdown.Menu>
								{Object.entries(examples).map(([id, example]) => (
									<Dropdown.Item
										key={id}
										value={id}
										onClick={this.onClickExample}
									>
										{example.name}
									</Dropdown.Item>
								))}
								{/*<Dropdown.Item onClick={this.onClickExample}>Canvas</Dropdown.Item>*/}
							</Dropdown.Menu>
						</Dropdown>

						<Menu.Item as="a" position="right" id="share">
							Share
						</Menu.Item>
					</Menu>

					{/* TODO: Find a way to use this.props.classes instead. */}
					<AceEditor
						style={styles.input}
						mode="wdte"
						theme="vibrant_ink"
						value={this.state.input}
						onChange={this.setVal('input')}
					/>

					<pre className={this.props.classes.output}>{this.state.output}</pre>
				</div>
			</div>
		)
	}
}

export default injectSheet(styles)(App)
