package test_prj

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SalesService_HealthCheck(t *testing.T) {
	resp, err := salesServiceClient.Healthcheck()
	if err != nil {
		t.Fatalf("Не удалось подключиться к сервису: %v", err)
	}
	t.Logf("Успешно отправили запрос статуса сервиса (Healthz). Статус ответа: %s", resp.Status)

	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode, "Healthcheck должен возвращать статус 200")
	t.Logf("Тест %s завершён успешно. Статус ответа: %s", t.Name(), resp.Status)
}
