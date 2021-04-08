## gin的快速应用脚手架

安装 egin
```bash
go get -u github.com/daodao97/egin
```

安装骨架包进行开发
```bash
cd /workspace
git clone github.com/daodao97/egin-skeleton my-app
cd my-app
goland . #编辑器打开  vscode .
air #开发模式启动
```

安装开发辅助工具
```bash
go get -u github.com/daodao97/egin-tools
```

#### 路由自动生成

### 配置的管理

- [x] utils/config.go env,json等配置方式的支持

### 日志

- [x] utils/logger.go 基于logrus的日志管理, 标准输出/文件(分割)
- [ ] 日志输出到 es/mongo/redis 

### 控制器的封装

基于单一数据模型的`CRUD`通用控制器方法

- [ ] 底层通用方法

### 数据模型的封装

基于单表的`CRUD`方法

- [x] db/model.go 连接池 
- [x] db/model.go crud 方法封装 

### 数据验证

- [x] 验证器
- [x] 接口参数合法性自动验证

### 权限验证

- [x] JWT
- [x] AKSK

### 缓存

- [x] redis的基础封装
- [ ] 全量方法完善

### 健康控制

- [ ] 接口频率控制
    - [x] 基于滑动窗口的频率限制
    - [x] ip维度
    - [ ] 用户维度
- [ ] prometheus打点
    - [x] api打点
    - [ ] db数据打点
    - [x] redis打点

### 微服务

- [x] grpc

### 其他

- [x] consul
- [x] 基于consul的配置中心
- [ ] 开关服务
- [ ] 规则引擎
- [x] RabbitMQ
- [x] kafka
- [x] nsq
- [x] 事件总线
- [ ] 配置解密
- [x] swagger
- [x] mongo
