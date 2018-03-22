const child_process = require('child_process')
const fs = require('fs')
const path = require('path')

module.exports = function(source) {
  const callback = this.async()
  const tmp = path.join(__dirname, 'tmp')
  const fname = path.basename(this.resourcePath, '.go')
  const cmd = `gopherjs build ${this.resourcePath} -o '${path.join(tmp, `${fname}.js`)}'`
  child_process.execSync(`rm -rf '${tmp}'`)
  child_process.execSync(`mkdir -p '${tmp}'`)

  child_process.exec(cmd, function(error, stdout, stderr) {
    if (error) {
			return callback(error, null)
		}

    const out = fs.readFileSync(path.join(tmp, `${fname}.js`), 'utf8')
		const map = fs.readFileSync(path.join(tmp, `${fname}.js.map`), 'utf8')
    callback(null, out)
  })
}
