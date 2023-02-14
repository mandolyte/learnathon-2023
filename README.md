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

This will extract the place data. Eventually will need to be two files.

```
$ cd extractPlaces/
$ go run extractPlaces.go -i ../properNames.tsv -o places.csv
```


