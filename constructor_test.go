package getdoc

import (
	"bytes"
	"io/ioutil"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConstructor(t *testing.T) {
	data, err := ioutil.ReadFile(path.Join("_testdata", "constructor.html"))
	if err != nil {
		t.Fatal(err)
	}

	v, err := ParseConstructor(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Constructor{
		Name:        "userProfilePhoto",
		Description: []string{"User profile photo."},
		Fields: map[string]string{
			"dc_id":       "DC ID where the photo is stored",
			"flags":       "Flags, see TL conditional fields",
			"has_video":   "Whether an animated profile picture is available for this user",
			"photo_big":   "Location of the file, corresponding to the big profile photo thumbnail",
			"photo_id":    "Identifier of the respective photoParameter added in Layer 2",
			"photo_small": "Location of the file, corresponding to the small profile photo thumbnail",
		},
	}
	require.Equal(t, expected, v)
}
