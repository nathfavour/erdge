package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func removeLines(filePath string, n, m int) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	totalLines := len(lines)
	if n+m >= totalLines {
		n, m = totalLines, 0
	}

	lines = lines[n : totalLines-m]
	output := strings.Join(lines, "\n")

	return ioutil.WriteFile(filePath, []byte(output), 0644)
}

func processPath(path string, n, m int) error {
	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err := removeLines(filePath, n, m); err != nil {
				return err
			}
		}

		return nil
	})
}

func main() {
	nPtr := flag.Int("n", 0, "Number of lines to remove from the start of the file")
	mPtr := flag.Int("m", 0, "Number of lines to remove from the end of the file")
	flag.Parse()

	args := flag.Args()
	path := "input"

	//     if len(args) != 1 {
	//         panic("Please specify a file or directory path")
	//     }

	//     if err := processPath(args[0], *nPtr, *mPtr); err != nil {
	//         panic(err)
	//     }
	// }
	if len(args) > 0 {
		path = args[0]
	}

	if err := processPath(path, *nPtr, *mPtr); err != nil {
		panic(err)
	}
}
