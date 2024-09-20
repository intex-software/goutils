package secret

var (
	prefix = "$|"
	suffix = "$"
)

type Secret []byte

func (k *Secret) Bytes() []byte {
	return *k
}

func (k *Secret) String() string {
	return string(*k)
}

func (k *Secret) StringPointer() *string {
	if k == nil {
		return nil
	}

	p := string(*k)
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
