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

	// open output file for main place table
	var w *csv.Writer
	if *output == "" {
		usage("Missing output filename")
	} else {
		fo, foerr := os.Create(*output)
		if foerr != nil {
			log.Fatal("os.Create() Error for main place file:" + foerr.Error())
		}
		defer fo.Close()
		w = csv.NewWriter(fo)
	}
	// open output file for significance table
	var s *csv.Writer
	fo, foerr := os.Create(strings.TrimSuffix(*output, ".csv") + "_significance.csv")
	if foerr != nil {
		log.Fatal("os.Create() Error for significance file:" + foerr.Error())
	}
	defer fo.Close()
	s = csv.NewWriter(fo)

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

	sheaders := []string{
		"UniqueName",
		"Qualifier",
		"Significance",
		"Strongs",
		"ESV Name",
		"References",
	}

	// write the header row first
	//strings.ReplaceAll(str, " ", "")
	herr := writeRow(w, headers)
	if herr != nil {
		log.Fatalf("writeRow() error on header row: \n%v\n", herr)
	}
	serr := writeRow(s, sheaders)
	if serr != nil {
		log.Fatalf("writeRow() error on significance header row: \n%v\n", herr)
	}

	const dataStart = 10330
	const placeMarker = "$========== PLACE"
	for row := 0; row < len(records); row++ {

		if records[row][0] == placeMarker && row > dataStart {
			// place data doesn't begin until row 10326
			// but there is intro data near the beginning
			// that we need to overlook
		} else {
			continue
		}

		// main place data file
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

		// significance file
		j := row + 2 // beginning after the main place data row
		for {
			if j == len(records) {
				break
			}
			cella := strings.TrimSpace(records[j][0])
			log.Printf("cella/%v/", cella)
			if cella == placeMarker {
				break
			}
			// at the end of the place data are empty rows
			if cella == "" {
				j++
				continue
			}
			// if it is a note
			if strings.HasPrefix(cella, "NOTES") {
				j++
				continue
			}
			// if the string is really long it is the note content
			if len(cella) > 21 {
				j++
				continue
			}
			if strings.HasPrefix(cella, "#==") {
				break
			}
			if strings.HasPrefix(cella, "$==") {
				break
			}
			// data misplace?
			if strings.HasPrefix(cella, "Rahab@Jos.2.1=H7343") {
				break
			}
			// columns for significance data:
			// 	- Significance,UniqueName,dStrong«eStrong=Heb/Grk,
			// ESV name (and KJV, NIV),STEPBible link for first,Refs
			// Will skip the step bible link
			uname := ""
			qualifier := ""
			significance := ""
			strongs := ""
			esvName := ""
			refs := ""

			for c, v := range records[j] {
				if c == 0 {
					significance = v
				}
				if c == 1 {
					// the unique name may have a qualifier.
					// if so, it will follow this pattern: "west|Arabia@2Sa.23.35"
					x := strings.Split(v, "|")
					if len(x) > 1 {
						qualifier = x[0]
						uname = x[1]
					} else {
						uname = v
					}
				}
				if c == 2 {
					strongs = v
				}
				if c == 3 {
					esvName = v
				}
				if c == 5 {
					refs = v
				}
			}
			var srow []string
			srow = append(srow,
				uname,        // unique name
				qualifier,    // qualifer
				significance, // significance
				strongs,      // strongs
				esvName,      // esv name
				refs,         // refs
			)
			serr := writeRow(s, srow)
			if serr != nil {
				log.Fatalf("writeRow() error on row %v: \n%v\n", srow, serr)
			}
			j++
		}
	}
	w.Flush()
	s.Flush()
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
	os.Exit(1)
}
