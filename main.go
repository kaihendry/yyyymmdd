package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"time"
)

var version = "dev"

type options struct {
	dryRun        bool
	yes           bool
	includeHidden bool
	minAge        time.Duration
}

type stats struct{ found, moved, skipped, conflicts int }

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }

func run(args []string, in io.Reader, out, errOut io.Writer) int {
	fs := flag.NewFlagSet("yyyymmdd", flag.ContinueOnError)
	fs.SetOutput(errOut)
	fs.Usage = func() {
		fmt.Fprintf(errOut, "Usage: yyyymmdd [options] [directory]\n\nOrganise files into YYYY-MM-DD folders by modification date.\nIf directory is omitted, ~/Downloads is used.\n\nOptions:\n")
		fs.PrintDefaults()
	}
	var opts options
	var help, showVersion bool
	fs.BoolVar(&help, "help", false, "show this help message")
	fs.BoolVar(&help, "h", false, "show this help message")
	fs.BoolVar(&showVersion, "version", false, "show version")
	fs.BoolVar(&showVersion, "v", false, "show version")
	fs.BoolVar(&opts.dryRun, "dry-run", false, "preview changes without modifying anything")
	fs.BoolVar(&opts.dryRun, "d", false, "preview changes without modifying anything")
	fs.BoolVar(&opts.yes, "yes", false, "move without asking for confirmation")
	fs.BoolVar(&opts.yes, "y", false, "move without asking for confirmation")
	fs.BoolVar(&opts.yes, "skip-confirmation", false, "move without asking for confirmation (legacy alias)")
	fs.BoolVar(&opts.includeHidden, "include-hidden", false, "include files whose names start with a dot")
	fs.DurationVar(&opts.minAge, "older-than", 0, "only organise files older than this (for example 10m or 2h)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if help {
		fs.Usage()
		return 0
	}
	if showVersion {
		if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "(devel)" {
			version = info.Main.Version
		}
		fmt.Fprintf(out, "yyyymmdd %s\n", version)
		return 0
	}
	if fs.NArg() > 1 {
		fmt.Fprintln(errOut, "error: expected at most one directory")
		fs.Usage()
		return 2
	}
	dir := ""
	if fs.NArg() == 1 {
		dir = fs.Arg(0)
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintf(errOut, "error: find home directory: %v\n", err)
			return 1
		}
		dir = filepath.Join(home, "Downloads")
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		fmt.Fprintf(errOut, "error: resolve %q: %v\n", dir, err)
		return 1
	}
	s, err := organise(absDir, opts, in, out, time.Now())
	if err != nil {
		fmt.Fprintf(errOut, "error: %v\n", err)
		return 1
	}
	action := "Moved"
	if opts.dryRun {
		action = "Would move"
	}
	fmt.Fprintf(out, "\n%s %d of %d files", action, s.moved, s.found)
	if s.conflicts > 0 || s.skipped > 0 {
		fmt.Fprintf(out, " (%d skipped, %d conflicts)", s.skipped, s.conflicts)
	}
	fmt.Fprintln(out, ".")
	return 0
}

func organise(dir string, opts options, in io.Reader, out io.Writer, now time.Time) (stats, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return stats{}, fmt.Errorf("read %s: %w", dir, err)
	}
	reader := bufio.NewReader(in)
	var s stats
	for _, entry := range entries {
		if entry.IsDir() || (!opts.includeHidden && strings.HasPrefix(entry.Name(), ".")) {
			continue
		}
		s.found++
		info, err := entry.Info()
		if err != nil {
			return s, fmt.Errorf("inspect %s: %w", entry.Name(), err)
		}
		if opts.minAge > 0 && now.Sub(info.ModTime()) < opts.minAge {
			s.skipped++
			continue
		}
		source := filepath.Join(dir, entry.Name())
		destinationDir := filepath.Join(dir, info.ModTime().Format("2006-01-02"))
		destination := filepath.Join(destinationDir, entry.Name())
		if _, err := os.Lstat(destination); err == nil {
			fmt.Fprintf(out, "skip  %s (already exists)\n", relativeMove(dir, source, destination))
			s.conflicts++
			continue
		} else if !errors.Is(err, os.ErrNotExist) {
			return s, fmt.Errorf("check destination %s: %w", destination, err)
		}
		if !opts.yes && !opts.dryRun {
			fmt.Fprintf(out, "move  %s? [y/N] ", relativeMove(dir, source, destination))
			answer, err := reader.ReadString('\n')
			if err != nil && !errors.Is(err, io.EOF) {
				return s, fmt.Errorf("read confirmation: %w", err)
			}
			if strings.ToLower(strings.TrimSpace(answer)) != "y" {
				s.skipped++
				continue
			}
		}
		if opts.dryRun {
			fmt.Fprintf(out, "move  %s\n", relativeMove(dir, source, destination))
			s.moved++
			continue
		}
		if err := os.MkdirAll(destinationDir, 0755); err != nil {
			return s, fmt.Errorf("create %s: %w", destinationDir, err)
		}
		if err := os.Rename(source, destination); err != nil {
			return s, fmt.Errorf("move %s: %w", entry.Name(), err)
		}
		fmt.Fprintf(out, "moved %s\n", relativeMove(dir, source, destination))
		s.moved++
	}
	return s, nil
}

func relativeMove(base, source, destination string) string {
	from, _ := filepath.Rel(base, source)
	to, _ := filepath.Rel(base, destination)
	return fmt.Sprintf("%s → %s", from, to)
}
