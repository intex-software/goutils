package internal

import (
	"encoding/base32"
	"encoding/base64"
)

var Base32 = base32.StdEncoding.WithPadding(base32.NoPadding)
var Base64 = base64.StdEncoding
