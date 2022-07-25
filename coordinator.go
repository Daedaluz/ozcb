package ozcb

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type SecurityLevel int

func (s SecurityLevel) String() string {
	switch s {
	case None:
		return "None"
	case MIC32:
		return "MIC32"
	case MIC64:
		return "MIC64"
	case MIC128:
		return "MIC128"
	case ENC:
		return "ENC"
	case ENC_MIC32:
		return "ENC_MIC32"
	case ENC_MIC64:
		return "ENC_MIC64"
	case ENC_MIC128:
		return "ENC_MIC128"
	default:
		return fmt.Sprintf("Unknown(%d)", int(s))
	}
}

// Security levels
const (
	None = SecurityLevel(iota)
	MIC32
	MIC64
	MIC128
	ENC
	ENC_MIC32
	ENC_MIC64
	ENC_MIC128
)

const FormatZigpy = "zigpy/open-coordinator-backup"

type Meta struct {
	Format   string         `json:"format"`
	Version  int            `json:"version"`
	Source   string         `json:"source"`
	Internal map[string]any `json:"internal"`
}

type ZStack struct {
	TCLKSeed string `json:"tclk_seed"`
}

type StackSpecific struct {
	ZStack ZStack `json:"zstack"`
}

type NetworkKey struct {
	Key     string `json:"key"`
	Seq     uint8  `json:"sequence_number"`
	Counter uint32 `json:"frame_counter"`
}

type LinkKey struct {
	Key       string `json:"key"`
	RXCounter uint32 `json:"rx_counter"`
	TXCounter uint32 `json:"tx_counter"`
}

type Backup struct {
	Metadata        Meta           `json:"metadata"`
	StackSpecific   *StackSpecific `json:"stack_specific,omitempty"`
	CoordinatorIEEE string         `json:"coordinator_ieee"`
	PanID           string         `json:"pan_id"`
	EXTPanID        string         `json:"extended_pan_id"`
	Channel         int            `json:"channel"`
	ChannelMask     []int          `json:"channel_mask"`
	SecurityLevel   SecurityLevel  `json:"security_level"`
	NWKUpdateID     int            `json:"nwk_update_id"`
	NetworkKey      NetworkKey     `json:"network_key"`
	Devices         []*Device      `json:"devices,omitempty"`
}

func NewBackup(applicationName string, meta map[string]any) *Backup {
	return &Backup{
		Metadata: Meta{
			Format:   FormatZigpy,
			Version:  1,
			Source:   applicationName,
			Internal: meta,
		},
		StackSpecific:   nil,
		CoordinatorIEEE: "",
		PanID:           "",
		EXTPanID:        "",
		Channel:         0,
		ChannelMask:     []int{},
		SecurityLevel:   0,
		NWKUpdateID:     0,
		NetworkKey:      NetworkKey{},
		Devices:         make([]*Device, 0, 100),
	}
}

func (b *Backup) SetPanIDBytes(panID []byte) *Backup {
	b.PanID = hex.EncodeToString(panID[0:2])
	return b
}

func (b *Backup) SetPanID(panID string) *Backup {
	b.PanID = fixHexString(panID, 4)
	return b
}
func (b *Backup) SetPanIDInt(addr uint16) *Backup {
	var x = [8]byte{}
	binary.BigEndian.PutUint16(x[:], addr)
	b.PanID = hex.EncodeToString(x[:])
	return b
}

func (b *Backup) GetPanIDBytes() []byte {
	x, _ := hex.DecodeString(b.PanID)
	return x
}

func (b *Backup) SetCoordinatorIEEEBytes(addr []byte) *Backup {
	b.CoordinatorIEEE = hex.EncodeToString(addr[0:8])
	return b
}

func (b *Backup) SetCoordinatorIEEE(addr string) *Backup {
	b.CoordinatorIEEE = fixHexString(addr, 16)
	return b
}
func (b *Backup) SetCoordinatorIEEEInt(addr uint64) *Backup {
	var x = [8]byte{}
	binary.BigEndian.PutUint64(x[:], addr)
	b.CoordinatorIEEE = hex.EncodeToString(x[:])
	return b
}

func (b *Backup) GetCoordinatorIEEEBytes() []byte {
	x, _ := hex.DecodeString(b.CoordinatorIEEE)
	return x
}

func (b *Backup) SetExtendedPanIDBytes(extPan []byte) *Backup {
	b.EXTPanID = hex.EncodeToString(extPan[0:8])
	return b
}

func (b *Backup) SetExtendedPanID(extPan string) *Backup {
	b.EXTPanID = fixHexString(extPan, 16)
	return b
}

func (b *Backup) SetExtendedPanIDInt(extPan uint64) *Backup {
	var x = [8]byte{}
	binary.BigEndian.PutUint64(x[:], extPan)
	b.EXTPanID = hex.EncodeToString(x[:])
	return b
}

func (b *Backup) GetExtendedPanIDBytes() []byte {
	x, _ := hex.DecodeString(b.EXTPanID)
	return x
}

func (b *Backup) SetChannel(ch int) *Backup {
	b.Channel = ch
	return b
}

func (b *Backup) AddChannels(chs ...int) *Backup {
	b.ChannelMask = append(b.ChannelMask, chs...)
	return b
}

func (b *Backup) SetNWKUpdateID(id int) *Backup {
	b.NWKUpdateID = id
	return b
}

func (b *Backup) SetSecurityLevel(level SecurityLevel) *Backup {
	b.SecurityLevel = level
	return b
}

func (b *Backup) SetNetworkKeyBytes(key []byte, sequence uint8, counter uint32) *Backup {
	b.NetworkKey.Key = hex.EncodeToString(key)
	b.NetworkKey.Counter = counter
	b.NetworkKey.Seq = sequence
	return b
}

func (b *Backup) SetNetworkKey(key string, sequence uint8, counter uint32) *Backup {
	b.NetworkKey.Key = fixHexString(key, 32)
	b.NetworkKey.Counter = counter
	b.NetworkKey.Seq = sequence
	return b
}

func (b *Backup) GetNetworkKeyBytes() []byte {
	x, _ := hex.DecodeString(b.NetworkKey.Key)
	return x
}

func (b *Backup) AddDevice(dev *Device) *Backup {
	b.Devices = append(b.Devices, dev)
	return b
}

func (b *Backup) Write(w io.Writer) error {
	// TODO: Validate structure
	enc := json.NewEncoder(w)
	enc.SetIndent("", "    ")
	return enc.Encode(b)
}

func (b *Backup) Encode() (string, error) {
	data, err := json.MarshalIndent(b, "", "    ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *Backup) WriteFile(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	return b.Write(file)
}

func Read(r io.Reader) (*Backup, error) {
	res := &Backup{}
	dec := json.NewDecoder(r)
	if err := dec.Decode(res); err != nil {
		return nil, err
	}
	// TODO: Validate
	return res, nil
}

func ReadFile(file string) (*Backup, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Read(f)
}
