package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
)

var version string

func main() {
	help := flag.Bool("help", false, "Show help message")
	versionFlag := flag.Bool("version", false, "Show version")
	dryRun := flag.Bool("dry-run", false, "Simulate moving files without making changes")
	skipConfirmation := flag.Bool("skip-confirmation", false, "Skip confirmation when moving files")

	flag.BoolVar(help, "h", false, "Show help message")
	flag.BoolVar(versionFlag, "v", false, "Show version")
	flag.BoolVar(dryRun, "d", false, "Simulate moving files without making changes")
	flag.BoolVar(skipConfirmation, "y", false, "Skip confirmation when moving files")

	flag.Parse()

	// Handle --help and --version flags before doing anything else
	if *help {
		flag.PrintDefaults()
		return
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatalf("debug.ReadBuildInfo() failed")
	}
	version = buildInfo.Main.Version

	if *versionFlag {
		fmt.Printf("%s, version %s\n", filepath.Base(os.Args[0]), version)
		return
	}

	// Ensure a directory is passed
	if len(flag.Args()) != 1 {
		log.Fatalf("Error: No directory specified\n\n")
		flag.PrintDefaults()
		return
	}

	dir := flag.Arg(0)

	// Check if the directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Fatalf("The directory %s does not exist.\n", dir)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatalf("Error reading directory: %v\n", err)
	}
	fmt.Printf("Found %d files in the directory.\n", len(files))

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		fullPath, err := filepath.Abs(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatalf("Error getting absolute path for %s: %v\n", file.Name(), err)
		}

		// Ignore files starting with .
		if strings.HasPrefix(file.Name(), ".") {
			continue
		}

		containingFolder := filepath.Base(filepath.Dir(fullPath))
		fileInfo, err := file.Info()
		if err != nil {
			log.Fatalf("Error getting info for %s: %v\n", fullPath, err)
		}

		yyyymmdd := fileInfo.ModTime().Format("2006-01-02")
		if yyyymmdd != containingFolder {
			newPath := filepath.Join(dir, yyyymmdd)
			newFullPath := filepath.Join(newPath, file.Name())

			if _, err := os.Stat(newPath); os.IsNotExist(err) {
				err := os.MkdirAll(newPath, 0755)
				if err != nil {
					log.Fatalf("Error creating directory %s: %v\n", newPath, err)
				}
			}

			if *skipConfirmation || promptUser(fullPath, newFullPath) {
				if *dryRun {
					fmt.Printf("[DRY RUN] Would move %s to %s\n", fullPath, newFullPath)
				} else {
					err := os.Rename(fullPath, newFullPath)
					if err != nil {
						log.Fatalf("Error moving %s to %s: %v\n", fullPath, newFullPath, err)
					}
					fmt.Printf("Moved %s to %s\n", fullPath, newFullPath)
				}
			}
		}
	}
}

// promptUser prompts the user for confirmation to move the file.
func promptUser(fullPath, newFullPath string) bool {
	fmt.Printf("Move %s to %s? (y/n) ", fullPath, newFullPath)
	var answer string
	_, err := fmt.Scanln(&answer)
	if err != nil {
		log.Fatalf("Error reading user input: %v\n", err)
	}
	return strings.ToLower(strings.TrimSpace(answer)) == "y"
}
