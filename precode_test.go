package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
func TestMainHandlerWhenOk(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)

	handler.ServeHTTP(responseRecorder, req)
	response := responseRecorder.Result()

	require.NotNil(t, req, "Erorr creating request")
	require.Equal(t, 200, response.StatusCode, "Unexpected status code")
	require.NotEmpty(t, response.Body, "Response body is empty")
}

// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
func TestMainHandlerWhereIsTheWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=4&city=wrongCity", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)

	handler.ServeHTTP(responseRecorder, req)
	response := responseRecorder.Result()

	assert.Equal(t, http.StatusBadRequest, response.StatusCode, "Unexepected status code")
	body, _ := io.ReadAll(response.Body)
	require.Contains(t, string(body), "wrong city value")
}

// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4
	req := httptest.NewRequest("GET", "/cafe?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(MainHandle)

	handler.ServeHTTP(responseRecorder, req)
	response := responseRecorder.Result()

	body, _ := io.ReadAll(response.Body)
	cafeList := strings.Split(string(body), ",")
	assert.Len(t, cafeList, totalCount, "Unexpected number of cafes")
}
