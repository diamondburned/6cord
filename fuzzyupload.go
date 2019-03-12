package main

import (
	"os"
	"strings"

	"github.com/sahilm/fuzzy"
)

func joinFolder(fs []string) (p string) {
	for _, f := range fs {
		p += f + "/"
	}

	return
}

func fuzzyUpload(text string) {
	fuzzied := [][]string{}

	if len(text) > 8 {
		text = strings.TrimPrefix(text, "/upload ")

		var (
			pathparts = strings.Split(text, "/")
			inputfile = pathparts[len(pathparts)-1]

			path = joinFolder(pathparts[:len(pathparts)-1])
		)

		f, err := os.Open(path)
		if err != nil {
			return
		}

		defer f.Close()

		files, err := f.Readdir(-1)
		if err != nil {
			return
		}

		filenames := make([]string, len(files))

		for i, f := range files {
			filename := f.Name()
			if f.IsDir() {
				filename += "/"
			}

			filenames[i] = filename
		}

		results := fuzzy.Find(inputfile, filenames)
		fuzzied = make([][]string, len(results))

		for i, r := range results {
			fuzzied[i] = []string{
				"[::u]" + path + "[::-]\u200B" + formatNeedle(r),
				path + filenames[r.Index],
			}
		}
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, filename := range fuzzied {
			autocomp.InsertItem(
				i,
				filename[0],
				"", 0, nil,
			)
		}

		rightflex.ResizeItem(autocomp, min(len(fuzzied), 10), 1)

		autofillfunc = func(i int) {
			input.SetText("/upload " + fuzzied[i][1])
			clearList()
			app.SetFocus(input)
		}

	} else {
		rightflex.ResizeItem(autocomp, 1, 1)
	}

	app.Draw()
}
