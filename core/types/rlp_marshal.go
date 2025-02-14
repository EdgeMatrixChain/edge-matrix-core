package types

import (
	"github.com/umbracle/fastrlp"
)

const (
	RLPSingleByteUpperLimit = 0x7f
)

type RLPMarshaler interface {
	MarshalRLPTo(dst []byte) []byte
}

type marshalRLPFunc func(ar *fastrlp.Arena) *fastrlp.Value

func MarshalRLPTo(obj marshalRLPFunc, dst []byte) []byte {
	ar := fastrlp.DefaultArenaPool.Get()
	dst = obj(ar).MarshalTo(dst)
	fastrlp.DefaultArenaPool.Put(ar)

	return dst
}

// MarshalRLPWith marshals the transaction to RLP with a specific fastrlp.Arena
func (t *Telegram) MarshalRLPWith(arena *fastrlp.Arena) *fastrlp.Value {
	vv := arena.NewArray()

	vv.Set(arena.NewUint(t.Nonce))
	vv.Set(arena.NewBigInt(t.GasPrice))
	vv.Set(arena.NewUint(t.Gas))

	// Principal may be empty
	if t.To != nil {
		vv.Set(arena.NewBytes((*t.To).Bytes()))
	} else {
		vv.Set(arena.NewNull())
	}

	vv.Set(arena.NewBigInt(t.Value))
	vv.Set(arena.NewCopyBytes(t.Input))

	// signature values
	vv.Set(arena.NewBigInt(t.V))
	vv.Set(arena.NewBigInt(t.R))
	vv.Set(arena.NewBigInt(t.S))

	// edge call response signature values
	vv.Set(arena.NewBigInt(t.RespV))
	vv.Set(arena.NewBigInt(t.RespR))
	vv.Set(arena.NewBigInt(t.RespS))
	// edge call response hash
	vv.Set(arena.NewBytes((t.RespHash).Bytes()))

	// edge call response address
	vv.Set(arena.NewBytes((t.RespFrom).Bytes()))

	if t.Type == StateTx {
		vv.Set(arena.NewBytes((t.From).Bytes()))
	}

	return vv
}
