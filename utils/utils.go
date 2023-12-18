package utils

import (
	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

const BATCH_SIZE = 5

func MapErrorCode(rpcResponse codes.Code) http.ConnState {
	switch rpcResponse {
	case codes.NotFound:
		return http.StatusNotFound
	case codes.Aborted:
		return http.StatusInternalServerError
	case codes.OK:
		return http.StatusOK
	case codes.AlreadyExists:
		return http.StatusNotAcceptable
	default:
		return http.StatusNotImplemented
	}
}

func FindDigitCount(num int) int {
	count := 0
	for {
		num /= 10
		count += 1
		if num == 0 {
			break
		}
	}
	return count
}

func NumberToString(num int, maxDigit any) string {
	maxDigitCount := FindDigitCount(int(maxDigit.(int64)))
	digitCount := FindDigitCount(num)
	var zeros string
	for i := 0; i < maxDigitCount-digitCount; i++ {
		zeros += "0"
	}
	return zeros + strconv.Itoa(num)
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

func ExtractNonEmptyFields(s any) []firestore.Update {

	structValues := reflect.ValueOf(s)

	var ret []firestore.Update

	for i := 0; i < structValues.NumField(); i++ {
		key := reflect.TypeOf(s).Field(i).Name
		if strings.Contains(strings.ToLower(key), "id") {
			continue
		}
		val := structValues.Field(i).Interface()
		_, isString := val.(string)
		valMap, isMap := val.(map[string]interface{})
		if (isString && val != "") || (isMap && len(valMap) != 0) {
			ret = append(ret, firestore.Update{
				Path:  strings.ToLower(string(key[0])) + key[1:],
				Value: val,
			})
		}
	}
	return ret
}
