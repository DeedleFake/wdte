import React, { Component } from 'react'

import AceEditor from 'react-ace'
import './brace'

import {
	Grid,
	Menu,
	Dropdown,
	Container,
} from 'semantic-ui-react'

import { connect } from 'react-redux'
import {
	runWDTE,
} from './store'

const styles = {
	main: {
		display: 'flex',
		flexDirection: 'column',

		padding: 8,

		position: 'absolute',
		top: 0,
		bottom: 0,
		left: 0,
		right: 0,
	},

	column: {
		display: 'flex',
		flexDirection: 'column',
	},

	scroll: {
		overflow: 'auto',
	},
}

class App extends Component {
	state = {
		input: '',
	}

	setVal = (k, f) => (val) => this.setState({
		[k]: (f || ((v) => v))(val),
	})

	onRun = () => {
		this.props.runWDTE(this.state.input)
	}

	render() {
		return (
			<div style={styles.main}>
				<Grid style={styles.grid} columns={2} divided stretched>
					<Grid.Column style={styles.scroll}>
					</Grid.Column>

					<Grid.Column style={styles.column}>
						<Menu inverted>
							<Menu.Item onClick={this.onRun}>Run</Menu.Item>

							<Dropdown item text='Examples'>
								<Dropdown.Menu>
									<Dropdown.Item>Fibonacci</Dropdown.Item>
									<Dropdown.Item>Stream</Dropdown.Item>
									<Dropdown.Item>Strings</Dropdown.Item>
									<Dropdown.Item>Lambda</Dropdown.Item>
									{/*<Dropdown.Item>Canvas</Dropdown.Item>*/}
								</Dropdown.Menu>
							</Dropdown>
						</Menu>

						<Grid columns={1} divided stretched>
							<Grid.Row>
								<Grid.Column style={styles.scroll}>
									<AceEditor
										mode='wdte'
										theme='vibrant_ink'
										value={this.state.input}
										onChange={this.setVal('input')}
									/>
								</Grid.Column>
							</Grid.Row>

							<Grid.Row>
								<Grid.Column style={styles.scroll}>
									<Container>
										<pre>{this.props.output}</pre>
									</Container>
								</Grid.Column>
							</Grid.Row>
						</Grid>
					</Grid.Column>
				</Grid>
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
