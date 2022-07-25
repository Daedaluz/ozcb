package ozcb

import (
	"encoding/binary"
	"encoding/hex"
)

type Device struct {
	IEEEAddress string   `json:"ieee_address"`
	NWKAddress  string   `json:"nwk_address"`
	IsChild     bool     `json:"is_child"`
	LinkKey     *LinkKey `json:"link_key,omitempty"`
}

func NewDevice() *Device {
	return &Device{
		IEEEAddress: "",
		NWKAddress:  "",
		IsChild:     true,
		LinkKey:     nil,
	}
}

func CreateDevice(ieeeAddr, nwkAddr string, child bool) *Device {
	return NewDevice().
		SetNWKAddress(nwkAddr).
		SetIEEEAddress(ieeeAddr).
		SetIsChild(child)
}

func (d *Device) SetIEEEAddressBytes(addr []byte) *Device {
	d.IEEEAddress = hex.EncodeToString(addr[0:8])
	return d
}

func (d *Device) SetIEEEAddress(addr string) *Device {
	d.IEEEAddress = fixHexString(addr, 16)
	return d
}

func (d *Device) SetIEEEAddressInt(addr uint64) *Device {
	var x = [8]byte{}
	binary.BigEndian.PutUint64(x[:], addr)
	d.IEEEAddress = hex.EncodeToString(x[:])
	return d
}

func (d *Device) GetIEEEAddressBytes() []byte {
	x, _ := hex.DecodeString(d.IEEEAddress)
	return x
}

func (d *Device) SetNWKAddressBytes(addr []byte) *Device {
	d.NWKAddress = hex.EncodeToString(addr[0:2])
	return d
}

func (d *Device) SetNWKAddress(addr string) *Device {
	d.NWKAddress = fixHexString(addr, 4)
	return d
}

func (d *Device) GetNWKAddressBytes() []byte {
	x, _ := hex.DecodeString(d.NWKAddress)
	return x
}

func (d *Device) SetNWKAddressInt(addr uint16) *Device {
	nwk := []byte{0, 0}
	binary.BigEndian.PutUint16(nwk, addr)
	d.NWKAddress = hex.EncodeToString(nwk)
	return d
}

func (d *Device) SetIsChild(yes bool) *Device {
	d.IsChild = yes
	return d
}

func (d *Device) SetLinkKey(key string, rx, tx uint32) *Device {
	d.LinkKey = &LinkKey{
		Key:       fixHexString(key, 32),
		RXCounter: rx,
		TXCounter: tx,
	}
	return d
}

func (d *Device) GetLinkKeyBytes() []byte {
	if d.LinkKey != nil {
		x, _ := hex.DecodeString(d.LinkKey.Key)
		return x
	}
	return nil
}
