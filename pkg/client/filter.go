package client

import (
	"fmt"
	"strings"
)

func buildEventFilter(events []string) string {
	xs := []string{}
	for _, event := range events {
		xs = append(xs, fmt.Sprintf(`r["_measurement"] == "%s"`, event))
	}
	if len(xs) == 0 {
		return ""
	}
	return fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(xs, " or "))
}

func buildFieldFilter(filters []string) string {
	xs := []string{}
	for _, f := range filters {
		fs := splitFilter(f)
		if len(fs) != 3 {
			continue
		}
		if strings.Contains(fs[2], "~") {
			xs = append(xs, fmt.Sprintf(`r["%s"] %s %s`, fs[0], fs[2], fs[1]))
		} else {
			xs = append(xs, fmt.Sprintf(`r["%s"] %s "%s"`, fs[0], fs[2], fs[1]))
		}
	}
	if len(xs) == 0 {
		return ""
	}
	return fmt.Sprintf(`|> filter(fn: (r) => %s)`, strings.Join(xs, " and "))
}

func splitFilter(f string) []string {
	operators := []string{"==", "!=", "<=", ">=", "=~", "!~", "<", ">"}
	for _, op := range operators {
		if strings.Contains(f, op) {
			xs := strings.Split(f, op)
			if !strings.HasPrefix(xs[0], "f_") {
				xs[0] = "f_" + xs[0]
			}
			return append(xs, op)
		}
	}
	return nil
}
