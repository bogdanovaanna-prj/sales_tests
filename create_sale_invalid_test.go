package test_prj

import (
	"fmt"
	"io"
	"net/http"
	"test_prj/fixtures"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SalesService_CreateSale_InvalidReq(t *testing.T) { // по БТ здесь ждём код 400
	type testCase struct {
		name           string
		saleName       string
		timezone       string
		schedule       fixtures.ScheduleMap
		expectedStatus int
	}
	tests := []testCase{
		{
			name:     "Ошибка: Пустое имя продажи",
			saleName: "",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"MONDAY": {"09:00", "18:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: Пустой timezone",
			saleName: "Empty timezone",
			timezone: "",
			schedule: fixtures.ScheduleMap{
				"MONDAY": {"09:00", "18:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: Невалидный день недели",
			saleName: "Impossible day",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"AWESOME-DAY": {"09:00", "18:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: Опечатка в названии дня недели",
			saleName: "Misspelled day",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"WENSDAY": {"09:00", "18:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: цифра вместо дня недели",
			saleName: "Number as day",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"5": {"09:00", "18:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: время в расписании больше 24",
			saleName: "Time more than 24 hours",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"THURSDAY": {"09:00", "25:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: время в расписании больше 24",
			saleName: "Minutes more than 59 minutes",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"THURSDAY": {"09:65", "15:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "Ошибка: валидное и невалидное время в расписании",
			saleName: "Both valid and invalid schedule",
			timezone: "Europe/Berlin",
			schedule: fixtures.ScheduleMap{
				"THURSDAY": {"09:00", "10:00"},
				"THURSDA":  {"10:00", "11:00"},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Тест # %d", i+1)
			t.Logf("Имя теста: \"%s\"", tc.name)
			t.Logf("Имя продажи: \"%s\"", tc.saleName)
			t.Logf("Зона: %s, расписание: %s", tc.timezone, tc.schedule)

			t.Run("Создание продажи", func(t *testing.T) {
				reqBody := fixtures.BuildSaleRequest(tc.saleName, tc.timezone, tc.schedule)

				resp, err := salesServiceClient.CreateSale(reqBody)
				if err != nil {
					t.Fatalf("Ошибка: %v", err)
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						t.Errorf("Ошибка при закрытии: %s", err.Error())
					}
				}(resp.Body)

				assert.Equal(t, tc.expectedStatus, resp.StatusCode, fmt.Sprintf("Ожидаем код: %d, фактический код: %d", tc.expectedStatus, resp.StatusCode))
				t.Logf("Тест # %d завершён. Статус ответа: %s\n\n", i+1, resp.Status)
			})
		})
	}
}
