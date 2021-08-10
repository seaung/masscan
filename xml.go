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

type ScanInfo struct {
	Type     string `xml:"type,attr" json:"type"`
	Protocol string `xml:"protocol,attr" json:"protocol"`
}

type Hosts struct {
	EndTime TimeStamp `xml:"endtime,attr" json:"end_time"`
	Address Address   `xml:"address" json:"address"`
	Ports   []Ports   `xml:"ports>port" json:"ports"`
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
	Elapsed string    `xml:"elapsed,attr" json:"elapsed"`
	Time    string    `xml:"time,attr" json:"time"`
	TimeStr TimeStamp `xml:"timeStr,attr" json:"time_str"`
}

type TimeStamp time.Time

func (t *TimeStamp) ParseTime(s string) error {
	timestamp, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}

	*t = TimeStamp(time.Unix(timestamp, 0))
	return nil
}

func (t TimeStamp) FormatTime() string {
	return strconv.FormatInt(time.Time(t).Unix(), 10)
}

func (t TimeStamp) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	if time.Time(t).IsZero() {
		return xml.Attr{}, nil
	}
	return xml.Attr{Name: name, Value: t.FormatTime()}, nil
}

func (t *TimeStamp) UnMarshalXMLAttr(attr xml.Attr) (err error) {
	return t.ParseTime(attr.Value)
}

func ParseXML(content []byte) (*MasscaRun, error) {
	r := &MasscaRun{
		rawXML: content,
	}
	err := xml.Unmarshal(content, r)
	return r, err
}
