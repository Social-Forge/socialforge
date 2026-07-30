package utils

import (
	"html"
	"strings"

	"github.com/microcosm-cc/bluemonday"
)

func SanitizeHTML(raw string) string {
	p := bluemonday.UGCPolicy()
	return p.Sanitize(raw)
}
func SanitizeNullableHTML(html *string) *string {
	if html == nil {
		return nil
	}
	sanitized := SanitizeHTML(*html)
	return &sanitized
}
func SanitizeCustomHTML(raw string) string {
	p := bluemonday.NewPolicy()
	// Izinkan tag dasar yang lebih luas (misalnya)
	p.AllowElements(
		"p", "br", "b", "i", "u", "strong", "em", "a", // Teks dasar
		"ul", "ol", "li", // List
		"h1", "h2", "h3", "h4", "h5", "h6", // Heading
		"blockquote",                                // Kutipan
		"img",                                       // Gambar
		"table", "thead", "tbody", "tr", "th", "td", // Tabel
		// Tambahkan tag lain jika kamu ingin mempertahankan (misalnya: figure, figcaption, div tertentu)
	)
	// Izinkan atribut tertentu untuk tag tertentu
	p.AllowAttrs("href").OnElements("a")
	p.AllowAttrs("src", "alt", "title").OnElements("img")

	// --- HATI-HATI DENGAN INI ---
	// Jika kamu perlu class/id untuk goquery parsing, tapi ini bisa jadi risiko XSS
	// Jika data ini akan ditampilkan di frontend, pastikan ada sanitasi lagi di frontend atau backend.
	// Jika hanya untuk parsing internal, risiko lebih rendah.
	// p.AllowAttrs("class", "id").Globally()
	// Pertimbangkan untuk hanya mengizinkan 'class'/'id' pada elemen spesifik jika memungkinkan.
	// Misalnya: p.AllowAttrs("class", "id").OnElements("div", "p", "h1", "h2", "h3", "table")
	// Atau, lebih baik lagi, setelah sanitasi, goquery masih bisa bekerja dengan struktur DOM yang bersih.
	// Bluemonday akan menghapus tag/attrs yang tidak diizinkan.

	// Allow rel="nofollow" for links (common in crawled content)
	p.AllowRelativeURLs(true) // Untuk URL relatif
	p.AllowDataAttributes()   // Jika ada data-attributes yang relevan

	return p.Sanitize(raw)
}
func EscapeNullableString(s *string) *string {
	if s == nil {
		return nil
	}
	escaped := html.EscapeString(*s)
	return &escaped
}
func SanitizeContentSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	sanitized := make([]string, len(slice))
	for i, s := range slice {
		sanitized[i] = SanitizeHTML(s)
	}
	return sanitized
}
func SanitizeStringSlice(slice []string) []string {
	if slice == nil {
		return nil
	}
	sanitized := make([]string, len(slice))
	for i, s := range slice {
		sanitized[i] = html.EscapeString(strings.TrimSpace(s))
	}
	return sanitized
}
func Sanitize2DStringSlice(slices [][]string) [][]string {
	if slices == nil {
		return nil
	}
	sanitized := make([][]string, len(slices))
	for i, slice := range slices {
		sanitized[i] = SanitizeStringSlice(slice)
	}
	return sanitized
}
