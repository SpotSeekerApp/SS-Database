package utils

import (
	"SS-Database/lib/types"
	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"net/http"
	"reflect"
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

func MatrixMul(mat1 [][]float64, mat2 [][]float64) [][]float64 {
	result := make([][]float64, len(mat1))
	for i := range result {
		result[i] = make([]float64, len(mat2[0]))
		for j := range result[i] {
			for k := range mat2 {
				result[i][j] += mat1[i][k] * mat2[k][j]
			}
		}
	}
	return result
}

func MatrixTranspose(mat [][]float64) [][]float64 {

	transposed := make([][]float64, len(mat[0]))
	for i := range transposed {
		transposed[i] = make([]float64, len(mat))
		for j := range transposed[i] {
			transposed[i][j] = mat[j][i]
		}
	}

	return transposed
}

func MapLists(key []float64, val []string) map[float64]string {
	res := make(map[float64]string)
	for idx, idxVal := range key {
		res[idxVal] = val[idx]
	}
	return res
}

func MatrixSum(mat1 [][]float64, mat2 [][]float64) [][]float64 {
	for rowIdx, row := range mat1 {
		for colIdx, _ := range row {
			mat1[rowIdx][colIdx] += mat2[rowIdx][colIdx]
		}
	}
	return mat1
}

func MatrixConstMul(mat [][]float64, num float64) [][]float64 {
	for rowIdx, row := range mat {
		for colIdx, _ := range row {
			mat[rowIdx][colIdx] *= num
		}
	}
	return mat
}

func MapToMatrix(strMap map[string]map[string]float64) ([][]float64, []string) {
	var matrixArray [][]float64
	var matrixKeys []string
	for parentKey, childMap := range strMap {
		matrixKeys = append(matrixKeys, parentKey)
		var matrixRow []float64
		for _, tag := range types.TAGLIST {
			matrixRow = append(matrixRow, childMap[tag])
		}
		matrixArray = append(matrixArray, matrixRow)
	}
	return matrixArray, matrixKeys
}

func ListToMatrix(strMap map[string]float64) [][]float64 {
	var matrixArray [][]float64
	for _, tag := range types.TAGLIST {
		matrixArray = append(matrixArray, []float64{strMap[tag]})
	}
	return matrixArray
}

func StringToList(strIn string) []string {
	var res []string

	strIn = strings.TrimLeft(strIn, "[")
	strIn = strings.TrimRight(strIn, "]")
	strList := strings.Split(strIn, ",")

	for _, str := range strList {
		str = strings.TrimSpace(str)
		str = strings.TrimLeft(str, "\"")
		str = strings.TrimRight(str, "\"")
		res = append(res, str)
	}
	return res
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
		_, isInt := val.(int)
		_, isString := val.(string)
		valMap, isMap := val.(map[string]interface{})
		if (isString && val != "") || (isMap && len(valMap) != 0) || isInt {
			ret = append(ret, firestore.Update{
				Path:  strings.ToLower(string(key[0])) + key[1:],
				Value: val,
			})
		}
	}
	return ret
}
