package sonus

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// Errors represents the XML structure returned by the SBC when an
// error is encountered
type Errors struct {
	XMLName   xml.Name `xml:"errors"`
	Text      string   `xml:",chardata"`
	Xmlns     string   `xml:"xmlns,attr"`
	ErrorMsgs []struct {
		Text         string `xml:",chardata"`
		ErrorTag     string `xml:"error-tag"`
		ErrorUrlpath string `xml:"error-urlpath"`
		ErrorMessage string `xml:"error-message"`
	} `xml:"error"`
}

// Error implements the error interface for the Sonus Errors.  It
// concatenates the array of errors into a string separated by '\n'.
func (e *Errors) Error() string {
	errors := make([]string, 0)
	for _, err := range e.ErrorMsgs {
		errors = append(errors, fmt.Sprintf("%s: %s", err.ErrorTag, err.ErrorMessage))
	}
	return strings.Join(errors, "\n")
}
