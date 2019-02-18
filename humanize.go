package main

// HumanizeStrings converts arrays to a friendly string
func HumanizeStrings(a []string) (s string) {
	switch len(a) {
	case 0:
	case 1:
		s = a[0]
	default:
		for i := 0; i < len(a)-2; i++ {
			s += a[i] + ", "
		}

		s += " and " + a[len(a)-1]
	}

	return
}
