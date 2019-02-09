package main

// HumanizeStrings converts arrays to a friendly string
func HumanizeStrings(a []string) (s string) {
	switch {
	case len(a) == 0:
		return ""
	case len(a) == 1:
		return a[0]
	default:
		for i := 0; i < len(a)-1; i++ {
			s += a[i] + ", "
		}

		s += " and " + a[len(a)-1]

		return
	}
}
