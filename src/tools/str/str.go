package str

import (
	"fmt"
	"regexp"
	"strings"
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
	first := fmt.Sprintf("%c", str[0])
	rest := str[1:]

	return fmt.Sprintf("%s%s", strings.ToUpper(first), strings.ToLower(rest))
}
