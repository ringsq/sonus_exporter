package sonus

import (
	"context"
	"crypto/tls"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"reflect"
	"time"

	log "github.com/ringsq/go-logger"
)

// Sonus URLs
const (
	// Gets the SBC system information
	systemInfoPath       = "/config/system"
	serverInfoPath       = "/operational/system/serverStatus/"
	contextListPath      = "/config/addressContext/"
	zoneStatusPath       = "/operational/addressContext/%s/zoneStatus/"
	ipInterfaceGroupPath = "/operational/addressContext/%s/ipInterfaceGroup/"
	sipStatsPath         = "/operational/addressContext/%s/zone/%s/sipCurrentStatistics/"
	fanStatusPath        = "/operational/system/fanStatus/"
	powerSupplyPath      = "/operational/system/powerSupplyStatus/"
	dspStatusPath        = "/operational/system/dspStatus/dspUsage/"
	tgStatusPath         = "/operational/global/globalTrunkGroupStatus/"
	tgConfigPath         = "/config/addressContext/%s/zone/%s/sipTrunkGroup/"
	callStatusPath       = "/operational/addressContext/%s/zone/%s/callCurrentStatistics/"
)

type system struct {
	Admin       admin   `json:"admin" xml:"admin"`
	ServerAdmin []admin `json:"serverAdmin" xml:"serverAdmin"`
}

type admin struct {
	Name string `json:"name" xml:"name"`
}

type AddressContexts struct {
	AddressContext []struct {
		Name     string `xml:"name"`
		DnsGroup struct {
			Name string `xml:"name"`
		} `xml:"dnsGroup"`
		IpInterfaceGroup []struct {
			Name string `xml:"name"`
		} `xml:"ipInterfaceGroup"`
		Zone []struct {
			Name string `xml:"name"`
		} `xml:"zone"`
	} `xml:"addressContext"`
}

// SBC represents a single Sonus session border controller
type SBC struct {
	target          string
	user            string
	password        string
	client          *http.Client
	System          string
	AddressContexts *AddressContexts
}

// NewSBC instantiates an SBC from the provided credentials
func NewSBC(address, user, password string) *SBC {
	ac := &AddressContexts{}
	sbc := &SBC{
		target:          address,
		user:            user,
		password:        password,
		AddressContexts: ac,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	sbc.client = &http.Client{Transport: tr}
	sys := &system{}
	ctx := context.Background()

	err := sbc.GetAndParse(ctx, sys, systemInfoPath)
	if err != nil {
		log.Errorf("Error calling SBC (%s): %v", systemInfoPath, err)
		return nil
	}
	sbc.System = sys.Admin.Name

	// Get the address contexts
	err = sbc.GetAndParse(ctx, ac, contextListPath)
	if err != nil {
		log.Errorf("Error getting context list: %v", err)
		return nil
	}
	sbc.AddressContexts = ac
	return sbc
}

// buildURL takes the given path, adds the base to the beginning, and applies any
// formatting arguments
func (s *SBC) buildURL(path string, args ...any) string {
	url := fmt.Sprintf("https://%s/api%s", s.target, path)
	return fmt.Sprintf(url, args...)
}

// GetAndParse builds the URL, does a GET against the SBC, and parses the XML response.
// Any errors are returned in error.
func (s *SBC) GetAndParse(ctx context.Context, response any, path string, args ...any) error {
	url := s.buildURL(path, args...)
	resp, err := s.callSBC(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Errorf("Error calling SBC (%s): %v", url, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Errorf("Error reading response body: %v", err)
		return err
	}
	err = xml.Unmarshal(body, response)
	if err != nil {
		log.Errorf("BODY: %v", body)
		log.Errorf("Failed to deserialize %s into %v: %v", url, reflect.TypeOf(response), err)
		return err
	}
	return nil
}

// callSBC is responsible for building the request object, sending it to the SBC, and checking the response status.
func (s *SBC) callSBC(ctx context.Context, method string, url string, body io.Reader) (*http.Response, error) {

	// startTime := time.Now()
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		log.Errorf("Error creating request to %s %s: %v", method, url, err)
		return nil, err
	}
	req.SetBasicAuth(s.user, s.password)
	var resp *http.Response
	for retries := 3; retries > 0; retries-- {
		resp, err = s.client.Do(req)
		if err != nil {
			log.Errorf("Error with SBC call %s %s: %v", method, url, err)
			return nil, err
		}

		if prob := checkResponse(resp); prob != nil {
			log.Errorf("Error response from %s %s: %v", method, url, prob)
			return nil, prob
		}
		if resp.StatusCode == 204 {
			sleepTime := rand.Intn(1000)
			log.Warnf("%s received for %s, retry %d in %dms...", resp.Status, url, 3-retries+1, sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Millisecond)
			continue
		}
		break
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Invalid response received: %s", resp.Status)
	}
	return resp, nil
}

// checkResponse inspects the response status codes and creates an appropriate problem.  If
// no issues are found `nil` is returned
func checkResponse(resp *http.Response) error {
	status := resp.StatusCode
	if status < 300 {
		return nil
	}
	if status < 400 {
		return fmt.Errorf("Redirect received from Sonus: %v", resp.Status)
	}
	if status < 500 {
		prob := errors.New(getError(resp))
		log.Warn(prob)
		return prob
	}
	if status < 600 {
		prob := errors.New(getError(resp))
		log.Error(prob)
		return prob
	}
	return nil
}

// getError reads the error message from the Sonus response
func getError(resp *http.Response) string {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return resp.Status
	}
	if len(body) == 0 {
		return resp.Status
	}
	sbcErr := Errors{}
	err = xml.Unmarshal(body, &sbcErr)
	if err != nil {
		return resp.Status
	}
	return sbcErr.Error()
}
