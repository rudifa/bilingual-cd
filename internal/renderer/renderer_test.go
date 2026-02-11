package renderer

import (
	"html/template"
	"strings"
	"testing"
)

func TestRender_BasicOutput(t *testing.T) {
	data := TemplateData{
		Title:       "Test Document",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<h1>Bonjour</h1>"),
				Target: template.HTML("<h1>Hola</h1>"),
			},
			{
				Source: template.HTML("<p>Le monde est beau.</p>"),
				Target: template.HTML("<p>El mundo es hermoso.</p>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Check basic structure
	checks := []string{
		"<!DOCTYPE html>",
		"<title>Test Document</title>",
		"French",
		"Spanish",
		"<h1>Bonjour</h1>",
		"<h1>Hola</h1>",
		"<p>Le monde est beau.</p>",
		"<p>El mundo es hermoso.</p>",
		"<table>",
		"</table>",
	}

	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("rendered HTML should contain %q", check)
		}
	}
}

func TestRender_EmptyPairs(t *testing.T) {
	data := TemplateData{
		Title:       "Empty",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs:       []BlockPair{},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	if !strings.Contains(html, "<tbody>") {
		t.Error("should contain tbody even with no pairs")
	}
}

func TestRender_ParagraphTextPresent(t *testing.T) {
	data := TemplateData{
		Title:       "Paragraph Check",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<h1>Titre</h1>"),
				Target: template.HTML("<h1>Título</h1>"),
			},
			{
				Source: template.HTML("<p>Ceci est un document de test en français.</p>"),
				Target: template.HTML("<p>Este es un documento de prueba en francés.</p>"),
			},
			{
				Source: template.HTML("<p>Le printemps est la saison du renouveau.</p>"),
				Target: template.HTML("<p>La primavera es la estación de la renovación.</p>"),
			},
			{
				Source: template.HTML(`<blockquote><p>La vie est belle.</p></blockquote>`),
				Target: template.HTML(`<blockquote><p>La vida es hermosa.</p></blockquote>`),
			},
			{
				Source: template.HTML("<div class=\"note\">\n<p>Paragraphe en <strong>HTML brut</strong>.</p>\n</div>"),
				Target: template.HTML("<div class=\"nota\">\n<p>Párrafo en <strong>HTML simple</strong>.</p>\n</div>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Every paragraph's text content must appear in the rendered HTML
	paragraphTexts := []string{
		"Ceci est un document de test en français.",
		"Este es un documento de prueba en francés.",
		"Le printemps est la saison du renouveau.",
		"La primavera es la estación de la renovación.",
		"La vie est belle.",
		"La vida es hermosa.",
		"Paragraphe en",
		"HTML brut",
		"Párrafo en",
		"HTML simple",
	}

	for _, text := range paragraphTexts {
		if !strings.Contains(html, text) {
			t.Errorf("rendered HTML should contain paragraph text %q", text)
		}
	}

	// Paragraph tags must not be empty
	if strings.Contains(html, "<p></p>") {
		t.Error("rendered HTML should not contain empty <p></p> tags")
	}

	// Must not contain raw-HTML-omitted comments
	if strings.Contains(html, "<!-- raw HTML omitted -->") {
		t.Error("rendered HTML should not contain '<!-- raw HTML omitted -->' comments")
	}
}

func TestRender_ListFormatPreserved(t *testing.T) {
	data := TemplateData{
		Title:       "List Check",
		SourceLabel: "French",
		TargetLabel: "Spanish",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<ul>\n<li>Pain frais</li>\n<li>Fromage de chèvre</li>\n<li>Vin rouge</li>\n</ul>"),
				Target: template.HTML("<ul>\n<li>pan fresco</li>\n<li>queso de cabra</li>\n<li>vino tinto</li>\n</ul>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Both source and target must have <li> items, not bare text in a <p>
	for _, item := range []string{
		"<li>Pain frais</li>",
		"<li>Fromage de chèvre</li>",
		"<li>Vin rouge</li>",
		"<li>pan fresco</li>",
		"<li>queso de cabra</li>",
		"<li>vino tinto</li>",
	} {
		if !strings.Contains(html, item) {
			t.Errorf("rendered HTML should contain list item %q", item)
		}
	}

	// Target list must not be rendered as a plain paragraph
	if strings.Contains(html, "<p>pan fresco") {
		t.Error("translated list should be <ul>/<li>, not a plain <p>")
	}
}

func TestRender_LinkPreserved(t *testing.T) {
	data := TemplateData{
		Title:       "Link Check",
		SourceLabel: "Français",
		TargetLabel: "Español",
		Pairs: []BlockPair{
			{
				Source: template.HTML(`<p>Visitez <a href="https://example.com">Example</a> pour chercher.</p>`),
				Target: template.HTML(`<p>Visite <a href="https://example.com">Example</a> para buscar.</p>`),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Both source and target must have the <a> link with href
	checks := []string{
		`<a href="https://example.com">Example</a>`,
		"Visitez",
		"Visite",
	}
	for _, check := range checks {
		if !strings.Contains(html, check) {
			t.Errorf("rendered HTML should contain %q", check)
		}
	}

	// The href must not be escaped
	if strings.Contains(html, "&lt;a href") {
		t.Error("link tags should not be HTML-escaped")
	}
}

func TestRender_HTMLEscaping(t *testing.T) {
	// The template.HTML type should NOT escape the content
	data := TemplateData{
		Title:       "Test",
		SourceLabel: "EN",
		TargetLabel: "FR",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<p>Hello <strong>world</strong></p>"),
				Target: template.HTML("<p>Bonjour <strong>monde</strong></p>"),
			},
		},
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should contain unescaped HTML tags
	if !strings.Contains(html, "<strong>world</strong>") {
		t.Error("template.HTML content should not be escaped")
	}
}

func TestRender_FontSizePresets(t *testing.T) {
	basePairs := []BlockPair{
		{
			Source: template.HTML("<p>Hello</p>"),
			Target: template.HTML("<p>Bonjour</p>"),
		},
	}

	tests := []struct {
		name  string
		fonts FontSizes
		body  string
		head  string
		code  string
		pre   string
	}{
		{"small", FontSizePresets["small"], "font-size: 9pt", "font-size: 10pt", "font-size: 8pt", "font-size: 7pt"},
		{"medium", FontSizePresets["medium"], "font-size: 10pt", "font-size: 11pt", "font-size: 9pt", "font-size: 8pt"},
		{"large", FontSizePresets["large"], "font-size: 11pt", "font-size: 12pt", "font-size: 10pt", "font-size: 9pt"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := TemplateData{
				Title:       "Font Size Test",
				SourceLabel: "EN",
				TargetLabel: "FR",
				Pairs:       basePairs,
				Fonts:       tc.fonts,
			}

			html, err := Render(data)
			if err != nil {
				t.Fatalf("Render failed: %v", err)
			}

			// Each font-size declaration should appear in the CSS
			for _, expected := range []string{tc.body, tc.head, tc.code, tc.pre} {
				if !strings.Contains(html, expected) {
					t.Errorf("rendered HTML should contain %q for %s preset", expected, tc.name)
				}
			}
		})
	}
}

func TestRender_DefaultFontSize(t *testing.T) {
	// When Fonts is zero-value, Render should apply the default (medium) preset
	data := TemplateData{
		Title:       "Default Font Size",
		SourceLabel: "EN",
		TargetLabel: "FR",
		Pairs: []BlockPair{
			{
				Source: template.HTML("<p>Hello</p>"),
				Target: template.HTML("<p>Bonjour</p>"),
			},
		},
		// Fonts intentionally omitted (zero value)
	}

	html, err := Render(data)
	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	// Should use medium sizes: body=10, head=11, code=9, pre=8
	for _, expected := range []string{"font-size: 10pt", "font-size: 11pt", "font-size: 9pt", "font-size: 8pt"} {
		if !strings.Contains(html, expected) {
			t.Errorf("default render should contain %q (medium preset)", expected)
		}
	}
}

func TestFontSizePresets_Values(t *testing.T) {
	// Verify the preset map contains exactly the expected keys and values
	expected := map[string]FontSizes{
		"small":  {Body: 9, Head: 10, Code: 8, Pre: 7},
		"medium": {Body: 10, Head: 11, Code: 9, Pre: 8},
		"large":  {Body: 11, Head: 12, Code: 10, Pre: 9},
	}

	if len(FontSizePresets) != len(expected) {
		t.Fatalf("expected %d presets, got %d", len(expected), len(FontSizePresets))
	}

	for name, want := range expected {
		got, ok := FontSizePresets[name]
		if !ok {
			t.Errorf("missing preset %q", name)
			continue
		}
		if got != want {
			t.Errorf("preset %q = %+v, want %+v", name, got, want)
		}
	}
}
