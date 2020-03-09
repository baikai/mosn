package integrate

import (
	"testing"
	"time"

	"mosn.io/mosn/pkg/mosn"
	"mosn.io/mosn/pkg/protocol"
	"mosn.io/mosn/pkg/types"
	testutil "mosn.io/mosn/test/util"
)

type tcpExtendCase struct {
	*TestCase
}

func (c *tcpExtendCase) Start(isRouteEntryMode bool) {
	c.AppServer.GoServe()
	appAddr := c.AppServer.Addr()
	meshAddr := testutil.CurrentMeshAddr()
	c.ClientMeshAddr = meshAddr
	cfg := testutil.CreateTCPProxyConfig(meshAddr, []string{appAddr}, isRouteEntryMode)
	mesh := mosn.NewMosn(cfg)
	mesh.Start()
	c.DeferFinishCase(func() {
		c.AppServer.Close()
		mesh.Close()
	})
	time.Sleep(1 * time.Second) //wait server and mesh start
}

const _NIL types.Protocol = "null"

func TestTCPProxy(t *testing.T) {
	testCases := []*tcpExtendCase{
		&tcpExtendCase{NewTestCase(t, protocol.HTTP1, _NIL, testutil.NewHTTPServer(t, nil))},
		&tcpExtendCase{NewTestCase(t, protocol.HTTP2, _NIL, testutil.NewUpstreamHTTP2WithAnyPort(t, nil))},
		&tcpExtendCase{NewTestCase(t, protocol.SofaRPC, _NIL, testutil.NewRPCServerWithAnyPort(t, testutil.Bolt1))},
	}
	for i, tc := range testCases {
		t.Logf("start case #%d\n", i)
		tc.Start(false)
		go tc.RunCase(1, 0)
		select {
		case err := <-tc.C:
			if err != nil {
				t.Errorf("[ERROR MESSAGE] #%d tcp proxy test failed, protocol: %s, error: %v\n", i, tc.AppProtocol, err)
			}
		case <-time.After(15 * time.Second):
			t.Errorf("[ERROR MESSAGE] #%d tcp proxy hang, protocol: %s\n", i, tc.AppProtocol)
		}
		tc.FinishCase()
	}
}
func TestTCPProxyRouteEntry(t *testing.T) {
	testCases := []*tcpExtendCase{
		&tcpExtendCase{NewTestCase(t, protocol.HTTP1, _NIL, testutil.NewHTTPServer(t, nil))},
		&tcpExtendCase{NewTestCase(t, protocol.HTTP2, _NIL, testutil.NewUpstreamHTTP2WithAnyPort(t, nil))},
		&tcpExtendCase{NewTestCase(t, protocol.SofaRPC, _NIL, testutil.NewRPCServerWithAnyPort(t, testutil.Bolt1))},
	}
	for i, tc := range testCases {
		t.Logf("start case #%d\n", i)
		tc.Start(true)
		go tc.RunCase(1, 0)
		select {
		case err := <-tc.C:
			if err != nil {
				t.Errorf("[ERROR MESSAGE] #%d tcp proxy route entry test failed, protocol: %s, error: %v\n", i, tc.AppProtocol, err)
			}
		case <-time.After(15 * time.Second):
			t.Errorf("[ERROR MESSAGE] #%d tcp proxy route entry hang, protocol: %s\n", i, tc.AppProtocol)
		}
		tc.FinishCase()
	}
}
