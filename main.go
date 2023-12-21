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

// using string for debugging purposes in errors
// int/iota would probably be more efficient
const (
	BYTES = "bytes"
	LINES = "lines"
	WORDS = "words"
	CHARS = "chars"
)

type CountOptionType string

// this is just me avoiding forloop + switch statement
var optionsTable = map[string]CountOptionType{
	"c": BYTES,
	"w": WORDS,
	"l": LINES,
	"m": CHARS,
}

// used for counting bytes when filepath is given
// faster to use this for counting bytes than bufio scanner
func getFileSize(filePath string) (int, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return -1, fmt.Errorf("error getting file info: %w", err)
	}
	return int(fileInfo.Size()), err
}

func count(path string, in io.ReadSeeker, option CountOptionType) (int, error) {
	_, err := in.Seek(0, io.SeekStart)
	if err != nil {
		return -1, fmt.Errorf("error seeking file end: %w", err)
	}

	scanner := bufio.NewScanner(in)
	switch option {
	case LINES:
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
		log.Fatalf("UNREACHABLE: %s", option)
	}

	var count int
	for scanner.Scan() {
		count += 1
	}
	return count, nil
}

func getResultsFromOptions(fpath string, in io.ReadSeeker, options []string) (string, error) {
	var result string
	for _, v := range options {
		option, ok := optionsTable[v]
		if !ok {
			return "", fmt.Errorf("invalid option: '%s'", v)
		}

		res, err := count(fpath, in, option)
		if err != nil {
			return "", fmt.Errorf("error counting %s: %w", option, err)
		}
		result += fmt.Sprintf("%d ", res)
	}
	return result, nil
}

func main() {
	args := os.Args[1:]

	var options []string
	var paths []string

	var file *os.File
	defer file.Close()

	// use this to get options and filepath from args
	for _, v := range args {
		if v[0] == '-' {
			options = append(options, strings.Split(v[1:], "")...)
		} else {
			paths = append(paths, v)
		}
	}

	if len(options) == 0 {
		options = append(options, "l", "w", "c")
	}

	var result string
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

		res, err := getResultsFromOptions("", bytes.NewReader(data), options)
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

			res, err := getResultsFromOptions(path, file, options)
			if err != nil {
				log.Fatalf("error getting results: %s", err.Error())
			}

			res += path
			result += fmt.Sprintln(res)
		}
	}
	fmt.Print(result)
}
