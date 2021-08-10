package masscan

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type MasscanRunner interface {
	Run() (result *MasscaRun, warnning []string, err error)
}

type MasscanScanner struct {
	cmd            *exec.Cmd
	args           []string
	masscanPath    string
	cxt            context.Context
	stderr, stdout bufio.Scanner
}

type Options func(*MasscanScanner)

func NewMasscanScanner(options ...Options) (*MasscanScanner, error) {
	masscanScanner := &MasscanScanner{}

	for _, option := range options {
		option(masscanScanner)
	}

	if masscanScanner.masscanPath == "" {
		var err error
		masscanScanner.masscanPath, err = exec.LookPath("masscan")
		if err != nil {
			return nil, MasscanNotInstalledError
		}
	}

	if masscanScanner.cxt == nil {
		masscanScanner.cxt = context.Background()
	}

	return masscanScanner, nil
}

func NewMasscanScannerWithBinaryPath(binaryPath string, options ...Options) (*MasscanScanner, error) {
	masscanScanner := &MasscanScanner{}

	for _, option := range options {
		option(masscanScanner)
	}

	if _, err := os.Stat(binaryPath); os.IsExist(err) {
		return nil, MasscanNotFoundError
	}

	masscanScanner.masscanPath = binaryPath

	fmt.Println("masscan path : ", masscanScanner.masscanPath)

	if masscanScanner.cxt == nil {
		masscanScanner.cxt = context.Background()
	}

	return masscanScanner, nil
}

func (m *MasscanScanner) Run() (result *MasscaRun, warnning []string, err error) {
	var (
		stdout, stderr bytes.Buffer
		resume         bool
	)

	args := m.args

	for _, arg := range args {
		if arg == "--resume" {
			resume = true
			break
		}
	}

	if !resume {
		args = append(args, "-oX")
		args = append(args, "-")
	}

	cmd := exec.Command(m.masscanPath, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Start()

	if err != nil {
		return nil, warnning, err
	}

	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-m.cxt.Done():
		_ = cmd.Process.Kill()
		return nil, warnning, MasscanScanTimeoutError
	case <-done:
		if stderr.Len() > 0 {
			warnning = strings.Split(strings.Trim(stderr.String(), "\n"), "\n")
		}

		result, err := ParseXML(stdout.Bytes())
		if err != nil {
			warnning = append(warnning, err.Error())
			return nil, warnning, MasscaScanResultParseError
		}
		return result, warnning, nil
	}
}

// anything on the command-line not prefixed with a '-' is assumed to be an IP address or range.
// There are three valid formats. The first is a single IPv4 address like "192.168.0.1".
// The second is a range like "10.0.0.1-10.0.0.100".
// The third is a CIDR address, like "0.0.0.0/0".
// At least one target must be specified. Multiple targets can be specified.
// This can be specified as multiple options separated by space,
// or can be separated by a comma as a single option, such as 10.0.0.0/8,192.168.0.1.
func WithTargets(targets ...string) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, targets...)
	}
}

// specifies the port(s) to be scanned.
func WithPorts(ports ...string) Options {
	portList := strings.Join(ports, ",")

	return func(m *MasscanScanner) {
		var flags int = -1
		for p, value := range m.args {
			if value == "-p" {
				flags = p
				break
			}
		}

		if flags >= 0 {
			portList = m.args[flags+1] + "," + portList
			m.args[flags+1] = portList
		} else {
			m.args = append(m.args, "-p")
			m.args = append(m.args, portList)
		}
	}
}

// specifies that banners should be grabbed, like HTTP server versions, HTML title fields, and so forth.
// Only a few protocols are supported.
func WithBanners() Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--banners")
	}
}

// specifies the desired rate for transmitting packets.
func WithRate(rate int) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--rate")
		m.args = append(m.args, fmt.Sprint(rate))
	}
}

// indicates that the scan should include an ICMP echo request.
func WithPingScan() Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--ping")
	}
}

// list the available network interfaces, and then exits.
func WithIfList() Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--iflist")
	}
}

// specifies the TTL of outgoing packets, defaults to 255.
func WithTTL(number int16) Options {
	return func(m *MasscanScanner) {
		if number < 0 || number > 255 {
			panic("")
		}

		m.args = append(m.args, "--ttl")
		m.args = append(m.args, fmt.Sprint(number))
	}
}

func WithConnectionTimeout(second int) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--connection-timeout")
		m.args = append(m.args, fmt.Sprint(second))
	}
}

func WithRetriesNumber(number int) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--retries")
		m.args = append(m.args, fmt.Sprint(number))
	}
}

func WithTargetExclusion(target string) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--exclude")
		m.args = append(m.args, target)
	}
}

func WithExclusionFile(fileName string) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--excludefile")
		m.args = append(m.args, fileName)
	}
}

func WithTopPorts(ports int) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "‐‐top-ports")
		m.args = append(m.args, fmt.Sprint(ports))
	}
}

func WithTopTen() Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "‐‐top-ten")
	}
}

func WithResumePreviousScan(resumeFile string) Options {
	return func(m *MasscanScanner) {
		m.args = append(m.args, "--resume")
		m.args = append(m.args, resumeFile)
	}
}

func WithContext(cxt context.Context) Options {
	return func(m *MasscanScanner) {
		m.cxt = cxt
	}
}
