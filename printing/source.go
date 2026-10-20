package printing

import "strings"

// Str is s as a Typst string literal placed in markup (#"..."): Typst shows
// it verbatim, so contract text and term values never become markup (no
// headings, emphasis, math or code from user input). Only \ and " need
// escaping inside the literal; newlines become \n.
func Str(s string) string { return "#" + Lit(s) }

// Lit is s as a Typst string literal for code ("..."), e.g. a function
// argument.
func Lit(s string) string {
	var b strings.Builder
	b.WriteString(`"`)
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
		case '\t':
			b.WriteString(`\t`)
		default:
			b.WriteRune(r)
		}
	}
	b.WriteString(`"`)
	return b.String()
}

// Content is s as Typst content ([#"..."]), for function arguments.
func Content(s string) string { return "[" + Str(s) + "]" }
