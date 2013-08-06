package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cespare/blackfriday"
)

// TODO: Syntax highlighting
// TODO: -s
// TODO: -w
// TODO: Allow for specifying the browser? (bcat has -b for this.)

var (
	serve = flag.Bool("s", false, "Open the output in your browser.")
	watch = flag.Bool("w", false, "Open the output in a browser and watch the input file for changes to reload.")
)

func init() { flag.Parse() }

func usage(status int) {
	fmt.Printf(`Usage:
  $ %s [OPTIONS] [MARKDOWN_FILE]
where OPTIONS are:
`, os.Args[0])
	flag.PrintDefaults()
	fmt.Println(`and MARKDOWN_FILE is some file containing markdown.
- If MARKDOWN_FILE is not given, markdownd will read markdown text from stdin.
- If you specify -w, you must specify MARKDOWN_FILE (stdin doesn't make sense).
- -w implies -s.
- If neither -w nor -s are given, the output is written to stdout.`)
	os.Exit(status)
}

// Render renders some markdown with syntax highlighting. It would be nicer if blackfriday.Markdown operated
// on io.Readers/Writers, but it uses []bytes so we need to fully buffer everything.
func render(input []byte) []byte {
	flags := 0
	flags |= blackfriday.HTML_GITHUB_BLOCKCODE
	renderer := blackfriday.HtmlRenderer(flags, "", "")

	extensions := 0
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_AUTOLINK

	return blackfriday.Markdown(input, renderer, extensions)
}

func fatal(args ...interface{}) {
	fmt.Println(args...)
	os.Exit(1)
}

func main() {
	if (flag.NArg() == 0 && *watch) || flag.NArg() > 1 {
		usage(1)
	}

	if *serve || *watch {
		panic("unimplemented")
	}
	if flag.NArg() > 0 {
		panic("unimplemented")
	}

	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fatal(err)
	}
	os.Stdout.Write(render(input))
}
