package main

import (
	"bytes"
	"compress/gzip"
	"regexp"
)

//Sanitize a string before use in a filename
func FileNameFilter(str string) string {
	alphafilter, _ := regexp.Compile("[^a-zA-Z0-9-_]+")
	str = alphafilter.ReplaceAllString(str, "")
	return str
}

//Gzip compress a byte array
func compressGzip(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	w.Write(data)
	w.Close()
	return b.Bytes()
}
