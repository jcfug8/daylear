package caldav

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

// Generic CalDAV/WebDAV response structures
type MultiStatusResponse struct {
	XMLName  xml.Name   `xml:"D:multistatus"`
	XMLNSD   string     `xml:"xmlns:D,attr"`
	XMLNSC   string     `xml:"xmlns:C,attr"`
	Response []Response `xml:"D:response"`
}

type Response struct {
	Href      string     `xml:"D:href"`
	Propstats []Propstat `xml:"D:propstat"`
}

type Propstat struct {
	Prop   interface{} `xml:"D:prop"`
	Status Status      `xml:"D:status"`
}

type Status struct {
	Status string `xml:"D:status"`
}

// Generic response builder that can handle any property type
type ResponseBuilder struct{}

func (rb ResponseBuilder) BuildMultiStatusResponse(responses []Response) MultiStatusResponse {
	return MultiStatusResponse{
		XMLNSD:   "DAV:",
		XMLNSC:   "urn:ietf:params:xml:ns:caldav",
		Response: responses,
	}
}

// Generic function to add property status
func (rb ResponseBuilder) AddPropertyStatus(response Response, prop interface{}, statusCode int) Response {
	statusText := fmt.Sprintf("HTTP/1.1 %d %s", statusCode, http.StatusText(statusCode))

	response.Propstats = append(response.Propstats, Propstat{
		Prop:   prop,
		Status: Status{Status: statusText},
	})

	return response
}

// Helper function to check if a prop struct has any non-zero values
func HasAnyProperties(prop interface{}) bool {
	// This is a simplified check - you might want to implement more sophisticated
	// property checking based on your specific needs
	return prop != nil
}
