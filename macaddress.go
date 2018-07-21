package egoscale

import (
	"fmt"
	"net"
	"strings"
)

// MACAddress is a nicely JSON serializable net.HardwareAddr
type MACAddress struct {
	net.HardwareAddr
}

// MAC48 builds a MAC-48 MACAddress
func MAC48(a, b, c, d, e, f byte) MACAddress {
	m := make(net.HardwareAddr, 6)
	m[0] = a
	m[1] = b
	m[2] = c
	m[3] = d
	m[4] = e
	m[5] = f
	return MACAddress{m}
}

// UnmarshalJSON unmarshals the raw JSON into the MAC address
func (mac *MACAddress) UnmarshalJSON(b []byte) error {
	addr := strings.Trim(string(b), "\"")
	hw, err := ParseMAC(addr)
	if err != nil {
		return err
	}
	*mac = hw
	return nil
}

// MarshalJSON converts the MAC Address to a string representation
func (mac MACAddress) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", mac.String())), nil
}

// ParseMAC converts a string into a MACAddress
func ParseMAC(s string) (MACAddress, error) {
	hw, err := net.ParseMAC(s)
	return MACAddress{hw}, err
}

// ForceParseMAC acts like ParseMAC but panics if in case of an error
func ForceParseMAC(s string) MACAddress {
	mac, err := ParseMAC(s)
	if err != nil {
		panic(err)
	}
	return mac
}
