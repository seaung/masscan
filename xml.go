package masscan

import (
	"encoding/xml"
)

type MasscaRun struct {
	XMLName          xml.Name   `xml:"nmaprun"`
	Scanner          string     `xml:"scanner,attr" json:"scanner"`
	Start            string     `xml:"start,attr" json:"start"`
	Version          string     `xml:"version,attr" json:"version"`
	XmlOutputVersion string     `xml:"xmloutputversion,attr" json:"xml_output_version"`
	ScanInfo         ScanInfo   `xml:"scaninfo" json:"scaninfo"`
	Hosts            []Hosts      `xml:"host" json:"host"`
	Finished         Finished   `xml:"finished" json:"finished"`
	HostRecord       HostRecord `xml:"hosts" json:"hosts"`
	rawXML           []byte
}

type ScanInfo struct {
	Type     string `xml:"type,attr" json:"type"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
}

type Hosts struct {
	EndTime string  `xml:"endtime,attr" json:"end_time"`
	Address Address `xml:"address" json:"address"`
	Ports   []Ports   `xml:"ports>port" json:"ports"`
}

type Ports []struct {
	ID       string `xml:"portid,attr" json:"port_id"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
	State    State  `xml:"state" json:"state"`
}

type Address struct {
	AddrType string `xml:"addrtype,attr" json:"addr_type"`
	Addr     string `xml:"addr,attr" json:"addr"`
}

type State struct {
	State     string `xml:"state,attr" json:"state"`
	Reason    string `xml:"reason,attr" json:"reason"`
	ReasonTTL string `xml:"reasonttl,attr" json:"reason_ttl"`
}

type HostRecord struct {
	Total string `xml:"total,attr" json:"total"`
	Up    string `xml:"up,attr" json:"up"`
	Down  string `xml:"down,attr" json:"down"`
}

type Finished struct {
	Elapsed string `xml:"elapsed,attr" json:"elapsed"`
	Time    string `xml:"time,attr" json:"time"`
	TimeStr string `xml:"timeStr,attr" json:"time_str"`
}

func ParseXML(content []byte) (*MasscaRun, error) {
	r := &MasscaRun{
		rawXML: content,
	}
	err := xml.Unmarshal(content, r)
	return r, err
}
