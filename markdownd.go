package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

// TODO: Allow for specifying the browser? (bcat has -b for this.)

var (
	serve   = flag.Bool("s", false, "Open the output in a browser")
	watch   = flag.Bool("w", false, "Open the output in a browser and watch the input file for changes to reload")
	verbose = flag.Bool("v", false, "Print some debugging information")

	sseHeaders = [][2]string{
		{"Content-Type", "text/event-stream"},
		{"Cache-Control", "no-cache"},
		{"Connection", "keep-alive"},
	}

	mu       sync.RWMutex // protects rendered
	rendered []byte

	baseDir  string // directory containing the markdown file, or cwd
	repoRoot string // root of the containing git repo, or ""
)

func main() {
	log.SetFlags(0)
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, `Usage:

  markdownd [flags] [file]

The flags are:

`)
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, `
Markdownd renders the provided file containing markdown text.
If no filename is given, markdownd reads markdown text from stdin.
If -w is used, a filename must also be given. The -w flag implies the -s flag.
If neither -w nor -s are given, the output is written to stdout.`)
	}
	flag.Parse()

	if flag.NArg() > 0 {
		absPath, err := filepath.Abs(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		baseDir = filepath.Dir(absPath)
	} else {
		var err error
		baseDir, err = os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
	}
	repoRoot = gitRepoRoot(baseDir)

	if (flag.NArg() == 0 && *watch) || flag.NArg() > 1 {
		flag.Usage()
		os.Exit(2)
	}

	if err := renderMarkdown(); err != nil {
		log.Fatal(err)
	}

	switch {
	case *watch:
		updates, err := updateListener(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		url := startServer(updates)
		fmt.Printf("Serving markdown rendered from %s at %s\n", flag.Arg(0), url)
		if err := bopen(url); err != nil {
			log.Fatal(err)
		}
		// Just sit and block infinitely.
		select {}
	case *serve:
		url := startServer(nil)
		if flag.NArg() > 0 {
			fmt.Printf("Serving markdown rendered from %s at %s\n", flag.Arg(0), url)
		} else {
			fmt.Printf("Serving markdown at %s\n", url)
		}
		if err := bopen(url); err != nil {
			log.Fatal(err)
		}
		select {}
	default:
		// Just write to stdout and we're done.
		os.Stdout.Write(rendered)
	}
}

func gitRepoRoot(dir string) string {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

var imageExts = map[string]struct{}{
	".apng": {},
	".avif": {},
	".bmp":  {},
	".gif":  {},
	".ico":  {},
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".svg":  {},
	".tif":  {},
	".tiff": {},
	".webp": {},
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if _, ok := imageExts[strings.ToLower(filepath.Ext(r.URL.Path))]; !ok {
		http.NotFound(w, r)
		return
	}

	rel := strings.TrimPrefix(r.URL.Path, "/")

	// Try relative to the markdown file's directory.
	fsPath := filepath.Join(baseDir, rel)
	if _, err := os.Stat(fsPath); err == nil {
		http.ServeFile(w, r, fsPath)
		return
	}

	// Try relative to the git repo root.
	if repoRoot != "" && repoRoot != baseDir {
		fsPath = filepath.Join(repoRoot, rel)
		if _, err := os.Stat(fsPath); err == nil {
			http.ServeFile(w, r, fsPath)
			return
		}
	}

	http.NotFound(w, r)
}

// render renders markdown text.
func render(input []byte) []byte {
	md := goldmark.New(goldmark.WithExtensions(extension.GFM))
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		// Only error should be on writes (there is no invalid
		// markdown).
		panic(err)
	}
	return buf.Bytes()
}

func renderFromFile(filename string) ([]byte, error) {
	input, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return render(input), nil
}

// bopen opens some (possibly file://) url in a browser.
func bopen(url string) error {
	cmd := exec.Command(openProgram, url)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func updateListener(filename string) (<-chan struct{}, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	c := make(chan struct{})
	go func() {
		for {
			select {
			case e := <-watcher.Events:
				if strings.TrimPrefix(e.Name, "./") != strings.TrimPrefix(filename, "./") {
					continue
				}
				if e.Op&fsnotify.Remove != 0 {
					// Need to reset the watch. No way to report errors
					// if this doesn't work for some reason.
					watcher.Remove(filename)
				}
				// Nonblocking send (if nobody's listening, skip it and move on).
				select {
				case c <- struct{}{}:
				default:
				}
			case err := <-watcher.Errors:
				fmt.Println("fsnotify error:", err)
			}
		}
	}()

	if err := watcher.Add(filepath.Dir(filename)); err != nil {
		return nil, err
	}
	return reRender(c), nil
}

// reRender re-renders the markdown on updates.
// It handles bursty events by batching slightly.
func reRender(updates <-chan struct{}) <-chan struct{} {
	out := make(chan struct{})
	t := time.NewTimer(0)
	<-t.C
	go func() {
		for range updates {
			t.Reset(50 * time.Millisecond)
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
			out <- struct{}{}
		}
	}()
	return out
}

func makeUpdateHandler(update <-chan struct{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			panic("HTTP server does not support Flusher needed for SSE")
		}
		for _, header := range sseHeaders {
			w.Header().Set(header[0], header[1])
		}
		done := r.Context().Done()

		for {
			select {
			case <-update:
				fmt.Fprint(w, "data:update\n\n")
				flusher.Flush()
			case <-done:
				// Exit now because the user closed the page.
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
		input, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		rendered = render(input)
	}

	// Embed the output in an HTML page with some nice CSS
	// unless we're printing the output directly to stdout.
	if *serve || *watch {
		rendered = append([]byte(htmlHeader), rendered...)
		rendered = append(rendered, []byte(htmlFooter)...)
	}
	return nil
}

// startServer serves output on a local webserver running at the returned URL.
func startServer(updates <-chan struct{}) (url string) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		mu.RLock()
		defer mu.RUnlock()
		if r.FormValue("nojs") != "true" {
			w.Write([]byte(jsHeader))
		}
		w.Write(rendered)
	})
	mux.HandleFunc("GET /updates", makeUpdateHandler(updates))
	mux.HandleFunc("GET /", imageHandler)
	return startLocalServer(mux)
}

func startLocalServer(handler http.Handler) (url string) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		if ln, err = net.Listen("tcp6", "[::1]:0"); err != nil {
			log.Fatalf("Failed to listen on a port: %v", err)
		}
	}
	s := &http.Server{Handler: handler}
	go s.Serve(ln)
	return "http://" + ln.Addr().String()
}
