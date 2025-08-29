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

type ResponseHref struct {
	Href string `xml:"D:href"`
}

type PrivilegeSet struct {
	Privileges []Privilege `xml:"D:privilege"`
}

type Privilege struct {
	Name string `xml:"D:privilege"`
}

func setCalDAVHeaders(w http.ResponseWriter) {
	w.Header().Set("DAV", "1, 2, calendar-access")
	w.Header().Set("CalDAV", "calendar-access")
}

func addXMLDeclaration(response []byte) []byte {
	// declaration := []byte(`<?xml version="1.0" encoding="utf-8"?>`)
	// return append(declaration, response...)
	return response
}
