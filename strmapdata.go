package rplanlib

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func GetInputStrStrMapFromFile(f string) (*map[string]string, error) {
	file, err := os.Open(f)
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, e
	}
	defer file.Close()

	ipsm := NewInputStringsMap()

	scanner := bufio.NewScanner(file)

	// Default scanner is bufio.ScanLines. Lets use ScanWords.
	// Could also use a custom function of SplitFunc type
	//scanner.Split(bufio.ScanWords)

	var key, val string
	// Scan for next token.
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Printf("Line: %s\n", line)
		if line != "" && line[0] == '#' {
			fmt.Printf("Skipping Line: %s\n", line)
			continue
		}
		tokens := strings.SplitAfter(line, "'")
		first := true
		havePair := false
		for _, token := range tokens {
			if strings.Index(token, ":") != -1 {
				continue
			}
			//fmt.Printf("token: %s\n", token)
			if first && strings.Index(token, "'") != -1 {
				key = strings.TrimRight(token, "'")
				key = strings.TrimSpace(key)
				//fmt.Printf("key: %s\n", key)
				first = false
			} else if strings.Index(token, "'") != -1 {
				val = strings.TrimRight(token, "'")
				val = strings.TrimSpace(val)
				//fmt.Printf("val: %s\n", val)
				havePair = true
				break
			}
		}
		if havePair {
			//fmt.Printf("key: %s, val: %s\n", key, val)
			err := setStringMapValueWithValue(&ipsm, key, val)
			if err != nil {
				e := fmt.Errorf("Error: %s", err)
				return nil, e
			}
		}
	}
	// False on error or EOF. Check error
	err = scanner.Err()
	if err != nil {
		e := fmt.Errorf("Error: %s", err)
		return nil, e
	}
	return &ipsm, nil
}
