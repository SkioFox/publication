## note
### 变量声明

如果没有显式为变量赋予初值，Go 编译器会为变量赋予这个类型的零值。 整数类型是0、浮点类型是0.0 、布尔是false、指针、接口、切片、channel、map和函数是nil
![img.png](go-first-course/img/img.png)
### 区别总结空结构和空接口

- struct{} 是一个空结构体，表示一个不包含任何字段的结构体；而 interface{} 是一个空接口，表示一个不包含任何方法的接口。
- 空结构体通常用于表示不需要存储任何数据的情况，而空接口通常用于表示不限定具体类型的情况。
- 空结构体在并发编程中常用作信号或事件的通知；空接口常用于需要接收任意类型的值的场景。

在实际编程中，空结构体和空接口都有其特定的用途，可以根据具体情况选择使用。

### byte类型 => 字节类型

在 gRPC 请求中，byte 类型通常用于表示二进制数据，例如文件内容、图像数据、视频数据等。gRPC 使用 Protocol Buffers 作为默认的序列化和反序列化格式，而 Protocol Buffers 支持将二进制数据表示为字节序列。

在 gRPC 中，您可能会使用 byte 类型来传输二进制数据。例如，如果您有一个 gRPC 服务，需要接收一个文件作为输入，您可以将文件内容表示为 byte 类型，并在 gRPC 请求中传输这些字节。
```
// 定义 gRPC 请求消息
message UploadFileRequest {
  bytes file_content = 1;
}
// 生成的 gRPC 服务端代码
func (s *Server) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
    fileContent := req.GetFileContent()
    // 处理文件内容
    return &pb.UploadFileResponse{}, nil
}

```