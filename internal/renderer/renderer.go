package renderer

import (
	"bytes"
	"html/template"
)

// BlockPair holds a source block and its translated counterpart as HTML.
type BlockPair struct {
	Source template.HTML
	Target template.HTML
}

// FontSizes holds the font sizes (in pt) for the HTML template.
type FontSizes struct {
	Body  int // body text
	Head  int // table header
	Code  int // inline code
	Pre   int // code blocks
}

// FontSizePresets maps size names to FontSizes.
var FontSizePresets = map[string]FontSizes{
	"small":  {Body: 9, Head: 10, Code: 8, Pre: 7},
	"medium": {Body: 10, Head: 11, Code: 9, Pre: 8},
	"large":  {Body: 11, Head: 12, Code: 10, Pre: 9},
}

// DefaultFontSize is the default font size preset name.
const DefaultFontSize = "medium"

// TemplateData holds all data for the HTML template.
type TemplateData struct {
	Title       string
	SourceLabel string
	TargetLabel string
	Pairs       []BlockPair
	Fonts       FontSizes
	Attribution bool
}

// Render produces a complete HTML document with a 2-column table layout.
func Render(data TemplateData) (string, error) {
	// Apply default font sizes if not set
	if data.Fonts == (FontSizes{}) {
		data.Fonts = FontSizePresets[DefaultFontSize]
	}

	tmpl, err := template.New("bilingual").Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
