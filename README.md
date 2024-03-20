[![Tests](https://github.com/kaatinga/ccase/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/kaatinga/ccase/actions/workflows/test.yml)
[![MIT license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kaatinga/ccase/blob/main/LICENSE)
[![lint workflow](https://github.com/kaatinga/ccase/actions/workflows/golangci-lint.yml/badge.svg)](https://github.com/kaatinga/ccase/actions?query=workflow%3Alinter)
[![Help wanted](https://img.shields.io/badge/Help%20wanted-True-yellow.svg)](https://github.com/kaatinga/ccase/issues?q=is%3Aopen+is%3Aissue+label%3A%22help+wanted%22)

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