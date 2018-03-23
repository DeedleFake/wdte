const path = require('path')

module.exports = {
	loaders: [
		{
			test: /\.txt$/,
			use: 'raw-loader',
		},
		{
			test: /\.go$/,
			use: [
				{
					loader: path.resolve('gopherjs-loader'),
				},
			],
		},
	],
}
