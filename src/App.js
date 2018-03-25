import React, { Component } from 'react'

import ReactMarkdown from 'react-markdown'

import AceEditor from 'react-ace'
import './brace'

import {
	Menu,
	Dropdown,
} from 'semantic-ui-react'

import initialDesc from './initialDesc.txt'
import * as examples from './examples'

import { connect } from 'react-redux'
import {
	runWDTE,
} from './store'

const styles = {
	main: {
		display: 'flex',
		flexDirection: 'row',
		//flexWrap: 'wrap-reverse',

		backgroundColor: '#EEEEEE',
		boxSizing: 'border-box',
		padding: 8,
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
	}
}

class App extends Component {
	state = {
		description: initialDesc,
		input: '',
	}

	setVal = (k, f) => (val) => this.setState({
		[k]: (f || ((v) => v))(val),
	})

	onRun = () => {
		this.props.runWDTE(this.state.input)
	}

	onClickExample = (ev, data) => {
		this.setState({
			description: examples[data.value].desc,
			input: examples[data.value].input,
		})
	}

	render() {
		return (
			<div style={styles.main}>
				<div style={styles.column}>
					<ReactMarkdown source={this.state.description} />
				</div>

				<div style={styles.column}>
					<Menu inverted>
						<Menu.Item onClick={this.onRun}>Run</Menu.Item>

						<Dropdown item text='Examples'>
							<Dropdown.Menu>
								{Object.entries(examples).map(([id, example]) => (
									<Dropdown.Item value={id} onClick={this.onClickExample}>{example.name}</Dropdown.Item>
								))}
								{/*<Dropdown.Item onClick={this.onClickExample}>Canvas</Dropdown.Item>*/}
							</Dropdown.Menu>
						</Dropdown>
					</Menu>

					<AceEditor
						style={styles.input}
						mode='wdte'
						theme='vibrant_ink'
						value={this.state.input}
						onChange={this.setVal('input')}
					/>

					<pre style={styles.output}>{this.props.output}</pre>
				</div>
			</div>
		)
	}
}

export default connect(
	(state) => ({
		output: state.wdte.output,
	}),

	{
		runWDTE,
	},
)(App)
