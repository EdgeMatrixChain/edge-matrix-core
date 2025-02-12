package application

import (
	"crypto/ecdsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	appAgent "github.com/emc-protocol/edge-matrix-core/core/application/proof/agent"
	"github.com/emc-protocol/edge-matrix-core/core/application/proof/helper"
	"github.com/emc-protocol/edge-matrix-core/core/crypto"
	"github.com/emc-protocol/edge-matrix-core/core/helper/rpc"
	"github.com/emc-protocol/edge-matrix-core/core/types"
	"github.com/emc-protocol/edge-matrix-core/core/versioning"
	"github.com/hashicorp/go-hclog"
	gostream "github.com/libp2p/go-libp2p-gostream"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
	"math/rand"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	// proto tag for p2phttp
	ProtoTagEcApp = "/em-app"
)

const (
	txSlotSize = 32 * 1024 // 32kB
)

const (
	DefaultAppStatusSyncDuration = 15 * time.Second
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
	appOrigin  string
	h          host.Host
	tag        string
	listener   net.Listener
	httpClient *rpc.FastHttpClient
	//signer     Signer
	privateKey *ecdsa.PrivateKey
	address    types.Address
	stream     *eventStream // Event subscriptions

	application *Application

	randomNum int

	latestBlockHeadHash string
	latestBlockNum      uint64

	isEdgeMode bool
}

// SubscribeEvents returns a application event subscription
func (b *Endpoint) SubscribeEvents() Subscription {
	return b.stream.subscribe()
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

// SetSigner sets the signer the endpint will use
// to validate a edge call response's signature.
//func (e *Endpoint) SetSigner(s Signer) {
//	e.signer = s
//}

func (e *Endpoint) GetEndpointApplication() *Application {
	return e.application
}

func (e *Endpoint) doAppNodeBind() error {
	agent := appAgent.NewAppAgent(e.appUrl)
	err := agent.BindAppNode(e.h.ID().String())
	if err != nil {
		return err
	}
	return nil
}

func (e *Endpoint) getAppOrigin() (error, string) {
	agent := appAgent.NewAppAgent(e.appUrl)
	err, appOrigin := agent.GetAppOrigin()
	if err != nil {
		return err, ""
	}
	return nil, appOrigin
}

func (e *Endpoint) getAppIdl() (error, string) {
	agent := appAgent.NewAppAgent(e.appUrl)
	err, appOrigin := agent.GetAppOrigin()
	if err != nil {
		return err, ""
	}
	return nil, appOrigin
}

func (e *Endpoint) validAppNode() (error, bool) {
	agent := appAgent.NewAppAgent(e.appUrl)
	err, nodeId := agent.GetAppNode()
	if err != nil {
		return err, false
	}
	if e.h.ID().String() == nodeId {
		return nil, true
	}
	return nil, false
}

func NewApplicationEndpoint(
	logger hclog.Logger,
	privateKey *ecdsa.PrivateKey,
	srvHost host.Host,
	name string,
	appUrl string,
	isEdgeMode bool) (*Endpoint, error) {
	endpoint := &Endpoint{
		logger:              logger.Named("app_endpoint"),
		name:                name,
		appUrl:              appUrl,
		appOrigin:           "",
		h:                   srvHost,
		tag:                 ProtoTagEcApp,
		stream:              &eventStream{},
		nonceCacheEnable:    false,
		latestBlockHeadHash: "",
		latestBlockNum:      0,
		isEdgeMode:          isEdgeMode,
	}
	rand.Seed(time.Now().Unix())
	endpoint.randomNum = rand.Intn(1000)
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
		Version:     versioning.Version + " Build" + versioning.Build,
	}

	// check app status
	if isEdgeMode {
		go func() {
			ticker := time.NewTicker(DefaultAppStatusSyncDuration)
			for {
				<-ticker.C
				event := &Event{}
				// bind app node
				err := endpoint.doAppNodeBind()
				if err != nil {
					endpoint.logger.Error("doAppNodeBind", "err", err.Error())
				}

				err, appOrigin := endpoint.getAppOrigin()
				if err != nil {
					endpoint.logger.Error("getAppOrigin", "err", err.Error())
				}

				endpoint.application.AppOrigin = appOrigin
				endpoint.application.Uptime = uint64(time.Now().UnixMilli()) - endpoint.application.StartupTime
				endpoint.application.MemInfo = helper.GetMemInfo()
				endpoint.application.GpuInfo = helper.GetGpuInfo()

				event.AddNewApp(endpoint.application)
				endpoint.stream.push(event)
				endpoint.logger.Debug("endpoint----> status", "AppOrigin", endpoint.application.AppOrigin, "Mac", endpoint.application.Mac, "CpuInfo", endpoint.application.CpuInfo, "GpuInfo", endpoint.application.GpuInfo, "MemInfo", endpoint.application.MemInfo)
			}
			ticker.Stop()
		}()
	}

	go func() {
		http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			var infoObj struct {
				Name        string `json:"name"`
				PeerID      string `json:"peerId"`
				Uptime      uint64 `json:"uptime"`
				StartupTime uint64 `json:"startupTime"`
				Version     string `json:"version"`
				Tag         string `json:"tag"`
				// ai model hash string
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

			info := make([]byte, 0)
			info, err := json.Marshal(infoObj)
			if err != nil {
				info = []byte("endpoint err: " + err.Error())
			}

			writeResponse(w, info, endpoint)
		})

		http.HandleFunc("/alive", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			resp := fmt.Sprintf("%s", time.Now().String())
			w.Write([]byte(resp))
		})

		http.HandleFunc("/idl", func(w http.ResponseWriter, r *http.Request) {
			defer r.Body.Close()
			err, appIdl := endpoint.getAppIdl()
			if err != nil {
				// TODO Fetch idl json text through GET #{appUrl}/getAppIdl
				endpoint.logger.Debug(fmt.Sprintf("/getAppIdl =>resp: %s", err.Error()))
				idlData, err := os.ReadFile("idl.json")
				if nil != err {
					idlData = []byte("[]")
				}
				writeResponse(w, idlData, endpoint)
			} else {
				if len(appIdl) > 0 {
					writeResponse(w, []byte(appIdl), endpoint)
				} else {
					writeResponse(w, []byte("[]"), endpoint)
				}
			}
		})

		server := &http.Server{}
		server.Serve(listener)
	}()

	return endpoint, nil
}

func writeResponse(w http.ResponseWriter, data []byte, endpoint *Endpoint) {
	resp := base64.StdEncoding.EncodeToString(data)

	endpoint.logger.Debug(fmt.Sprintf("/api =>resp size: %d", len(resp)))

	w.Write(data)
}
