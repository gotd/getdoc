package getdoc

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseMethod(t *testing.T) {
	data, err := os.ReadFile(path.Join("_testdata", "method.html"))
	if err != nil {
		t.Fatal(err)
	}

	v, err := ParseMethod(bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}

	expected := &Method{
		Name:        "langpack.getDifference",
		Description: []string{"Get new strings in languagepack"},
		Parameters: map[string]ParamDescription{
			"from_version": {
				Name:        "from_version",
				Description: "Previous localization pack version",
			},
			"lang_code": {
				Name:        "lang_code",
				Description: "Language code",
			},
			"lang_pack": {
				Name:        "lang_pack",
				Description: "Language pack",
			},
		},
		Errors: []Error{
			{Code: 400, Type: "LANG_PACK_INVALID", Description: "The provided language pack is invalid"},
		},
	}
	require.Equal(t, expected, v)
}
