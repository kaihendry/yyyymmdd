package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const version = "0.0.1"

func main() {

	fmt.Printf("%s, version %s\n", filepath.Base(os.Args[0]), version)

	if len(os.Args) != 2 {
		log.Fatalf("No directory specified")
	}

	dir := os.Args[1]
	log.Printf("Organising %s", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath, err := filepath.Abs(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}

		// ignore files starting with .
		if strings.HasPrefix(file.Name(), ".") {
			log.Printf("Ignoring dot file %s", fullPath)
			continue
		}

		containingFolder := filepath.Base(filepath.Dir(fullPath))

		fileInfo, err := file.Info()
		if err != nil {
			log.Fatal(err)
		}

		yyyymmdd := fileInfo.ModTime().Format("2006-01-02")
		if yyyymmdd != containingFolder {
			log.Printf("%s is not in the right folder: %s", fullPath, yyyymmdd)
		}

		// new path in "yyyy-mm-dd" folder
		newPath := filepath.Join(dir, yyyymmdd)
		newFullPath := filepath.Join(newPath, file.Name())

		// prompt user to move file
		fmt.Printf("Move %s to %s? (y/n) ", fullPath, newFullPath)
		var answer string
		fmt.Scanln(&answer)
		if strings.ToLower(strings.TrimSpace(answer)) == "y" {
			err = os.MkdirAll(newPath, 0755)
			if err != nil {
				log.Fatal(err)
			}

			err = os.Rename(fullPath, newFullPath)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
