# markdownd

markdownd is a markdown renderer for your command line. It can also serve the
rendered markdown to your browser and reload it automatically when the file
changes.

## Installation

For now, you need to install from source. With Go 1.21+, run

    go install github.com/cespare/markdownd@latest

## Tips

Here are two lines I use in my [Neovim config] to make it quick to open the
current markdown file in a browser and have it auto-refresh when the file
changes:

``` vimscript
command! Markdownd call jobstart(['markdownd', '-w', @%])
noremap <leader>m :Markdownd<CR>
```

Note that when you use `markdownd -w` to watch a file, the server will exit as
soon as you close the web browser tab so you don't need to worry about extra
markdownd servers hanging around.

[Neovim config]: https://github.com/cespare/nvim-config/blob/main/init.vim
