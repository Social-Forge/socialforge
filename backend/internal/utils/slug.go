package utils

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/mozillazg/go-unidecode"
)

func GenerateSlug(title string) string {
	slug := strings.ToLower(title)
	slug = strings.Map(func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == ' ' || r == '-' {
			return r
		}
		return -1
	}, slug)

	return strings.ReplaceAll(strings.TrimSpace(slug), " ", "-")
}
func GenerateArticleSlug(s string) string {
	s = strings.ToLower(s)

	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	s = reg.ReplaceAllString(s, "")

	s = strings.ReplaceAll(s, " ", "-")

	reg = regexp.MustCompile(`-+`)
	s = reg.ReplaceAllString(s, "-")

	s = strings.Trim(s, "-")

	if s == "" {
		str := GenerateRandomToken()
		return fmt.Sprintf("untitled-%s", str)
	}

	return s
}
func GenerateSlugUnicode(s string) string {
	s = strings.ToLower(s)

	f := func(r rune) rune {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			return r
		}
		return ' '
	}
	s = strings.Map(f, s)

	reSpace := regexp.MustCompile(`[\p{Z}\p{P}\p{S}]+`)
	s = reSpace.ReplaceAllString(s, " ")

	s = strings.ReplaceAll(s, " ", "-")

	s = strings.Trim(s, "- ")

	if s == "" {
		str := generateRandomSlug()
		return fmt.Sprintf("untitled-%s", str)
	}

	return s
}
func GenerateSlugUnicodeV2(s string) string {
	s = strings.TrimSpace(s)
	s = strings.ToLower(s)

	s = unidecode.Unidecode(s)

	var result strings.Builder
	var lastChar rune = 0
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z': // Huruf kecil
			result.WriteRune(r)
			lastChar = r
		case r >= '0' && r <= '9': // Angka
			result.WriteRune(r)
			lastChar = r
		case unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r):
			// Ganti dengan dash jika karakter sebelumnya bukan dash
			if lastChar != '-' && result.Len() > 0 {
				result.WriteRune('-')
				lastChar = '-'
			}
		}
	}

	slug := result.String()

	slug = regexp.MustCompile(`-+`).ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	if slug == "" || !isValidSlug(slug) {
		return generateRandomSlug()
	}

	return slug
}
func isValidSlug(s string) bool {
	if len(s) < 3 || len(s) > 200 {
		return false
	}

	if !regexp.MustCompile(`[a-z]`).MatchString(s) {
		return false
	}

	return regexp.MustCompile(`^[a-z0-9-]+$`).MatchString(s)
}
func generateRandomSlug() string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 8)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return "untitled-" + string(b)
}
func Transliterate(s string) string {
	transliterations := map[rune]string{
		// Latin Extended
		'à': "a", 'á': "a", 'â': "a", 'ã': "a", 'ä': "a", 'å': "a",
		'è': "e", 'é': "e", 'ê': "e", 'ë': "e",
		'ì': "i", 'í': "i", 'î': "i", 'ï': "i",
		'ò': "o", 'ó': "o", 'ô': "o", 'õ': "o", 'ö': "o", 'ø': "o",
		'ù': "u", 'ú': "u", 'û': "u", 'ü': "u",
		'ñ': "n", 'ç': "c", 'ß': "ss",

		// Cyrillic (Rusia)
		'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e",
		'ё': "yo", 'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k",
		'л': "l", 'м': "m", 'н': "n", 'о': "o", 'п': "p", 'р': "r",
		'с': "s", 'т': "t", 'у': "u", 'ф': "f", 'х': "kh", 'ц': "ts",
		'ч': "ch", 'ш': "sh", 'щ': "shch", 'ъ': "", 'ы': "y", 'ь': "",
		'э': "e", 'ю': "yu", 'я': "ya",

		// Arab (transliterasi dasar)
		'ا': "a", 'ب': "b", 'ت': "t", 'ث': "th", 'ج': "j",
		'ح': "h", 'خ': "kh", 'د': "d", 'ذ': "dh", 'ر': "r",
		'ز': "z", 'س': "s", 'ش': "sh", 'ص': "s", 'ض': "d",
		'ط': "t", 'ظ': "z", 'ع': "a", 'غ': "gh", 'ف': "f",
		'ق': "q", 'ك': "k", 'ل': "l", 'م': "m", 'ن': "n",
		'ه': "h", 'و': "w", 'ي': "y",
	}

	var result strings.Builder
	for _, r := range s {
		if replacement, ok := transliterations[r]; ok {
			result.WriteString(replacement)
		} else if r > 127 {
			// Skip karakter non-ASCII yang tidak ada di mapping
			continue
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}
