package fixtures

import "test_prj/client"

type ScheduleMap map[string][2]string

func BuildSaleRequest(name string, timezone string, scheduleData ScheduleMap) client.CreateSaleRequest {
	var schedule []client.ScheduleItem

	for day, hours := range scheduleData {
		item := client.ScheduleItem{
			Day:  day,
			From: hours[0],
			To:   hours[1],
		}
		schedule = append(schedule, item)
	}

	return client.CreateSaleRequest{
		Name:     name,
		Timezone: timezone,
		Schedule: schedule,
	}
}
