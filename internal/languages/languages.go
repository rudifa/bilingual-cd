package languages

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// LangInfo holds display information for a language code.
type LangInfo struct {
	Code string
	Name string
}

// supported is the table of principal supported language codes.
var supported = map[string]string{
	"af": "Afrikaans",
	"ar": "Arabic",
	"bg": "Bulgarian",
	"bn": "Bengali",
	"ca": "Catalan",
	"cs": "Czech",
	"da": "Danish",
	"de": "German",
	"el": "Greek",
	"en": "English",
	"es": "Spanish",
	"et": "Estonian",
	"fa": "Persian",
	"fi": "Finnish",
	"fr": "French",
	"he": "Hebrew",
	"hi": "Hindi",
	"hr": "Croatian",
	"hu": "Hungarian",
	"id": "Indonesian",
	"it": "Italian",
	"ja": "Japanese",
	"ko": "Korean",
	"lt": "Lithuanian",
	"lv": "Latvian",
	"ms": "Malay",
	"nl": "Dutch",
	"no": "Norwegian",
	"pl": "Polish",
	"pt": "Portuguese",
	"ro": "Romanian",
	"ru": "Russian",
	"sk": "Slovak",
	"sl": "Slovenian",
	"sr": "Serbian",
	"sv": "Swedish",
	"th": "Thai",
	"tr": "Turkish",
	"uk": "Ukrainian",
	"vi": "Vietnamese",
	"zh": "Chinese",
}

// nativeNames maps language codes to their names in the native language.
var nativeNames = map[string]string{
	"af": "Afrikaans",
	"ar": "العربية",
	"bg": "Български",
	"bn": "বাংলা",
	"ca": "Català",
	"cs": "Čeština",
	"da": "Dansk",
	"de": "Deutsch",
	"el": "Ελληνικά",
	"en": "English",
	"es": "Español",
	"et": "Eesti",
	"fa": "فارسی",
	"fi": "Suomi",
	"fr": "Français",
	"he": "עברית",
	"hi": "हिन्दी",
	"hr": "Hrvatski",
	"hu": "Magyar",
	"id": "Bahasa Indonesia",
	"it": "Italiano",
	"ja": "日本語",
	"ko": "한국어",
	"lt": "Lietuvių",
	"lv": "Latviešu",
	"ms": "Bahasa Melayu",
	"nl": "Nederlands",
	"no": "Norsk",
	"pl": "Polski",
	"pt": "Português",
	"ro": "Română",
	"ru": "Русский",
	"sk": "Slovenčina",
	"sl": "Slovenščina",
	"sr": "Српски",
	"sv": "Svenska",
	"th": "ไทย",
	"tr": "Türkçe",
	"uk": "Українська",
	"vi": "Tiếng Việt",
	"zh": "中文",
}

// Validate checks if a language code is in the supported list.
func Validate(code string) error {
	if _, ok := supported[code]; !ok {
		return fmt.Errorf("unsupported language code: %q (use --list-languages to see supported codes)", code)
	}
	return nil
}

// Supported returns a sorted list of supported languages.
func Supported() []LangInfo {
	langs := make([]LangInfo, 0, len(supported))
	for code, name := range supported {
		langs = append(langs, LangInfo{Code: code, Name: name})
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].Code < langs[j].Code
	})
	return langs
}

// Name returns the display name for a language code, or the code itself if unknown.
func Name(code string) string {
	if name, ok := supported[code]; ok {
		return name
	}
	return code
}

// NativeName returns the language name in its own language, or falls back to the English name.
func NativeName(code string) string {
	if name, ok := nativeNames[code]; ok {
		return name
	}
	return Name(code)
}

// PrintSupported writes the supported languages table to the given writer.
func PrintSupported(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "Code\tLanguage")
	fmt.Fprintln(tw, "----\t--------")
	for _, lang := range Supported() {
		fmt.Fprintf(tw, "%s\t%s\n", lang.Code, lang.Name)
	}
	tw.Flush()
}
