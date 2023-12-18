# ERDGE

ERDGE is a command-line tool written in Go that allows you to delete a specific number of lines from the start and end of a file. This tool is particularly useful when dealing with very large files where manual editing is not feasible.

## Features

- Remove a specified number of lines from the start of a file.
- Remove a specified number of lines from the end of a file.
- Process individual files or entire directories.
- Cross-platform support (Windows, macOS, Linux).

## Installation

Download the appropriate binary for your operating system from the releases page. The binary files are named `erdgeWin` for Windows, `erdgeMac` for macOS, and `erdgeLin` for Linux.

## Usage

You can run ERDGE from the command line using the following syntax:

### Windows

```pwsh
./erdgeWin -n <num_start> -m <num_end> <path_to_file_or_directory>```

### macOS

```./erdgeMac -n <num_start> -m <num_end> <path_to_file_or_directory>```

### Linux

```./erdgeLin -n <num_start> -m <num_end> <path_to_file_or_directory>>```

Replace <num_start> with the number of lines you want to remove from the start of the file, <num_end> with the number of lines you want to remove from the end of the file, and <path_to_file_or_directory> with the path to the file or directory you want to process.

If no path is specified, ERDGE will use the "input" folder in the current working directory.

## Contributing
Contributions are welcome! Please feel free to submit a pull request.

## License