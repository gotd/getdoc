package getdoc

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseType(t *testing.T) {
	data, err := ioutil.ReadFile(path.Join("_testdata", "type.html"))
	if err != nil {
		t.Fatal(err)
	}

	v, err := ParseType(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Type{
		Name:        "messages.Chats",
		Description: []string{"Object contains list of chats with auxiliary data."},
	}
	require.Equal(t, expected, v)
}
