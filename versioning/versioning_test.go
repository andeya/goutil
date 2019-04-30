package versioning

import (
	"reflect"
	"testing"
)

func Test(t *testing.T) {
	cases := []struct {
		ver    string
		semVer SemVer
	}{
		{"1.0.0-alpha.1",
			SemVer{
				"1",
				"0",
				"0",
				"-alpha.1",
				[3]uint32{1, 0, 0},
			},
		},
		{"1.0.2-alpha",
			SemVer{
				"1",
				"0",
				"2",
				"-alpha",
				[3]uint32{1, 0, 2},
			},
		},
		{"1.0.0+20130313144700",
			SemVer{
				"1",
				"0",
				"0",
				"+20130313144700",
				[3]uint32{1, 0, 0},
			},
		},
		{"1.0.0rc",
			SemVer{
				"1",
				"0",
				"0",
				"rc",
				[3]uint32{1, 0, 0},
			},
		},
	}

	for _, c := range cases {
		semVer := Parse(c.ver)
		if !reflect.DeepEqual(*semVer, c.semVer) {
			t.Fatal(c.ver)
		}
		if semVer.String() != c.ver || c.semVer.Compare(semVer) != 0 {
			t.Fatalf("expect:%s, got:%s", c.ver, semVer.String())
		}
	}
}
