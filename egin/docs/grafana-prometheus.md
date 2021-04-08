`mac`下安装

### 安装 prometheus
[prometheus/download](https://prometheus.io/download/)下载指定版本

`mac` 下载 `prometheus-2.20.1.darwin-386.tar.gz` 即可

```bash
tar zxvf prometheus-2.20.1.darwin-386.tar.gz
```

`prometheus.yml`

```yml
global:
  scrape_interval: 10s
  scrape_timeout: 10s
  evaluation_interval: 10m
scrape_configs:
  - job_name: 服务名
    scrape_interval: 5s
    scrape_timeout: 5s
    metrics_path: /prometheus 
    scheme: http
    basic_auth:
      username: user
      password: pwd
    static_configs:
      - targets:
        - 127.0.0.1:8888 # 服务地址
```

```bash
./prometheus --config.file="prometheus.yml"
```

[prometheus后台地址](http://127.0.0.1:9090)

### 安装 grafana

```bash
brew install grafana
brew service start grafana
```

[grafana后台地址](http://127.0.0.1:3000)


### grafana 添加数据源

http://127.0.0.1:3000/datasources/new?gettingstarted

![IbIART](https://cdn.jsdelivr.net/gh/daodao97/FigureBed@master/uPic/IbIART.png)

![a8oNOT](https://cdn.jsdelivr.net/gh/daodao97/FigureBed@master/uPic/a8oNOT.png)

### 导入数据面板 

[面板模板](https://grafana.com/grafana/dashboards)

![63uonZ](https://cdn.jsdelivr.net/gh/daodao97/FigureBed@master/uPic/63uonZ.png)


![Jk4P4m](https://cdn.jsdelivr.net/gh/daodao97/FigureBed@master/uPic/Jk4P4m.png)

![5LqUsy](https://cdn.jsdelivr.net/gh/daodao97/FigureBed@master/uPic/5LqUsy.png)

[参考](https://github.com/chenjiandongx/ginprom)

