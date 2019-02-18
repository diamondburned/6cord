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

		for i, r := range results {
			if i == 10 {
				break
			}

			fuzzied = append(
				fuzzied,
				[]string{
					"[::u]" + path + "[::-]" + formatNeedle(r, filenames[r.Index]),
					path + filenames[r.Index],
				},
			)
		}
	}

	clearList()

	if len(fuzzied) > 0 {
		for i, filename := range fuzzied {
			autocomp.InsertItem(
				i,
				filename[0],
				"",
				rune(0x31+i),
				nil,
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
