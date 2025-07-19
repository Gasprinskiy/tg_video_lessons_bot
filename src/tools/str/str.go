package str

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	FullNameRegexp   = regexp.MustCompile(`^(([А-ЯЁA-Z][а-яёa-z]+)\s+([А-ЯЁA-Z][а-яёa-z]+))$`)
	BirthDateRegexp  = regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	EmptySpaceRegexp = regexp.MustCompile(`\s+`)
)

func SplitStringByEmptySpace(str string) []string {
	return EmptySpaceRegexp.Split(str, -1)
}

func CapFirstLowerRest(str string) string {
	runes := []rune(str)
	first := unicode.ToUpper(runes[0])
	rest := strings.ToLower(string(runes[1:]))

	return fmt.Sprintf("%c%s", first, rest)
}
