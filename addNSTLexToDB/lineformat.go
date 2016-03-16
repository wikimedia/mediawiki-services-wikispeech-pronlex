package main

import "github.com/stts-se/pronlex/line"

// TODO: HL work in progress

var lineFmt line.Format

func initLineFmt() line.Format {
	var nilFmt line.Format
	if line.Equals(lineFmt, nilFmt) {
		tests := []line.FormatTest{
			// TODO: update to actual format
			line.FormatTest{"hannas;PM;GEN;hannas;;;eng;;;;;\"\" h a . n a s;;;swe;hanna_01",
				map[line.Field]string{
					line.Orth:       "hannas",
					line.Pos:        "PM",
					line.Morph:      "GEN",
					line.Decomp:     "hannas",
					line.WordLang:   "eng",
					line.Trans1:     "\"\" h a . n a s",
					line.Translang1: "swe",
					line.Lemma:      "hanna_01"},
			},
		}
		f, err := line.NewFormat(
			"NST",
			";",
			map[line.Field]int{
				line.Orth:           0,
				line.Pos:            1,
				line.Morph:          2,
				line.Decomp:         3,
				line.WordLang:       6,
				line.Trans1:         11,
				line.Translang1:     14,
				line.Trans2:         15,
				line.Translang2:     18,
				line.Trans3:         19,
				line.Translang3:     22,
				line.Trans4:         23,
				line.Translang4:     26,
				line.Lemma:          32,
				line.InflectionRule: 33,
			},
			54,
			tests,
		)
		if err != nil {
			panic(err) // TODO
		}
		lineFmt = f
	}
	return lineFmt
}
