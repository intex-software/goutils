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

func (k *Secret) Pointer() *string {
	p := string(*k)
	return &p
}

func (k *Secret) Empty() bool {
	return len(*k) == 0
}
