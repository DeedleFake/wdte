import React, { Component } from 'react'

import Typography from 'material-ui/Typography'
import Paper from 'material-ui/Paper'
import Grid from 'material-ui/Grid'
import { withStyles } from 'material-ui/styles'

const styles = (theme) => ({
	root: {
		position: 'absolute',
		top: theme.spacing.unit,
		bottom: theme.spacing.unit,
		left: theme.spacing.unit,
		right: theme.spacing.unit,
	},
})

class App extends Component {
	onRun = () => {
		alert('Not implemented.')
	}

	render() {
		return (
			<React.Fragment>
				<div className={this.props.classes.root}>
					<Grid container>
						<Grid item xs>
							<Paper>This is a test.</Paper>
						</Grid>

						<Grid container direction='column' item xs>
							<Grid item>
								<Typography>
									Toolbar
								</Typography>
							</Grid>

							<Grid item xs>
								<Paper>This is also a test.</Paper>
							</Grid>

							<Grid item xs>
								<Paper>This might not be.</Paper>
							</Grid>
						</Grid>
					</Grid>
				</div>
			</React.Fragment>
		)
	}
}

export default withStyles(styles)(App)
