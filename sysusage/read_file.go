package sysusage

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"unsafe"
)

var (
	empty = [][]byte{}
)

func init() {
	for i := 0; i < 10; i++ {
		empty = append(empty, []byte{'0'})
	}
}

// HeadLineSplitOfFile 读取文件的第一行
func HeadLineSplitOfFile(filename string) ([][]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	line, isPrefix, err := buf.ReadLine()
	f.Close()
	if isPrefix || (err != nil && err != io.EOF) {
		return nil, err
	}
	splitLine := bytes.Fields(line)
	return splitLine, nil
}

// byteToUint64
func byteToUint(b []byte) uint64 {
	u64, err := strconv.ParseUint(*(*string)(unsafe.Pointer(&b)), 10, 64)
	if err != nil {
		return 0
	}
	return u64
}
