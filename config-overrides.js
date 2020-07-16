// @format

module.exports = (config, env) => ({
	...config,
	devServer: {
		...config.devServer,
		mimeTypes: {
			...(config.devServer || {}).mimeTypes,
			wasm: 'application/wasm',
		},
	},
})
