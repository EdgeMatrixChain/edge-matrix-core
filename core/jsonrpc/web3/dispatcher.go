package web3

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/EdgeMatrixChain/edge-matrix-core/core/jsonrpc"
	"github.com/libp2p/go-libp2p/core/host"
	"math"
	"reflect"
	"strings"
	"unicode"

	"github.com/hashicorp/go-hclog"
)

type serviceData struct {
	sv      reflect.Value
	funcMap map[string]*funcData
}

type funcData struct {
	inNum int
	reqt  []reflect.Type
	fv    reflect.Value
	isDyn bool
}

func (f *funcData) numParams() int {
	return f.inNum - 1
}

type endpoints struct {
	Edge *Edge
	//TelePool *TelePool
}

// Dispatcher handles all json rpc requests by delegating
// the execution flow to the corresponding service
type Dispatcher struct {
	logger            hclog.Logger
	serviceMap        map[string]*serviceData
	nodeFilterManager *NodeFilterManager
	endpoints         endpoints
	params            *dispatcherParams
	host              host.Host
}

type dispatcherParams struct {
	networkID   uint64
	networkName string

	jsonRPCBatchLengthLimit uint64
}

func newDispatcher(
	logger hclog.Logger,
	store JSONRPCStore,
	params *dispatcherParams,
) *Dispatcher {
	d := &Dispatcher{
		logger: logger.Named("dispatcher"),
		params: params,
	}

	if store != nil {
		d.nodeFilterManager = NewNodeFilterManager(logger, store)
		go d.nodeFilterManager.Run()

		d.host = store.GetHost()
	}

	d.registerEndpoints(store)

	return d
}

func (d *Dispatcher) registerEndpoints(store JSONRPCStore) {
	d.endpoints.Edge = &Edge{
		d.logger,
		store,
		d.params.networkID,
		d.nodeFilterManager,
	}

	d.registerService("edge", d.endpoints.Edge)
}

func (d *Dispatcher) getFnHandler(req Request) (*serviceData, *funcData, jsonrpc.Error) {
	callName := strings.SplitN(req.Method, "_", 2)
	if len(callName) != 2 {
		return nil, nil, jsonrpc.NewMethodNotFoundError(req.Method)
	}

	serviceName, funcName := callName[0], callName[1]

	service, ok := d.serviceMap[serviceName]
	if !ok {
		return nil, nil, jsonrpc.NewMethodNotFoundError(req.Method)
	}

	fd, ok := service.funcMap[funcName]

	if !ok {
		return nil, nil, jsonrpc.NewMethodNotFoundError(req.Method)
	}

	return service, fd, nil
}

type wsConn interface {
	WriteMessage(messageType int, data []byte) error
	GetFilterID() string
	SetFilterID(string)
}

// as per https://www.jsonrpc.org/specification, the `id` in JSON-RPC 2.0
// can only be a string or a non-decimal integer
func formatFilterResponse(id interface{}, resp string) (string, jsonrpc.Error) {
	switch t := id.(type) {
	case string:
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":"%s","result":"%s"}`, t, resp), nil
	case float64:
		if t == math.Trunc(t) {
			return fmt.Sprintf(`{"jsonrpc":"2.0","id":%d,"result":"%s"}`, int(t), resp), nil
		} else {
			return "", jsonrpc.NewInvalidRequestError("Invalid json request")
		}
	case nil:
		return fmt.Sprintf(`{"jsonrpc":"2.0","id":null,"result":"%s"}`, resp), nil
	default:
		return "", jsonrpc.NewInvalidRequestError("Invalid json request")
	}
}

func (d *Dispatcher) handleEdgeSubscribe(req Request, conn wsConn) (string, jsonrpc.Error) {
	var params []interface{}
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return "", jsonrpc.NewInvalidRequestError("Invalid json request")
	}

	if len(params) == 0 {
		return "", jsonrpc.NewInvalidParamsError("Invalid params")
	}

	subscribeMethod, ok := params[0].(string)
	if !ok {
		return "", jsonrpc.NewSubscriptionNotFoundError(subscribeMethod)
	}

	var filterID string
	if subscribeMethod == "node" {
		if len(params) < 2 {
			return "", jsonrpc.NewInvalidRequestError("params[1] is not exist")
		}
		nodeQuery, err := DecodeNodeQueryFromInterface(params[1])
		if err != nil {
			return "", jsonrpc.NewInternalError(err.Error())
		}
		filterID = d.nodeFilterManager.NewNodeFilter(nodeQuery, conn)
	} else {
		return "", jsonrpc.NewSubscriptionNotFoundError(subscribeMethod)
	}

	return filterID, nil
}

func (d *Dispatcher) RemoveFilterByWs(conn wsConn) {
	d.nodeFilterManager.RemoveFilterByWs(conn)
}

func (d *Dispatcher) HandleWs(reqBody []byte, conn wsConn) ([]byte, error) {
	var req Request
	if err := json.Unmarshal(reqBody, &req); err != nil {
		return NewRPCResponse(req.ID, "2.0", nil, jsonrpc.NewInvalidRequestError("Invalid json request")).Bytes()
	}

	// if the request method is edge_subscribe we need to create a
	// new filter with ws connection
	// example:
	//{
	//  "jsonrpc": "2.0",
	//  "id": 1,
	//  "method": "edge_subscribe",
	//  "params": [
	//    "node",
	//    {}
	//  ]
	//}
	if req.Method == "edge_subscribe" {
		filterID, err := d.handleEdgeSubscribe(req, conn)
		if err != nil {
			return NewRPCResponse(req.ID, "2.0", nil, err).Bytes()
		}

		resp, err := formatFilterResponse(req.ID, filterID)

		if err != nil {
			return NewRPCResponse(req.ID, "2.0", nil, err).Bytes()
		}

		return []byte(resp), nil
	}

	// its a normal query that we handle with the dispatcher
	resp, err := d.handleReq(req)
	if err != nil {
		return nil, err
	}

	return NewRPCResponse(req.ID, "2.0", resp, err).Bytes()
}

func (d *Dispatcher) Handle(reqBody []byte) ([]byte, error) {
	x := bytes.TrimLeft(reqBody, " \t\r\n")
	if len(x) == 0 {
		return NewRPCResponse(nil, "2.0", nil, jsonrpc.NewInvalidRequestError("Invalid json request")).Bytes()
	}

	if x[0] == '{' {
		var req Request
		if err := json.Unmarshal(reqBody, &req); err != nil {
			return NewRPCResponse(nil, "2.0", nil, jsonrpc.NewInvalidRequestError("Invalid json request")).Bytes()
		}

		if req.Method == "" {
			return NewRPCResponse(req.ID, "2.0", nil, jsonrpc.NewInvalidRequestError("Invalid json request")).Bytes()
		}

		resp, err := d.handleReq(req)

		return NewRPCResponse(req.ID, "2.0", resp, err).Bytes()
	}

	// handle batch requests
	var requests []Request
	if err := json.Unmarshal(reqBody, &requests); err != nil {
		return NewRPCResponse(
			nil,
			"2.0",
			nil,
			jsonrpc.NewInvalidRequestError("Invalid json request"),
		).Bytes()
	}

	// if not disabled, avoid handling long batch requests
	if d.params.jsonRPCBatchLengthLimit != 0 && len(requests) > int(d.params.jsonRPCBatchLengthLimit) {
		return NewRPCResponse(
			nil,
			"2.0",
			nil,
			jsonrpc.NewInvalidRequestError("Batch request length too long"),
		).Bytes()
	}

	responses := make([]Response, 0)

	for _, req := range requests {
		var response, err = d.handleReq(req)
		if err != nil {
			errorResponse := NewRPCResponse(req.ID, "2.0", nil, err)
			responses = append(responses, errorResponse)

			continue
		}

		resp := NewRPCResponse(req.ID, "2.0", response, nil)
		responses = append(responses, resp)
	}

	respBytes, err := json.Marshal(responses)
	if err != nil {
		return NewRPCResponse(nil, "2.0", nil, jsonrpc.NewInternalError("Internal error")).Bytes()
	}

	return respBytes, nil
}

func (d *Dispatcher) handleReq(req Request) ([]byte, jsonrpc.Error) {
	d.logger.Debug("request", "method", req.Method, "id", req.ID)

	service, fd, ferr := d.getFnHandler(req)
	if ferr != nil {
		return nil, ferr
	}

	inArgs := make([]reflect.Value, fd.inNum)
	inArgs[0] = service.sv

	inputs := make([]interface{}, fd.numParams())

	for i := 0; i < fd.inNum-1; i++ {
		val := reflect.New(fd.reqt[i+1])
		inputs[i] = val.Interface()
		inArgs[i+1] = val.Elem()
	}

	if fd.numParams() > 0 {
		if err := json.Unmarshal(req.Params, &inputs); err != nil {
			return nil, jsonrpc.NewInvalidParamsError("Invalid Params")
		}
	}

	output := fd.fv.Call(inArgs)
	if err := getError(output[1]); err != nil {
		d.logInternalError(req.Method, err)

		return nil, jsonrpc.NewInvalidRequestError(err.Error())
	}

	var (
		data []byte
		err  error
	)

	if res := output[0].Interface(); res != nil {
		data, err = json.Marshal(res)
		if err != nil {
			d.logInternalError(req.Method, err)

			return nil, jsonrpc.NewInternalError("Internal error")
		}
	}

	return data, nil
}

func (d *Dispatcher) logInternalError(method string, err error) {
	d.logger.Error("failed to dispatch", "method", method, "err", err)
}

func (d *Dispatcher) registerService(serviceName string, service interface{}) {
	if d.serviceMap == nil {
		d.serviceMap = map[string]*serviceData{}
	}

	if serviceName == "" {
		panic("jsonrpc: serviceName cannot be empty")
	}

	st := reflect.TypeOf(service)
	if st.Kind() == reflect.Struct {
		panic(fmt.Sprintf("jsonrpc: service '%s' must be a pointer to struct", serviceName))
	}

	funcMap := make(map[string]*funcData)

	for i := 0; i < st.NumMethod(); i++ {
		mv := st.Method(i)
		if mv.PkgPath != "" {
			// skip unexported methods
			continue
		}

		name := lowerCaseFirst(mv.Name)
		funcName := serviceName + "_" + name
		fd := &funcData{
			fv: mv.Func,
		}

		var err error

		if fd.inNum, fd.reqt, err = validateFunc(funcName, fd.fv, true); err != nil {
			panic(fmt.Sprintf("jsonrpc: %s", err))
		}
		// check if last item is a pointer
		if fd.numParams() != 0 {
			last := fd.reqt[fd.numParams()]
			if last.Kind() == reflect.Ptr {
				fd.isDyn = true
			}
		}

		funcMap[name] = fd
	}

	d.serviceMap[serviceName] = &serviceData{
		sv:      reflect.ValueOf(service),
		funcMap: funcMap,
	}
}

func validateFunc(funcName string, fv reflect.Value, _ bool) (inNum int, reqt []reflect.Type, err error) {
	if funcName == "" {
		err = fmt.Errorf("funcName cannot be empty")

		return
	}

	ft := fv.Type()
	if ft.Kind() != reflect.Func {
		err = fmt.Errorf("function '%s' must be a function instead of %s", funcName, ft)

		return
	}

	inNum = ft.NumIn()

	if outNum := ft.NumOut(); ft.NumOut() != 2 {
		err = fmt.Errorf("unexpected number of output arguments in the function '%s': %d. Expected 2", funcName, outNum)

		return
	}

	if !isErrorType(ft.Out(1)) {
		err = fmt.Errorf(
			"unexpected type for the second return value of the function '%s': '%s'. Expected '%s'",
			funcName,
			ft.Out(1),
			errt,
		)

		return
	}

	reqt = make([]reflect.Type, inNum)
	for i := 0; i < inNum; i++ {
		reqt[i] = ft.In(i)
	}

	return
}

var errt = reflect.TypeOf((*error)(nil)).Elem()

func isErrorType(t reflect.Type) bool {
	return t.Implements(errt)
}

func getError(v reflect.Value) error {
	if v.IsNil() {
		return nil
	}

	extractedErr, ok := v.Interface().(error)
	if !ok {
		return errors.New("invalid type assertion, unable to extract error")
	}

	return extractedErr
}

func lowerCaseFirst(str string) string {
	for i, v := range str {
		return string(unicode.ToLower(v)) + str[i+1:]
	}

	return ""
}
