package relay

import (
	"context"
	"errors"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/application"
	appProto "github.com/EdgeMatrixChain/edge-matrix-core/core/application/proto"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/network/common"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/network/grpc"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/relay/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/multiformats/go-multiaddr"
	"net"
	"regexp"
	"sync"
)

const (
	EdgeAliveProto = "/alive/0.2"
)

// networkingServer defines the base communication interface between
// any networking server implementation and the AliveService
type networkingServer interface {

	// GetPeerAddrInfo fetches the AddrInfo of a peer
	GetPeerAddrInfo(peerID peer.ID) peer.AddrInfo

	// PublishApplicationStatus publish application status
	PublishApplicationStatus(status *appProto.AppStatus)

	GetRelayProxyAddr() *net.TCPAddr

	GetRandomBootnode() *peer.AddrInfo
}

// BOOTNODE QUERIES //
// AliveService is a service that finds other peers in the network
// and connects them to the current running node
type AliveService struct {
	proto.UnimplementedAliveServer
	pendingPeerConnections sync.Map // Map that keeps track of the pending status of peers; peerID -> bool

	baseServer networkingServer // The interface towards the base networking server
	logger     hclog.Logger     // The AliveService logger

	syncAppPeerClient application.SyncAppPeerClient
}

// NewAliveService creates a new instance of the alive service
func NewAliveService(
	server networkingServer,
	logger hclog.Logger,
	syncAppPeerClient application.SyncAppPeerClient,
) *AliveService {
	return &AliveService{
		logger:            logger.Named("AliveService"),
		baseServer:        server,
		syncAppPeerClient: syncAppPeerClient,
	}
}

func (d *AliveService) Hello(ctx context.Context, status *proto.AliveStatus) (*proto.AliveStatusResp, error) {
	// Extract the requesting peer ID from the gRPC context
	grpcContext, ok := ctx.(*grpc.Context)
	if !ok {
		return nil, errors.New("invalid type assertion")
	}

	from := grpcContext.PeerID
	addr := ""
	innerIp := false
	addrInfo := d.baseServer.GetPeerAddrInfo(from)
	if len(addrInfo.Addrs) > 0 {
		addr = addrInfo.Addrs[0].String()
		innerIp = isInnerIp(addrInfo.Addrs[0])
	}
	d.logger.Debug("-------->Alive status", "from", from, "name", status.Name, "app_origin", status.AppOrigin, "addr", addr, "relay", status.Relay)

	// get RelayProxyPort and put it into AppStatus
	proxyPort := 0
	proxyAddr := d.baseServer.GetRelayProxyAddr()
	if proxyAddr != nil {
		proxyPort = proxyAddr.Port
	}

	if !innerIp || status.Relay != "" {
		d.baseServer.PublishApplicationStatus(&appProto.AppStatus{
			Name:           status.Name,
			NodeId:         from.String(),
			Uptime:         status.Uptime,
			StartupTime:    status.StartupTime,
			Relay:          status.Relay,
			Addr:           addr,
			AppOrigin:      status.AppOrigin,
			Mac:            status.Mac,
			CpuInfo:        status.CpuInfo,
			GpuInfo:        status.GpuInfo,
			MemInfo:        status.MemInfo,
			ModelHash:      status.ModelHash,
			AveragePower:   status.AveragePower,
			Version:        status.Version,
			RelayProxyPort: uint64(proxyPort),
		})
	}

	newRelayNode := d.baseServer.GetRandomBootnode()
	discovery := ""
	if newRelayNode != nil {
		discovery = common.AddrInfoToString(newRelayNode)
	}
	return &proto.AliveStatusResp{
		Success:   true,
		Discovery: discovery,
	}, nil
}

func isInnerIp(ma multiaddr.Multiaddr) (innerIp bool) {
	innerIp = false
	ip4Addr, err := ma.ValueForProtocol(multiaddr.P_IP4)
	if err != nil {
		ip4Addr = ""
	}

	re := regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)\.(\d+)`)
	submatches := re.FindStringSubmatch(ip4Addr)
	if len(submatches) > 0 {
		// 127.0.0.1
		// 10.0.0.0/8
		// 172.16.0.0/12
		// 169.254.0.0/16
		// 192.168.0.0/16
		if submatches[0] == "127.0.0.1" ||
			submatches[1] == "10" ||
			(submatches[1] == "172" && submatches[2] >= "16" && submatches[2] <= "31") ||
			(submatches[1] == "169" && submatches[2] == "254") ||
			(submatches[1] == "192" && submatches[2] == "168") {
			innerIp = true
		}
	}
	return innerIp
}
