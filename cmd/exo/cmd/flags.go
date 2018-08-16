package cmd

import (
	"fmt"
	"strconv"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

type uint8PtrValue struct {
	*uint8
}

func (v *uint8PtrValue) Set(val string) error {
	r, err := strconv.ParseUint(val, 10, 8)
	if err != nil {
		return err
	}
	res := uint8(r)
	v.uint8 = &res
	return nil
}

func (v *uint8PtrValue) Type() string {
	return "uint8"
}

func (v *uint8PtrValue) String() string {
	if v.uint8 == nil {
		return "nil"
	}
	return strconv.FormatUint(uint64(*v.uint8), 10)
}

func getUint8CustomFlag(cmd *cobra.Command, name string) (uint8PtrValue, error) {
	it := cmd.Flags().Lookup(name)
	if it != nil {
		r := it.Value.(*uint8PtrValue)
		if r != nil {
			return *r, nil
		}
	}
	return uint8PtrValue{}, fmt.Errorf("unable to get flag %q", name)
}

type int64PtrValue struct {
	*int64
}

func (v *int64PtrValue) Set(val string) error {
	r, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return err
	}
	v.int64 = &r
	return nil
}

func (v *int64PtrValue) Type() string {
	return "int64"
}

func (v *int64PtrValue) String() string {
	if v.int64 == nil {
		return "nil"
	}
	return strconv.FormatInt(*v.int64, 10)
}

func getInt64CustomFlag(cmd *cobra.Command, name string) (int64PtrValue, error) {
	it := cmd.Flags().Lookup(name)
	if it != nil {
		r := it.Value.(*int64PtrValue)
		if r != nil {
			return *r, nil
		}
	}
	return int64PtrValue{}, fmt.Errorf("unable to get flag %q", name)
}

// EXO API flags

// UUID egoscale.UUID embeded
type UUID struct {
	UUID **egoscale.UUID
}

// Set set
func (v *UUID) Set(val string) error {
	r, err := egoscale.ParseUUID(val)
	if err != nil {
		return err
	}
	*(v.UUID) = r

	return nil
}

//Type type
func (v *UUID) Type() string {
	return "UUID"
}

func (v *UUID) String() string {
	if v.UUID == nil || *(v.UUID) == nil {
		return "nil"
	}

	return (*(v.UUID)).String()
}

func getUUIDFlag(cmd *cobra.Command, name string) (UUID, error) {
	it := cmd.Flags().Lookup(name)
	if it != nil {
		r := it.Value.(*UUID)
		if r != nil {
			return *r, nil
		}
	}
	return UUID{}, fmt.Errorf("unable to get flag %q", name)
}

// CIDR egoscale.UUID embeded
type CIDR struct {
	CIDR **egoscale.CIDR
}

// Set set
func (v *CIDR) Set(val string) error {
	r, err := egoscale.ParseCIDR(val)
	if err != nil {
		return err
	}
	*(v.CIDR) = r

	return nil
}

//Type type
func (v *CIDR) Type() string {
	return "CIDR"
}

func (v *CIDR) String() string {
	if v.CIDR == nil || *(v.CIDR) == nil {
		return "nil"
	}

	return (*(v.CIDR)).String()
}

func getCIDRFlag(cmd *cobra.Command, name string) (CIDR, error) {
	it := cmd.Flags().Lookup(name)
	if it != nil {
		r := it.Value.(*CIDR)
		if r != nil {
			return *r, nil
		}
	}
	return CIDR{}, fmt.Errorf("unable to get flag %q", name)
}
