export default {
	run: async (input) => {
		return await new Promise((resolve, reject) => {
			window.WDTE.run(input, (err, output) => {
				if (err != null) {
					return reject(err)
				}

				resolve(output)
			})
		})
	},
}
