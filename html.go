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
:root {
  --standard-radius: 4px;
}
body {
  font: 14px / 20px "Helvetica Neue", "Lucida Grande", Helvetica, Arial, Verdana, sans-serif;
  background-color: #ccc;
}

a {
  color: #4183C4;
  text-decoration: none;

  &:hover {
    text-decoration: underline;
  }
}

pre, code {
  font-family: "Ubuntu Mono", Courier, monospace;
}

code {
  background-color: #f5f5f4;
  padding: 2px 4px;
  border-radius: var(--standard-radius);
}

pre {
  background-color: #f5f5f4;
  padding: 10px 10px;
  border-radius: var(--standard-radius);
  overflow: auto;

  code {
    padding: 0;
    background: none;
    border-radius: 0;
  }
}

blockquote {
  margin-left: 0;
  padding-left: 12px;
  border-left: 4px solid #ddd;
  color: #555;
}

table {
  border-collapse: collapse;

  th, td {
    padding: 4px 8px;
    border: 1px solid #ddd;
  }

  th {
    background-color: #f5f5f4;
    font-weight: 600;
  }

  tbody tr:nth-child(even) {
    background-color: #fafafa;
  }
}

hr {
  border: none;
  border-top: 2px solid #ccc;
  margin: 12px 0;
}

img {
  max-width: 100%;
}

dl {
  dt {
    font-weight: 600;
  }

  dd {
    margin-bottom: 8px;
  }
}

#wrapper {
  max-width: 800px;
  margin: 50px auto;
  padding: 16px;
  background-color: white;
  border-radius: var(--standard-radius);
}
</style>
<div id="wrapper">
`

const htmlFooter = `
</div>
</body>
`
