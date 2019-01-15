package main

const jsHeader = `
<head>
<script type="application/javascript">
var source = new EventSource("/updates");
var handler = function(e) {
	var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function() {
		if (xhr.readyState != 4 || xhr.status != 200) {
			return;
		}
		document.body.innerHTML = xhr.responseText;
	}
	xhr.open("GET", "?nojs=true", true);
	xhr.send();
}
source.addEventListener("message", handler, false);
</script>
</head>
`

const htmlHeader = `
<body>
<style type="text/css">
a {
  color: #4183C4;
  text-decoration: none;
}
a:hover {
  text-decoration: underline;
}
h1 {
  border-bottom: 3px solid #ccc;
  padding-bottom: 10px;
}
body {
  font: 14px / 20px "Helvetica Neue", "Lucida Grande", Helvetica, Arial, Verdana, sans-serif;
}
pre, code {
  font-family: "Ubuntu Mono", Courier, monospace;
  background-color: #F0EEEA;
  padding: 2px;
  overflow: auto;
}
.highlight pre {
  padding-left: 6px;
}
#wrapper {
  max-width: 800px;
  margin: 50px auto;
  border: 3px solid #ccc;
  padding: 0px 15px;
}
</style>
<div id="wrapper">
`

const htmlFooter = `
</div>
</div>
</body>
`
