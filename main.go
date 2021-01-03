package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/rjeczalik/notify"
)

var Version string

func main() {
	bin, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting", bin, Version)

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dldir := path.Join(home, "Downloads")

	n := make(chan notify.EventInfo, 1)

	err = notify.Watch(dldir, n, notify.All)
	if err != nil {
		log.Fatal(err)
	}

	defer notify.Stop(n)

	for {
		switch ei := <-n; ei.Event() {
		case notify.Create, notify.Rename:
			log.Println(ei)
			from := ei.Path()

			t := time.Now().Local()
			yyyymmdd := path.Join(dldir, t.Format("2006-01-02"))

			if path.Base(from) == yyyymmdd {
				log.Println("ignoring files already with", yyyymmdd)
				continue
			}

			if strings.HasPrefix(path.Base(from), ".com.google.Chrome") {
				log.Println("ignoring .com.google.Chrome file", from)
				continue
			}

			if strings.HasSuffix(path.Base(from), ".crdownload") {
				log.Println("ignoring *.crdownload file", from)
				continue
			}

			if err := os.Mkdir(yyyymmdd, 0777); err != nil {
				log.Println(err)
			}

			to := path.Join(yyyymmdd, path.Base(from))
			log.Println("Symlinking", from, "to", to)

			if err := os.Symlink(from, to); err != nil {
				log.Println(err)
			}

		default:
			log.Println(ei)
		}
	}

}
