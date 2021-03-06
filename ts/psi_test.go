//    Copyright 2017 drillbits
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package ts

import (
	"bytes"
	"testing"
)

func TestPSITableID(t *testing.T) {
	p := PSI([]byte{0x01, 0x02, 0x03, 0x04, 0x05})
	exp := TableID(0x01)
	got := p.TableID()
	if got != exp {
		t.Errorf("PSI.TableID() => 0x%02X, want 0x%02X", got, exp)
	}
}

func TestPSICRC32(t *testing.T) {
	p := PSI([]byte{0x01, 0x02, 0x03, 0x04, 0x05})
	exp := []byte{0x02, 0x03, 0x04, 0x05}
	got := p.CRC32()
	if !bytes.Equal(got, exp) {
		t.Errorf("PSI.CRC32() => 0x%02X, want 0x%02X", got, exp)
	}
}

func TestPAT(t *testing.T) {
	for i, tc := range []struct {
		b          []byte
		tableID    TableID
		secInd     byte
		secLen     int
		crc32      CRC32
		tsID       TransportStreamID
		ver        int
		curNextInd byte
		secNum     byte
		lastSecNum byte
		nPID       PID
		pPIDMap    map[ProgramNumber]PID
		err        error
	}{
		{
			b: []byte{
				0x00, 0xB0, 0x1D, 0x7F, 0xE5, 0xED, 0x00, 0x00, 0x00, 0x00,
				0xE0, 0x10, 0x04, 0x28, 0xE4, 0x28, 0x04, 0x29, 0xE4, 0x29,
				0x04, 0x2A, 0xE4, 0x2A, 0x05, 0xA8, 0xFF, 0xC8, 0x8E, 0xFD,
				0xB2, 0xA4},
			tableID:    TableID(0x00),
			secInd:     0x01,
			secLen:     29,
			crc32:      CRC32([]byte{0x8E, 0xFD, 0xB2, 0xA4}),
			tsID:       0x7FE5,
			ver:        22,
			curNextInd: 1,
			secNum:     0x00,
			lastSecNum: 0x00,
			err:        nil,
		},
		{
			b: []byte{
				0x00, 0xB0, 0x1D, 0x7F, 0xE5, 0xED, 0x00, 0x00, 0x00, 0x00, 0xE0},
			err: ErrTooShort,
		},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			pat, err := NewPAT(tc.b)
			if tc.err == nil {
				if err != nil {
					t.Errorf("%0d: NewPAT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
				if PSI(pat).TableID() != tc.tableID {
					t.Errorf("%0d: PAT(0x%04X).TableID() => 0x%04X, want 0x%04X", i, tc.b, PSI(pat).TableID(), tc.tableID)
				}
				if PSI(pat).SectionSyntaxIndicator() != tc.secInd {
					t.Errorf("%0d: PAT(0x%04X).SectionSyntaxIndicator() => 0x%02X, want 0x%02X", i, tc.b, PSI(pat).SectionSyntaxIndicator(), tc.secInd)
				}
				if PSI(pat).SectionLength() != tc.secLen {
					t.Errorf("%0d: PAT(0x%04X).SectionLength() => %d, want %d", i, tc.b, PSI(pat).SectionLength(), tc.secLen)
				}
				if !bytes.Equal(PSI(pat).CRC32(), tc.crc32) {
					t.Errorf("%0d: PAT(0x%04X).CRC32() => 0x%04X, want 0x%04X", i, tc.b, PSI(pat).CRC32(), tc.crc32)
				}
				if pat.TransportStreamID() != tc.tsID {
					t.Errorf("%0d: PAT(0x%04X).TransportStreamID() => 0x%04X, want 0x%04X", i, tc.b, pat.TransportStreamID(), tc.tsID)
				}
				if pat.VersionNumber() != tc.ver {
					t.Errorf("%0d: PAT(0x%04X).VersionNumber() => %d, want %d", i, tc.b, pat.VersionNumber(), tc.ver)
				}
				if pat.CurrentNextIndicator() != tc.curNextInd {
					t.Errorf("%0d: PAT(0x%04X).CurrentNextIndicator() => 0x%04X, want 0x%04X", i, tc.b, pat.CurrentNextIndicator(), tc.curNextInd)
				}
				if pat.SectionNumber() != tc.secNum {
					t.Errorf("%0d: PAT(0x%04X).SectionNumber() => 0x%04X, want 0x%04X", i, tc.b, pat.SectionNumber(), tc.secNum)
				}
				if pat.LastSectionNumber() != tc.lastSecNum {
					t.Errorf("%0d: PAT(0x%04X).LastSectionNumber() => 0x%04X, want 0x%04X", i, tc.b, pat.LastSectionNumber(), tc.lastSecNum)
				}
			} else {
				if err != tc.err {
					t.Errorf("%0d: NewPAT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
			}
		})
	}
}

func TestPATPID(t *testing.T) {
	for i, tc := range []struct {
		b       []byte
		nPID    PID
		pPIDMap map[ProgramNumber]PID
		err     error
	}{
		{
			b: []byte{
				0x00, 0xB0, 0x1D, 0x7F, 0xE5, 0xED, 0x00, 0x00, 0x00, 0x00,
				0xE0, 0x10, 0x04, 0x28, 0xE4, 0x28, 0x04, 0x29, 0xE4, 0x29,
				0x04, 0x2A, 0xE4, 0x2A, 0x05, 0xA8, 0xFF, 0xC8, 0x8E, 0xFD,
				0xB2, 0xA4},
			nPID: 0x0010,
			pPIDMap: map[ProgramNumber]PID{
				0x0428: 0x0428,
				0x0429: 0x0429,
				0x042A: 0x042A,
				0x05A8: 0x1FC8,
			},
			err: nil,
		},
		{
			b: []byte{
				0x00, 0xB0, 0x19, 0x7F, 0xE5, 0xED, 0x00, 0x00, 0x04, 0x28,
				0xE4, 0x28, 0x04, 0x29, 0xE4, 0x29, 0x04, 0x2A, 0xE4, 0x2A,
				0x05, 0xA8, 0xFF, 0xC8, 0x8E, 0xFD, 0xB2, 0xA4},
			nPID: 0x0010,
			pPIDMap: map[ProgramNumber]PID{
				0x0428: 0x0428,
				0x0429: 0x0429,
				0x042A: 0x042A,
				0x05A8: 0x1FC8,
			},
			err: ErrNoNetworkID,
		},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			pat, err := NewPAT(tc.b)
			if err != nil {
				t.Fatal(err)
			}
			nPID, err := pat.NetworkPID()
			if tc.err == nil {
				if err != nil {
					t.Errorf("%0d: PAT(0x%04X).NetworkPID() \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
				if nPID != tc.nPID {
					t.Errorf("%0d: PAT(0x%04X).NetworkPID() => 0x%04X, want 0x%04X", i, tc.b, nPID, tc.nPID)
				}
				for n, exp := range tc.pPIDMap {
					if got, ok := pat.ProgramPIDMap()[n]; !ok || got != exp {
						t.Errorf("%0d: PAT(0x%04X).ProgramPIDMap()[0x%04X] => 0x%04X, want 0x%04X", i, tc.b, n, exp, got)
					}
				}
			} else {
				if err != tc.err {
					t.Errorf("%0d: PAT(0x%04X).NetworkPID() \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
			}
		})
	}
}

func TestCAT(t *testing.T) {
	for i, tc := range []struct {
		b           []byte
		tableID     TableID
		secInd      byte
		secLen      int
		crc32       CRC32
		ver         int
		curNextInd  byte
		secNum      byte
		lastSecNum  byte
		descriptors []Descriptor
		err         error
	}{
		{
			b: []byte{0x01, 0xB0, 0x10, 0xFF, 0xFF, 0xF9, 0x00, 0x00, 0xF6,
				0x05, 0x00, 0x0E, 0xE0, 0x71, 0x01, 0x04, 0xCC, 0x5F, 0xAB},
			tableID:    TableID(0x01),
			secInd:     0x01,
			secLen:     16,
			crc32:      CRC32([]byte{0x04, 0xCC, 0x5F, 0xAB}),
			ver:        28,
			curNextInd: 1,
			secNum:     0x00,
			lastSecNum: 0x00,
			descriptors: []Descriptor{
				{0xF6, 0x05, 0x00, 0x0E, 0xE0, 0x71, 0x01},
			},
			err: nil,
		},
		{
			b: []byte{
				0x01, 0xB0, 0x10, 0xFF, 0xFF, 0xF9, 0x00, 0x00, 0xF6, 0x05, 0x00},
			err: ErrTooShort,
		},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			cat, err := NewCAT(tc.b)
			if tc.err == nil {
				if err != nil {
					t.Errorf("%0d: NewCAT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
				if PSI(cat).TableID() != tc.tableID {
					t.Errorf("%0d: CAT(0x%04X).TableID() => 0x%04X, want 0x%04X", i, tc.b, PSI(cat).TableID(), tc.tableID)
				}
				if PSI(cat).SectionSyntaxIndicator() != tc.secInd {
					t.Errorf("%0d: CAT(0x%04X).SectionSyntaxIndicator() => 0x%02X, want 0x%02X", i, tc.b, PSI(cat).SectionSyntaxIndicator(), tc.secInd)
				}
				if PSI(cat).SectionLength() != tc.secLen {
					t.Errorf("%0d: CAT(0x%04X).SectionLength() => %d, want %d", i, tc.b, PSI(cat).SectionLength(), tc.secLen)
				}
				if !bytes.Equal(PSI(cat).CRC32(), tc.crc32) {
					t.Errorf("%0d: CAT(0x%04X).CRC32() => 0x%04X, want 0x%04X", i, tc.b, PSI(cat).CRC32(), tc.crc32)
				}
				if cat.VersionNumber() != tc.ver {
					t.Errorf("%0d: CAT(0x%04X).VersionNumber() => %d, want %d", i, tc.b, cat.VersionNumber(), tc.ver)
				}
				if cat.CurrentNextIndicator() != tc.curNextInd {
					t.Errorf("%0d: CAT(0x%04X).CurrentNextIndicator() => 0x%04X, want 0x%04X", i, tc.b, cat.CurrentNextIndicator(), tc.curNextInd)
				}
				if cat.SectionNumber() != tc.secNum {
					t.Errorf("%0d: CAT(0x%04X).SectionNumber() => 0x%04X, want 0x%04X", i, tc.b, cat.SectionNumber(), tc.secNum)
				}
				if cat.LastSectionNumber() != tc.lastSecNum {
					t.Errorf("%0d: CAT(0x%04X).LastSectionNumber() => 0x%04X, want 0x%04X", i, tc.b, cat.LastSectionNumber(), tc.lastSecNum)
				}
				if len(cat.Descriptors()) != len(tc.descriptors) {
					t.Errorf("%0d: CAT(0x%04X).Descriptors() => len: %d, want %d", i, tc.b, len(cat.Descriptors()), len(tc.descriptors))
				} else {
					for j, exp := range tc.descriptors {
						got := cat.Descriptors()[j]
						if !bytes.Equal(got, exp) {
							t.Errorf("%0d: CAT(0x%04X).Descriptors()[%d] => 0x%04X, want 0x%04X", i, tc.b, j, got, exp)
						}
					}
				}
			} else {
				if err != tc.err {
					t.Errorf("%0d: NewCAT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
			}
		})
	}
}

func TestPMT(t *testing.T) {
	for i, tc := range []struct {
		b           []byte
		tableID     TableID
		secInd      byte
		secLen      int
		crc32       CRC32
		pgNum       ProgramNumber
		ver         int
		curNextInd  byte
		secNum      byte
		lastSecNum  byte
		pcrPID      PID
		pgInfoLen   int
		descriptors []Descriptor
		elemInfo    []ProgramElementInfo
		err         error
	}{
		{
			b: []byte{0x02, 0xB0, 0x63, 0x05, 0xA8, 0xED, 0x00, 0x00, 0xE1,
				0x01, 0xF0, 0x06, 0xC1, 0x01, 0x88, 0xDE, 0x01, 0xEF, 0x1B,
				0xE1, 0x81, 0xF0, 0x03, 0x52, 0x01, 0x81, 0x0F, 0xE1, 0x82,
				0xF0, 0x03, 0x52, 0x01, 0x83, 0x06, 0xE1, 0x84, 0xF0, 0x08,
				0x52, 0x01, 0x87, 0xFD, 0x03, 0x00, 0x12, 0xAD, 0x0D, 0xF0,
				0x30, 0xF0, 0x0F, 0x52, 0x01, 0x80, 0xFD, 0x0A, 0x00, 0x0D,
				0x3F, 0x2F, 0x00, 0x0C, 0x00, 0x00, 0xFF, 0xBF, 0x0D, 0xF0,
				0x39, 0xF0, 0x03, 0x52, 0x01, 0x89, 0x0D, 0xF0, 0x3A, 0xF0,
				0x03, 0x52, 0x01, 0x8A, 0x0D, 0xF0, 0x3B, 0xF0, 0x0A, 0x52,
				0x01, 0x8B, 0xFD, 0x05, 0x00, 0x0D, 0x1F, 0xFF, 0xBF, 0xFE,
				0x9B, 0xEB, 0xD9},
			tableID:    TableID(0x02),
			secInd:     0x01,
			secLen:     99,
			crc32:      CRC32([]byte{0xFE, 0x9B, 0xEB, 0xD9}),
			pgNum:      0x05A8,
			ver:        22,
			curNextInd: 1,
			secNum:     0x00,
			lastSecNum: 0x00,
			pcrPID:     0x0101,
			pgInfoLen:  6,
			descriptors: []Descriptor{
				{0xC1, 0x01, 0x88},
				{0xDE, 0x01, 0xEF},
			},
			elemInfo: []ProgramElementInfo{
				{0x1B, 0xE1, 0x81, 0xF0, 0x03, 0x52, 0x01, 0x81},
				{0x0F, 0xE1, 0x82, 0xF0, 0x03, 0x52, 0x01, 0x83},
				{0x06, 0xE1, 0x84, 0xF0, 0x08, 0x52, 0x01, 0x87, 0xFD, 0x03, 0x00, 0x12, 0xAD},
				{0x0D, 0xF0, 0x30, 0xF0, 0x0F, 0x52, 0x01, 0x80, 0xFD, 0x0A, 0x00, 0x0D, 0x3F, 0x2F, 0x00, 0x0C, 0x00, 0x00, 0xFF, 0xBF},
				{0x0D, 0xF0, 0x39, 0xF0, 0x03, 0x52, 0x01, 0x89},
				{0x0D, 0xF0, 0x3A, 0xF0, 0x03, 0x52, 0x01, 0x8A},
				{0x0D, 0xF0, 0x3B, 0xF0, 0x0A, 0x52, 0x01, 0x8B, 0xFD, 0x05, 0x00, 0x0D, 0x1F, 0xFF, 0xBF},
			},
			err: nil,
		},
		{
			b: []byte{0x02, 0xB0, 0x63, 0x05, 0xA8, 0xED, 0x00, 0x00, 0xE1,
				0x01, 0xF0, 0x06, 0xC1, 0x01, 0x88},
			err: ErrTooShort,
		},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			pmt, err := NewPMT(tc.b)
			if tc.err == nil {
				if err != nil {
					t.Errorf("%0d: NewPMT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
				if PSI(pmt).TableID() != tc.tableID {
					t.Errorf("%0d: PMT(0x%04X).TableID() => 0x%04X, want 0x%04X", i, tc.b, PSI(pmt).TableID(), tc.tableID)
				}
				if PSI(pmt).SectionSyntaxIndicator() != tc.secInd {
					t.Errorf("%0d: PMT(0x%04X).SectionSyntaxIndicator() => 0x%02X, want 0x%02X", i, tc.b, PSI(pmt).SectionSyntaxIndicator(), tc.secInd)
				}
				if PSI(pmt).SectionLength() != tc.secLen {
					t.Errorf("%0d: PMT(0x%04X).SectionLength() => %d, want %d", i, tc.b, PSI(pmt).SectionLength(), tc.secLen)
				}
				if !bytes.Equal(PSI(pmt).CRC32(), tc.crc32) {
					t.Errorf("%0d: PMT(0x%04X).CRC32() => 0x%04X, want 0x%04X", i, tc.b, PSI(pmt).CRC32(), tc.crc32)
				}
				if pmt.ProgramNumber() != tc.pgNum {
					t.Errorf("%0d: PMT(0x%04X).ProgramNumber() => 0x%04X, want 0x%04X", i, tc.b, pmt.ProgramNumber(), tc.pgNum)
				}
				if pmt.VersionNumber() != tc.ver {
					t.Errorf("%0d: PMT(0x%04X).VersionNumber() => %d, want %d", i, tc.b, pmt.VersionNumber(), tc.ver)
				}
				if pmt.CurrentNextIndicator() != tc.curNextInd {
					t.Errorf("%0d: PMT(0x%04X).CurrentNextIndicator() => 0x%04X, want 0x%04X", i, tc.b, pmt.CurrentNextIndicator(), tc.curNextInd)
				}
				if pmt.SectionNumber() != tc.secNum {
					t.Errorf("%0d: PMT(0x%04X).SectionNumber() => 0x%04X, want 0x%04X", i, tc.b, pmt.SectionNumber(), tc.secNum)
				}
				if pmt.LastSectionNumber() != tc.lastSecNum {
					t.Errorf("%0d: PMT(0x%04X).LastSectionNumber() => 0x%04X, want 0x%04X", i, tc.b, pmt.LastSectionNumber(), tc.lastSecNum)
				}
				if pmt.PCRPID() != tc.pcrPID {
					t.Errorf("%0d: PMT(0x%04X).PCRPID() => 0x%04X, want 0x%04X", i, tc.b, pmt.PCRPID(), tc.pcrPID)
				}
				if pmt.ProgramInfoLength() != tc.pgInfoLen {
					t.Errorf("%0d: PMT(0x%04X).ProgramInfoLength() => %d, want %d", i, tc.b, pmt.ProgramInfoLength(), tc.pgInfoLen)
				}
				if len(pmt.Descriptors()) != len(tc.descriptors) {
					t.Errorf("%0d: PMT(0x%04X).Descriptors() => len: %d, want %d", i, tc.b, len(pmt.Descriptors()), len(tc.descriptors))
				} else {
					for j, exp := range tc.descriptors {
						got := pmt.Descriptors()[j]
						if !bytes.Equal(got, exp) {
							t.Errorf("%0d: PMT(0x%04X).Descriptors()[%d] => 0x%04X, want 0x%04X", i, tc.b, j, got, exp)
						}
					}
				}
				if len(pmt.ElementInfo()) != len(tc.elemInfo) {
					t.Errorf("%0d: PMT(0x%04X).ElementInfo() => len: %d, want %d", i, tc.b, len(pmt.ElementInfo()), len(tc.elemInfo))
				} else {
					// TODO: other tests
				}
			} else {
				if err != tc.err {
					t.Errorf("%0d: NewPMT(0x%04X) \ncauses %s, want %s", i, tc.b, err, tc.err)
				}
			}
		})
	}
}

func TestProgramElementInfo(t *testing.T) {
	for i, tc := range []struct {
		b           []byte
		streamType  byte
		elemPID     PID
		esLen       int
		descriptors []Descriptor
	}{
		{
			b:          []byte{0x1B, 0xE1, 0x81, 0xF0, 0x03, 0x52, 0x01, 0x81},
			streamType: 0x1B,
			elemPID:    0x0181,
			esLen:      3,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x81},
			},
		},
		{
			b:          []byte{0x0F, 0xE1, 0x82, 0xF0, 0x03, 0x52, 0x01, 0x83},
			streamType: 0x0F,
			elemPID:    0x0182,
			esLen:      3,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x83},
			},
		},
		{
			b:          []byte{0x06, 0xE1, 0x84, 0xF0, 0x08, 0x52, 0x01, 0x87, 0xFD, 0x03, 0x00, 0x12, 0xAD},
			streamType: 0x06,
			elemPID:    0x0184,
			esLen:      8,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x87},
				{0xFD, 0x03, 0x00, 0x12, 0xAD},
			},
		},
		{
			b:          []byte{0x0D, 0xF0, 0x30, 0xF0, 0x0F, 0x52, 0x01, 0x80, 0xFD, 0x0A, 0x00, 0x0D, 0x3F, 0x2F, 0x00, 0x0C, 0x00, 0x00, 0xFF, 0xBF},
			streamType: 0x0D,
			elemPID:    0x1030,
			esLen:      15,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x80},
				{0xFD, 0x0A, 0x00, 0x0D, 0x3F, 0x2F, 0x00, 0x0C, 0x00, 0x00, 0xFF, 0xBF},
			},
		},
		{
			b:          []byte{0x0D, 0xF0, 0x39, 0xF0, 0x03, 0x52, 0x01, 0x89},
			streamType: 0x0D,
			elemPID:    0x1039,
			esLen:      3,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x89},
			},
		},
		{
			b:          []byte{0x0D, 0xF0, 0x3A, 0xF0, 0x03, 0x52, 0x01, 0x8A},
			streamType: 0x0D,
			elemPID:    0x103A,
			esLen:      3,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x8A},
			},
		},
		{
			b:          []byte{0x0D, 0xF0, 0x3B, 0xF0, 0x0A, 0x52, 0x01, 0x8B, 0xFD, 0x05, 0x00, 0x0D, 0x1F, 0xFF, 0xBF},
			streamType: 0x0D,
			elemPID:    0x103B,
			esLen:      10,
			descriptors: []Descriptor{
				{0x52, 0x01, 0x8B},
				{0xFD, 0x05, 0x00, 0x0D, 0x1F, 0xFF, 0xBF},
			},
		},
	} {
		i, tc := i, tc
		t.Run("", func(t *testing.T) {
			t.Parallel()

			info := ProgramElementInfo(tc.b)
			if info.StreamType() != tc.streamType {
				t.Errorf("%0d: ProgramElementInfo(0x%04X).StreamType() => 0x%04X, want 0x%04X", i, tc.b, info.StreamType(), tc.streamType)
			}
			if info.ElementaryPID() != tc.elemPID {
				t.Errorf("%0d: ProgramElementInfo(0x%04X).ElementaryPID() => 0x%04X, want 0x%04X", i, tc.b, info.ElementaryPID(), tc.elemPID)
			}
			if info.ESInfoLength() != tc.esLen {
				t.Errorf("%0d: ProgramElementInfo(0x%04X).ESInfoLength() => %d, want %d", i, tc.b, info.ESInfoLength(), tc.esLen)
			}
			if len(info.Descriptors()) != len(tc.descriptors) {
				t.Errorf("%0d: ProgramElementInfo(0x%04X).Descriptors() => len: %d, want %d", i, tc.b, len(info.Descriptors()), len(tc.descriptors))
			} else {
				for j, exp := range tc.descriptors {
					got := info.Descriptors()[j]
					if !bytes.Equal(got, exp) {
						t.Errorf("%0d: ProgramElementInfo(0x%04X).Descriptors()[%d] => 0x%04X, want 0x%04X", i, tc.b, j, got, exp)
					}
				}
			}
		})
	}
}
