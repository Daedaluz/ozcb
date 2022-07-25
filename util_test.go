package ozcb

import (
	"testing"
)

type test struct {
	input, expect string
	pad           int
}

var tests = []test{
	{
		pad:    10,
		input:  "000102030A",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "0x00 0x01 0x02 0x03 0x0a",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "00 01 02 03 0A",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "00:01:02:03:0A",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "0x00:0x01:0x02:0x03:0x0A",
		expect: "000102030a",
	},
	{
		pad: 10,
		input: "0x00: 0x01 : 0x02  : 0x03	: 0x0A",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "0x000102030A",
		expect: "000102030a",
	},
	{
		pad:    10,
		input:  "0x102030A",
		expect: "000102030a",
	},
}

func Test_fixHexString(t *testing.T) {
	for _, x := range tests {
		if res := fixHexString(x.input, x.pad); res != x.expect {
			t.Logf("fix(%s) != %s; got %s", x.input, x.expect, res)
			t.Fail()
		}
	}
}
