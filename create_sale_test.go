package test_prj

import (
	"fmt"
	"net/http"
	"test_prj/client"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SalesService_CreateSale_HappyPath(t *testing.T) {
	type testCase struct {
		name             string // имя тест-кейса
		saleName         string
		timezone         string
		schedule         []client.ScheduleItem
		targetActiveTime string
		saleIsActive     bool
	}

	tests := []testCase{
		{
			//внутри интервала
			name:     "Два дня, зона: Берлин, акция: активна, время: внутри интервала дня 1",
			saleName: "Two days, same time, Berlin, active, middle",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-22T11:30:00+02:00",
			saleIsActive:     true,
		},
		{
			//внутри интервала
			name:     "Два дня, зона: Берлин, акция: активна, время: внутри интервала дня 2",
			saleName: "Two days, same time, Berlin, active, middle",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-24T11:30:00+02:00",
			saleIsActive:     true,
		},
		{
			//точная левая граница
			name:     "Два дня, зона: Берлин, акция: активна, время: точная левая граница дня 1",
			saleName: "Two days, same time, Berlin, active, left border",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-22T09:00:00+02:00",
			saleIsActive:     true,
		},
		{
			//точная левая граница
			name:     "Два дня, зона: Берлин, акция: активна, время: точная левая граница дня 2",
			saleName: "Two days, same time, Berlin, active, left border",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-24T10:00:00+02:00",
			saleIsActive:     true,
		},
		{
			//точная правая граница
			name:     "Два дня, зона: Берлин, акция: неактивна, время: точная правая граница дня 1",
			saleName: "Two days, same time, Berlin, inactive, right border",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-22T18:00:00+02:00",
			saleIsActive:     false,
		},
		{
			//точная правая граница
			name:     "Два дня, зона: Берлин, акция: неактивна, время: точная правая граница дня 2",
			saleName: "Two days, same time, Berlin, inactive, right border",
			timezone: "Europe/Berlin",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
				{Day: "WEDNESDAY", From: "10:00", To: "19:00"},
			},
			targetActiveTime: "2026-06-24T19:00:00+02:00",
			saleIsActive:     false,
		},
		{
			//внутри интервала
			name:     "Один день, зона: Москва, акция: активна, время: внутри интервала",
			saleName: "One day Moscow active",
			timezone: "Europe/Moscow",
			schedule: []client.ScheduleItem{
				{Day: "TUESDAY", From: "10:00", To: "15:00"},
			},
			targetActiveTime: "2026-06-23T11:00:00+02:00",
			saleIsActive:     true,
		},
		{
			//вне интервала
			name:     "Один день, зона: Москва, акция: неактивна, время: другой день",
			saleName: "One day Moscow inactive",
			timezone: "Europe/Moscow",
			schedule: []client.ScheduleItem{
				{Day: "TUESDAY", From: "10:00", To: "15:00"},
			},
			targetActiveTime: "2026-06-27T01:00:00+02:00",
			saleIsActive:     false,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: активна, время: внутри интервала 1",
			saleName: "One day Tokyo, active, combined interval 1",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T13:00:00+02:00",
			saleIsActive:     true,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: активна, время: внутри интервала 2",
			saleName: "One day Tokyo, active, combined interval 2",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T19:20:00+02:00",
			saleIsActive:     true,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: неактивна, время: между интервалами",
			saleName: "One day Tokyo, inactive, combined interval 3",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T14:01:00+02:00",
			saleIsActive:     false,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: неактивна, время: между интервалами",
			saleName: "One day Tokyo, inactive, combined interval 4",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T17:59:00+02:00",
			saleIsActive:     false,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: неактивна, время: до интервала 1",
			saleName: "One day Tokyo, inactive, combined interval 5",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T11:59:00+02:00",
			saleIsActive:     false,
		},
		{
			//сложный интервал
			name:     "Один день, сложный интервал, зона: Токио, акция: неактивна, время: после интервала 2",
			saleName: "One day Tokyo, inactive, combined interval 6",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "FRIDAY", From: "12:00", To: "14:00"},
				{Day: "FRIDAY", From: "18:00", To: "21:00"},
			},
			targetActiveTime: "2026-06-26T21:01:00+02:00",
			saleIsActive:     false,
		},
		{
			// интервал через полночь
			name:     "Один день, интервал через полночь, зона: Токио, акция: активна, время: внутри интервала до полуночи",
			saleName: "One day Tokyo, active, extended interval",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "SATURDAY", From: "12:00", To: "02:00"},
			},
			targetActiveTime: "2026-06-27T23:59:00+02:00",
			saleIsActive:     true,
		},
		{
			// интервал через полночь, здесь акция по расписанию начинается в 12 часов дня субботы, а запрос идёт на 00 00 часов субботы, когда суббота только наступила
			// в это время акция ещё должна быть неактивна
			name:     "Один день, интервал через полночь, зона: Токио, акция: НЕактивна, время: полночь до начала акции",
			saleName: "One day Tokyo, inactive, extended interval 1",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "SATURDAY", From: "12:00", To: "02:00"},
			},
			targetActiveTime: "2026-06-27T00:00:00+02:00",
			saleIsActive:     false,
		},
		{
			// интервал через полночь, здесь по расписанию мы захватываем полночь следующего дня и она должна войти в акцию по ТЗ:
			// "Интервал может переходить через полночь. Например, FRIDAY 22:00–02:00 означает:
			//в пятницу с 22:00 до 23:59 продажа активна;
			//в субботу с 00:00 до 01:59 продажа активна;
			//в субботу в 02:00 продажа уже не активна."
			name:     "Один день, интервал через полночь, зона: Токио, акция: активна, время: полночь",
			saleName: "One day Tokyo, active, extended interval 2",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "SATURDAY", From: "12:00", To: "02:00"},
			},
			targetActiveTime: "2026-06-28T00:00:00+02:00",
			saleIsActive:     true,
		},
		{
			// интервал через полночь
			name:     "Один день, интервал через полночь, зона: Токио, акция: активна, время: за полночь",
			saleName: "One day Tokyo, active, extended interval 3",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "SATURDAY", From: "12:00", To: "02:00"},
			},
			targetActiveTime: "2026-06-28T00:01:00+02:00",
			saleIsActive:     true,
		},
		{
			// интервал через полночь. Правая граница интервала не включается по ТЗ
			name:     "Один день, интервал через полночь, зона: Токио, акция: неактивна, время: правая граница акции",
			saleName: "One day Tokyo, active, extended interval 4",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "SATURDAY", From: "12:00", To: "02:00"},
			},
			targetActiveTime: "2026-06-28T02:00:00+02:00",
			saleIsActive:     false,
		},
		{
			// пустое расписание. По ТЗ: "Пустой schedule означает, что продажа никогда не активна."
			name:     "Пустое расписание, акция всегда неактивна",
			saleName: "Empty schedule, Moscow, inactive",
			timezone: "Europe/Moscow",
			schedule: []client.ScheduleItem{
				{},
			},
			targetActiveTime: "2026-01-01T00:00:00+02:00",
			saleIsActive:     false,
		},
		{
			// проверка работы с таймзонами
			//В UTC - воскресенье, а в таймзоне продажи - понедельник (смещение зоны в +)
			name:     "Вокресенье в UTC, зона: Токио, акция: активна, время: внутри интервала",
			saleName: "UTC plus",
			timezone: "Asia/Tokyo",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
			},
			targetActiveTime: "2026-06-22T09:30:00+09:00", //9:30 утра пн в Токио, в UTC это воскресенье, 2026-06-21T23:30:00Z
			saleIsActive:     true,
		},
		{
			// проверка работы с таймзонами
			//В UTC - воскресенье, а в таймзоне продажи - суббота (смещение зоны в -)
			name:     "Вокресенье в UTC, зона: Америка, акция: активна, время: внутри интервала",
			saleName: "UTC minus",
			timezone: "America/New_York",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
			},
			targetActiveTime: "2026-06-22T17:30:00-04:00", //17:30 в Нью_йорке, в UTC это вторник, 2026-06-23T21:30:00Z
			saleIsActive:     true,
		},
		{
			// проверка работы с таймзонами
			//смещение на нецелое количество часов
			name:     "Смещение на нецелое количество часов, зона: Индия, акция: активна, время: внутри интервала",
			saleName: "UTC with half an hour",
			timezone: "Asia/Kolkata",
			schedule: []client.ScheduleItem{
				{Day: "MONDAY", From: "09:00", To: "18:00"},
			},
			targetActiveTime: "2026-06-22T10:30:00+05:30", //10:30 в Индии, в UTC это понедельник, 5 утра, 2026-06-22T05:00:00Z
			saleIsActive:     true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Имя теста: %s", t.Name())
			t.Logf("Имя продажи: %s, таймзона: %s, расписание: %s", tc.saleName, tc.timezone, tc.schedule)

			var createdID string

			t.Run("Создание продажи", func(t *testing.T) {
				reqCreateSale := client.CreateSaleRequest{
					Name:     tc.saleName,
					Timezone: tc.timezone,
					Schedule: tc.schedule,
				}
				respCreateSale, statusCreateSale, errCreateSale := salesServiceClient.CreateSaleAndParse(reqCreateSale)
				if errCreateSale != nil {
					t.Fatalf("Ошибка: %v", errCreateSale)
				}
				assert.Equal(t, http.StatusCreated, statusCreateSale, fmt.Sprintf("Ожидаем код: %d, фактический код: %d", http.StatusCreated, statusCreateSale))
				assert.NotEmpty(t, respCreateSale.ID, "ID продажи должен быть непустым")

				createdID = respCreateSale.ID
			})

			t.Run("Проверка статуса продажи", func(t *testing.T) {
				respActiveCheck, statusActiveCheck, errActiveCheck := salesServiceClient.IsSaleActiveAndParse(createdID, tc.targetActiveTime)
				if errActiveCheck != nil {
					t.Fatalf("Ошибка при проверке статуса продажи: %v", errActiveCheck)
				}
				assert.Equal(t, http.StatusOK, statusActiveCheck)
				assert.Equal(t, tc.saleIsActive, respActiveCheck.Active, fmt.Sprintf("Продажа должна иметь статус %t", tc.saleIsActive))
				t.Logf("Тест %s завершён успешно. Статус акции: %t", t.Name(), respActiveCheck.Active)
			})

		})
	}
}
