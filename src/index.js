// @format

import React from 'react'
import ReactDOM from 'react-dom'

import 'semantic-ui-css/semantic.min.css'

import App from './App'

import './assets/wasm_exec'

let go = new window.Go()
WebAssembly.instantiateStreaming(fetch('./wdte.wasm'), go.importObject).then(
	(wasm) => go.run(wasm.instance),
)

ReactDOM.render(<App />, document.getElementById('root'))
