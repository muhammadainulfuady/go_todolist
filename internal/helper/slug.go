package helper

import (
	"context"
	"strconv"
	"strings"
	"unicode"
)

// GenerateSlug mengubah judul menjadi slug URL: huruf kecil, tanpa aksara
// khusus, spasi/tanda hubung digabung menjadi satu "-".
func GenerateSlug(title string) string {
	title = strings.ToLower(strings.TrimSpace(title))

	var b strings.Builder
	lastDash := false
	for _, r := range title {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			lastDash = false
		case r == ' ' || r == '-':
			if !lastDash {
				b.WriteByte('-')
				lastDash = true
			}
		}
	}

	return strings.Trim(b.String(), "-")
}

type SlugChecker interface {
	SlugExists(ctx context.Context, slug string) (bool, error)
}

// EnsureUniqueSlug memastikan slug belum dipakai; jika sudah, diberi suffix
// -1, -2, dst (misal "mengerjakan-pr-1").
func EnsureUniqueSlug(ctx context.Context, checker SlugChecker, slug string) (string, error) {
	exists, err := checker.SlugExists(ctx, slug)
	if err != nil {
		return "", err
	}
	if !exists {
		return slug, nil
	}

	for i := 1; ; i++ {
		candidate := slug + "-" + strconv.Itoa(i)
		exists, err := checker.SlugExists(ctx, candidate)
		if err != nil {
			return "", err
		}
		if !exists {
			return candidate, nil
		}
	}
}
