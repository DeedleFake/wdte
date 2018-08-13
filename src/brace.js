// @format

import brace from 'brace'
import 'brace/theme/vibrant_ink'

brace.define('ace/mode/wdte', (acequire, exports, module) => {
	let oop = acequire('ace/lib/oop')
	let TextMode = acequire('ace/mode/text').Mode
	let TextHighlightRules = acequire('ace/mode/text_highlight_rules')
		.TextHighlightRules

	let HighlightRules = function() {
		this.$rules = {
			start: [
				{
					token: 'comment',
					regex: '^#.*$',
				},
				{
					token: 'constant',
					regex: '[0-9]+|[0-9]+\\.[0-9]+|\\.[0-9]+',
				},
				{
					token: 'keyword',
					regex: '(\\b(switch|memo|let|import)\\b)',
				},
				{
					token: 'keyword.operator',
					regex: '\\.|->|\\{|\\}|\\[|\\]|\\(|\\)|=>|;|:|--|\\(@',
				},
				{
					token: 'string',
					regex: '((["\']).*\\2)',
				},
			],
		}
	}
	oop.inherits(HighlightRules, TextHighlightRules)

	let Mode = function() {
		this.HighlightRules = HighlightRules
	}
	oop.inherits(Mode, TextMode)

	exports.Mode = Mode
})
