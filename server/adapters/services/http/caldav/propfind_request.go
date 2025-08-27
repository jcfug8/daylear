package caldav

import (
	"encoding/xml"
	"errors"
	"io"
)

type PropFindRequestType string

const (
	PropFindRequestTypeProp     PropFindRequestType = "prop"
	PropFindRequestTypeAllProp  PropFindRequestType = "allprop"
	PropFindRequestTypePropName PropFindRequestType = "propname"
)

type PropFindRequest struct {
	XMLName  xml.Name  `xml:"propfind"`
	Prop     *Prop     `xml:"prop,omitempty"`
	AllProp  *struct{} `xml:"allprop,omitempty"`
	PropName *struct{} `xml:"propname,omitempty"`
}

type Prop struct {
	XMLName xml.Name      `xml:"prop"`
	Raw     []RawXMLValue `xml:",any"`
}

type RawXMLValue struct {
	XMLName  xml.Name
	Attrs    []xml.Attr    `xml:",any,attr"`
	Content  string        `xml:",chardata"`
	Children []RawXMLValue `xml:",any"`
}

func NewPropFindRequestFromReader(reader io.Reader) (PropFindRequest, error) {
	content, err := io.ReadAll(reader)
	if err != nil {
		return PropFindRequest{}, err
	}

	return NewPropFindRequestFromBytes(content)
}

func NewPropFindRequestFromBytes(bytes []byte) (PropFindRequest, error) {
	var propFindRequest PropFindRequest
	err := xml.Unmarshal(bytes, &propFindRequest)
	if err != nil {
		return propFindRequest, err
	}

	if err := propFindRequest.Validate(); err != nil {
		return propFindRequest, err
	}

	return propFindRequest, nil
}

func (p PropFindRequest) Validate() error {
	count := 0
	if p.Prop != nil {
		count++
	}
	if p.AllProp != nil {
		count++
	}
	if p.PropName != nil {
		count++
	}

	if count == 0 {
		return errors.New("PROPFIND request must specify one of: prop, allprop, or propname")
	}

	if count > 1 {
		return errors.New("PROPFIND request can only specify one of: prop, allprop, or propname")
	}

	return nil
}

// GetRequestType returns what type of request this is
func (p PropFindRequest) GetRequestType() PropFindRequestType {
	if p.Prop != nil {
		return PropFindRequestTypeProp
	}
	if p.AllProp != nil {
		return PropFindRequestTypeAllProp
	}
	if p.PropName != nil {
		return PropFindRequestTypePropName
	}
	return ""
}

// // GetRequestedProperties extracts property names from the request
// func (p Prop) GetRequestedProperties() []string {
// 	var properties []string
// 	for _, raw := range p.Raw {
// 		// Convert XML name to string representation
// 		propName := raw.XMLName.Local
// 		if raw.XMLName.Space != "" {
// 			propName = raw.XMLName.Space + ":" + raw.XMLName.Local
// 		}
// 		properties = append(properties, propName)
// 	}
// 	return properties
// }

// // HasProperty checks if a specific property was requested
// func (p Prop) HasProperty(namespace, local string) bool {
// 	for _, raw := range p.Raw {
// 		if raw.XMLName.Space == namespace && raw.XMLName.Local == local {
// 			return true
// 		}
// 	}
// 	return false
// }

// // HasPropertyByXMLName checks if a property was requested by its XML name
// func (p Prop) HasPropertyByXMLName(xmlName xml.Name) bool {
// 	for _, raw := range p.Raw {
// 		if raw.XMLName == xmlName {
// 			return true
// 		}
// 	}
// 	return false
// }

// // HasPropertyByString checks if a property was requested by its string representation
// func (p Prop) HasPropertyByString(propString string) bool {
// 	for _, raw := range p.Raw {
// 		var currentProp string
// 		if raw.XMLName.Space != "" {
// 			currentProp = raw.XMLName.Space + ":" + raw.XMLName.Local
// 		} else {
// 			currentProp = raw.XMLName.Local
// 		}

// 		if currentProp == propString {
// 			return true
// 		}
// 	}
// 	return false
// }
