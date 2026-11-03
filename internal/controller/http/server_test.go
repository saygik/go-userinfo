package http

import (
	"context"
	"fmt"
	"net"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func freePort(t *testing.T) int {
	t.Helper()

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to get free port: %v", err)
	}
	defer l.Close()

	return l.Addr().(*net.TCPAddr).Port
}

func TestServerStartShutsDownOnContextCancel(t *testing.T) {
	gin.SetMode(gin.TestMode)

	s := NewServer("local", logrus.New())
	port := freePort(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan error, 1)
	go func() {
		done <- s.Start(ctx, fmt.Sprint(port))
	}()

	url := fmt.Sprintf("http://127.0.0.1:%d/", port)
	deadline := time.Now().Add(5 * time.Second)
	for {
		resp, err := nethttp.Get(url)
		if err == nil {
			resp.Body.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("server did not start: %v", err)
		}
		time.Sleep(20 * time.Millisecond)
	}

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Start returned error on shutdown: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("server did not stop after context cancel")
	}
}
