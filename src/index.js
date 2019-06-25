// @format

import React from 'react'
import ReactDOM from 'react-dom'

import 'semantic-ui-css/semantic.min.css'

import Go from './assets/wasm_exec'

import App from './App'

let go = new Go()
WebAssembly.instantiateStreaming(fetch('./wdte.wasm'), go.importObject).then(
	(wasm) => go.run(wasm.instance),
)

ReactDOM.render(<App />, document.getElementById('root'))
