# learnathon-2023

## Convert Proper Names TSV

This is just a quick verification that the TSV parses OK.

```
$ cd convertToCsv/
$ go run convertToCsv.go -i ../properNames.tsv -o properNames.csv
$ 
```

Looks good.

## Extract Places

This will extract the place data. Two output files are produced.

If `places.csv` is the given output, then it will also create 
`places_significance.csv`.

Note: there is a many to one relation between the places and the
"significance" (and the associated references of each significance)
data.

Run and import into database:
```
$ cd extractPlaces/
$ go run extractPlaces.go -i ../properNames.tsv -o places.csv
$ ls
extractPlaces.go  places.csv places_significance.csv
$ sqlite3 ../properNames.db
SQLite version 3.34.1 2021-01-20 14:10:07
Enter ".help" for usage hints.
sqlite> .mode csv
sqlite> .import places.csv places
sqlite> select count(*) from places;
1017
sqlite> .import places_significance.csv significance
sqlite> .schema
CREATE TABLE IF NOT EXISTS "places"(
  "UniqueName" TEXT,
  "OpenBible" TEXT,
  "Founder" TEXT,
  "People Group" TEXT,
  "GoogleMap URL" TEXT,
  "Palopenmaps URL" TEXT
);
CREATE TABLE IF NOT EXISTS "significance"(
  "UniqueName" TEXT,
  "Significance" TEXT,
  "Strongs" TEXT,
  "ESV Name" TEXT,
  "References" TEXT
);
sqlite> select count(*) from significance ;
2414
sqlite> 
```
