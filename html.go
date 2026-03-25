package main

const jsHeader = `
<head>
<script type="application/javascript">
const source = new EventSource("/updates");
source.addEventListener("message", async () => {
  const response = await fetch("?nojs=true");
  if (response.ok) {
    document.body.innerHTML = await response.text();
  }
});
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
</body>
`
