package main

import (
	"github.com/hanleo-6/glance-go/internal/server"

	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	// Parse CLI arguments
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: glance <path>")
		os.Exit(1)
	}

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// verify path exists on disk
	_, err = os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	server := server.NewServer(absPath)
	fmt.Println("Listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}


