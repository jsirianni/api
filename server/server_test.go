package server

import (
	"testing"

	"github.com/jsirianni/server/logging"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestWithBindAddress(t *testing.T) {
	cases := []struct {
		name      string
		addr      string
		port      uint
		expect    string
		expectErr string
	}{
		{
			"port-only",
			"",
			8000,
			":8000",
			"",
		},
		{
			"ivp4-and-port",
			"192.168.1.40",
			10000,
			"192.168.1.40:10000",
			"",
		},
		{
			"ipv4-localhost",
			"127.0.0.1",
			8000,
			"127.0.0.1:8000",
			"",
		},
		{
			"ipv6-and0-port",
			"fe80::fc18:f93a:9769:3fcc",
			9090,
			"[fe80::fc18:f93a:9769:3fcc]:9090",
			"",
		},
		{
			"localhost",
			"localhost",
			8000,
			"",
			"failed to parse 'localhost' as an IP address",
		},
		{
			"invalid-ip",
			"0.0.0.",
			8000,
			"",
			"failed to parse '0.0.0.' as an IP address",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f := WithBindAddress(tc.addr, tc.port)
			s := &Server{}
			err := f(s)
			if tc.expectErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.expectErr)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expect, s.server.Addr)
		})
	}
}

func TestNew(t *testing.T) {
	cases := []struct {
		name      string
		ops       []Option
		expectErr bool
		errStr    string
	}{
		{
			"valid",
			[]Option{
				WithBindAddress("", 10000),
				WithMemoryStore(false),
			},
			false,
			"",
		},
		{
			"invalid",
			[]Option{
				WithBindAddress("x.x.x", 0),
				WithMemoryStore(false),
			},
			true,
			"failed to parse 'x.x.x' as an IP address",
		},
		{
			"missing-store",
			[]Option{
				WithBindAddress("", 10000),
			},
			true,
			"server must be configured with a storage backend",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			server, err := New(testLogger(t), tc.ops...)

			if tc.expectErr {
				require.Error(t, err)
				require.ErrorContains(t, err, tc.errStr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, server)
		})
	}
}

func TestNewNoLogger(t *testing.T) {
	_, err := New(nil)
	require.Error(t, err, "an error is expected when a nil logger is passed")
}

func testLogger(t *testing.T) *zap.Logger {
	logger, err := logging.New(logging.DebugLevel)
	if err != nil {
		t.Errorf("%s: expected logging.New to return a logger without an error. this indicated an issue with the internal/logger package.", err)
		t.Fail()
	}
	return logger
}
