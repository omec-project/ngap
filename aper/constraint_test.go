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

// A full 16K fragment must be followed by a final (possibly zero) length determinant (X.691 #11.9.3.8.3).
func TestStringFragmentTerminator(t *testing.T) {
	type octets struct {
		A OctetString `aper:"sizeLB:1"`
		B int64       `aper:"valueLB:0,valueUB:255"`
	}
	for _, tc := range []struct {
		n    int
		want func(data []byte) []byte // encoding of A, without B
	}{
		{100, func(d []byte) []byte { return append([]byte{0x64}, d...) }},
		{16384, func(d []byte) []byte { return append(append([]byte{0xC1}, d...), 0x00) }},
		{16385, func(d []byte) []byte { return append(append([]byte{0xC1}, d[:16384]...), 0x01, 0x5A) }},
		{32768, func(d []byte) []byte { return append(append([]byte{0xC2}, d...), 0x00) }},
		{65536, func(d []byte) []byte { return append(append([]byte{0xC4}, d...), 0x00) }},
	} {
		data := bytes.Repeat([]byte{0x5A}, tc.n)
		roundTrip(t, octets{OctetString(data), 7}, append(tc.want(data), 0x07))
	}

	// BIT STRING of exactly 16384 bits: 0xC1, 2048 octets, 0x00 terminator.
	type bits struct {
		A BitString `aper:"sizeLB:1,sizeUB:131072"`
		B int64     `aper:"valueLB:0,valueUB:255"`
	}
	bitData := bytes.Repeat([]byte{0xA5}, 2048)
	exp := append(append([]byte{0xC1}, bitData...), 0x00, 0x07)
	roundTrip(t, bits{BitString{bitData, 16384}, 7}, exp)
}

// A value below the SIZE lower bound must be rejected on encode and on decode.
func TestStringBelowLowerBound(t *testing.T) {
	type octets struct {
		Value OctetString `aper:"sizeLB:1"`
	}
	type bits struct {
		Value BitString `aper:"sizeLB:1"`
	}
	type bigBits struct {
		Value BitString `aper:"sizeLB:1,sizeUB:131072"`
	}
	type bigOctets struct {
		Value OctetString `aper:"sizeLB:1,sizeUB:131072"`
	}
	for _, in := range []any{octets{}, bits{}, bigBits{}, bigOctets{}} {
		if got, err := MarshalWithParams(in, ""); err == nil {
			t.Errorf("marshal %T: expected error, got %x", in, got)
		}
		out := reflect.New(reflect.TypeOf(in))
		if err := Unmarshal([]byte{0x00}, out.Interface()); err == nil {
			t.Errorf("unmarshal %T: expected error for empty value", in)
		}
	}
}

// The declared SIZE bounds of an unconstrained-count SEQUENCE OF are enforced while decoding.
func TestSequenceOfDecodeBounds(t *testing.T) {
	// Zero-bit components: two 64K fragments exceed SIZE(0..65536).
	type bounded struct {
		List []int64 `aper:"sizeLB:0,sizeUB:65536,valueLB:0,valueUB:0"`
	}
	var b bounded
	if err := Unmarshal([]byte{0xC4, 0xC4, 0x00}, &b); err == nil {
		t.Errorf("expected error for %d components over the upper bound", len(b.List))
	}

	// No declared bound: the cumulative cap applies (17 * 64K > 1M).
	type unbounded struct {
		List []int64 `aper:"valueLB:0,valueUB:0"`
	}
	in := append(bytes.Repeat([]byte{0xC4}, 17), 0x00)
	var u unbounded
	if err := Unmarshal(in, &u); err == nil {
		t.Errorf("expected error for %d components over the cap", len(u.List))
	}

	// Below the declared minimum.
	type minimum struct {
		List []bool `aper:"sizeLB:1,sizeUB:65536"`
	}
	var m minimum
	if err := Unmarshal([]byte{0x00}, &m); err == nil {
		t.Errorf("expected error for an empty list with sizeLB:1")
	}
}

// A declared minimum of 64K or more must still be enforced (the lower bound is not offset, but it is checked).
func TestSequenceOfLargeLowerBound(t *testing.T) {
	type list struct {
		List []bool `aper:"sizeLB:65536,sizeUB:131072"`
	}
	if got, err := MarshalWithParams(list{}, ""); err == nil {
		t.Errorf("marshal: expected error for an empty list, got %x", got)
	}
	var out list
	if err := Unmarshal([]byte{0x00}, &out); err == nil {
		t.Errorf("unmarshal: expected error for an empty list")
	}
}

// Zero-bit components may end exactly at the end of the data, including in a final fragment.
func TestSequenceOfZeroBitComponentsAtEnd(t *testing.T) {
	type list struct {
		List []int64 `aper:"valueLB:0,valueUB:0"`
	}
	roundTrip(t, list{List: make([]int64, 3)}, []byte{0x03})
	roundTrip(t, list{List: make([]int64, 16385)}, []byte{0xC1, 0x01})
	roundTrip(t, list{List: make([]int64, 16384)}, []byte{0xC1, 0x00})
}

// Truncated input must still be rejected for components that need bits.
func TestSequenceOfTruncated(t *testing.T) {
	type list struct {
		List []bool `aper:"sizeLB:1,sizeUB:65536"`
	}
	var out list
	if err := Unmarshal([]byte{0x08, 0xFF}, &out); err != nil {
		t.Fatalf("complete input: %v", err)
	}
	if want := (list{List: []bool{true, true, true, true, true, true, true, true}}); !reflect.DeepEqual(out, want) {
		t.Fatalf("complete input = %+v, want %+v", out, want)
	}
	if err := Unmarshal([]byte{0x09, 0xFF}, &out); err == nil {
		t.Errorf("expected error for 9 components in 8 bits")
	}
	if err := Unmarshal(nil, &out); err == nil {
		t.Errorf("expected error for empty input")
	}
}

// An extensible SEQUENCE OF below the root minimum uses the extension encoding (bit 1, general length).
func TestSequenceOfExtensibleBelowRootMinimum(t *testing.T) {
	type list struct {
		List []bool `aper:"sizeExt,sizeLB:65536,sizeUB:131072"`
	}
	roundTrip(t, list{List: []bool{}}, []byte{0x80, 0x00})

	type small struct {
		List []bool `aper:"sizeExt,sizeLB:2,sizeUB:4"`
	}
	// Root: ext bit 0, count-lb in 2 bits, then the components.
	roundTrip(t, small{List: []bool{true, false}}, []byte{0x10})
	// Below the root: ext bit 1, general length 1, one component.
	roundTrip(t, small{List: []bool{true}}, []byte{0x80, 0x01, 0x80})
	// Above the root.
	roundTrip(t, small{List: []bool{true, true, true, true, true}}, []byte{0x80, 0x05, 0xF8})
}

// An empty encoding is not a valid aligned PER encoding, even for a zero-bit value.
func TestEmptyEncodingRejected(t *testing.T) {
	type zero struct {
		Value int64 `aper:"valueLB:0,valueUB:0"`
	}
	var z zero
	if err := Unmarshal(nil, &z); err == nil {
		t.Errorf("expected error for an empty buffer")
	}
	roundTrip(t, zero{}, []byte{0x00})
}

// An empty open type is rejected, and a zero-bit open type value is sent as one zero octet.
func TestOpenTypeEmpty(t *testing.T) {
	type zeroBit struct {
		Value int64 `aper:"valueLB:0,valueUB:0"`
	}
	type openChoice struct {
		Present int
		Zero    zeroBit `aper:"referenceFieldValue:1"`
	}
	type withOpen struct {
		ID    int64      `aper:"valueLB:0,valueUB:255"`
		Value openChoice `aper:"openType,referenceFieldName:ID"`
	}
	roundTrip(t, withOpen{1, openChoice{1, zeroBit{}}}, []byte{0x01, 0x01, 0x00})

	var out withOpen
	if err := Unmarshal([]byte{0x01, 0x00}, &out); err == nil {
		t.Errorf("expected error for an empty open type")
	}
}

// Constrained values and counts must lie within their declared range, even when the bit
// pattern read is representable (e.g. 255 in 8 bits for 255 values).
func TestConstrainedValueOutOfRange(t *testing.T) {
	type zeroBitList struct {
		List []int64 `aper:"sizeLB:0,sizeUB:254,valueLB:0,valueUB:0"`
	}
	var l zeroBitList
	if err := Unmarshal([]byte{0xFF}, &l); err == nil {
		t.Errorf("expected error for %d components over the upper bound", len(l.List))
	}

	type small struct {
		Value int64 `aper:"valueLB:0,valueUB:4"`
	}
	var s small
	if err := Unmarshal([]byte{0xE0}, &s); err == nil {
		t.Errorf("expected error for value 7 in INTEGER (0..4), got %d", s.Value)
	}
	roundTrip(t, small{4}, []byte{0x80})

	type enum struct {
		Value Enumerated `aper:"valueLB:0,valueUB:2"`
	}
	var e enum
	if err := Unmarshal([]byte{0xC0}, &e); err == nil {
		t.Errorf("expected error for value 3 in ENUMERATED (0..2)")
	}
}

// An empty open type is rejected even when its reference value is unknown and it would be skipped.
func TestOpenTypeEmptyUnknownID(t *testing.T) {
	type zeroBit struct {
		Value int64 `aper:"valueLB:0,valueUB:0"`
	}
	type openChoice struct {
		Present int
		Zero    zeroBit `aper:"referenceFieldValue:1"`
	}
	type withOpen struct {
		ID    int64      `aper:"valueLB:0,valueUB:255"`
		Value openChoice `aper:"openType,referenceFieldName:ID"`
	}
	var out withOpen
	if err := Unmarshal([]byte{0x02, 0x00}, &out); err == nil {
		t.Errorf("expected error for an empty skipped open type")
	}
	if err := Unmarshal([]byte{0x02, 0x01, 0x00}, &out); err != nil || out.Value.Present != 0 {
		t.Errorf("non-empty unknown open type: err=%v present=%d", err, out.Value.Present)
	}
}

// ENUMERATED extension values are encoded as normally small numbers; a zero-length one is malformed.
func TestEnumeratedExtensionEncoding(t *testing.T) {
	type enum struct {
		Value Enumerated `aper:"valueExt,valueLB:0,valueUB:2"`
	}
	roundTrip(t, enum{1}, []byte{0x20})
	// Extension bit 1, then the "small" flag 0 and 6 bits holding value-3.
	roundTrip(t, enum{3}, []byte{0x80})
	roundTrip(t, enum{66}, []byte{0xBF})
	// 67-3 = 64 does not fit in 6 bits: "not small" flag, then a length and the value.
	roundTrip(t, enum{67}, []byte{0xC0, 0x01, 0x40})

	// Extension bit 1, "not small" flag 1, then a length of zero octets.
	var out enum
	if err := Unmarshal([]byte{0xC0, 0x00}, &out); err == nil {
		t.Errorf("expected error for a zero-length extension value, got %d", out.Value)
	}
	// One octet holding 5 decodes to 3+5.
	if err := Unmarshal([]byte{0xC0, 0x01, 0x05}, &out); err != nil || out.Value != 8 {
		t.Errorf("unmarshal = %d, %v, want 8", out.Value, err)
	}
}
