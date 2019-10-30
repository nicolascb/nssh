package config

import (
	"log"
	"regexp"
	"strings"
)

// matchStr check regexp match
func matchStr(rgxp string, compare string) bool {
	r, err := regexp.Compile(rgxp)
	if err != nil {
		log.Fatalf("invalid regexp: %s", rgxp)
	}
	return r.MatchString(strings.ToLower(compare))
}
