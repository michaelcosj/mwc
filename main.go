package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const (
	BYTES = iota
	LINES
	WORDS
	CHARS
	MAX_LINE
)

type CountOptionType int

// used for counting bytes when filepath is given
// apparently more efficient than manually counting
func getFileSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -1, fmt.Errorf("error getting file info: %w", err)
	}
	return int(fileInfo.Size()), err
}

func count(path string, in io.ReadSeeker, opt CountOptionType) (int, error) {
	_, err := in.Seek(0, io.SeekStart)
	if err != nil {
		return -1, fmt.Errorf("error seeking file end: %w", err)
	}

	scanner := bufio.NewScanner(in)
	switch opt {
	case LINES, MAX_LINE:
		scanner.Split(bufio.ScanLines)
	case WORDS:
		scanner.Split(bufio.ScanWords)
	case CHARS:
		scanner.Split(bufio.ScanRunes)
	case BYTES:
		if len(path) > 0 {
			return getFileSize(path)
		}
		scanner.Split(bufio.ScanBytes)
	default:
		log.Fatalf("UNREACHABLE: %d", opt)
	}

	var count int
	for scanner.Scan() {
		if opt == MAX_LINE {
			lineLen := len(scanner.Text())
			if count < lineLen {
				count = lineLen
			}
			continue
		}
		count += 1
	}
	return count, nil
}

// evaluates the file based on the options and returns the results
func getResults(fpath string, in io.ReadSeeker, options []string) (string, error) {
	var result string
	for _, option := range options {
		var res int
		var err error

		switch option {
		case "c":
			res, err = count(fpath, in, BYTES)
		case "l":
			res, err = count(fpath, in, LINES)
		case "L":
			res, err = count(fpath, in, MAX_LINE)
		case "w":
			res, err = count(fpath, in, WORDS)
		case "m":
			res, err = count(fpath, in, CHARS)
		}

		if err != nil {
			return "", fmt.Errorf("error counting: %w", err)
		}
		result += fmt.Sprintf("%d ", res)
	}

	result += fmt.Sprintln(fpath)
	return result, nil
}

func parseArgs(args []string) ([]string, []string) {
	var options []string
	var paths []string
	for _, arg := range args {
		if arg[0] == '-' {
			// TODO: remove duplicates
			options = append(options, strings.Split(arg[1:], "")...)
		} else {
			paths = append(paths, arg)
		}
	}

	if len(options) == 0 {
		options = append(options, "l", "w", "c")
	}
	return options, paths
}

func main() {
	options, paths := parseArgs(os.Args[1:])
	var file *os.File
	var result string

	// if no paths, it means we were probably passed in data from stdin
	if len(paths) == 0 {
		file = os.Stdin
		defer file.Close()

		// apparently you can't seek stdin so i'm reading all the data
		// and then creating a reader seeker from that
		// not sure if there's a better way though
		data, err := io.ReadAll(file)
		if err != nil {
			log.Fatalf("error reading stdin: %s", err.Error())
		}

		res, err := getResults("", bytes.NewReader(data), options)
		if err != nil {
			log.Fatalf("error getting results: %s", err.Error())
		}
		result = fmt.Sprintln(res)
	} else {
		for _, path := range paths {
			file, err := os.Open(path)
			if err != nil {
				log.Fatalf("error opening file: %s", err.Error())
			}
			defer file.Close()

			res, err := getResults(path, file, options)
			if err != nil {
				log.Fatalf("error getting results: %s", err.Error())
			}
			result += fmt.Sprint(res)
		}
	}
	fmt.Print(result)
}
