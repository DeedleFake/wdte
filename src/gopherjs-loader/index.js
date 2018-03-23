const child_process = require('child_process')
const fs = require('fs')
const path = require('path')

module.exports = function(source) {
  const callback = this.async()

  const tmp = path.join(__dirname, 'tmp')
  const fname = path.basename(this.resourcePath, '.go')
  const cmd = `gopherjs build -m ${this.resourcePath} -o '${path.join(tmp, `${fname}.js`)}'`

  child_process.execSync(`rm -fr '${tmp}'`)
  child_process.execSync(`mkdir -p '${tmp}'`)

  child_process.exec(cmd, function(error, stdout, stderr) {
    if (error) {
			return callback(error)
		}

    const out = fs.readFileSync(path.join(tmp, `${fname}.js`), 'utf8')
		const map = fs.readFileSync(path.join(tmp, `${fname}.js.map`), 'utf8')
		child_process.execSync(`rm -fr '${tmp}'`)

    callback(null, out, map)
  })
}
