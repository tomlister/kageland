var errors = []
var panel = null

var image_1 = "gopher"
var image_2 = "gopher"
var image_3 = "gopher"
var image_4 = "gopher"

const searchParams = new URLSearchParams(window.location.search)
const params = Object.fromEntries(searchParams.entries())

CodeMirror.registerHelper("lint", "go", text => {
	console.log("linting:", text);
	return errors
})

var editor = CodeMirror(document.querySelector("#editor"), {
	value: `// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build ignore

//kage:unit pixels

package main

var Time float
var Cursor vec2

func Fragment(dstPos vec4, srcPos vec2, color vec4) vec4 {
	pos := (dstPos.xy - imageDstOrigin()) / imageDstSize()
	pos += Cursor / imageDstSize() / 4
	clr := 0.0
	clr += sin(pos.x*cos(Time/15)*80) + cos(pos.y*cos(Time/15)*10)
	clr += sin(pos.y*sin(Time/10)*40) + cos(pos.x*sin(Time/25)*40)
	clr += sin(pos.x*sin(Time/5)*10) + sin(pos.y*sin(Time/35)*80)
	clr *= sin(Time/10) * 0.5
	return vec4(clr, clr*0.5, sin(clr+Time/3)*0.75, 1)
}`,
	mode:  "go",
	lineNumbers: true
})

if (params.remix != null) {
	axios.get(`api/shader?id=${params.remix}`).catch(err => {
		if (err.response.status == 404) {
			document.querySelector("#loaderContent").innerHTML = "<h1>404 Shader not found. 😔</h1>"
		}
		if (err.response.status == 500) {
			document.querySelector("#loaderContent").innerHTML = "<h1>500 Yeah nah yeah nah, she'll be right mate! 😓</h1>"
		}
		throw err
	}).then(res => {
		document.querySelector("#loader").remove()
		document.querySelector("#shaderName").value = res.data.name
		image_1 = res.data.image_1
		image_2 = res.data.image_2
		image_3 = res.data.image_3
		image_4 = res.data.image_4
		document.querySelector("#image1select").value = image_1
		document.querySelector("#image2select").value = image_2
		document.querySelector("#image3select").value = image_3
		document.querySelector("#image4select").value = image_4
		editor.setValue(res.data.frag_shader)
		document.querySelector("#content").style.visibility = 'visible'
	})
} else {
	document.querySelector("#loader").remove()
	document.querySelector("#content").style.visibility = 'visible'
}

const runShader = () => {
	const code = editor.getValue('\n')
	const response = document.querySelector("#viewer").contentWindow.compileShader(code, image_1, image_2, image_3, image_4)
	console.log(response)
	const firstError = document.querySelector("#errors").firstChild
	if (firstError !== null) {
		firstError.remove()
	}
	if (response !== null) {
		const parseError = response.match(/(\d+):(\d+):\s/)
		if (parseError !== null) {
			const split = parseError[0].split(':')
			const lineNumber = split[0]
			console.log(lineNumber)
			const charNumber = parseInt(split[1])
			var node = document.createElement("div")
			node.className = "error"
			node.textContent = response
			document.querySelector("#errors").appendChild(node)
		} else {
			var node = document.createElement("div")
			node.className = "error"
			node.textContent = response
			document.querySelector("#errors").appendChild(node)
		}
	}
}

const resetTime = () => {
	document.querySelector("#viewer").contentWindow.resetTime()
}

document.querySelector("#runCode").addEventListener('click', runShader)

document.querySelector("#resetTime").addEventListener('click', resetTime)

document.querySelector("#shareCode").addEventListener('click', () => {
	if (!this.disabled) {
		const code = editor.getValue('\n')
		const shaderName = document.querySelector("#shaderName").value
		const btn = document.querySelector("#shareCode")
		btn.innerHTML = `<span class="spinner-border spinner-border-perfect" role="status" aria-hidden="true"></span>`
		btn.disabled = true
		axios.post('api/shader', {
			name: shaderName,
			frag_shader: code,
			image_1: image_1,
			image_2: image_2,
			image_3: image_3,
			image_4: image_4,
		}).catch(err => {
			throw err
		}).then(res => {
			window.location.href = 'shader?id=' + res.data.id
		})
	}
})



document.querySelector("#image1select").addEventListener('change', event => {image_1 = event.target.value; runShader()})
document.querySelector("#image2select").addEventListener('change', event => {image_2 = event.target.value; runShader()})
document.querySelector("#image3select").addEventListener('change', event => {image_3 = event.target.value; runShader()})
document.querySelector("#image4select").addEventListener('change', event => {image_4 = event.target.value; runShader()})