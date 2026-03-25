# Heading 1

## Heading 2

### Heading 3

#### Heading 4

##### Heading 5

###### Heading 6

## Paragraphs and Inline Formatting

This is a plain paragraph with some **bold text**, some *italic text*, and some
***bold italic text***. Here is some `inline code` in the middle of a sentence.
You can also use ~~strikethrough~~ to indicate deleted text.

This is a second paragraph to verify spacing between paragraphs. It contains a
[link to example.com](http://example.com) and a bare autolink: https://example.com.

## Blockquotes

> This is a blockquote. It can contain **bold**, *italic*, and `code`.
>
> It can also span multiple paragraphs.

> > And blockquotes can be nested.

## Unordered Lists

- Item one
- Item two
  - Nested item A
  - Nested item B
    - Deeply nested item
- Item three

## Ordered Lists

1. First item
2. Second item
   1. Nested first
   2. Nested second
3. Third item

## Task Lists

- [x] Completed task
- [ ] Incomplete task
- [x] Another completed task
  - [ ] Nested incomplete task
  - [x] Nested completed task

## Code

Inline `code` was shown above. Here is a fenced code block:

```
package main

import "fmt"

func main() {
    fmt.Println("Hello, world!")
    for i := 0; i < 10; i++ {
        fmt.Printf("i = %d\n", i)
    }
}
```

And an indented code block:

    def hello():
        print("This is an indented code block")
        return True

## Tables

| Left-aligned | Center-aligned | Right-aligned |
| :----------- | :------------: | ------------: |
| Row 1 Col 1  | Row 1 Col 2    | Row 1 Col 3   |
| Row 2 Col 1  | Row 2 Col 2    | Row 2 Col 3   |
| Row 3 Col 1  | Row 3 Col 2    | Row 3 Col 3   |

A simple table with no alignment:

| Name    | Value | Description         |
| ------- | ----- | ------------------- |
| alpha   | 1     | The first item      |
| beta    | 2     | The second item     |
| gamma   | 3     | The third item      |

## Horizontal Rules

Content above the rule.

---

Content between rules.

***

Content below the rules.

## Images

![Alt text for a missing image](missing.png)

## Links

- [Regular link](http://example.com)
- [Link with title](http://example.com "Example Title")
- Autolink: https://example.com
- Email autolink: user@example.com

## HTML Passthrough

<details>
<summary>Click to expand</summary>

This is content inside a `<details>` element.

</details>

<dl>
  <dt>Definition Term</dt>
  <dd>Definition description.</dd>
  <dt>Another Term</dt>
  <dd>Another description.</dd>
</dl>

## Long Content

This paragraph has a very long line that should test how the wrapper and text handle overflow within the max-width container, since the body font is set to 14px with a 20px line height and the wrapper is constrained to 800px max-width with 15px of horizontal padding.

```
This is a long line inside a code block that should trigger horizontal scrolling: aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa
```
