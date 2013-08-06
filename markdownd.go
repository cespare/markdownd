package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/cespare/blackfriday"
)

// TODO: Syntax highlighting
// TODO: -w
// TODO: Allow for specifying the browser? (bcat has -b for this.)

const tempfile = "/tmp/markdownd_tempfile.html"

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

// bopen opens some (possibly file://) url in a browser.
func bopen(url string) error {
	cmd := exec.Command(openProgram, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func main() {
	if (flag.NArg() == 0 && *watch) || flag.NArg() > 1 {
		usage(1)
	}

	f := os.Stdin

	if flag.NArg() > 0 {
		var err error
		f, err = os.Open(flag.Arg(0))
		if err != nil {
			fatal(err)
		}
		defer f.Close()
	}

	input, err := ioutil.ReadAll(f)
	if err != nil {
		fatal(err)
	}
	rendered := render(input)

	// Embed the output in an HTML page with some nice CSS, unless we're printing the output directly to stdout.
	if *serve || *watch {
		rendered = append([]byte(htmlHeader), rendered...)
		rendered = append(rendered, []byte(htmlFooter)...)
	}

	switch {
	case *watch:
		panic("unimplemented")
	case *serve:
		// Write to a temp file and open it in a browser, then exit.
		temp, err := os.Create(tempfile)
		if err != nil {
			fatal("Could not create a tempfile:", err)
		}
		if _, err := temp.Write(rendered); err != nil {
			fatal(err)
		}
		if err := bopen(temp.Name()); err != nil {
			fatal(err)
		}
	default:
		// Just write to stdout and we're done.
		os.Stdout.Write(rendered)
	}
}
