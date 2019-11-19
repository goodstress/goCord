// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    superProp, err := UnmarshalSuperProp(bytes)
//    bytes, err = superProp.Marshal()

package main

import "encoding/json"

func UnmarshalSuperProp(data []byte) (SuperProp, error) {
	var r SuperProp
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SuperProp) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type SuperProp struct {
	OS                     string      `json:"os"`
	Browser                string      `json:"browser"`
	Device                 string      `json:"device"`
	BrowserUserAgent       string      `json:"browser_user_agent"`
	BrowserVersion         string      `json:"browser_version"`
	OSVersion              string      `json:"os_version"`
	Referrer               string      `json:"referrer"`
	ReferringDomain        string      `json:"referring_domain"`
	ReferrerCurrent        string      `json:"referrer_current"`
	ReferringDomainCurrent string      `json:"referring_domain_current"`
	ReleaseChannel         string      `json:"release_channel"`
	ClientBuildNumber      int64       `json:"client_build_number"`
	ClientEventSource      interface{} `json:"client_event_source"`
}
