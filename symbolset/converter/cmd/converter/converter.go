package main

import (
	"fmt"
	"log"
	"os"

	"github.com/stts-se/pronlex/symbolset"
	"github.com/stts-se/pronlex/symbolset/converter"
)

func main() {
	var usage = `converter <SYMBOL SET FOLDER> <CONVERTER FILE or FOLDER>`

	if len(os.Args) != 3 {
		fmt.Println(usage)
		os.Exit(1)
	}

	sSets, err := symbolset.LoadSymbolSetsFromDir(os.Args[1])
	if err != nil {
		log.Printf("Couldn't load symbol set files from dir %s: %s", os.Args[1], err)
		os.Exit(1)
	}

	convFile := os.Args[2]

	fi, err := os.Stat(convFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	switch mode := fi.Mode(); {
	case mode.IsDir():
		_, testRes, err := converter.LoadFromDir(sSets, convFile)
		if err != nil {
			log.Printf("Couldn't load converter files from dir %s: %s", os.Args[1], err)
			os.Exit(1)
		}
		for name, res := range testRes {
			if !res.OK {
				for _, err := range res.Errors {
					log.Printf("%s: %s", name, err)
				}
			}
		}
		// do directory stuff
		fmt.Println("directory")
	case mode.IsRegular():
		conv, res, err := converter.LoadFile(sSets, convFile)
		if err != nil {
			log.Printf("Couldn't load converter file %s: %s", os.Args[1], err)
			os.Exit(1)
		}
		if !res.OK {
			for _, err := range res.Errors {
				log.Printf("%s: %s", conv.Name, err)
			}
		}
	}

}
