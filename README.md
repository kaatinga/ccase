# ccase

ccase is a command line tool which helps to convert ".go" file names to a name that follows the conventions of the Go
programming language.

## Installation

```bash
go install github.com/kaatinga/ccase@latest
```

## Usage

To get help run:

```bash
ccase -h
```

To rename files in the pkg/service folder run:

```bash
ccase --path pkg/service
```

or just:

```bash
ccase
```

in the folder you want to rename files.

Note that all the subfolders will be processed recursively.