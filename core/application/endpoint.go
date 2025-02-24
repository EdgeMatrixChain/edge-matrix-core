package application

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix-core/core/application/proof"
	"github.com/emc-protocol/edge-matrix-core/core/application/proof/helper"
	"github.com/emc-protocol/edge-matrix-core/core/crypto"
	"github.com/emc-protocol/edge-matrix-core/core/helper/rpc"
	"github.com/emc-protocol/edge-matrix-core/core/types"
	"github.com/hashicorp/go-hclog"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	ProtoTagEcApp = "/em-app"
)

const (
	DefaultAppStatusSyncDuration = 30 * time.Second
)

type Endpoint struct {
	logger hclog.Logger

	// gauge for measuring app capacity
	gauge slotGauge
	sync.Mutex
	nextNonce        uint64
	nonceCacheEnable bool

	name       string
	appUrl     string
	appPort    uint64
	appOrigin  string
	h          host.Host
	tag        string
	listener   net.Listener
	httpClient *rpc.FastHttpClient
	signer     proof.Signer
	privateKey *ecdsa.PrivateKey
	address    types.Address
	stream     *eventStream // Event subscriptions

	application *Application

	httpHandler *EndpointHandler
}

// SubscribeEvents returns a application event subscription
func (e *Endpoint) SubscribeEvents() Subscription {
	return e.stream.subscribe()
}

func (e *Endpoint) getID() peer.ID {
	return e.h.ID()
}

func (e *Endpoint) getProtocolOption() p2phttp.Option {
	return p2phttp.ProtocolOption(protocol.ID(e.tag))
}

func (e *Endpoint) Close() {
	e.listener.Close()
	e.h.Close()
}

// SetSigner sets the signer the endpoint will use
// to validate an edge call response's signature.
func (e *Endpoint) SetSigner(s proof.Signer) {
	e.signer = s
}

func (e *Endpoint) GetEndpointApplication() *Application {
	return e.application
}

func (e *Endpoint) GetHandlerList() []string {
	keys := make([]string, 0, len(e.httpHandler.routes))
	for k := range e.httpHandler.routes {
		keys = append(keys, k)
	}
	return keys
}

func (e *Endpoint) AddHandler(url string, handler func(w http.ResponseWriter, r *http.Request)) {
	e.httpHandler.AddHandler(url, handler)
}

func (e *Endpoint) SetAppOrigin(appOrigin string) {
	e.appOrigin = appOrigin
	e.application.AppOrigin = appOrigin
}

func NewApplicationEndpoint(
	logger hclog.Logger,
	privateKey *ecdsa.PrivateKey,
	srvHost host.Host,
	name string,
	appUrl string,
	appPort uint64,
	version string) (*Endpoint, error) {
	endpoint := &Endpoint{
		logger:           logger.Named("app_endpoint"),
		name:             name,
		appUrl:           appUrl,
		appPort:          appPort,
		appOrigin:        "",
		h:                srvHost,
		tag:              ProtoTagEcApp,
		stream:           &eventStream{},
		nonceCacheEnable: false,
		httpHandler:      &EndpointHandler{routes: make(map[string]func(w http.ResponseWriter, r *http.Request))},
	}
	endpoint.httpClient = rpc.NewDefaultHttpClient()
	listener, err := gostream.Listen(srvHost, ProtoTagEcApp)
	if err != nil {
		return nil, err
	}
	endpoint.listener = listener

	address, err := crypto.GetAddressFromKey(privateKey)
	if err != nil {
		endpoint.logger.Error("unable to extract key, error: %v", err.Error())
		return nil, err
	}
	endpoint.address = address
	endpoint.privateKey = privateKey
	// Push the initial event to the stream
	endpoint.stream.push(&Event{})

	// init application metric
	mac, _ := helper.GetLocalMac()
	endpoint.application = &Application{
		Name:        name,
		PeerID:      srvHost.ID(),
		StartupTime: uint64(time.Now().UnixMilli()),
		Uptime:      0,
		AppOrigin:   "",
		GuageHeight: 0,
		GuageMax:    200,
		Mac:         mac,
		CpuInfo:     helper.GetCpuInfo(),
		GpuInfo:     helper.GetGpuInfo(),
		MemInfo:     helper.GetMemInfo(),
		Version:     version,
	}

	// update app status
	go func() {
		ticker := time.NewTicker(DefaultAppStatusSyncDuration)
		defer ticker.Stop()
		for {
			<-ticker.C
			event := &Event{}
			endpoint.application.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			endpoint.application.MemInfo = helper.GetMemInfo()
			endpoint.application.GpuInfo = helper.GetGpuInfo()

			event.AddNewApp(endpoint.application)
			endpoint.stream.push(event)
			endpoint.logger.Debug("endpoint----> status", "AppOrigin", endpoint.application.AppOrigin, "Mac", endpoint.application.Mac, "CpuInfo", endpoint.application.CpuInfo, "GpuInfo", endpoint.application.GpuInfo, "MemInfo", endpoint.application.MemInfo)
		}
	}()

	go func() {
		endpoint.httpHandler.AddHandler("/health", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var health struct {
				PeerID      string `json:"peerId"`
				Uptime      uint64 `json:"uptime"`
				StartupTime uint64 `json:"startupTime"`
				Version     string `json:"version"`
				Time        string `json:"time"`
			}
			health.PeerID = endpoint.application.PeerID.String()
			health.Version = endpoint.application.Version
			health.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			health.StartupTime = endpoint.application.StartupTime
			health.Time = time.Now().String()
			resp, err := json.Marshal(health)
			if err != nil {
				resp = []byte("err: " + err.Error())
			}

			w.Write(resp)
		})

		endpoint.httpHandler.AddHandler("/info", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var infoObj struct {
				Name        string `json:"name"`
				PeerID      string `json:"peerId"`
				Uptime      uint64 `json:"uptime"`
				StartupTime uint64 `json:"startupTime"`
				Version     string `json:"version"`
				Tag         string `json:"tag"`
				// AI model hash string
				ModelHash string `json:"model_hash"`
				// mac addr
				Mac string `json:"mac"`
				// memory info
				MemInfo string `json:"mem_info"`
				// cpu info
				CpuInfo string `json:"cpu_info"`
				// average e power
				AveragePower float32 `json:"average_power"`
				// gpu info
				GpuInfo string `json:"gpu_info"`
			}
			infoObj.PeerID = endpoint.application.PeerID.String()
			infoObj.Version = endpoint.application.Version
			infoObj.Tag = endpoint.application.AppOrigin
			infoObj.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
			infoObj.StartupTime = endpoint.application.StartupTime
			infoObj.Name = endpoint.application.Name
			infoObj.CpuInfo = endpoint.application.CpuInfo
			infoObj.GpuInfo = endpoint.application.GpuInfo
			infoObj.MemInfo = endpoint.application.MemInfo
			infoObj.Mac = endpoint.application.Mac
			infoObj.ModelHash = endpoint.application.ModelHash
			infoObj.AveragePower = endpoint.application.AveragePower

			info, err := json.Marshal(infoObj)
			if err != nil {
				info = []byte("endpoint err: " + err.Error())
			}

			WriteSignedResponse(w, info, endpoint)
		})

		http.Handle("/", endpoint.httpHandler)
		server := &http.Server{}
		server.Serve(listener)
	}()

	return endpoint, nil
}

func WriteSignedResponse(w http.ResponseWriter, info []byte, endpoint *Endpoint) {
	resp := base64.StdEncoding.EncodeToString(info)
	edgeResp := &proof.EdgeResponse{
		RespString: resp,
	}
	endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(edgeResp.RespString)))

	signedResp, err := endpoint.signer.SignEdgeResp(edgeResp, endpoint.privateKey)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
	provider, err := endpoint.signer.Provider(signedResp)
	if err != nil {
		return
	}
	signedResp.From = provider
	signedResp.Hash = endpoint.signer.Hash(edgeResp)

	w.Write(signedResp.MarshalRLP())
}
