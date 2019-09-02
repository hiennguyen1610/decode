package decode

import (
	"fmt"
	"strconv"
	"math"
	"os"
	"io"
	"bufio"
)

var messages = []string{}
var message = ""
var isNewMessage = true
var isReadingHeader = true
var header = ""
var isNewSegment = true
var segmentKey = ""
var segmentKeyLen = 0
var key = ""

func Decode(filePath string) []string {
	f, _ := os.Open(filePath)
	if f != nil {
		reader := bufio.NewReader(f)
		for {
			if c, _, err := reader.ReadRune(); err != nil {
				if err == io.EOF {
					break
				} else {
					fmt.Println(err)
				}
			} else {
				check(c)
			}
		}
	} else {
		fmt.Println("File not found: ", filePath)
	}
	return messages
}

func check(c rune) {
	// Ignore all line break character
	if isNewMessage && isEOL(c) {
		return
	}
	isNewMessage = false
	if isReadingHeader {
		writeHeader(c)
	} else {
		if !isEOL(c) {
			if isNewSegment {
				writeSegmentKey(c)
			} else {
				writeKey(c)
			}
		}
	}
}

func writeHeader(c rune) {
	if isEOL(c) {
		isReadingHeader = false
		fmt.Printf("HEADER: %v\n", header)
	} else {
		header += string(c)
	}
}

func writeSegmentKey(c rune) {
	segmentKey += string(c)
	if len(segmentKey) == 3 {
		v, _ := strconv.ParseInt(segmentKey, 2, 0)
		segmentKeyLen = int(v)
		if segmentKeyLen == 0 {
			fmt.Printf("MESSAGE DECODE DONE: %v \n\n\n", segmentKey)
			messages = append(messages, message)
			isNewMessage = true
			message = ""
			header = ""
			isReadingHeader = true
			resetSegment()
		} else {
			fmt.Printf("Segment start: %v %v\n", segmentKey, segmentKeyLen)
			isNewSegment = false
		}
	}
}

func writeKey(c rune) {
	key += string(c)
	if len(key) == segmentKeyLen {

		if isEndSegment(key) {
			fmt.Printf("Segment end: %v\n", key)
			resetSegment()
		} else {
			value := header[indexOf(key):indexOf(key) + 1]
			fmt.Printf("Key: %v index: %v value: %v\n", key, indexOf(key), value)
			message += value
		}
		key = ""
	}
}

func resetSegment() {
	isNewSegment = true
	segmentKey = ""
	segmentKeyLen = 0
}

func isEOL(c rune) bool {
	if c == 10 || c == 13 {
		return true
	}
	return false
}

/*
	If the key contains a sequence of "1", so that is end of segment.
*/
func isEndSegment(key string) bool {
	for _, char := range key {
		if char != rune('1') {
			return false
		}
	}
	return true
}

/*
	Return index of character in header by key
	Actual index = start index + value of key
*/
func indexOf(key string) int {
	n := len(key)
	start := 0
	for i := 1; i < n; i++ {
		start += int(math.Pow(2, float64(i))) - 1
	}
	v, _ := strconv.ParseInt(key, 2, 0)
	return start + int(v)
}

