package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Print("enter the names of files you want to rename (i.e file.txt file2.txt...)\n>>> ")
	input, err := bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err != nil {
		log.Fatal("error while reading user input): ", err)
	}

	rename(input)
}

func rename(input []byte) {
	input = trimSpaces(input)
	start := 0
	end := 0
	var tokenList TokenList

	for end != len(input)-1 {
		if input[end] == ' ' || input[end] == '\n' {
			token := input[start:end]
			tokenList.Append(token)
			start = end + 1
		}
		end++
	}

	tokenList.Rename()
}

func trimSpaces(input []byte) []byte {
	start := 0
	for input[start] == ' ' {
		start++
	}

	end := len(input) - 1
	for input[end] == ' ' {
		end--
	}

	input = append(input, '\n')
	return input
}
