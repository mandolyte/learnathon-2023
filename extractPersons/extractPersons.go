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
Here is the documentation on the "columns" for person data:
	UnifiedName=uStrong	Description	Parents Male+Female	Siblings	Partners	Offspring	>Tribe/Nation of father	#Summary description
	- Significance	UniqueName	dStrong«eStrong=Heb/Grk	ESV name (and KJV, NIV)	STEPBible link for first Refs                                 /	All Refs

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

	// open output file for main person table
	var w *csv.Writer
	if *output == "" {
		usage("Missing output filename")
	} else {
		fo, foerr := os.Create(*output)
		if foerr != nil {
			log.Fatal("os.Create() Error for main person file:" + foerr.Error())
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
	// UnifiedName=uStrong	Description	Parents Male+Female	Siblings	Partners	Offspring	>Tribe/Nation of father	#Summary description

	headers := []string{
		"UniqueName", // sans the Strong's number
		"UnifiedName",
		"Description",
		"Parents",
		"Siblings",
		"Partners",
		"Offspring",
		"TribeNation",
		"Summary",
	}

	// - Significance	UniqueName	dStrong«eStrong=Heb/Grk	ESV name (and KJV, NIV)	STEPBible link for first Refs                                 /	All Refs

	sheaders := []string{
		"UniqueName",
		"Qualifier",
		"Significance",
		"Strongs",
		"EsvName",
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

	const dataStart = 112
	const personMarker = "$========== PERSON(s)"
	for row := 0; row < len(records); row++ {

		if records[row][0] == personMarker && row > dataStart {
			// person data doesn't begin until row 112
			// but there is intro data near the beginning
			// that we need to overlook
		} else {
			continue
		}

		// main person data file
		var arow []string
		// fill up the row
		// "UniqueName", // sans the Strong's number
		// "UnifiedName",
		// "Description",
		// "Parents",
		// "Siblings",
		// "Partners",
		// "Offspring",
		// "Tribe/Nation",
		// "Summary",
		unifiedName := records[row+1][0]
		x := strings.Split(unifiedName, "=")

		arow = append(arow,
			x[0],
			unifiedName,
			records[row+1][1], // desc
			records[row+1][2], // parents
			records[row+1][3], // sibs
			records[row+1][4], // partners
			records[row+1][5], // offspring
			records[row+1][6], // tribe/nation
			records[row+1][7], // Summary
		)
		werr := writeRow(w, arow)
		if werr != nil {
			log.Fatalf("writeRow() error on row %v: \n%v\n", row, werr)
		}

		// significance file
		j := row + 2 // beginning after the main person data row
		for {
			if j == len(records) {
				break
			}
			cella := strings.TrimSpace(records[j][0])
			log.Printf("cella/%v/", cella)
			if cella == personMarker {
				break
			}
			// at the end of the person data are empty rows
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
