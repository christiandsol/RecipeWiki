package main

import (
	"fmt"
	"os"
)

func isLocal(args []string) bool {
	// fmt.Printf("OS args[0]: %v", args)
	hasArg := false
	for _, arg := range args {
		if arg == "--help" {
			fmt.Printf("\t--help: get information\n")
			fmt.Printf("\t--local: develop local (default)\n")
			fmt.Printf("\t--cloud: run cloud\n")
			hasArg = true
		} else if arg == "--cloud" {
			fmt.Printf("\t Running via cloud")
			hasArg = true
			return false
		}
	}
	if hasArg == false {
		fmt.Printf("\t--help: get information\n")
		fmt.Println()
	}
	return true
}

func parseToEqual(str string) (string, string) {
	key := ""
	value := ""
	hitEquals := false
	for i := range str {
		if string(str[i]) == "=" {
			hitEquals = true
			continue
		}
		if !hitEquals {
			key += string(str[i])
		} else {
			value += string(str[i])
		}
	}
	return key, value
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
