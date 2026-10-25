package translate
// Translate markdown source into HTML using goldmark's parser and renderer
// Default style is Github Flavoured Markdown (GFM).

import (
	"github.com/yuin/goldmark/v2/parser"
	"github.com/yuin/goldmark/v2/renderer/html"
	"github.com/yuin/goldmark/v2/extension"

	"bytes"
)
var (
	mdParser = parser.New(
		parser.WithExtensions(extension.GFMParser),
	)
	htmlRenderer = html.New(
		html.WithExtensions(extension.GFMHTMLRenderer),
	)
)

// Markdown: returns Markdown source converted to HTML as a string
func Markdown(source []byte) (string, error) {
	doc := mdParser.Parse(source)
	var buf bytes.Buffer

	if err := htmlRenderer.Render(&buf, source, doc); err != nil {
		return "", err
	}

	return buf.String(), nil
}

