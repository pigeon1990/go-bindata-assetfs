package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

const outfile = "bindata.go"

func main() {
	if _, err := exec.LookPath("go-bindata"); err != nil {
		fmt.Println("Cannot find go-bindata executable in path")
		fmt.Println("Maybe you need: go get github.com/elazarl/go-bindata-assetfs/...")
		os.Exit(1)
	}
	cmd := exec.Command("go-bindata", os.Args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}
	in, err := os.Open(outfile)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot read 'bindata.go'", err)
		return
	}
	out, err := os.Create("bindata_assetfs.go")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Cannot write 'bindata_assetfs.go'", err)
		return
	}
	r := bufio.NewReader(in)
	done := false
	for line, isPrefix, err := r.ReadLine(); err == nil; line, isPrefix, err = r.ReadLine() {
		line = append(line, '\n')
		if _, err := out.Write(line); err != nil {
			fmt.Fprintln(os.Stderr, "Cannot write to 'bindata_assetfs.go'", err)
			return
		}
		if !done && !isPrefix && bytes.HasPrefix(line, []byte("import (")) {
			fmt.Fprintln(out, "\t\"github.com/elazarl/go-bindata-assetfs\"")
			done = true
		}
	}
	fmt.Fprintln(out, `
func assetFS() *assetfs.AssetFS {
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, Prefix: k}
	}
	panic("unreachable")
}`)
	// Close files BEFORE remove calls (don't use defer).
	in.Close()
	out.Close()
	if err := os.Remove(outfile); err != nil {
		fmt.Fprintln(os.Stderr, "Cannot remove", outfile, err)
	}
}
