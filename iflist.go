package masscan

import (
	"bytes"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type InterfaceList struct {
	Interfaces []*Interface `json:"interfaces"`
}

type Interface struct {
	Index       int
	IFace       string
	Description string
}

func (m *MasscanScanner) GetInterfaceList() (ifaces *InterfaceList, err error) {
	var stdout, stderr bytes.Buffer

	args := append(m.args, "--iflist")

	cmd := exec.Command(m.masscanPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	if err != nil {
		return nil, err
	}

	ifaces = parseInterfaces(stdout.Bytes())
	return ifaces, nil
}

func parseInterfaces(content []byte) *InterfaceList {
	ifaces := InterfaceList{
		Interfaces: make([]*Interface, 0),
	}

	output := string(content)

	lines := strings.Split(output, "\n")

	for i, line := range lines {
		if match, _ := regexp.MatchString("^[0-9]+$", line); match {
			for _, li := range lines[i+2:] {
				if iface := converInterface(li); iface != nil {
					ifaces.Interfaces = append(ifaces.Interfaces, iface)
				}
			}
		}
	}

	return &ifaces
}

func converInterface(line string) *Interface {
	fields := strings.Fields(line)

	if len(fields) < 3 {
		return nil
	}

	iface := &Interface{
		IFace:       fields[1],
		Description: fields[2],
	}

	if value, err := strconv.Atoi(fields[0]); err != nil {
		iface.Index = value
	}

	return iface
}
