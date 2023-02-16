#!/bin/sh

# run the command to extract the data in normalized form
echo "Run the exract"
go run extractPersons.go -i ../properNames.tsv -o persons.csv

echo "Run the database imports"
# NOTE!! cannot indent the sql commands
sqlite3 ../properNames.db <<EoF
.echo on
drop table if exists persons;
drop table if exists persons_significance;

.mode csv
.import persons.csv persons
select count(*) from persons;
.import persons_significance.csv persons_significance
select count(*) from persons_significance;
.exit
EoF



