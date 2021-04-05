package toolkit

import (
	"strings"
	"fmt"
	"regexp"
)

func MatchAndReplace(expr, reg string, target string) (string, error) {
	if t, err := regexp.Compile(reg); err == nil {
		r := target
		s := t.FindStringSubmatch(expr)
		if s != nil && len(s) > 1 {
			fmt.Println("find match", s)
			for i := 1; i < len(s); i++ {
				// replace $n in the target 
				old := fmt.Sprintf("$%d", i)
				r = strings.Replace(r, old, s[i], -1) 
			}
		}
		return r, nil
	} else { 
		return "", err
	}	
}
