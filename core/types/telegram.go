package types

import (
	"fmt"
	"github.com/umbracle/fastrlp"
	"math/big"
	"sync/atomic"

	"github.com/EdgeMatrixChain/edge-matrix-core/core/helper/keccak"
)

type TeleType byte

var marshalArenaPool fastrlp.ArenaPool

const (
	LegacyTx TeleType = 0x0
	StateTx  TeleType = 0x7f
)

func txTypeFromByte(b byte) (TeleType, error) {
	tt := TeleType(b)

	switch tt {
	case LegacyTx, StateTx:
		return tt, nil
	default:
		return tt, fmt.Errorf("unknown transaction type: %d", b)
	}
}

func (t TeleType) String() (s string) {
	switch t {
	case LegacyTx:
		return "LegacyTx"
	case StateTx:
		return "StateTx"
	}

	return
}

type Telegram struct {
	Nonce    uint64
	GasPrice *big.Int
	Gas      uint64
	To       *Address
	Value    *big.Int
	Input    []byte
	V        *big.Int
	R        *big.Int
	S        *big.Int

	Hash Hash
	From Address

	RespV    *big.Int
	RespR    *big.Int
	RespS    *big.Int
	RespHash Hash
	RespFrom Address

	Type TeleType

	// Cache
	size atomic.Value
}

// IsContractCreation checks if tx is contract creation
func (t *Telegram) IsContractCreation() bool {
	return t.To == nil
}

// ComputeHash computes the hash of the transaction
func (t *Telegram) ComputeHash() *Telegram {
	ar := marshalArenaPool.Get()
	hash := keccak.DefaultKeccakPool.Get()

	v := t.MarshalRLPWith(ar)
	hash.WriteRlp(t.Hash[:0], v)

	marshalArenaPool.Put(ar)
	keccak.DefaultKeccakPool.Put(hash)

	return t
}

func (t *Telegram) Copy() *Telegram {
	tt := new(Telegram)
	*tt = *t

	tt.GasPrice = new(big.Int)
	if t.GasPrice != nil {
		tt.GasPrice.Set(t.GasPrice)
	}

	tt.Value = new(big.Int)
	if t.Value != nil {
		tt.Value.Set(t.Value)
	}

	if t.R != nil {
		tt.R = new(big.Int)
		tt.R = big.NewInt(0).SetBits(t.R.Bits())
	}

	if t.S != nil {
		tt.S = new(big.Int)
		tt.S = big.NewInt(0).SetBits(t.S.Bits())
	}

	tt.Input = make([]byte, len(t.Input))
	copy(tt.Input[:], t.Input[:])

	return tt
}

// Cost returns gas * gasPrice + value
func (t *Telegram) Cost() *big.Int {
	total := new(big.Int).Mul(t.GasPrice, new(big.Int).SetUint64(t.Gas))
	total.Add(total, t.Value)

	return total
}

func (t *Telegram) Size() uint64 {
	if size := t.size.Load(); size != nil {
		sizeVal, ok := size.(uint64)
		if !ok {
			return 0
		}

		return sizeVal
	}

	size := uint64(len(t.MarshalRLP()))
	t.size.Store(size)

	return size
}

func (t *Telegram) MarshalRLP() []byte {
	return t.MarshalRLPTo(nil)
}

func (t *Telegram) MarshalRLPTo(dst []byte) []byte {
	if t.Type != LegacyTx {
		dst = append(dst, byte(t.Type))
	}

	return MarshalRLPTo(t.MarshalRLPWith, dst)
}

func (t *Telegram) UnmarshalRLP(input []byte) error {
	t.Type = LegacyTx
	offset := 0

	if len(input) > 0 && input[0] <= RLPSingleByteUpperLimit {
		var err error
		if t.Type, err = txTypeFromByte(input[0]); err != nil {
			return err
		}

		offset = 1
	}

	return UnmarshalRlp(t.unmarshalRLPFrom, input[offset:])
}

// unmarshalRLPFrom unmarshals a Transaction in RLP format
func (t *Telegram) unmarshalRLPFrom(p *fastrlp.Parser, v *fastrlp.Value) error {
	elems, err := v.GetElems()
	if err != nil {
		return err
	}

	if len(elems) < 9 {
		return fmt.Errorf("incorrect number of elements to decode transaction, expected 9 but found %d", len(elems))
	}

	p.Hash(t.Hash[:0], v)

	// nonce
	if t.Nonce, err = elems[0].GetUint64(); err != nil {
		return err
	}

	// gasPrice
	t.GasPrice = new(big.Int)
	if err = elems[1].GetBigInt(t.GasPrice); err != nil {
		return err
	}

	// gas
	if t.Gas, err = elems[2].GetUint64(); err != nil {
		return err
	}

	// to
	if vv, _ := v.Get(3).Bytes(); len(vv) == 20 {
		// address
		addr := BytesToAddress(vv)
		t.To = &addr
	} else {
		// reset To
		t.To = nil
	}

	// value
	t.Value = new(big.Int)
	if err = elems[4].GetBigInt(t.Value); err != nil {
		return err
	}

	// input
	if t.Input, err = elems[5].GetBytes(t.Input[:0]); err != nil {
		return err
	}

	// V
	t.V = new(big.Int)
	if err = elems[6].GetBigInt(t.V); err != nil {
		return err
	}

	// R
	t.R = new(big.Int)
	if err = elems[7].GetBigInt(t.R); err != nil {
		return err
	}

	// S
	t.S = new(big.Int)
	if err = elems[8].GetBigInt(t.S); err != nil {
		return err
	}

	if len(elems) > 9 {
		// edge call response signature values
		t.RespV = new(big.Int)
		if err = elems[9].GetBigInt(t.RespV); err != nil {
			return err
		}

		t.RespR = new(big.Int)
		if err = elems[10].GetBigInt(t.RespR); err != nil {
			return err
		}

		t.RespS = new(big.Int)
		if err = elems[11].GetBigInt(t.RespS); err != nil {
			return err
		}

		// edge call response hash
		if vv, _ := v.Get(12).Bytes(); len(vv) == 32 {
			// address
			respHash := BytesToHash(vv)
			t.RespHash = respHash
		} else {
			// reset To
			t.RespHash = ZeroHash
		}

		// edge call response address
		if vv, _ := v.Get(13).Bytes(); len(vv) == 20 {
			// address
			respAddr := BytesToAddress(vv)
			t.RespFrom = respAddr
		} else {
			// reset To
			t.RespFrom = ZeroAddress
		}
	}
	if t.Type == StateTx {
		// set From with default value
		t.From = ZeroAddress

		// We need to set From field for state transaction,
		// because we are using unique, predefined address, for sending such transactions
		// From
		if len(elems) > 14 {
			if vv, err := v.Get(14).Bytes(); err == nil && len(vv) == AddressLength {
				// address
				t.From = BytesToAddress(vv)
			}
		}
	}

	return nil
}
