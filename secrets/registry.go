package secrets

import "time"

type Encryption interface {
	CreateNonce() ([]byte, error)
	Decrypt([]byte) ([]byte, error)
	Encrypt([]byte) ([]byte, error)
	Close()
}

type Registry struct {
	encrypterMap map[string]Encryption
}

// TODO(sszuecs): get rid of global
var r *Registry

// NewRegistry returns a singleton Registry
func NewRegistry() *Registry {
	if r != nil {
		return r
	}
	e := make(map[string]Encryption)
	r = &Registry{
		encrypterMap: e,
	}
	return r
}

func (r *Registry) NewEncrypter(refreshInterval time.Duration, file string) (Encryption, error) {
	if e, ok := r.encrypterMap[file]; ok {
		return e, nil
	}

	e, err := newEncrypter(file)
	if err != nil {
		return nil, err
	}

	if refreshInterval > 0 {
		err := e.runCipherRefresher(refreshInterval)
		if err != nil {
			return nil, err
		}

	}
	r.encrypterMap[file] = e
	return e, nil
}

// Close will close all Encryption of the Registry
func (r *Registry) Close() {
	for _, v := range r.encrypterMap {
		v.Close()
	}
}
