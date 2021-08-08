package storage

import (
	"io"
	"os"
	"unsafe"
)

type fstoreSize = uint16

var fstoreSS = unsafe.Sizeof(fstoreSize(0))

type FStoreValue struct {
	Key []byte
	Val []byte
}

type FStore struct {
	f    *os.File
	data []FStoreValue
}

func OpenFStore(file string) (*FStore, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_RDWR, 0777)
	if err != nil {
		return nil, err
	}

	fstore := FStore{
		f:    f,
		data: make([]FStoreValue, 0, 8),
	}

	var ss = make([]byte, fstoreSS)

	for {
		_, err = f.Read(ss[:]) // key size
		if err != nil {
			break
		}
		key := make([]byte, *(*fstoreSize)(unsafe.Pointer(&ss)))
		_, err = f.Read(key) // key value
		if err != nil {
			break
		}

		_, err = f.Read(ss[:]) // val size
		if err != nil {
			break
		}
		val := make([]byte, *(*fstoreSize)(unsafe.Pointer(&ss)))
		_, err = f.Read(val) // val value
		if err != nil {
			break
		}

		fstore.data = append(fstore.data, FStoreValue{key, val})
	}

	if err == io.EOF {
		err = nil
	}

	return &fstore, err
}
