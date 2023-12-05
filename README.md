# `dukCSV`

`dukCSV` is a library that attempts to read CSV files efficiently using``
random-access patterns. This is a project born from a slight itch caused by
trying to read 2GB+ CSV files from a database at work.

There are no benchmarks as of now, as this library was never meant to be
super performant. Just needed to do its job.

## Prior Art

There are other projects out there that mmap CSV files (and this is heavily
inspired by [one of them](https://github.com/carbocation/genomisc/tree/master/ramcsv)),
but they were either lacking a key feature, or weren't exactly what I was
looking for.

## Key Features

### Multi-line Records

The main feature that this CSV reader has over others is the ability to read
multi-line records. Most CSV readers treat every `\n` character as a new row
in the CSV file, but sometimes that's annoyingly not true. `dukCSV` reads in
`"` characters, keeping track of what is between quotes and what is not and
allows `\n` characters within records/cells (see
[testdata.csv](./testdata/test.csv)).

## TODO List

Once there's actual to-do items, I will add them here.