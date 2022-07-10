## *读取数据到固定长度的缓存中*

```golang
package part2

import (
	"io"
	"math/rand"
	"net"
	"testing"
)

func TestReadIntoBuff(t *testing.T) {
	payload := make([]byte, 1<<24) // 16MB data
	_, err := rand.Read(payload)
	if err != nil {
		t.Fatal(err)
	}

	listener, err := net.Listen("tcp", "127.0.0.1:")
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Log(err)
			return
		}
		defer conn.Close()

		_, err = conn.Write(payload)
		if err != nil {
			t.Error(err)
		}
	}()

	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1<<19) // 512KB read buf
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Error(err)
			}
			break
		}
		t.Logf("read %d bytes", n) // buf[:n] is the data read from conn
	}
}
```

输出结果：

```
$ go test -v
=== RUN   TestReadIntoBuff
    4-read-data_test.go:51: read 65536 bytes
    4-read-data_test.go:51: read 65536 bytes
    ......
    4-read-data_test.go:51: read 65483 bytes
    4-read-data_test.go:51: read 45064 bytes
--- PASS: TestReadIntoBuff (0.07s)
PASS
ok      bekars.github.com/gonetwork/v2/part2    0.075s
```