package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/stts-se/pronlex/lex"
	"github.com/stts-se/pronlex/line"
	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/vrules"
)

var sucTags = map[string]bool{
	"AB": true,
	"DT": true,
	"HA": true,
	"HD": true,
	"HP": true,
	"HS": true,
	"IE": true,
	"IN": true,
	"JJ": true,
	"KN": true,
	"NN": true,
	"PC": true,
	"PL": true,
	"PM": true,
	"PN": true,
	"PP": true,
	"PS": true,
	"RG": true,
	"RO": true,
	"SN": true,
	"UO": true,
	"VB": true,
}

var langCodes = map[string]string{
	"aar": "aa",
	"abk": "ab",
	"afr": "af",
	"aka": "ak",
	"alb": "sq",
	"amh": "am",
	"ara": "ar",
	"arg": "an",
	"arm": "hy",
	"asm": "as",
	"ava": "av",
	"ave": "ae",
	"aym": "ay",
	"aze": "az",
	"bak": "ba",
	"bam": "bm",
	"baq": "eu",
	"bel": "be",
	"ben": "bn",
	"bih": "bh",
	"bis": "bi",
	"bod": "bo",
	"bos": "bs",
	"bre": "br",
	"bul": "bg",
	"bur": "my",
	"cat": "ca",
	"ces": "cs",
	"cha": "ch",
	"che": "ce",
	"chi": "zh",
	"chu": "cu",
	"chv": "cv",
	"cor": "kw",
	"cos": "co",
	"cre": "cr",
	"cym": "cy",
	"cze": "cs",
	"dan": "da",
	"deu": "de",
	"div": "dv",
	"dut": "nl",
	"dzo": "dz",
	"ell": "el",
	"eng": "en",
	"epo": "eo",
	"est": "et",
	"eus": "eu",
	"ewe": "ee",
	"fao": "fo",
	"fas": "fa",
	"fij": "fj",
	"fin": "fi",
	"fra": "fr",
	"fre": "fr",
	"fry": "fy",
	"ful": "ff",
	"geo": "ka",
	"ger": "de",
	"gla": "gd",
	"gle": "ga",
	"glg": "gl",
	"glv": "gv",
	"gre": "el",
	"grn": "gn",
	"guj": "gu",
	"hat": "ht",
	"hau": "ha",
	"heb": "he",
	"her": "hz",
	"hin": "hi",
	"hmo": "ho",
	"hrv": "hr",
	"hun": "hu",
	"hye": "hy",
	"ibo": "ig",
	"ice": "is",
	"ido": "io",
	"iii": "ii",
	"iku": "iu",
	"ile": "ie",
	"ina": "ia",
	"ind": "id",
	"ipk": "ik",
	"isl": "is",
	"ita": "it",
	"jav": "jv",
	"jpn": "ja",
	"kal": "kl",
	"kan": "kn",
	"kas": "ks",
	"kat": "ka",
	"kau": "kr",
	"kaz": "kk",
	"khm": "km",
	"kik": "ki",
	"kin": "rw",
	"kir": "ky",
	"kom": "kv",
	"kon": "kg",
	"kor": "ko",
	"kua": "kj",
	"kur": "ku",
	"lao": "lo",
	"lat": "la",
	"lav": "lv",
	"lim": "li",
	"lin": "ln",
	"lit": "lt",
	"ltz": "lb",
	"lub": "lu",
	"lug": "lg",
	"mac": "mk",
	"mah": "mh",
	"mal": "ml",
	"mao": "mi",
	"mar": "mr",
	"may": "ms",
	"mkd": "mk",
	"mlg": "mg",
	"mlt": "mt",
	"mon": "mn",
	"mri": "mi",
	"msa": "ms",
	"mya": "my",
	"nau": "na",
	"nav": "nv",
	"nbl": "nr",
	"nde": "nd",
	"ndo": "ng",
	"nep": "ne",
	"nld": "nl",
	"nno": "nn",
	"nob": "nb",
	"nor": "no",
	"nya": "ny",
	"oci": "oc",
	"oji": "oj",
	"ori": "or",
	"orm": "om",
	"oss": "os",
	"pan": "pa",
	"per": "fa",
	"pli": "pi",
	"pol": "pl",
	"por": "pt",
	"pus": "ps",
	"que": "qu",
	"roh": "rm",
	"ron": "ro",
	"rum": "ro",
	"run": "rn",
	"rus": "ru",
	"sag": "sg",
	"san": "sa",
	"sin": "si",
	"slk": "sk",
	"slo": "sk",
	"slv": "sl",
	"sme": "se",
	"smo": "sm",
	"sna": "sn",
	"snd": "sd",
	"som": "so",
	"sot": "st",
	"spa": "es",
	"sqi": "sq",
	"srd": "sc",
	"srp": "sr",
	"ssw": "ss",
	"sun": "su",
	"swa": "sw",
	"swe": "sv",
	"tah": "ty",
	"tam": "ta",
	"tat": "tt",
	"tel": "te",
	"tgk": "tg",
	"tgl": "tl",
	"tha": "th",
	"tib": "bo",
	"tir": "ti",
	"ton": "to",
	"tsn": "tn",
	"tso": "ts",
	"tuk": "tk",
	"tur": "tr",
	"twi": "tw",
	"uig": "ug",
	"ukr": "uk",
	"urd": "ur",
	"uzb": "uz",
	"ven": "ve",
	"vie": "vi",
	"vol": "vo",
	"wel": "cy",
	"wln": "wa",
	"wol": "wo",
	"xho": "xh",
	"yid": "yi",
	"yor": "yo",
	"zha": "za",
	"zho": "zh",
	"zul": "zu",
	"for": "foreign",
}

func validPos(pos string) bool {
	if pos == "" {
		return true
	}
	_, ok := sucTags[pos]
	if ok {
		return true
	}
	return false
}

func mapLanguage(lang string) (string, error) {
	if lang == "" {
		return lang, nil
	}
	l, ok := langCodes[strings.ToLower(lang)]
	if ok {
		return l, nil
	}
	return lang, fmt.Errorf("couldn't map language <%v>", lang)
}

func main() {

	if len(os.Args) != 4 {
		fmt.Fprintln(os.Stderr, "<INPUT NST LEX FILE> <LEX2IPA MAPPER> <IPA2SAMPA MAPPER>")
		fmt.Fprintln(os.Stderr, "\tsample invokation:  go run convertNST2WS.go swe030224NST.pron.utf8 sv-se_nst-xsampa.tab sv-se_ws-sampa.tab ")
		return
	}

	nstFileName := os.Args[1]
	ssFileName1 := os.Args[2]
	ssFileName2 := os.Args[3]

	mapper, err := symbolset.LoadMapperFromFile("SAMPA", "SYMBOL", ssFileName1, ssFileName2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't load mappers: %v\n", err)
		return
	}
	ssRuleTo := vrules.SymbolSetRule{mapper.SymbolSet2.To}

	nstFile, err := os.Open(nstFileName)
	defer nstFile.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't open lexicon file: %v\n", err)
		return
	}

	nstFmt, err := line.NewNST()
	if err != nil {
		log.Fatal(err)
	}
	wsFmt, err := line.NewWS()
	if err != nil {
		log.Fatal(err)
	}

	nst := bufio.NewScanner(nstFile)
	n := 0
	for nst.Scan() {
		n++
		if err := nst.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "failed reading line %v : %v\n", n, err)
		}
		line := nst.Text()

		e, err := nstFmt.ParseToEntry(line)
		fmt.Println(e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to convert line %v to entry : %v\n", n, err)
			fmt.Fprintf(os.Stderr, "failing line: %v\n", line)
		}

		e.EntryStatus.Name = "imported"
		e.EntryStatus.Source = "nst"
		e.Language, err = mapLanguage(e.Language)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to map language <%v>\n", err)
		}
		if !validPos(e.PartOfSpeech) {
			fmt.Fprintf(os.Stderr, "invalid pos tag <%v>\n", e.PartOfSpeech)
		}

		err = mapper.MapTranscriptions(&e)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to map transcription symbols : %v\n", err)
		} else {
			for _, r := range ssRuleTo.Validate(e) {
				panic(r) // shouldn't happen
			}

			res, err := wsFmt.Entry2String(e)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to convert entry to string : %v\n", err)
			} else {
				fmt.Printf("%v\n", res)
			}
		}
	}

	_ = lex.Entry{}
}
