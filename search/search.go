package search

import (
	"bufio"
	"config"
	"fmt"
	"os"
	"strings"
)

var DIR = config.DIR

/* takes to params: search term and filename (both strings) and returns true
if string was found in filename. helper function for search handler */
func Search(searchterm, fname string) bool {
	file, err := os.Open(DIR + fname)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// fmt.Println("Opening", file.Name()) // just for debugging

	defer file.Close()
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchterm) {
			// fmt.Fprintln(w, scanner.Text())
			// fmt.Println(scanner.Text()) // just for debugging
			return true
		}
	}
	return false
}
