// @format

import React from 'react'
import ReactDOM from 'react-dom'

import 'semantic-ui-css/semantic.min.css'

import jss from 'jss'
import preset from 'jss-preset-default'

import App from './App'

jss.setup(preset())

ReactDOM.render(<App />, document.getElementById('root'))
