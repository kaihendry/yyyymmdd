package main

import (
	"log"
	"os"
	"path"
	"time"

	"github.com/rjeczalik/notify"
)

func main() {

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
		case notify.Create:
			log.Println(ei)
			from := ei.Path()

			t := time.Now().Local()
			yyyymmdd := path.Join(dldir, t.Format("2006-01-02"))

			if path.Base(from) == yyyymmdd {
				// Ignore YYYY-MM-DD files
				continue
			}

			if err := os.Mkdir(yyyymmdd, 0777); err != nil {
				log.Println(err)
			}

			to := path.Join(yyyymmdd, path.Base(from))
			log.Println("Moving", from, "to", to)

			if err := os.Rename(from, to); err != nil {
				log.Fatal(err)
			}

		default:
			log.Println(ei)
		}
	}

}
