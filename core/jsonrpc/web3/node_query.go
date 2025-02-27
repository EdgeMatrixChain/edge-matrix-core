package web3

import (
	"encoding/json"
	"github.com/emc-protocol/edge-matrix-core/core/application"
)

// NodeQuery is a query to filter node
type NodeQuery struct {
	Name    string `json:"name"`
	Tag     string `json:"tag"`
	Id      string `json:"id"`
	Version string `json:"version"`
}

func DecodeNodeQueryFromInterface(i interface{}) (*NodeQuery, error) {
	// once the node filter is decoded as map[string]interface we cannot use unmarshal json
	raw, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}

	query := &NodeQuery{}
	if err := json.Unmarshal(raw, &query); err != nil {
		return nil, err
	}

	return query, nil
}

func (q *NodeQuery) Match(rm *application.Application) bool {
	if q.Tag != "" {
		match := false
		if q.Tag == rm.Tag {
			match = true
		}
		if !match {
			return false
		}
	}
	// check name
	if q.Name != "" {
		match := false
		if rm.Name == q.Name {
			match = true
		}

		if !match {
			return false
		}
	}

	if q.Id != "" {
		match := false
		if rm.PeerID.String() == q.Id {
			match = true
		}

		if !match {
			return false
		}
	}

	if q.Version != "" {
		match := false
		if rm.Version == q.Version {
			match = true
		}

		if !match {
			return false
		}
	}

	return true
}
