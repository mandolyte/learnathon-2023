package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

/*
Here is the documentation on the "columns" for place data:
	UniqueName=uStrong	OpenBible name=Near	Founder	People living there	GoogleMap URL	Palopenmaps URL	>Geographical area
	- Significance	UniqueName	dStrong«eStrong=Heb/Grk	ESV name (and KJV, NIV)	STEPBible link for first Refs                 /	All Refs

There is one line for first row, which is "key" data
Then one or more rows for the second row, which refers to
the various kinds of "sigificances".

Most of these appear to start with a dash and the most common
is "- Named"

*/

func main() {
	input := flag.String("i", "", "Input CSV filename")
	output := flag.String("o", "", "Output CSV filename")
	flag.Parse()

	// open output file
	var w *csv.Writer
	if *output == "" {
		usage("Missing output filename")
	} else {
		fo, foerr := os.Create(*output)
		if foerr != nil {
			log.Fatal("os.Create() Error:" + foerr.Error())
		}
		defer fo.Close()
		w = csv.NewWriter(fo)
	}

	// open input file
	var r *csv.Reader
	if *input == "" {
		usage("Missing input filename")
	} else {
		fi, fierr := os.Open(*input)
		if fierr != nil {
			log.Fatal("os.Open() Error:" + fierr.Error())
		}
		defer fi.Close()
		r = csv.NewReader(fi)
	}

	// ignore expectations of fields per row
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	r.Comma = '\t'
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	// read loop for CSV
	//	UniqueName=uStrong	OpenBible name=Near	Founder	People living there	GoogleMap URL	Palopenmaps URL	>Geographical area

	headers := []string{
		"UniqueName",
		"OpenBible",
		"Founder",
		"People Group",
		"GoogleMap URL",
		"Palopenmaps URL",
	}

	// write the header row first
	//strings.ReplaceAll(str, " ", "")
	herr := writeRow(w, headers)
	if herr != nil {
		log.Fatalf("writeRow() error on header row: \n%v\n", herr)
	}
	const dataStart = 10330
	for row := 0; row < len(records); row++ {

		if records[row][0] == "$========== PLACE" && row > dataStart {
			// place data doesn't begin until row 10326
			// but there is intro data near the beginning
			// that we need to overlook
		} else {
			continue
		}

		var arow []string
		// fill up the row
		arow = append(arow,
			records[row+1][0],
			records[row+1][1],
			records[row+1][2],
			records[row+1][3],
			strings.ReplaceAll(records[row+1][4], " ", ""), // data has spaces where they should not be
			records[row+1][5])
		werr := writeRow(w, arow)
		if werr != nil {
			log.Fatalf("writeRow() error on row %v: \n%v\n", row, werr)
		}
	}
	w.Flush()
}

func writeRow(w *csv.Writer, cells []string) error {
	err := w.Write(cells)
	if err != nil {
		return err
	}
	return nil
}

func usage(msg string) {
	fmt.Println(msg + "\n")
	fmt.Print("Usage: parseProperNames -i input.csv -o output.csv\n")
	flag.PrintDefaults()
}
