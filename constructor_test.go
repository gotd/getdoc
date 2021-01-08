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
		Fields: map[string]ParamDescription{
			"dc_id": {
				Name:        "dc_id",
				Description: "DC ID where the photo is stored",
			},
			"flags": {
				Name:        "flags",
				Description: "Flags, see TL conditional fields¹",
				Links:       []string{"https://core.telegram.org/mtproto/TL-combinators#conditional-fields"},
			},
			"has_video": {
				Name:        "has_video",
				Description: "Whether an animated profile picture¹ is available for this user",
				Links:       []string{"https://core.telegram.org/api/files#animated-profile-pictures"},
			},
			"photo_big": {
				Name:        "photo_big",
				Description: "Location of the file, corresponding to the big profile photo thumbnail",
			},
			"photo_id": {
				Name:        "photo_id",
				Description: "Identifier of the respective photoParameter added in Layer 2¹",
				Links:       []string{"https://core.telegram.org/api/layers#layer-2"},
			},
			"photo_small": {
				Name:        "photo_small",
				Description: "Location of the file, corresponding to the small profile photo thumbnail",
			},
		},
	}
	require.Equal(t, expected, v)
}
