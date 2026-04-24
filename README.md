# File Processor with History

A command-line tool that processes files using interchangeable plugins and tracks already processed files using a content hash to avoid redundant work.

## Features

- **Multiple processing modes** – `lines`, `words`, `hash`, `compress`
- **Content‑based deduplication** – computes SHA‑256 hash of the file and skips files that have been processed before (in‑memory history)
- **Plugin architecture** – new processors can be added by implementing the `Processor` interface
- **Clear report output** – shows file name and processing result

## Installation

Ensure you have Go 1.16 or later installed. Clone the repository and build:

```bash
go build -o fileproc .
```
## Usage

```bash
./fileproc -mode=<mode> <file>
```
## Arguments

Flag	Description	                                         Required
-mode	Processing mode: lines, words, hash, or compress	 Yes
<file>	Path to the file to process	                         Yes

## Modes

lines – counts the number of lines in the file

words – counts the number of words (whitespace‑separated tokens)

hash – computes the SHA‑256 checksum of the file content

compress – compresses the file content using gzip and reports the size of the compressed data

## Examples
```bash
# Count lines in a text file
./fileproc -mode=lines sample.txt

# Count words
./fileproc -mode=words README.md

# Compute content hash
./fileproc -mode=hash data.bin

# Test compression
./fileproc -mode=compress large.log
```