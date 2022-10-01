
function colorToggle(target) {
	let color = "black"
	return function() {
		if (color == "black") { color = "red" }
		else { color = "black" }
		target.style.color = color
	}
}

hello = document.getElementById("hello")
hello.onclick = colorToggle(hello)
