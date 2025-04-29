package secret

var (
	prefix1 = "$|"
	prefix2 = "$:"

	prefix = "$"
	suffix = "$"
)

type Secret []byte

func (k *Secret) Bytes() []byte {
	if k == nil {
		return nil
	}

	return *k
}

func (k *Secret) String() string {
	if k == nil {
		return ""
	}

	return string(*k)
}

func (k *Secret) StringPointer() *string {
	var p string
	if k == nil {
		return &p
	}

	p = string(*k)
	return &p
}

func (k *Secret) Empty() bool {
	return k == nil || len(*k) == 0
}

func (k *Secret) NotEmpty() bool {
	return k != nil && len(*k) > 0
}

func FromString(raw string) (secret *Secret) {
	sec := Secret(raw)
	return &sec
}

func FromBytes(raw []byte) (secret *Secret) {
	sec := Secret(raw)
	return &sec
}

func Parse(raw string) (secret *Secret, err error) {
	if sec, err := newSecret(raw); err != nil {
		return nil, err
	} else {
		secret = &sec
	}
	return
}
