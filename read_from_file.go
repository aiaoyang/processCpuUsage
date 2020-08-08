package main

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
)

func spliLine(line []byte, sep []byte) {}

// HeadLineSplitOfFile 读取文件的第一行
func HeadLineSplitOfFile(filename string) [][]byte {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	buf := bufio.NewReader(f)
	line, _, err := buf.ReadLine()
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	splitLine := bytes.Fields(line)
	return splitLine
}
