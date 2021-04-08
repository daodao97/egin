## egin开发辅助工具

安装
```bash
go get -u github.com/daodao97/egin-tools
```

常用命令
```bash
# 查看帮助
egin-tools -h

# 生成 swagger 文件 并启动 ui
egin-tools -swagger -ui

# 根据 controller/* 文件 自动生成 gin 路由注册代码
egin-tools -route

# 生成数据库模型文件
egin-tools -model -database hyperf_admin -table reports

# 生成数据库模型文件
egin-tools -model -database hyperf_admin
```
