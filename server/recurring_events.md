# Update recurring
## This event
- add the recurring event time it overrides
- add the exclude to the parent
- create new event data row for this event
## All future events
- add the stop to the parent's rrule
- delete unneeded clones and any event data rows for clones that don't just use the parent's event data
- create a new event with the new rrule and create new clones

# Delete recurring
## This event
- add an exclude to the parent
- delete the event and any event data rows for this event
## All future events
- add the stop to the parent's rrule
- delete unneeded clones and any event data rows for clones that don't just use the parent's event data


### IDEAS
- maybe create a event_metadata table that could assist in individual event data overrides? Not sure yet
- maybe instead of duplicating events for reccuring events, we just introduce two fields "query_start_time" and "query_end_time". For non recurring events, we just use the start and end time. For recurring events, we set the start to the beginning of the first event and the end to the end of the last event (if it exists). Then when we query a data range we'll get the single recurring event and we can just pass the to the UI or we could generate them on the fly.


## Modifying a single event
Great question! CalDAV expects to see two separate VEVENT components when you modify a single instance within an RRULE:
CalDAV Expectation for Modified Instances:
1. The Original Recurring Event (with EXDATE):
```
BEGIN:VEVENT
UID:event123@example.com
DTSTART:20240101T090000Z
DTEND:20240101T100000Z
SUMMARY:Weekly Meeting
RRULE:FREQ=WEEKLY;COUNT=10
EXDATE:20240108T090000Z
END:VEVENT
```
2. The Modified Instance (with RECURRENCE-ID):
```
BEGIN:VEVENT
UID:event123@example.com
RECURRENCE-ID:20240108T090000Z
DTSTART:20240108T140000Z
DTEND:20240108T150000Z
SUMMARY:Weekly Meeting (Rescheduled)
END:VEVENT
```
Key Points:
- Same UID: Both events share the same UID to link them together
- EXDATE: The original event excludes the modified occurrence
- RECURRENCE-ID: The modified instance specifies which occurrence it's overriding
- Separate VEVENTs: They are two distinct calendar components, not one combined event