# Chapter 4: 传送TCP数据

我们将开始学习从一个网络连接中读取数据的方法。我们将设计一个简单的传输协议，可以在不同的节点之间传输不固定大小的payload负载。我们也将扩展net.Conn接口以实现各种网络操作。我们也会深入探讨在TCP网络编程中会碰到的各种问题。

## 使用net.Conn接口

大部分代码都会使用net.Conn接口实现，因为它提供了大多数实例所需要的功能，通过net.Conn接口我们可以写出非常健壮的跨平台通信系统。

两个经常被用到的函数是Read和Write，他们实现了io.Reader和io.Writer接口，这两个接口在Go的标准库中非常常见。使用这两个接口可以写出功能非常强大的网络应用程序。

我们使用net.Conn的Close方法来关闭网络链接，这个方法返回nil表示链接被正常关闭，否则返回error。SetReadDeadline和SetWriteDeadline方法接收一个time.Time的对象，通过设置一个绝对时间可以在read和write方法中返回超时错误。SetDeadline会同时为read和write设置结束时间。

## 发送和接收数据

### *读取数据到可变长缓存中*

```golang
package main

func main() {
    go send_data()
    return
}
```