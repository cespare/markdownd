# markdownd

markdownd is a markdown renderer for your command-line. It can also serve the rendered markdown to your
browser and reload it automatically when the file changes.

markdownd is based on [markdown_doctor](https://github.com/dmacdougall/markdown_doctor) (to which I am also a
contributor). However, markdownd is written in Go and is more easily distributed (no need for Ruby and
Rubygems).

## Installation

For now, you need to install from source. Install Go 1.1+, then

    $ go get github.com/cespare/markdownd

## Protips

Here are two lines I use in my `.vimrc` to make it really quick to open the current markdown file in a browser
and have it auto-refresh when the file changes:

``` vimscript
command! Markdownd !markdownd -w % >/dev/null &
noremap <leader>m :Markdownd<cr><cr>
```

Note that when you use `markdownd -w` to watch a file, the server will exit as soon as you close the web
browser tab so you don't need to worry about extra markdownd servers hanging around.
