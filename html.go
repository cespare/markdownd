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
  width: 900px;
  margin: 50px auto;
  border: 3px solid #ccc;
  padding: 0px 15px;
}
pre .hll {
  background-color: #ffffcc;
}
pre .c {
  color: #408080;
  font-style: italic;
}
pre .k {
  color: #5d2846;
  font-weight: bold;
}
pre .o {
  color: #666666;
}
pre .cm {
  color: #408080;
  font-style: italic;
}
pre .cp {
  color: #bc7a00;
}
pre .c1 {
  color: #408080;
  font-style: italic;
}
pre .cs {
  color: #408080;
  font-style: italic;
}
pre .gd {
  color: #a00000;
}
pre .ge {
  font-style: italic;
}
pre .gr {
  color: #ff0000;
}
pre .gh {
  color: #000080;
  font-weight: bold;
}
pre .gi {
  color: #00a000;
}
pre .go {
  color: #808080;
}
pre .gp {
  color: #000080;
  font-weight: bold;
}
pre .gs {
  font-weight: bold;
}
pre .gu {
  color: #800080;
  font-weight: bold;
}
pre .gt {
  color: #0040d0;
}
pre .kc {
  color: #5d2846;
  font-weight: bold;
}
pre .kd {
  color: #5d2846;
  font-weight: bold;
}
pre .kn {
  color: #5d2846;
  font-weight: bold;
}
pre .kp {
  color: #5d2846;
}
pre .kr {
  color: #5d2846;
  font-weight: bold;
}
pre .kt {
  color: #b00040;
}
pre .m {
  color: #666666;
}
pre .s {
  color: #4eb25a;
}
pre .na {
  color: #7d9029;
}
pre .nb {
  color: #5d2846;
}
pre .nc {
  color: #3333a0;
  font-weight: bold;
}
pre .no {
  color: #28732c;
}
pre .nd {
  color: #aa22ff;
}
pre .ni {
  color: #999999;
  font-weight: bold;
}
pre .ne {
  color: #d2413a;
  font-weight: bold;
}
pre .nf {
  color: #3333a0;
}
pre .nl {
  color: #a0a000;
}
pre .nn {
  color: #3333a0;
  font-weight: bold;
}
pre .nt {
  color: #5d2846;
  font-weight: bold;
}
pre .nv {
  color: #353c92;
}
pre .ow {
  color: #aa22ff;
  font-weight: bold;
}
pre .w {
  color: #bbbbbb;
}
pre .mf {
  color: #666666;
}
pre .mh {
  color: #666666;
}
pre .mi {
  color: #666666;
}
pre .mo {
  color: #666666;
}
pre .sb {
  color: #4eb25a;
}
pre .sc {
  color: #4eb25a;
}
pre .sd {
  color: #4eb25a;
  font-style: italic;
}
pre .s2 {
  color: #4eb25a;
}
pre .se {
  color: #bb6622;
  font-weight: bold;
}
pre .sh {
  color: #4eb25a;
}
pre .si {
  color: #bb6688;
  font-weight: bold;
}
pre .sx {
  color: #5d2846;
}
pre .sr {
  color: #bb6688;
}
pre .s1 {
  color: #4eb25a;
}
pre .ss {
  color: #353c92;
}
pre .bp {
  color: #5d2846;
}
pre .vc {
  color: #353c92;
}
pre .vg {
  color: #353c92;
}
pre .vi {
  color: #353c92;
}
pre .il {
  color: #666666;
}
</style>
<div id="wrapper">
`

const htmlFooter = `
</div>
</div>
</body>
`
