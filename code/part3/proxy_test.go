package part3

import (
	"io"
	"net"
	"sync"
	"testing"
)

func mockServer(l net.Listener, t *testing.T) {
	buf := make([]byte, 1024)

	//for {
	conn, err := l.Accept()
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	n, err := conn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	t.Logf("SERVER read %d bytes, %s", n, string(buf[:n]))

	m := copy(buf[n:], []byte(" ^_^"))
	n, err = conn.Write(buf[:n+m])
	if err != nil {
		t.Error(err)
	}
	t.Logf("SERVER write %d bytes", n)
	//}
}

func handleConn(client net.Conn, serverAddr string, t *testing.T) {
	defer client.Close()

	server, err := net.Dial("tcp", serverAddr)
	if err != nil {
		t.Error(err)
	}
	defer server.Close()

	var wg sync.WaitGroup

	doProxy := func(dst net.Conn, src net.Conn, isClient bool) {
		_, err := io.Copy(dst, src)
		t.Logf("%v %s => %s, %v", isClient, src.LocalAddr(), dst.LocalAddr(), err)
		wg.Done()
	}

	wg.Add(1)
	go doProxy(server, client, true)
	wg.Add(1)
	go doProxy(client, server, false)

	wg.Wait()
}

func mockProxy(l net.Listener, serverAddr string, t *testing.T) {
	for {
		conn, err := l.Accept()
		if err != nil {
			t.Error(err)
			return
		}

		go handleConn(conn, serverAddr, t)
	}
}

func TestTcpProxy(t *testing.T) {

	serverListener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	proxyListener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go mockServer(serverListener, t)
	go mockProxy(proxyListener, serverListener.Addr().String(), t)

	// client
	conn, err := net.Dial("tcp", proxyListener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	n := copy(buf, []byte("Proxy Testing"))
	n, err = conn.Write(buf[:n])
	if err != nil {
		t.Error(err)
	}
	t.Logf("CLIENT send %d bytes", n)

	n, err = conn.Read(buf)
	if err != nil {
		t.Error(err)
	}
	t.Logf("CLIENT read %d bytes, %s", n, string(buf[:n]))
}
