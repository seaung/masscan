package masscan

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMasscanGetInterfaceList(t *testing.T) {
	scanner, err := NewMasscanScanner(WithMasscanBinaryPath("tests/scripts/fake_masscan_iflist.sh"))

	assert.NoError(t, err)

	ifaces, err := scanner.GetInterfaceList()

	assert.NoError(t, err)
	assert.NotNil(t, ifaces)

	assert.Len(t, ifaces.Interfaces, 1)
}

func TestConverInterface(t *testing.T) {
	iface := converInterface("0  eth0        (No description available)")

	assert.Equal(t, "0", iface.Index)
	assert.Equal(t, "eth0", iface.IFace)
	assert.Equal(t, "(No description available)", iface.Description)
}
