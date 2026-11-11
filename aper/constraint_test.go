// Copyright (C) 2026 Intel Corporation
// SPDX-License-Identifier: Apache-2.0

package aper

import (
	"bytes"
	"reflect"
	"testing"
)

// roundTrip marshals in, compares with want, then unmarshals into a new value of the same type.
func roundTrip(t *testing.T, in any, want []byte) {
	t.Helper()
	got, err := MarshalWithParams(in, "")
	if err != nil {
		t.Fatalf("marshal %T: %v", in, err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("marshal %T = %x, want %x", in, got, want)
	}
	out := reflect.New(reflect.TypeOf(in))
	if err := Unmarshal(want, out.Interface()); err != nil {
		t.Fatalf("unmarshal %T: %v", in, err)
	}
	if !reflect.DeepEqual(in, out.Elem().Interface()) {
		t.Fatalf("unmarshal %T = %+v, want %+v", in, out.Elem().Interface(), in)
	}
}

// INTEGER (0..86399): the range needs 3 octets, so the length field is 2 bits (X.691 #11.5.7.4).
func TestConstrainedIntegerLengthField(t *testing.T) {
	type intLarge struct {
		Value int64 `aper:"valueLB:0,valueUB:86399"`
	}
	roundTrip(t, intLarge{86399}, []byte{0x80, 0x01, 0x51, 0x7F})
	roundTrip(t, intLarge{1}, []byte{0x00, 0x01})
}

// INTEGER (1..40000000): the offset from the lower bound, not the value, is sized.
func TestConstrainedIntegerUsesOffsetFromLowerBound(t *testing.T) {
	type intOffset struct {
		Value int64 `aper:"valueLB:1,valueUB:40000000"`
	}
	// 65537 - 1 = 0x010000 needs 3 octets; length field is 2 bits holding 3-1.
	roundTrip(t, intOffset{65537}, []byte{0x80, 0x01, 0x00, 0x00})
	// 256 - 1 = 0xFF fits in 1 octet even though 256 needs 2.
	roundTrip(t, intOffset{256}, []byte{0x00, 0xFF})
}

// INTEGER (-640000..640000): negative values are encoded as an offset from the lower bound.
func TestConstrainedIntegerNegativeLowerBound(t *testing.T) {
	type intNeg struct {
		Value int64 `aper:"valueLB:-640000,valueUB:640000"`
	}
	// -1 + 640000 = 639999 = 0x09C3FF
	roundTrip(t, intNeg{-1}, []byte{0x80, 0x09, 0xC3, 0xFF})
}

// SIZE(1..MAX): the length determinant is the length itself, not length - lb (X.691 #11.9.4.2).
func TestUnboundedSizeIgnoresLowerBound(t *testing.T) {
	type octets struct {
		Value OctetString `aper:"sizeLB:1"`
	}
	roundTrip(t, octets{OctetString("AB")}, []byte{0x02, 0x41, 0x42})

	type bits struct {
		Value BitString `aper:"sizeLB:1"`
	}
	roundTrip(t, bits{BitString{[]byte{0xA0}, 3}}, []byte{0x03, 0xA0})

	// ub >= 64K is also an unconstrained length.
	type bigBits struct {
		Value BitString `aper:"sizeLB:1,sizeUB:131072"`
	}
	roundTrip(t, bigBits{BitString{[]byte{0xA0}, 3}}, []byte{0x03, 0xA0})
}

// SEQUENCE OF with ub = 64K is not a constrained count: it uses a general length determinant.
func TestSequenceOfUnconstrainedCount(t *testing.T) {
	type list struct {
		List []bool `aper:"sizeLB:1,sizeUB:65536"`
	}
	for _, n := range []int{1, 127, 128, 200, 16383} {
		in := list{List: make([]bool, n)}
		for i := range in.List {
			in.List[i] = true
		}
		got, err := MarshalWithParams(in, "")
		if err != nil {
			t.Fatalf("n=%d marshal: %v", n, err)
		}
		var prefix []byte
		if n < 128 {
			prefix = []byte{byte(n)}
		} else {
			prefix = []byte{0x80 | byte(n>>8), byte(n)}
		}
		if !bytes.HasPrefix(got, prefix) || len(got) != len(prefix)+(n+7)/8 {
			t.Fatalf("n=%d marshal = %x...", n, got[:4])
		}
		var out list
		if err := Unmarshal(got, &out); err != nil {
			t.Fatalf("n=%d unmarshal: %v", n, err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("n=%d round trip mismatch", n)
		}
	}
}

// Counts of 16K or more are sent as fragments (X.691 #10.9.3.8).
func TestSequenceOfFragmentedCount(t *testing.T) {
	type list struct {
		List []bool `aper:"sizeLB:1,sizeUB:65536"`
	}
	for _, n := range []int{16384, 20000, 32768, 65536} {
		in := list{List: make([]bool, n)}
		for i := range in.List {
			in.List[i] = i%3 == 0
		}
		got, err := MarshalWithParams(in, "")
		if err != nil {
			t.Fatalf("n=%d marshal: %v", n, err)
		}
		var out list
		if err := Unmarshal(got, &out); err != nil {
			t.Fatalf("n=%d unmarshal: %v", n, err)
		}
		if !reflect.DeepEqual(in, out) {
			t.Fatalf("n=%d round trip mismatch", n)
		}
	}
}
