package application

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	p2phttp "github.com/libp2p/go-libp2p-http"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/emc-protocol/edge-matrix-core/core/versioning"
	"github.com/hashicorp/go-hclog"
)

type serverType int

const (
	serverIPC serverType = iota
	serverHTTP
	serverWS
)

func (s serverType) String() string {
	switch s {
	case serverIPC:
		return "ipc"
	case serverHTTP:
		return "http"
	case serverWS:
		return "ws"
	default:
		panic("BUG: Not expected")
	}
}

// TransparentProxy is an API consensus
type TransparentProxy struct {
	logger hclog.Logger
	config *Config
}

// TransparentProxyStore defines all the methods required
// by all the proxy endpoints
type TransparentProxyStore interface {
	GetRelayHost() host.Host
	GetNetworkHost() host.Host
}

type Config struct {
	Store                    TransparentProxyStore
	Addr                     *net.TCPAddr
	NetworkID                uint64
	ChainName                string
	AccessControlAllowOrigin []string
}

// NewTransportProxy returns the TransparentProxy http server
func NewTransportProxy(logger hclog.Logger, config *Config) (*TransparentProxy, error) {
	srv := &TransparentProxy{
		logger: logger.Named("transport-proxy"),
		config: config,
	}

	// start http server
	if err := srv.setupHTTP(); err != nil {
		return nil, err
	}

	return srv, nil
}

func (j *TransparentProxy) setupHTTP() error {
	j.logger.Info("http server started", "addr", j.config.Addr.String())

	lis, err := net.Listen("tcp", j.config.Addr.String())
	if err != nil {
		return err
	}

	// NewServeMux must be used, as it disables all debug features.
	// For some strange reason, with DefaultServeMux debug/vars is always enabled (but not debug/pprof).
	// If pprof need to be enabled, this should be DefaultServeMux
	mux := http.NewServeMux()

	// The middleware factory returns a handler, so we need to wrap the handler function properly.
	proxyHandler := http.HandlerFunc(j.handle)
	mux.Handle("/", middlewareFactory(j.config)(proxyHandler))

	//mux.HandleFunc("/edge_ws", j.handleWs)

	srv := http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 60 * time.Second,
	}

	go func() {
		if err := srv.Serve(lis); err != nil {
			j.logger.Error("closed http connection", "err", err)
		}
	}()

	return nil
}

// The middlewareFactory builds a middleware which enables CORS using the provided config.
func middlewareFactory(config *Config) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			for _, allowedOrigin := range config.AccessControlAllowOrigin {
				if allowedOrigin == "*" {
					w.Header().Set("Access-Control-Allow-Origin", "*")

					break
				}

				if allowedOrigin == origin {
					w.Header().Set("Access-Control-Allow-Origin", origin)

					break
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (j *TransparentProxy) handle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set(
		"Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization",
	)

	switch req.Method {
	case "POST":
		j.handlePostRequest(w, req)
	case "GET":
		j.handleGetRequest(w)
	case "OPTIONS":
		// nothing to return
	default:
		_, _ = w.Write([]byte("method " + req.Method + " not allowed"))
	}
}

type EdgePath struct {
	NodeID       string `json:"node_id"`
	Port         int    `json:"port"`
	InterfaceURL string `json:"interface_url"`
}

type TransparentForward struct {
	EdgePath EdgePath `json:"edge_path"`
	Payload  string   `json:"payload"`
}

func ParseEdgePath(req *http.Request) (*EdgePath, error) {
	path := req.URL.Path
	parts := strings.Split(path, "/")

	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid path format: expected at least 4 parts, got %d", len(parts))
	}

	nodeID := parts[2]
	port := parts[3]
	interfaceURL := strings.Join(parts[4:], "/")

	// 解码URL编码的部分
	decodedNodeID, err := url.QueryUnescape(nodeID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode nodeID: %w", err)
	}

	decodedPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("failed to decode port: %w", err)
	}

	decodedInterfaceURL, err := url.QueryUnescape(interfaceURL)
	if err != nil {
		return nil, fmt.Errorf("failed to decode interfaceURL: %w", err)
	}

	return &EdgePath{
		NodeID:       decodedNodeID,
		Port:         decodedPort,
		InterfaceURL: decodedInterfaceURL,
	}, nil
}

func (j *TransparentProxy) handlePostRequest(w http.ResponseWriter, req *http.Request) {
	pathInfo, err := ParseEdgePath(req)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	j.logger.Info("handle", "NodeID", pathInfo.NodeID, "Port", pathInfo.Port, "InterfaceURL", pathInfo.InterfaceURL)

	// TODO verify NodeID by whitelist

	defer req.Body.Close()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}
	// log request
	j.logger.Debug("handle", "request", string(body))

	//var randReader io.Reader
	//randReader = rand.Reader
	//prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, randReader)
	//if err != nil {
	//	panic(err)
	//}

	// TODO replace by j.config.Store.GetRelayHost()
	//listen, _ := ma.NewMultiaddr("/ip4/127.0.0.1/tcp/10001")
	//clientHost, err := libp2p.New(
	//	libp2p.ListenAddrs(listen),
	//	libp2p.Security(noise.ID, noise.New),
	//	libp2p.Identity(prvKey),
	//	//libp2p.EnableRelay(),
	//	//libp2p.EnableAutoRelayWithStaticRelays([]peer.AddrInfo{*relayinfo}, autorelay.WithNumRelays(1)),
	//)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer clientHost.Close()
	clientHost := j.config.Store.GetRelayHost()

	// TODO query pathInfo.NodeID in PeerStore
	// TODO query relayer's peerID by pathInfo.NodeID
	targetRelayInfo, err := peer.AddrInfoFromString(fmt.Sprintf("%s/p2p/%s/p2p-circuit/p2p/%s", j.config.Store.GetRelayHost().Addrs()[0].String(), j.config.Store.GetRelayHost().ID().String(), pathInfo.NodeID))
	if err != nil {
		log.Fatal(err)
		return
	}
	clientHost.Peerstore().AddAddrs(targetRelayInfo.ID, targetRelayInfo.Addrs, peerstore.PermanentAddrTTL)

	tr := &http.Transport{}
	tr.RegisterProtocol("libp2p", p2phttp.NewTransport(clientHost, p2phttp.ProtocolOption(ProtoTagEcApp)))
	client := &http.Client{Transport: tr}

	transparentForwardData := &TransparentForward{
		EdgePath: *pathInfo,
		Payload:  string(body),
	}
	data, err := json.Marshal(transparentForwardData)
	if err != nil {
		return
	}

	targetURL := fmt.Sprintf("libp2p://%s%s", pathInfo.NodeID, TransparentRewardUrl)
	request, err := http.NewRequest(req.Method, targetURL, bytes.NewBufferString(string(data)))
	if err != nil {
		http.Error(w, "Failed to create p2p request", http.StatusInternalServerError)
		return
	}
	// forward headers
	for key, values := range req.Header {
		for _, value := range values {
			request.Header.Add(key, value)
		}
	}

	// do request
	resp, err := client.Do(request)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	for key, value := range resp.Header {
		w.Header().Set(key, value[0])
	}
	w.WriteHeader(resp.StatusCode)
	if resp.Header.Get("Content-Type") == "text/event-stream" {
		reader := bufio.NewReader(resp.Body)
		for {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					j.logger.Debug("handlePostRequest", "msg", "SSE stream closed by server")
					return
				}
				j.logger.Warn("handlePostRequest", "err", fmt.Sprintf("Error reading SSE stream: %v\n", err))
				return
			}

			_, err = w.Write(line)
			if err != nil {
				j.logger.Warn("handlePostRequest", "err", fmt.Sprintf("Error writing to client: %v\n", err))
				return
			}

			w.(http.Flusher).Flush()
		}
	} else {
		io.Copy(w, resp.Body)
	}
}

type GetResponse struct {
	Name    string `json:"name"`
	ChainID uint64 `json:"chain_id"`
	Version string `json:"version"`
}

func (j *TransparentProxy) handleGetRequest(writer io.Writer) {
	data := &GetResponse{
		Name:    j.config.ChainName,
		ChainID: j.config.NetworkID,
		Version: versioning.Version,
	}

	resp, err := json.Marshal(data)
	if err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}

	if _, err = writer.Write(resp); err != nil {
		_, _ = writer.Write([]byte(err.Error()))
	}
}
