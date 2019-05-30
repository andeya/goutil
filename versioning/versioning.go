// Package versioning is a version comparison tool that conforms to semantic version 2.0.0
//
// Copyright 2019 henrylee2cn Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package versioning

import (
	"errors"
	"strconv"
)

// SemVer semantic version object
// via https://semver.org/
type SemVer struct {
	major    string
	minor    string
	patch    string
	metadata string
	nums     [3]uint32
}

// Create creates a semantic version object.
func Create(major, minor, patch uint32, metadata string) *SemVer {
	return &SemVer{
		major:    uint32ToString(major),
		minor:    uint32ToString(minor),
		patch:    uint32ToString(patch),
		metadata: metadata,
		nums: [3]uint32{
			major, minor, patch,
		},
	}
}

// Parse parses the semantic version string to object.
// NOTE:
// If metadata part exists, the separator must not be a number.
func Parse(semVer string) (*SemVer, error) {
	a := [4][]rune{}
	var i int
	for _, r := range semVer {
		switch {
		case i == 3 || (r >= '0' && r <= '9'):
			a[i] = append(a[i], r)
		case r == '.':
			if i < 3 {
				i++
			}
		default:
			i = 3
			a[i] = append(a[i], r)
		}
	}
	for _, s := range a[:3] {
		if len(s) == 0 {
			return nil, errors.New("invalid semantic version 2: " + semVer)
		}
	}
	return &SemVer{
		major:    string(a[0]),
		minor:    string(a[1]),
		patch:    string(a[2]),
		metadata: string(a[3]),
		nums: [3]uint32{
			runeToUint32(a[0]),
			runeToUint32(a[1]),
			runeToUint32(a[2]),
		},
	}, nil
}

// Compare compares 'a' and 'b'.
// The result will be 0 if a==b, -1 if a < b, and +1 if a > b.
// If compareMetadata==nil, will not compare metadata.
func Compare(a, b string, compareMetadata func(aMeta, bMeta string) int) (int, error) {
	ver1, err := Parse(a)
	if err != nil {
		return 0, err
	}
	ver2, err := Parse(b)
	if err != nil {
		return 0, err
	}
	return ver1.Compare(ver2, compareMetadata), nil
}

// Compare compares whether 's' and 'semVer'.
// The result will be 0 if s==semVer, -1 if s < semVer, and +1 if s > semVer.
// If compareMetadata==nil, will not compare metadata.
func (s *SemVer) Compare(semVer *SemVer, compareMetadata func(sMeta, semVerMeta string) int) int {
	for k, v := range s.nums {
		v2 := semVer.nums[k]
		if v < v2 {
			return -1
		} else if v > v2 {
			return 1
		}
	}
	if compareMetadata != nil {
		return compareMetadata(s.Metadata(), semVer.Metadata())
	}
	return 0
}

// Major returns the version major.
func (s *SemVer) Major() string {
	return s.major
}

// Minor returns the version minor.
func (s *SemVer) Minor() string {
	return s.minor
}

// Patch returns the version patch.
func (s *SemVer) Patch() string {
	return s.patch
}

// Metadata returns the version metadata.
// Examples:
//  1.0.0-alpha+001 => -alpha+001
//  1.0.0+20130313144700 => +20130313144700
//  1.0.0-beta+exp.sha.5114f85 => -beta+exp.sha.5114f85
//  1.0.0rc => rc
func (s *SemVer) Metadata() string {
	return s.metadata
}

// String returns the version string.
func (s *SemVer) String() string {
	var ver = s.major
	if s.minor != "" {
		ver += "." + s.minor
	}
	if s.patch != "" {
		ver += "." + s.patch
	}
	return ver + s.metadata
}

func uint32ToString(u uint32) string {
	return strconv.FormatUint(uint64(u), 10)
}

func runeToUint32(r []rune) uint32 {
	var p = 1
	var u int
	for i := len(r) - 1; i >= 0; i-- {
		u += int(r[i]-'0') * p
		p *= 10
	}
	return uint32(u)
}
