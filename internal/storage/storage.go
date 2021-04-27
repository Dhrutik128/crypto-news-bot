package storage

type Storable interface {
	Key() []byte
}
