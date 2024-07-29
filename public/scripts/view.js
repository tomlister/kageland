var image_1 = "gopher"
var image_2 = "gopher"
var image_3 = "gopher"
var image_4 = "gopher"
var liked = false

var editor = CodeMirror(document.querySelector("#editor"), {
	value: "",
	mode:  "go",
	lineNumbers: true
})

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

const searchParams = new URLSearchParams(window.location.search)
const params = Object.fromEntries(searchParams.entries())

const waitLoad = () => {
    const el = document.getElementById('viewer');
    if (el?.contentWindow && el?.contentWindow?.compileShader) {
        runShader();
        return true;
    }
    return false;
}

axios.get(`api/shader?id=${params.id}`).catch(err => {
	if (err.response.status == 404) {
		document.querySelector("#loaderContent").innerHTML = "<h1>404 Shader not found. 😔</h1>"
	}
	if (err.response.status == 500) {
		document.querySelector("#loaderContent").innerHTML = "<h1>500 Yeah nah yeah nah, she'll be right mate! 😓</h1>"
	}
	throw err
}).then(res => {
	document.querySelector("#loader").remove()
	document.querySelector("#shaderTitle").textContent = res.data.name
	image_1 = res.data.image_1
	image_2 = res.data.image_2
	image_3 = res.data.image_3
	image_4 = res.data.image_4
	document.querySelector("#image1select").value = image_1
	document.querySelector("#image2select").value = image_2
	document.querySelector("#image3select").value = image_3
	document.querySelector("#image4select").value = image_4
	document.querySelector("#likeCount").innerHTML = res.data.likes
	document.querySelector("#viewCount").innerHTML = res.data.views
	editor.setValue(res.data.frag_shader)
	document.querySelector("#content").style.visibility = 'visible'
	let i = setInterval(() => {
        if (waitLoad()) clearInterval(i);
    }, 100)
})

document.querySelector("#likeButton").addEventListener('click', () => {
	document.querySelector("#likeButton").classList.toggle("react-group-active")
	if (!liked) {
		liked = true
		axios.post(`api/like?id=${params.id}`)
		document.querySelector("#likeButtonIcon").innerHTML = "favorite"
		const count = document.querySelector("#likeCount").innerHTML
		document.querySelector("#likeCount").innerHTML = parseInt(count) + 1
		return
	}
	liked = false
	axios.post(`api/unlike?id=${params.id}`)
	document.querySelector("#likeButtonIcon").innerHTML = "favorite_border"
	const count = document.querySelector("#likeCount").innerHTML
	document.querySelector("#likeCount").innerHTML = Math.max(0, parseInt(count) - 1)
})

document.querySelector("#remixCode").addEventListener('click', () => {
	document.location.href = `/edit?remix=${params.id}`
})

document.querySelector("#image1select").addEventListener('change', event => {image_1 = event.target.value; runShader()})
document.querySelector("#image2select").addEventListener('change', event => {image_2 = event.target.value; runShader()})
document.querySelector("#image3select").addEventListener('change', event => {image_3 = event.target.value; runShader()})
document.querySelector("#image4select").addEventListener('change', event => {image_4 = event.target.value; runShader()})