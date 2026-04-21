package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
)

type Token struct {
	token []byte
	next  *Token
}

type TokenList struct {
	head *Token
	tail *Token
	num  uint
}

func newToken(token []byte) *Token {
	return &Token{
		token: token,
		next:  nil,
	}
}

func (t *TokenList) Append(token []byte) {
	newToken := newToken(token)
	if t.tail == nil {
		t.head = newToken
		t.tail = newToken
		t.num = 1
		return
	}

	t.tail.next = newToken
	t.tail = newToken
	t.num++
}

func (t *TokenList) Rename() {
	_, fileName := fileNames(t)

	cmd := exec.Command(os.Getenv("EDITOR"), fileName)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("the following files will be renamed")
	newNames, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	temp := t.head
	newFileName := getFileNames(newNames)
	i := 0
	for temp != nil {
		fmt.Println(string(temp.token), "->", string(newFileName[i]))
		i++
		temp = temp.next
	}
	fmt.Print("is this okay [y/n]:")
	var input string
	if _, err := fmt.Scanf("%s", &input); err != nil {
		log.Fatal(err)
	}
	if input == "n" {
		if err := os.Remove(fileName); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}

	temp = t.head
	i = 0
	for temp != nil {
		cmd := exec.Command("mv", string(temp.token), string(newFileName[i]))
		cmd.Stderr = os.Stderr
		_ = cmd.Run()

		i++
		temp = temp.next
	}

	if err := os.Remove(fileName); err != nil {
		log.Fatal(err)
	}
}

func getFileNames(newNames []byte) [][]byte {
	start := 0
	end := 0
	var newFileNames [][]byte
	for end != len(newNames) {
		if newNames[end] == '\n' {
			newFileNames = append(newFileNames, newNames[start:end])
			start = end + 1
		}
		end++
	}
	return newFileNames
}

func fileNames(t *TokenList) (*os.File, string) {
	fileName := strconv.Itoa(rand.Int())
	renameFile, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	temp := t.head
	for temp != nil {
		if _, err := renameFile.Write(temp.token); err != nil {
			log.Fatal(err)
		}
		if _, err := renameFile.Write([]byte("\n")); err != nil {
			log.Fatal(err)
		}

		temp = temp.next
	}

	if err := renameFile.Close(); err != nil {
		log.Fatal(err)
	}

	return renameFile, fileName
}

func (t *TokenList) Print() {
	if t.head == nil {
		fmt.Println("there are not tokens")
	}

	temp := t.head
	for temp != nil {
		fmt.Println(temp)
		temp = temp.next
	}
}
