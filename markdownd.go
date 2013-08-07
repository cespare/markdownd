package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/cespare/blackfriday"
	"github.com/howeyc/fsnotify"
)

// TODO: Allow for specifying the browser? (bcat has -b for this.)

const tempfile = "/tmp/markdownd_tempfile.html"

var (
	serve = flag.Bool("s", false, "Open the output in your browser.")
	watch = flag.Bool("w", false, "Open the output in a browser and watch the input file for changes to reload.")
	verbose = flag.Bool("v", false, "Print some debugging information.")

	sseHeaders = [][2]string{
		{"Content-Type", "text/event-stream"},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
	}

	pygmentize     string
	validLanguages = make(map[string]struct{})

	mu       sync.RWMutex // protects rendered
	rendered []byte
)

type debug struct{}
var dbg = debug{}
func (debug) Println(args ...interface{}) {
	if *verbose {
		fmt.Fprintln(os.Stderr, append([]interface{}{"DEBUG:"}, args...)...)
	}
}

func init() {
	flag.Parse()

	var err error
	pygmentize, err = findPygments()
	if err != nil {
		fatal("Pygments could not be loaded:", err)
	}

	rawLexerList, err := exec.Command(pygmentize, "-L", "lexers").Output()
	if err != nil {
		fatal(err)
	}
	for _, line := range bytes.Split(rawLexerList, []byte("\n")) {
		if len(line) == 0 || line[0] != '*' {
			continue
		}
		for _, l := range bytes.Split(bytes.Trim(line, "* :"), []byte(",")) {
			lexer := string(bytes.TrimSpace(l))
			if len(lexer) != 0 {
				validLanguages[lexer] = struct{}{}
			}
		}
	}
}

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

func syntaxHighlight(out io.Writer, in io.Reader, language string) {
	_, ok := validLanguages[language]
	if !ok || language == "" {
		language = "text"
	}
	pygmentsCmd := exec.Command(pygmentize, "-l", language, "-f", "html", "-P", "encoding=utf-8")
	pygmentsCmd.Stdin = in
	pygmentsCmd.Stdout = out
	var stderr bytes.Buffer
	pygmentsCmd.Stderr = &stderr
	if err := pygmentsCmd.Run(); err != nil {
		fatal(err)
	}
}

// Render renders some markdown with syntax highlighting. It would be nicer if blackfriday.Markdown operated
// on io.Readers/Writers, but it uses []bytes so we need to fully buffer everything.
func render(input []byte) []byte {
	flags := 0
	flags |= blackfriday.HTML_GITHUB_BLOCKCODE
	renderer := blackfriday.HtmlRenderer(flags, "", "")
	renderer.SetBlockCodeProcessor(syntaxHighlight)

	extensions := 0
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_NO_EMPTY_LINE_BEFORE_BLOCK

	return blackfriday.Markdown(input, renderer, extensions)
}

func renderFromFile(filename string) ([]byte, error) {
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return render(input), nil
}

func fatal(args ...interface{}) {
	fmt.Fprintln(os.Stderr, args...)
	os.Exit(1)
}

// bopen opens some (possibly file://) url in a browser.
func bopen(url string) error {
	cmd := exec.Command(openProgram, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func updateListener(filename string) (<-chan bool, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	c := make(chan bool)
	go func() {
		for {
			select {
			case e := <-watcher.Event:
				if strings.TrimPrefix(e.Name, "./") != strings.TrimPrefix(filename, "./") {
					continue
				}
				if e.IsDelete() {
					// Need to reset the watch. No way to report errors if this doesn't work for some reason.
					watcher.RemoveWatch(filename)
					watcher.Watch(filename)
				}
				// Nonblocking send (if nobody's listening, skip it and move on).
				select {
				case c <- true:
				default:
				}
			case err := <-watcher.Error:
				fmt.Println("fsnotify error:", err)
			}
		}
	}()

	if err := watcher.Watch(filepath.Dir(filename)); err != nil {
		return nil, err
	}
	return reRender(c), nil
}

// reRender re-renders the markdown on updates. It handles bursty events by batching slightly.
func reRender(updates <-chan bool) <-chan bool {
	out := make(chan bool)
	go func() {
		for _ = range updates {
			t := time.NewTimer(50 * time.Millisecond)
		loop:
			for {
				select {
				case <-updates:
					// drop
				case <-t.C:
					break loop
				}
			}
			if err := renderMarkdown(); err != nil {
				fmt.Println("Warning:", err)
			}
			out <- true
		}
	}()
	return out
}

// Unification of http.ResponseWriter, http.Flusher, and http.CloseNotifier
type HTTPWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(int)
	Flush()
	CloseNotify() <-chan bool
}

func makeUpdateHandler(update <-chan bool) http.HandlerFunc {
	return func(writer http.ResponseWriter, r *http.Request) {
		w, ok := writer.(HTTPWriter)
		if !ok {
			panic("HTTP server does not support Flusher and/or CloseNotifier needed for SSE.")
		}
		closed := w.CloseNotify()
		for _, header := range sseHeaders {
			w.Header().Set(header[0], header[1])
		}

		for {
			select {
			case <-update:
				fmt.Fprint(w, "data:update\n\n")
				w.Flush()
			case <-closed:
				// We're ready to exit now, because the user closed the page.
				os.Exit(0)
			}
		}
	}
}

func renderMarkdown() error {
	mu.Lock()
	defer mu.Unlock()
	if flag.NArg() > 0 {
		var err error
		rendered, err = renderFromFile(flag.Arg(0))
		if err != nil {
			return err
		}
	} else {
		input, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		rendered = render(input)
	}

	// Embed the output in an HTML page with some nice CSS, unless we're printing the output directly to stdout.
	if *serve || *watch {
		rendered = append([]byte(htmlHeader), rendered...)
		rendered = append(rendered, []byte(htmlFooter)...)
	}
	return nil
}

// serve serves output on a local webserver. The caller should Close the server it gets back.
func makeServer(updates <-chan bool) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		if r.FormValue("nojs") != "true" {
			w.Write([]byte(jsHeader))
		}
		w.Write(rendered)
	})
	mux.HandleFunc("/updates", makeUpdateHandler(updates))
	return httptest.NewServer(mux)
}

func main() {
	if (flag.NArg() == 0 && *watch) || flag.NArg() > 1 {
		usage(1)
	}

	if err := renderMarkdown(); err != nil {
		fatal(err)
	}

	switch {
	case *watch:
		updates, err := updateListener(flag.Arg(0))
		if err != nil {
			fatal(err)
		}
		server := makeServer(updates)
		defer server.Close()
		fmt.Printf("Serving markdown rendered from %s at %s\n", flag.Arg(0), server.URL)
		if err := bopen(server.URL); err != nil {
			fatal(err)
		}
		// Just sit and block infinitely.
		select {}
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
