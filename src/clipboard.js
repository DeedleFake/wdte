// @format

export const copy = (text) => {
	const el = document.createElement('textarea')
	el.value = text

	try {
		document.body.appendChild(el)
		el.focus()
		el.select()

		if (!document.execCommand('copy')) {
			throw new Error('copy failed')
		}
	} finally {
		document.body.removeChild(el)
	}
}
