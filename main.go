package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/neovim/go-client/nvim"
)

func instances(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("error fetching nvim instances: %w", err)
	}

	return matches, nil
}

func main() {
	socketPattern := flag.String(
		"socket-pattern",
		filepath.Join(os.TempDir(), "nvim*"),
		"Pattern where nvim sockets are located",
	)
	flag.Parse()

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s: [options] COMMAND\n", os.Args[0])
		flag.PrintDefaults()
	}

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	if os.Args[1] == "--help" {
		flag.Usage()
		os.Exit(0)
	}

	command := os.Args[1]

	list, err := instances(*socketPattern)
	if err != nil {
		log.Fatalln(err)
	}

	var addr string
	for _, i := range list {
		addr = filepath.Join(i, "0")

		v, err := nvim.Dial(addr)
		if err != nil {
			log.Fatal(err)
		}

		defer v.Close()

		b := v.NewBatch()
		b.Command(command)

		if err := b.Execute(); err != nil {
			log.Fatal(err)
		}
	}
}
