package parser

import (
	"testing"
)

func TestConvertJiraToTgMarkup(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "plain text",
			input: "This is plain text.",
			want:  "This is plain text\\.",
		},
		{
			name:  "bold formatting",
			input: "Start *bold* end.",
			// Bold tokens are represented with "*" in Telegram MarkdownV2.
			want: "Start *bold* end\\.",
		},
		{
			name:  "italic formatting",
			input: "This _is_ italic.",
			// Italic tokens use "_" as both markers.
			want: "This _is_ italic\\.",
		},
		{
			name:  "strike formatting",
			input: "This -struck- text.",
			// Strike token marker "-" turns into "~" in Telegram MarkdownV2.
			want: "This ~struck~ text\\.",
		},
		{
			name:  "underline formatting",
			input: "This +underline+ text.",
			// Underline token marker "+" turns into "__" in Telegram MarkdownV2.
			want: "This __underline__ text\\.",
		},
		{
			name:  "citation formatting",
			input: "A ??citation?? here.",
			// Citation token "??" uses "_" as Telegram Markdown marker.
			want: "A _citation_ here\\.",
		},
		{
			name:  "escaped marker",
			input: `Escaped \*not bold* remains.`,
			want:  "Escaped \\*not bold* remains\\.*",
		},
		{
			name:  "missing closing marker",
			input: "Unclosed *bold text",
			want:  "Unclosed *bold text*",
		},
		{
			name:  "nested formatting",
			input: "This is *bold and _italic_* text.",
			want:  "This is *bold and _italic_* text\\.",
		},
		{
			name:  "nested formatting italic strike bold ",
			input: "This is _italic -strike- *bold*_ text.",
			want:  "This is _italic ~strike~ *bold*_ text\\.",
		},
		{
			name:  "monospace formatting",
			input: "This is {{monospace}} text.",
			want:  "This is `monospace` text\\.",
		},
		{
			name:  "monospace formatting with text formatting",
			input: "This is {{mon*osp*ace}} text.",
			want:  "This is `mon*osp*ace` text\\.",
		},
		{
			name:  "text formatting with monospace",
			input: "This is *bold {{mon*osp*ace}}* _text_.",
			want:  "This is *bold \\{\\{mon*osp*ace\\}\\}* _text_\\.",
		},
		{
			name:  "images and attachments",
			input: "This is ! image.png! and !attachment.pdf!",
			want:  "This is \\! image\\.png\\! and ",
		},
		{
			name:  "color formatting",
			input: "This is {color:red}*_tetx_*{color}.",
			want:  "This is *_tetx_*\\.",
		},
		{
			name:  "quote formatting",
			input: "This is {quote}*_tetx_*{quote}.",
			want: `This is 
>*_tetx_*
\.`,
		},
		{
			name: "quote formatting",
			input: `This is 
{quote}
*_tetx_*
{quote}.
`,
			want: `This is 

>
>*_tetx_*
>
\.
`,
		},
		{
			name:  "anchor formatting",
			input: "This is {anchor:myanchor} adasd.",
			want:  "This is  adasd\\.",
		},
		{
			name:  "file link formatting",
			input: "This is [file://path/to/file] adasd.",
			want:  "This is  adasd\\.",
		},
		{
			name:  "user link formatting",
			input: "This is [~username] adasd.",
			want:  "This is  adasd\\.",
		},
		{
			name:  "attachment link formatting",
			input: "This is [^attachment.txt] adasd.",
			want:  "This is  adasd\\.",
		},
		{
			name:  "anchor link formatting",
			input: "This is [#anchor] adasd.",
			want:  "This is  adasd\\.",
		},
		{
			name:  "link formatting",
			input: "This is [https://example.com] adasd.",
			want:  "This is [https://example\\.com](https://example.com) adasd\\.",
		},
		{
			name:  "link formatting with text",
			input: "This is [Ex:*am*ple |https://example.com] adasd.",
			want:  "This is [Ex:*am*ple ](https://example.com) adasd\\.",
		},
		{
			name:  "link formatting without end",
			input: "This is [Ex:*am*ple |https://example.com",
			want:  "This is [Ex:*am*ple ](https://example.com)",
		},
		{
			name:  "link formatting without end 2",
			input: "This is [https://example.com adasd",
			want:  "This is [https://example\\.com adasd](https://example.com)",
		},
		{
			name: "preformatted text",
			input: `asdad {noformat} asdasd
		sdad
			asdasd {noformat}`,
			want: "asdad ``` asdasd\n\t\tsdad\n\t\t\tasdasd ```",
		},
		{
			name: "preformatted text 2",
			input: `asdad {noformat} asdasd
		*_sdad_*
			asdasd {noformat}`,
			want: "asdad ``` asdasd\n\t\t*_sdad_*\n\t\t\tasdasd ```",
		},
		{
			name: "code block",
			input: `asdad {code}class HelloWorld {
    public static void main(String[] args) {
        System.out.println("Hello, World!"); 
    }
}{code}`,
			want: "asdad \n```java\nclass HelloWorld {\n    public static void main(String[] args) {\n        System.out.println(\"Hello, World!\"); \n    }\n}```",
		},
		{
			name: "code block",
			input: `asdad {code:go}
func main() {
    fmt.Println("Hello, World!")
}{code}`,
			want: "asdad \n```go\n\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}```",
		},
		{
			name: "list formatting",
			input: `* Item 1
* Item 2
*# Subitem 1
*# Subitem 2
* Item 3
** Subitem 1
** Subitem 2`,
			want: "* Item 1\n* Item 2\n*# Subitem 1\n*# Subitem 2\n* Item 3\n** Subitem 1\n** Subitem 2",
		},
	}

	for _, tt := range tests {
		got := ConvertJiraToTgMarkup(tt.input)
		if got != tt.want {
			t.Errorf("%s: ParseInline(%q) = %q, want %q", tt.name, tt.input, got, tt.want)
		}
	}
}

func TestNestedFormatting(t *testing.T) {
	// Nested formatting: Bold wrapping italic.
	input := "*bold and _italic_* text"
	// Expected: Each token is replaced by its Telegram Markdown equivalent.
	// Bold token with "*" and italic with "_".
	want := "*bold and _italic_* text"
	got := ConvertJiraToTgMarkup(input)
	if got != want {
		t.Errorf("Nested formatting: ParseInline(%q) = %q, want %q", input, got, want)
	}
}
