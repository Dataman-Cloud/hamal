package utils

import (
	"io/ioutil"
	"net/http"
	"unsafe"
)

// Byte2str array byte parse to string
func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// ReadRequestBody read request body
func ReadRequestBody(request *http.Request) ([]byte, error) {
	defer request.Body.Close()
	return ioutil.ReadAll(request.Body)
}

// ReadResponseBody read response body
func ReadResponseBody(resp *http.Response) ([]byte, error) {
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func StageInterval() {
}
