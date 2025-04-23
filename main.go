package main

import (
	"fmt"
	"fusefs/fsys"
	"log"
	"os"
	"os/signal"
	"syscall"

	"bazil.org/fuse"
	"bazil.org/fuse/fs"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./memtreefs <mountpoint>")
		os.Exit(1)
	}

	mountpoint := os.Args[1]
	root := fsys.NewDir("/")
	fuseFS := fsys.NewFS(root)

	c, err := fuse.Mount(mountpoint)
	if err != nil {
		log.Fatal(err)
	}

	// Signal handling
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sig
		fmt.Println("\nUnmounting and exiting...")
		fuse.Unmount(mountpoint)
		c.Close()
		os.Exit(0)
	}()

	if err := fs.Serve(c, fuseFS); err != nil {
		log.Fatal(err)
	}
}
