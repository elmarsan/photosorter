# Photosorter

[![Go Report Card](https://goreportcard.com/badge/github.com/elmarsan/photosorter)](https://goreportcard.com/report/github.com/elmarsan/photosorter) 
![Coverage](https://img.shields.io/badge/Coverage-90.8%25-brightgreen)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Reference](https://pkg.go.dev/badge/github.com/elmarsan/photosorter.svg)](https://pkg.go.dev/github.com/elmarsan/photosorter)

Photosorter is a CLI tool designed for sorting your photos by their original creatuion date.
Powered by [Cobra](https://github.com/spf13/cobra) and [Go-exif](https://github.com/dsoprea/go-exif).

## Usage

![Demo gif](./doc/photosorter.gif)

The `sort` command requires two arguments: the source directory `src` containing the target photos, and the destination directory `dst` where the sorted photos will be output.
Additionally, you can use the `--format` flag to choose the output structure directory. It accepts `month` and `year` as options.