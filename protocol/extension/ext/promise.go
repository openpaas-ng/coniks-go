package ext

import "errors"

type Promise interface{}

type newPromiseSigner func() (PromiseSigner, error)

// PromiseSigner is a generic interface for registration promises.
type PromiseSigner interface {
	Issue(string, []byte, []byte, []byte) Promise
	Lookup(string) (Promise, error)
	Clear()
}

var (
	ErrPromiseNotFound = errors.New("[ext-promise] Promise not found")
)

var registeredPromiseSigners map[string]newPromiseSigner

// RegisterPromiseSigner is must to be called to register a promise signer.
func RegisterPromiseSigner(name string, ps newPromiseSigner) {
	if registeredPromiseSigners == nil {
		registeredPromiseSigners = make(map[string]newPromiseSigner)
	}
	registeredPromiseSigners[name] = ps
}

// NewPromiseSigner returns a registerd promise signer.
func NewPromiseSigner(name string) (PromiseSigner, error) {
	ps, ok := registeredPromiseSigners[name]
	if !ok {
		return nil, ErrExtensionNotFound
	}
	psInst, err := ps()
	if err != nil {
		return nil, err
	}
	return psInst, nil
}
