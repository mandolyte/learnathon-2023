#!/bin/sh

# run the command to extract the data in normalized form
echo "Run the exract"
go run extractPlaces.go -i ../properNames.tsv -o places.csv

echo "Run the database imports"
# NOTE!! cannot indent the sql commands
sqlite3 ../properNames.db <<EoF
.echo on
drop table if exists places;
drop table if exists places_significance;

.mode csv
.import places.csv places
select count(*) from places;
.import places_significance.csv places_significance
select count(*) from places_significance;
.exit
EoF



