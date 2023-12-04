---
layout: 'sumologic'
page_title: 'SumoLogic: sumologic_muting_schedule'
description: |-
  Provides the ability to create, read, delete, and update muting schedule.
---

# sumologic_muting_schedule

Provides the ability to create, read, delete, and update [MutingSchedule][1].

## Example One-time Muting Schedule From 12:00 AM To 1:00 AM On 2023-08-05 For All monitor

```hcl
resource "sumologic_muting_schedule" "muting_schedule" {
  name = "Muting Schedule For one time"
  description = "This is an example for one time Muting schedule for all monitor"
  type = "MutingSchedulesLibraryMutingSchedule"
  content_type = "MutingSchedule"
  monitor {
	all = true
  }
  schedule  {
	timezone = "America/Los_Angeles"
	start_date = "2023-08-05"
	start_time = "00:00"
	duration = 60
  }
}
```

## Example One-time Muting Schedule From 12:00 AM To 1:00 AM On 2023-08-05 For Specifc Monitor/Folder ids

```hcl
resource "sumologic_muting_schedule" "muting_schedule" {
  name = "Muting Schedule For one time"
  description = "This is an example for one time Muting schedule for all monitor"
  type = "MutingSchedulesLibraryMutingSchedule"
  content_type = "MutingSchedule"
  monitor {
	ids = ["0000000000200B92"]
  }
  schedule  {
	timezone = "America/Los_Angeles"
	start_date = "2023-08-05"
	start_time = "00:00"
	duration = 60
  }
}
```

## Example Daily Muting Schedule From 9:00 AM to 9:30 and 10:00 AM to 10:30 AM Since 2023-08-05 For All monitor

```hcl
resource "sumologic_muting_schedule" "muting_schedule" {
  name = "Muting Schedule For one time"
  description = "This is an example for one time Muting schedule for all monitor"
  type = "MutingSchedulesLibraryMutingSchedule"
  content_type = "MutingSchedule"
  monitor {
	all = true
  }
  schedule  {
	timezone = "America/Los_Angeles"
	start_date = "2023-08-05"
	start_time = "00:00"
	duration = 30
	rrule = "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10"
  }
}
```

## Example Daily Muting Schedule From 9:00 AM to 9:30 and 10:00 AM to 10:30 AM Since 2023-08-05 For Specifc Monitor/Folder ids 

```hcl
resource "sumologic_muting_schedule" "muting_schedule" {
  name = "Muting Schedule For one time"
  description = "This is an example for one time Muting schedule for all monitor"
	type = "MutingSchedulesLibraryMutingSchedule"
	content_type = "MutingSchedule"
	monitor {
		ids = ["0000000000200B92"]
	  }
	schedule  {
		timezone = "America/Los_Angeles"
		start_date = "2023-08-05"
		start_time = "00:00"
		duration = 30
    	rrule = "FREQ=DAILY;INTERVAL=1;BYHOUR=9,10"
	}
}
```

## Argument reference

The following arguments are supported:

- `type` - (Optional) The type of object model. Valid value:
  - `MutingSchedulesLibraryMutingSchedule`
- `name` - (Required) The name of the muting schedule. The name must be alphanumeric.
- `description` - (Optional) The description of the muting schedule.
- `content_type` - (Optional) The type of the content object. Valid value:
  - `MutingSchedule`
- `monitor` - (Optional) The monitors which need to put in the muting schedule. see `monitor_scope_type`:
- `schedule` - (Required) The schedule information. see `schedule_type`.
- `notification_groups` -(Optinal) The muting schedule group supporting key and values. see `notification_group_type`

#### schedule_type
  - `timezone` - (Required) Time zone for the schedule per
            [IANA Time Zone Database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones#List).
  - `start_date` - (Required) Schedule start date in the format of `yyyy-mm-dd`
  - `start_time` - (Required) Schedule start time in the format of `hh:mm`
  - `duration` - (Required) Duration of the muting in minutes
  - `rrule` - (Optional) RRule (Recurrence Rule) Below are some examples of how to represent recurring events using the RRULE format:
  A rule occurring on the third Sunday of April would be as follows: `FREQ=YEARLY;BYMONTH=4;BYDAY=SU;BYSETPOS=3`
  An event occurring on the first and second Monday of October would be specified by the rule: `FREQ=YEARLY;BYMONTH=10;BYDAY=MO;BYSETPOS=1,2`
  Event that repeats monthly: every 29th of every other month! `FREQ=MONTHLY;INTERVAL=2;BYMONTHDAY=29`
  (https://freetools.textmagic.com/rrule-generator)

#### monitor_scope_type
  - `ids` - (Optional) List of monitor Ids in hex. Must be empty if `all` is true.
  - `all` - (Optional) True if the schedule applies to all monitors

#### notification_group_type
  - `group_key` - (Required) the monitor notification group key .
  - `group_values` - (Required) List of monitor notification group values.

[1]: https://help.sumologic.com/docs/alerts/monitors/muting-schedules/