package masscan

import (
	"encoding/xml"
	"strconv"
	"time"
)

type MasscaRun struct {
	XMLName          xml.Name   `xml:"nmaprun"`
	Scanner          string     `xml:"scanner,attr" json:"scanner"`
	Start            string     `xml:"start,attr" json:"start"`
	Version          string     `xml:"version,attr" json:"version"`
	XmlOutputVersion string     `xml:"xmloutputversion,attr" json:"xml_output_version"`
	ScanInfo         ScanInfo   `xml:"scaninfo" json:"scaninfo"`
	Hosts            []Hosts    `xml:"host" json:"host"`
	Finished         Finished   `xml:"finished" json:"finished"`
	HostRecord       HostRecord `xml:"hosts" json:"hosts"`
	rawXML           []byte
}

func (m MasscaRun) String() string {
	return formatter(m.Start)
}

type ScanInfo struct {
	Type     string `xml:"type,attr" json:"type"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
}

type Hosts struct {
	EndTime string  `xml:"endtime,attr" json:"end_time"`
	Address Address `xml:"address" json:"address"`
	Ports   []Ports `xml:"ports>port" json:"ports"`
}

func (h Hosts) String() string {
	return formatter(h.EndTime)
}

type Ports struct {
	ID       string `xml:"portid,attr" json:"port_id"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
	State    State  `xml:"state" json:"state"`
}

type Address struct {
	AddrType string `xml:"addrtype,attr" json:"addr_type"`
	Addr     string `xml:"addr,attr" json:"addr"`
}

func (a Address) String() string {
	return a.Addr
}

type State struct {
	State     string `xml:"state,attr" json:"state"`
	Reason    string `xml:"reason,attr" json:"reason"`
	ReasonTTL string `xml:"reasonttl,attr" json:"reason_ttl"`
}

func (s State) Status() string {
	return s.State
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

func (f Finished) String() string {
	return formatter(f.Time)
}

func formatter(timeStr string) string {
	tm, _ := strconv.ParseInt(timeStr, 0, 64)
	return time.Unix(tm, 0).Format("2006-01-02 15:04:05")
}

func ParseXML(content []byte) (*MasscaRun, error) {
	r := &MasscaRun{
		rawXML: content,
	}
	err := xml.Unmarshal(content, r)
	return r, err
}
