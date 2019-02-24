package md

func isFormatEnter(e bool, p string) string {
	if e {
		return "[::" + p + "]"
	}

	return "[::-]"
}
