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

func main() {
	if len(os.Args) <= 1 {
		log.Printf("Please supply one or more files to checksum\n")
		os.Exit(1)
	}
	for i, p := range os.Args {
		if i == 0 {
			continue
		}
		sum, err := HashFile(p)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%s\t%s", hex.EncodeToString(sum)[:15], p)
	}

}
