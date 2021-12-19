package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		log.Fatal("No directory specified")
	}

	dir := flag.Arg(0)
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
		if answer == "y" {

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
