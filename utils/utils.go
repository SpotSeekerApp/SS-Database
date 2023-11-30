package utils

import (
	"google.golang.org/grpc/codes"
	"net/http"
	"strings"
)

func MapErrorCode(rpcResponse codes.Code) http.ConnState {
	switch rpcResponse {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Aborted:
		return http.StatusInternalServerError
	case codes.OK:
		return http.StatusOK
	default:
		return http.StatusNotImplemented
	}
}
func SubstrInList(str string, list string) bool {
	subStrList := strings.Split(list, "|")
	for _, s := range subStrList {
		if strings.Contains(str, s) {
			return true
		}
	}
	return false
}
