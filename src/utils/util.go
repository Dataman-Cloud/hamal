package utils

import (
	"io/ioutil"
	"net/http"
	"reflect"
	"unsafe"
)

// Byte2str array byte parse to string
func Byte2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func StringToByte(v string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&v))
	return *(*[]byte)(unsafe.Pointer(sh))
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
