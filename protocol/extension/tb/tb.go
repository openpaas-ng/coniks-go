// Package tb contains
// This would be used (copy from tb.go)
package tb

import (
	"bytes"

	"github.com/coniks-sys/coniks-go/protocol/extension/ext"
)

func init() {
	ext.RegisterPromiseSigner(PromiseSignerID, New)
}

type temporaryBinding struct {
	Index     []byte
	Value     []byte
	Signature []byte
}

const PromiseSignerID = "TemporaryBinding"

type TBPool struct {
	tbs map[string]*temporaryBinding
}

var _ ext.PromiseSigner = (*TBPool)(nil)

func New() (ext.PromiseSigner, error) {
	return &TBPool{
		tbs: make(map[string]*temporaryBinding),
	}, nil
}

func (signer *TBPool) Issue(uname string, index, value, sig []byte) ext.Promise {
	tb := &temporaryBinding{
		Index:     index,
		Value:     value,
		Signature: sig,
	}
	signer.tbs[uname] = tb
	return tb
}

func (signer *TBPool) Lookup(uname string) (ext.Promise, error) {
	if tb, ok := signer.tbs[uname]; ok {
		return tb, nil
	}
	return nil, ext.ErrPromiseNotFound
}

func (signer *TBPool) Clear() {
	for key := range signer.tbs {
		delete(signer.tbs, key)
	}
}

func (tb *temporaryBinding) Serialize(strSig []byte) []byte {
	var tbBytes []byte
	tbBytes = append(tbBytes, strSig...)
	tbBytes = append(tbBytes, tb.Index...)
	tbBytes = append(tbBytes, tb.Value...)
	return tbBytes
}

func (tb *temporaryBinding) Verify(index, value []byte) bool {
	return bytes.Equal(tb.Index, index) &&
		(value != nil && bytes.Equal(tb.Value, value))
}
