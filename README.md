# Installer Bundler
A tool to bundle executable files into a single installer binary.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
<br>

## Setup
Download the [latest release](https://github.com/MatthiasHarzer/installer-bundler/releases) and add the executable to your `PATH`.

### Usage Requirements 
> These are even required when using the pre-built binaries, since this tool will build another binary from code it generates.
- [Go](https://golang.org/dl/) (version 1.25 or higher) in order to build the installer from the generated source code.
- [Make](https://www.gnu.org/software/make/) in order to run the Makefile to simplify the build process.

## Usage
1. Create a resource file containing a title and a download URL for each file to bundle. It should be a comma-separated values with one pair per line, e.g.:
   ```
   File 1,https://example.com/file1.exe
   File 2,https://example.com/file2.exe
   ```
2. Run `installer-bundler bundle -f <resource-file> -o <output-file>` to create the installer binary.
    - To include the binaries directly, instead of downloading them at runtime, use the `--embed` / `-e` flag. This will greatly increase the size of the installer binary.
3. Use the generated installer binary:
   - `output.exe list` to list all included files.
   - `output.exe extract -d <output-dir>` to download all files, or copy them from the embedded resources if the `--embed` / `-e` flag was used during bundling.
   - `output.exe install` to extract and install all files (if applicable). Use the `--parallel` / `-p` flag to install all files in parallel.