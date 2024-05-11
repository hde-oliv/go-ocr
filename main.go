package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type Tesseract struct {
	Binary     string
	InputFile  string
	OutputFile *os.File
	Extension  string
	Language   string
}

func (t Tesseract) readOutput() []byte {
	f, err := os.Open(t.OutputFile.Name() + ".txt")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()
	defer os.Remove(t.OutputFile.Name() + ".txt")

	reader := bufio.NewReader(f)
	data, err := io.ReadAll(reader)

	if err != nil {
		log.Fatal(err)
	}

	return data
}

func (t Tesseract) runTesseract() []byte {
	var args = make([]string, 0)

	args = append(args, "-l", t.Language, t.InputFile, t.OutputFile.Name())

	cmd := exec.Command(t.Binary, args...)

	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}

	return t.readOutput()
}

func (t Tesseract) ImageToString() string {
	var err error

	t.OutputFile, err = os.CreateTemp("", "tess_")

	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(t.OutputFile.Name())

	t.InputFile, err = filepath.Abs(t.InputFile)

	if err != nil {
		log.Fatal(err)
	}

	data := t.runTesseract()

	return string(data)
}

func main() {
	var t Tesseract

	t.Binary = "tesseract.exe"

	if len(os.Args) == 2 {
		t.InputFile = os.Args[1]
	} else if len(os.Args) == 4 && os.Args[1] == "-l" {
		t.Language = os.Args[2]
		t.InputFile = os.Args[3]
	} else {
		fmt.Fprintf(os.Stderr, "Usage: go-ocr [-l lang] input_file\n")
		os.Exit(1)
	}

	text := t.ImageToString()

	print(text)
}
