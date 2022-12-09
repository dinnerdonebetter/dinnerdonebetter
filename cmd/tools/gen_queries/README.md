# Query Generator

This is a WIP tool to generate queries for the database library. The idea is that since so many of the queries share
the precise ordering of columns (i.e. the columns in a select query for table A are the same for queries that simply join A),
we can just keep one bona fide list of the columns and generate the queries from that.
