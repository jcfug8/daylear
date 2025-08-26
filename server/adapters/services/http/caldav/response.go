package caldav

import "net/http"

type ResourceType struct {
	Collection *Collection `xml:"D:collection,omitempty"`
	Principal  *Principal  `xml:"D:principal,omitempty"`
	Calendar   *Calendar   `xml:"C:calendar,omitempty"`
}

type Collection struct{}

type Principal struct{}

type Calendar struct{}

type Href struct {
	Href string `xml:"D:href"`
}

type Status struct {
	Status string `xml:",chardata"`
}

func setCalDAVHeaders(w http.ResponseWriter) {
	w.Header().Set("DAV", "1, 2, calendar-access, calendar-schedule")
	w.Header().Set("CalDAV", "calendar-access, calendar-schedule")
}

func addXMLDeclaration(response []byte) []byte {
	// declaration := []byte(`<?xml version="1.0" encoding="utf-8"?>`)
	// return append(declaration, response...)
	return response
}
