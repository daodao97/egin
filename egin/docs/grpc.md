
### 安装 protobuf
```bash
 brew install protobuf
```

### 安装 protobuf golang 插件
```bash
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```

### 生成代码
```bash
protoc --go_out=plugins=grpc:. hello.proto
```

### 整体流程
1. 定义服务描述文件 ***.proto
2. 生成服务定义代码 ***.pb.go
3. 完成接口实现  ***Serve
4. 绑定具体实现到接口定义
5. 启动 rpc 服务
