package main

import (
	"crypto"
	_ "crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
)

func HashFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		return []byte{}, err
	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("Failed to read file %s: %s", path, err)
	}
	hasher := crypto.SHA256.HashFunc().New()
	hasher.Write(bytes)
	sum := hasher.Sum(make([]byte, 0))
	return sum, nil
}

func NewFileEntry(path string) (FileEntry, error) {
	hash, err := HashFile(path)
	if err != nil {
		return FileEntry{}, err
	}
	h := hex.EncodeToString(hash)
	return FileEntry{path, h}, nil
}

type FileEntry struct {
	path string
	hash string
}

func RecursiveHash(root string, p chan FileEntry) error {
	if stat, err := os.Stat(root); err != nil {
		return err
	} else if stat.IsDir() {
		entries, err := os.ReadDir(root)
		if err != nil {
			return err
		}
		results := make(chan error)
		for _, entry := range entries {
			go func() {
				err := RecursiveHash(root+"/"+entry.Name(), p)
				results <- err
			}()
		}

		for range entries {
			if err := <-results; err != nil {
				return err
			}
		}
	} else {
		fe, err := NewFileEntry(root)
		if err != nil {
			return err
		}
		p <- fe
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Printf("Please supply one file or directory to checksum\n")
		os.Exit(1)
	}

	p := make(chan FileEntry, 5)
	root := os.Args[1]
	if root[len(root)-1] == '/' {
		root = root[:len(root)-1]
	}

	go func() {
		err := RecursiveHash(root, p)
		if err != nil {
			panic(err)
		}
		close(p)
	}()

	entries := make(map[string][]string)

	for e := range p {
		entries[e.hash] = append(entries[e.hash], e.path)
	}
	for h, p := range entries {
		if len(p) == 1 {
			continue
		}
		fmt.Println("Duplicate Files Detected")
		for _, p := range p {
			fmt.Printf("\t%s\t%s\n", h[:15], p)
		}
		fmt.Printf("\n")
	}

}
