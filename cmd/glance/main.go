package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	// 0. get cli args
	flag.Parse()
	args := flag.Args()

	// 1. check <path> arg
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: glance <path>")
		os.Exit(1)
	}

	// 2. resolve <path> arg to absolute filepath
	absPath, err := filepath.Abs(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// 3. verify path exists on disk
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	// 4. print resolved path
	fmt.Fprintf(os.Stdout, "path: %s\n", absPath)
	fmt.Printf("info: %s\n", info.Name())
}


