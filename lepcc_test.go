package lepcc

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestLerc(t *testing.T) {

	f, err := os.Open("./testdata/0.bin.pccxyz")
	defer f.Close()

	if err != nil {
		t.Error("error")
	}
	src, err := ioutil.ReadAll(f)
	if err != nil {
		t.Error("error")
	}

	ctx := NewContext()

	xyz, err := DecodeXYZ(ctx, src)
	if err != nil {
		t.Error("error")
	}

	if len(xyz) == 0 {
		t.Error("error")
	}

}
