package caldav

import (
	"encoding/xml"
	"errors"
	"io"
)

type ReportRequestType string

const (
	ReportRequestTypeCalendarQuery    ReportRequestType = "calendar-query"
	ReportRequestTypeCalendarMultiget ReportRequestType = "calendar-multiget"
	ReportRequestTypeSyncCollection   ReportRequestType = "sync-collection"
)

type ReportRequest struct {
	CalendarQuery    *CalendarQueryReport
	CalendarMultiget *CalendarMultigetReport
	SyncCollection   *SyncCollectionReport
}

type CalendarQueryReport struct {
	XMLName xml.Name      `xml:"calendar-query"`
	Prop    *Prop         `xml:"prop,omitempty"`
	Filter  *Filter       `xml:"filter,omitempty"`
	Raw     []RawXMLValue `xml:",any"`
}

type CalendarMultigetReport struct {
	XMLName xml.Name      `xml:"calendar-multiget"`
	Prop    *Prop         `xml:"prop,omitempty"`
	Hrefs   []string      `xml:"href,omitempty"`
	Raw     []RawXMLValue `xml:",any"`
}

type SyncCollectionReport struct {
	XMLName   xml.Name      `xml:"sync-collection"`
	Prop      *Prop         `xml:"prop,omitempty"`
	SyncToken *string       `xml:"sync-token,omitempty"`
	SyncLevel *string       `xml:"sync-level,omitempty"`
	Raw       []RawXMLValue `xml:",any"`
}

type Filter struct {
	XMLName    xml.Name      `xml:"filter"`
	CompFilter *CompFilter   `xml:"comp-filter,omitempty"`
	Raw        []RawXMLValue `xml:",any"`
}

type CompFilter struct {
	XMLName    xml.Name      `xml:"comp-filter"`
	Name       string        `xml:"name,attr,omitempty"`
	CompFilter *CompFilter   `xml:"comp-filter,omitempty"`
	TimeRange  *TimeRange    `xml:"time-range,omitempty"`
	Raw        []RawXMLValue `xml:",any"`
}

type TimeRange struct {
	XMLName xml.Name      `xml:"time-range"`
	Start   string        `xml:"start,attr,omitempty"`
	End     string        `xml:"end,attr,omitempty"`
	Raw     []RawXMLValue `xml:",any"`
}

// NewReportRequestFromReader parses a REPORT request by first detecting the root element,
// then parsing into the appropriate struct
func NewReportRequestFromReader(reader io.Reader) (ReportRequest, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return ReportRequest{}, err
	}

	return NewReportRequestFromBytes(content)
}

func NewReportRequestFromBytes(bytes []byte) (ReportRequest, error) {
	root := struct {
		XMLName xml.Name
	}{}
	err := xml.Unmarshal(bytes, &root)
	if err != nil {
		return ReportRequest{}, err
	}

	reportRequest := ReportRequest{}

	switch root.XMLName.Local {
	case string(ReportRequestTypeCalendarQuery):
		err = xml.Unmarshal(bytes, &reportRequest.CalendarQuery)
		if err != nil {
			return ReportRequest{}, err
		}
	case string(ReportRequestTypeCalendarMultiget):
		err = xml.Unmarshal(bytes, &reportRequest.CalendarMultiget)
		if err != nil {
			return ReportRequest{}, err
		}
	case string(ReportRequestTypeSyncCollection):
		err = xml.Unmarshal(bytes, &reportRequest.SyncCollection)
		if err != nil {
			return ReportRequest{}, err
		}
	default:
		return ReportRequest{}, errors.New("unknown report type: " + root.XMLName.Local)
	}

	return reportRequest, nil
}

// GetRequestType returns what type of report request this is
func (r ReportRequest) GetRequestType() ReportRequestType {
	if r.CalendarQuery != nil {
		return ReportRequestTypeCalendarQuery
	}
	if r.CalendarMultiget != nil {
		return ReportRequestTypeCalendarMultiget
	}
	if r.SyncCollection != nil {
		return ReportRequestTypeSyncCollection
	}
	return ""
}
