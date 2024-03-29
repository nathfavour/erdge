package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
    "strconv"
    "math/rand"
)

func combineFiles(inputDir, outputFilePath string) error {
    var files []string

    err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            files = append(files, path)
        }
        return nil
    })

    if err != nil {
        return err
    }

        // Create the output directory if it doesn't exist
        outputDir := filepath.Dir(outputFilePath)
        if _, err := os.Stat(outputDir); os.IsNotExist(err) {
            err = os.MkdirAll(outputDir, 0755)
            if err != nil {
                return err
            }
        }

    outputFile, err := os.Create(outputFilePath)
    if err != nil {
        return err
    }
    defer outputFile.Close()

    writer := bufio.NewWriter(outputFile)

    for _, file := range files {
        data, err := ioutil.ReadFile(file)
        if err != nil {
            return err
        }

        _, err = writer.WriteString(string(data) + "\n")
        if err != nil {
            return err
        }
    }

    return writer.Flush()
}



func splitFiles(filePath string, linesPerFile int) error {
    file, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)

    var lines []string
    fileCount := 1

    for scanner.Scan() {
        lines = append(lines, scanner.Text())

        if len(lines) == linesPerFile {
            err = writeLinesToFile(lines, "split_"+strconv.Itoa(fileCount)+".txt")
            if err != nil {
                return err
            }

            lines = []string{}
            fileCount++
        }
    }

    if len(lines) > 0 {
        err = writeLinesToFile(lines, "split_"+strconv.Itoa(fileCount)+".txt")
        if err != nil {
            return err
        }
    }

    return scanner.Err()
}

func writeLinesToFile(lines []string, filePath string) error {
    file, err := os.Create(filePath)
    if err != nil {
        return err
    }
    defer file.Close()

    for _, line := range lines {
        _, err := file.WriteString(line + "\n")
        if err != nil {
            return err
        }
    }

    return nil
}


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
    linesPerFilePtr := flag.Int("x", 10, "Number of lines per file")
    outputPtr := flag.String("p", "", "Path of the output file")

    flag.Parse()

    args := flag.Args()
    if len(args) < 1 {
        panic("Please specify a file or directory path")
    }
    path := args[0]

	fileInfo, err := os.Stat(path)
    if err != nil {
        panic(err)
    }

    if *nPtr > 0 || *mPtr > 0 {
        // Perform processPath workflow
        if err := processPath(path, *nPtr, *mPtr); err != nil {
            panic(err)
        }


        } else if fileInfo.IsDir() {
            // Perform combineFiles workflow
            if *outputPtr == "" {
                fmt.Print("Enter the desired pathname of the final combined file: ")
                fmt.Scanln(outputPtr)
    
                if *outputPtr == "" {
                    *outputPtr = fmt.Sprintf("combined_%d.txt", rand.Int())
                    fmt.Println("No name provided. The output file will be named:", *outputPtr)
                }
            }
    
            err := combineFiles(path, *outputPtr)
            if err != nil {
                fmt.Println("Error combining files:", err)
            } else {
                fmt.Println("Files combined successfully into:", *outputPtr)
            }
        } else {
            // Perform splitFiles workflow
            err := splitFiles(path, *linesPerFilePtr)
            if err != nil {
                fmt.Println("Error splitting file:", err)
            } else {
                fmt.Println("File split successfully")
            }
        }
    }