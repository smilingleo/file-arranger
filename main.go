package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

func statMTime(path string) (mYear int, mMonth time.Month, mDay int) {
	var st syscall.Stat_t
	if err := syscall.Stat(path, &st); err != nil {
		log.Fatal(err)
	}
	mYear, mMonth, mDay = time.Unix(st.Mtimespec.Sec, 0).Date()
	return mYear, mMonth, mDay
}

func mkdirIfNotExist(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0700)
	}
}

func removeTailingSlash(str *string) string {
	s := *str
	if s[len(s)-1:] == "/" {
		s = s[:len(s)-1]
	}
	return s
}

func parseArguments() (fromDir, toBaseDir string) {
	searchDir := flag.String("from", "", "The folder you want to re-arrange")
	targetBase := flag.String("to", "", "The target folder the files will be put under, it will be sibling folder as source by default.")

	flag.Parse()

	if *searchDir == "" {
		fmt.Println("You have to specify a 'from' argument which is the source folder you want to re-arrange.")
		os.Exit(1)
	}
	fromDir = removeTailingSlash(searchDir)
	if _, err := os.Stat(fromDir); os.IsNotExist(err) {
		fmt.Printf("'%s' does not exists, please specify the source folder.\n", fromDir)
		os.Exit(2)
	}

	if *targetBase == "" {
		toBaseDir = fromDir[:strings.LastIndex(fromDir, "/")]
	} else {
		toBaseDir = removeTailingSlash(targetBase)
	}
	return fromDir, toBaseDir
}

func main() {
	fromDir, toBaseDir := parseArguments()

	mkdirIfNotExist(toBaseDir)

	// <year, <month, <day, fileName>>>
	var m map[int]map[time.Month]map[int][]string = make(map[int]map[time.Month]map[int][]string)

	err := filepath.Walk(fromDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			return nil
		}
		mYear, mMonth, mDay := statMTime(path)
		if m[mYear] == nil {
			m[mYear] = make(map[time.Month]map[int][]string)
		}
		if m[mYear][mMonth] == nil {
			m[mYear][mMonth] = make(map[int][]string)
		}
		m[mYear][mMonth][mDay] = append(m[mYear][mMonth][mDay], path)

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	// Year
	for year, monthMap := range m {
		yearFolder := fmt.Sprintf("%s/%d", toBaseDir, year)
		mkdirIfNotExist(yearFolder)

		for month, dayMap := range monthMap {
			monthFolder := fmt.Sprintf("%s/%02d", yearFolder, month)
			mkdirIfNotExist(monthFolder)

			for day, files := range dayMap {
				dayFolder := fmt.Sprintf("%s/%d", monthFolder, day)
				mkdirIfNotExist(dayFolder)

				for _, file := range files {
					newPath := fmt.Sprintf("%s/%s", dayFolder, file[strings.LastIndex(file, "/")+1:])
					os.Rename(file, newPath)
				}
			}
		}
	}
}
