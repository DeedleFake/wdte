// @format

import React, { Component } from 'react'

import ReactMarkdown from 'react-markdown'

import AceEditor from 'react-ace'
import './brace'

import injectSheet from 'react-jss'

import { Menu, Dropdown } from 'semantic-ui-react'

import initialDesc from './initialDesc'
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
}

class App extends Component {
	state = {
		description: initialDesc,
		input: '',
		output: '',
	}

	setVal = (k, f) => (val) =>
		this.setState({
			[k]: (f || ((v) => v))(val),
		})

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
				<div className={this.props.classes.column}>
					<ReactMarkdown source={this.state.description} />
				</div>

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
