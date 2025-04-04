package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

/*
 * Complete the 'customSorting' function below.
 *
 * The function is expected to return a STRING_ARRAY.
 * The function accepts STRING_ARRAY strArr as parameter.
 */
func customSorting(strArr []string) []string {
	// Use sort.Slice with a custom comparison function
	sort.Slice(strArr, func(i, j int) bool {
		// Compare based on the length parity first (odd lengths before even)
		lenI, lenJ := len(strArr[i]), len(strArr[j])

		if lenI%2 != lenJ%2 {
			// Odd-length strings (len % 2 != 0) should come first
			return lenI%2 != 0
		}

		// If both are odd or both are even, compare by length (shorter before longer for odd, longer before shorter for even)
		if lenI != lenJ {
			if lenI%2 == 1 {
				// For odd lengths, the shorter one should come first
				return lenI < lenJ
			}
			// For even lengths, the longer one should come first
			return lenI > lenJ
		}

		// If lengths are equal, compare alphabetically
		return strArr[i] < strArr[j]
	})

	return strArr
}

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	stdout, err := os.Create(os.Getenv("OUTPUT_PATH"))
	checkError(err)

	defer stdout.Close()

	writer := bufio.NewWriterSize(stdout, 16*1024*1024)

	strArrCount, err := strconv.ParseInt(strings.TrimSpace(readLine(reader)), 10, 64)
	checkError(err)

	var strArr []string

	for i := 0; i < int(strArrCount); i++ {
		strArrItem := readLine(reader)
		strArr = append(strArr, strArrItem)
	}

	result := customSorting(strArr)

	for i, resultItem := range result {
		fmt.Fprintf(writer, "%s", resultItem)

		if i != len(result)-1 {
			fmt.Fprintf(writer, "\n")
		}
	}

	fmt.Fprintf(writer, "\n")

	writer.Flush()
}

func readLine(reader *bufio.Reader) string {
	str, _, err := reader.ReadLine()
	if err == io.EOF {
		return ""
	}

	return strings.TrimRight(string(str), "\r\n")
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
