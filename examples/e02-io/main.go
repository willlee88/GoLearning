package main

import (
	"io"
	"os"
	"strings"
)

func main() {
	_, _ = io.Copy(os.Stdout, strings.NewReader("hello io.Reader\n"))
}
