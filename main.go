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




func removeLines(filePath string, n, m, q, r int, s, t []string) error {
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

    var newLines []string
    for i, line := range lines {
        if (i+1)%q != 0 && (i+1)%r != 0 {
            newLines = append(newLines, line)
        } else if (i+1)%r == 0 {
            if len(s) > 0 {
                for _, str := range s {
                    if strings.Contains(line, str) {
                        newLines = append(newLines, line)
                        break
                    }
                }
            }
            if len(t) > 0 {
                for _, str := range t {
                    if !strings.Contains(line, str) {
                        newLines = append(newLines, line)
                        break
                    }
                }
            }
        }
    }

    // Remove the first n lines and the last m lines
    if n > 0 {
        newLines = newLines[n:]
    }
    if m > 0 {
        newLines = newLines[:len(newLines)-m]
    }

    output := strings.Join(newLines, "\n")
    return ioutil.WriteFile(filePath, []byte(output), 0644)
}





// func processPath(path string, n, m int) error {
func processPath(path string, n, m, p, q int, r, s []string) error {

	return filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err := removeLines(filePath, n, m, p, q, r, s); err != nil {
				return err
			}
		}

		return nil
	})
}




func main() {
    nPtr := flag.Int("n", 0, "Number of lines to remove from the start of the file")
    mPtr := flag.Int("m", 0, "Number of lines to remove from the end of the file")
    qPtr := flag.Int("q", 0, "Line number to be deleted")
    rPtr := flag.Int("r", 0, "Line number to be deleted relative to -q")
    sPtr := flag.String("s", "", "Lines containing these strings will be deleted")
    tPtr := flag.String("t", "", "Lines not containing these strings will be deleted")
    linesPerFilePtr := flag.Int("x", 10, "Number of lines per file")
    outputPtr := flag.String("p", "", "Path of the output file")
    hPtr := flag.Bool("h", false, "Print help message")

    flag.Parse()

    if *hPtr {
        flag.Usage()
        os.Exit(0)
    }

    s := strings.Split(*sPtr, ",")
    t := strings.Split(*tPtr, ",")

    args := flag.Args()
    if len(args) < 1 {
        panic("Please specify a file or directory path")
    }
    path := args[0]

    fileInfo, err := os.Stat(path)
    if err != nil {
        panic(err)
    }

    if *nPtr > 0 || *mPtr > 0 || *qPtr > 0 || *rPtr > 0 || *sPtr != "" || *tPtr != "" {
        // Perform removeLines workflow
        if err := removeLines(path, *nPtr, *mPtr, *qPtr, *rPtr, s, t); err != nil {
            panic(err)
        }





        // if fileInfo.IsDir() {
        //     // Perform removeLines workflow on each file in the directory
        //     if err := processPath(path, *nPtr, *mPtr, *qPtr, *rPtr, s, t); err != nil {
        //         panic(err)
        //     }
        // } else if *nPtr > 0 || *mPtr > 0 || *qPtr > 0 || *rPtr > 0 || *sPtr != "" || *tPtr != "" {
        //     // Perform removeLines workflow on the single file
        //     if err := removeLines(path, *nPtr, *mPtr, *qPtr, *rPtr, s, t); err != nil {
        //         panic(err)
        //     }
        // }



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