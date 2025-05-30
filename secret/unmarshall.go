package secret

import (
	"encoding/json"
	"math"
	"strings"

	"github.com/intex-software/goutils/internal"
	"github.com/intex-software/goutils/obfuscate"
	"gopkg.in/yaml.v3"
)

func (k *Secret) UnmarshalJSON(data []byte) (err error) {
	raw := string(*k)
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	*k, err = newSecret(raw)
	return
}

func (k *Secret) UnmarshalYAML(data *yaml.Node) (err error) {
	raw, err := unbinary(data)
	if err != nil {
		return
	}
	*k, err = newSecret(raw)
	return
}

func newSecret(raw string) (secret Secret, err error) {
	if strings.HasPrefix(raw, prefix2) && strings.HasSuffix(raw, suffix) {
		raw = raw[len(prefix2) : len(raw)-len(suffix)]
	} else if strings.HasPrefix(raw, prefix1) && strings.HasSuffix(raw, suffix) {
		raw = raw[len(prefix1) : len(raw)-len(suffix)]
	} else if strings.HasPrefix(raw, prefix) && strings.HasSuffix(raw, suffix) {
		raw = raw[len(prefix) : len(raw)-len(suffix)]
	} else {
		secret = Secret(raw)
		return
	}

	bytes, err := obfuscate.DeobfuscateString(raw)
	if err != nil {
		return nil, err
	}
	secret = Secret(bytes)

	return
}

func unbinary(data *yaml.Node) (value string, err error) {
	value = data.Value

	if data.ShortTag() == "!!binary" {
		if binary, e := internal.Base64.DecodeString(value); e != nil {
			err = e
			return
		} else {
			value = string(binary)
		}
	}

	return
}

func split(s string, width int) (result []string) {
	length := len(s)

	slices := math.Ceil(float64(length) / float64(width))
	result = make([]string, 0, int(slices))

	for start, end := 0, min(width, length); start < length; start, end = start+width, min(end+width, length) {
		result = append(result, s[start:end])
	}
	return
}
