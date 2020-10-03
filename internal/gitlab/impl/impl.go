// Package impl implements services for gitlab
package impl

func hideToken(s string) string {
	if s == "" {
		return s
	}

	if len(s) < 3 {
		return s
	}

	return s[:3] + "..."
}
