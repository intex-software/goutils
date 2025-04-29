package secret

import (
	"encoding/json"
	"strings"

	"github.com/intex-software/goutils/obfuscate"
)

func (k *Secret) MarshalJSON() (data []byte, err error) {
	secret, err := k.Marshal()
	if err != nil {
		return
	}
	return json.Marshal(secret)
}

func (k *Secret) MarshalYAML() (data any, err error) {
	secret, err := k.Marshal()
	if err != nil {
		return
	}

	data = strings.Join(split(secret, 64), "\n")
	return
}

func (k *Secret) Marshal() (raw string, err error) {
	raw = obfuscate.ObfuscateString(*k)
	raw = prefix + raw + suffix
	return
}
