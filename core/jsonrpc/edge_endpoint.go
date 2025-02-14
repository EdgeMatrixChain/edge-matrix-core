package jsonrpc

import (
	"fmt"
	"github.com/emc-protocol/edge-matrix-core/core/types"
	"github.com/hashicorp/go-hclog"
)

//type edgeTelePoolStore interface {
//	// AddTele adds a new telegram to the telegram pool
//	AddTele(tx *types.Telegram) (string, error)
//
//	// GetPendingTx gets the pending transaction from the transaction pool, if it's present
//	GetPendingTele(txHash types.Hash) (*types.Telegram, bool)
//
//	// GetNonce returns the next nonce for this address
//	GetNonce(addr types.Address) uint64
//}

type edgeTelePoolStore interface {
	// AddTele adds a new telegram to the telegram pool
	AddTele(tx *types.Telegram) (string, error)
}

// edgeStore provides access to the methods needed by edge endpoint
type edgeStore interface {
	edgeTelePoolStore
}

// Edge is the edge jsonrpc endpoint
type Edge struct {
	logger        hclog.Logger
	store         edgeStore
	chainID       uint64
	filterManager *NodeFilterManager
}

// ChainId returns the chain id of the client
//
//nolint:stylecheck
func (e *Edge) ChainId() (interface{}, error) {
	return argUintPtr(e.chainID), nil
}

// SendRawTelegram sends a raw telegram
func (e *Edge) SendRawTelegram(buf argBytes) (interface{}, error) {
	tele := &types.Telegram{}
	if err := tele.UnmarshalRLP(buf); err != nil {
		return nil, err
	}
	e.logger.Debug(fmt.Sprintf("SendRawTelegram To: %s, Nonce:%d", tele.To.String(), tele.Nonce))
	tele.ComputeHash()

	teleResp, teleErr := e.store.AddTele(tele)
	if teleErr != nil {
		return nil, teleErr
	}
	resp := fmt.Sprintf(`{"telegram_hash":"%s","response":"%s"}`, tele.Hash.String(), teleResp)
	return resp, nil
}

// GetFilterChanges is a polling method for a filter, which returns an array of logs which occurred since last poll.
func (e *Edge) GetFilterChanges(id string) (interface{}, error) {
	return e.filterManager.GetFilterChanges(id)
}

// UninstallFilter uninstalls a filter with given ID
func (e *Edge) UninstallFilter(id string) (bool, error) {
	return e.filterManager.Uninstall(id), nil
}

// Unsubscribe uninstalls a filter in a websocket
func (e *Edge) Unsubscribe(id string) (bool, error) {
	return e.filterManager.Uninstall(id), nil
}
