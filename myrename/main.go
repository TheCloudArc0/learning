package main

import (
	"fmt"
	"log"
	"math/rand/v2"
	"os"
	"os/exec"
	"slices"
	"strconv"
)

func main() {
	if len(os.Args) <= 1 {
		rename()
	} else if len(os.Args) >= 1 {
		if os.Args[1] == "help" {
			fmt.Println("Usage: rename\n\t\tthis will rename file in the current folder")
		}
	}
}

func rename() {
	renameFile := "." + strconv.Itoa(rand.Int())
	tempFile, err := os.Create(renameFile)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command("ls", "-F")
	cmd.Stdout = tempFile
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	oldFileNames := fileNames(renameFile)
	openEditor(renameFile)
	newFileNames := fileNames(renameFile)
	compareAndRename(oldFileNames, newFileNames, renameFile)

	removeRenameFile(renameFile)
}

func fileNames(renameFile string) [][]byte {
	s, err := os.ReadFile(renameFile)
	if err != nil {
		log.Fatal(err)
	}
	return fileList(s)
}

func openEditor(renameFile string) {
	cmd := exec.Command(os.Getenv("EDITOR"), renameFile)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func fileList(s []byte) [][]byte {
	start := 0
	end := 0
	var fileNames [][]byte
	for end != len(s) {
		if s[end] == '\n' {
			fileNames = append(fileNames, s[start:end])
			start = end + 1
		}
		end++
	}
	return fileNames
}

func compareAndRename(oldName, newName [][]byte, renameFile string) {
	var oldFiles [][]byte
	var newFiles [][]byte
	for i := range oldName {
		if slices.Equal(oldName[i], newName[i]) {
			continue
		} else {
			oldFiles = append(oldFiles, oldName[i])
			newFiles = append(newFiles, newName[i])
		}
	}

	if newFiles == nil {
		fmt.Println("nothin to rename")
		removeRenameFile(renameFile)
		os.Exit(0)
	}

	fmt.Println("the following will be renamed:\n==============================")
	for i := range oldFiles {
		fmt.Println(string(oldFiles[i]), "->", string(newFiles[i]))
	}
	fmt.Print("\nis this OK [y/n]: ")
	var input string
	if _, err := fmt.Scanf("%s", &input); err != nil {
		log.Fatal(err)
	}

	switch input {
	case "n":
		removeRenameFile(renameFile)
		os.Exit(0)
	case "y":
		for i := range oldFiles {
			cmd := exec.Command("mv", string(oldFiles[i]), string(newFiles[i]))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}
	}
}

func removeRenameFile(renameFile string) {
	if err := os.Remove(renameFile); err != nil {
		log.Fatal(err)
	}
}
