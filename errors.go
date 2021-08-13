package masscan

import (
	"errors"
)

var (
	MasscanNotInstalledError   = errors.New("masscan is not installed")
	MasscanNotFoundError       = errors.New("path does not exist")
	MasscanScanTimeoutError    = errors.New("scan timeout")
	MasscaScanResultParseError = errors.New("parse result error")
)
