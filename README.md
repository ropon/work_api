### 后端快速Api脚手架

#### 1、sdk

- web:gin
- db:gorm
- redis:go-redis
- etcd:go.etcd.io
- kafka:sarama
- swagger
- jaeger

#### 2、usage

```shell
#第一个参数服务名称 第二个参数监听端口
bash new_project.sh testpro 8866
```

#### 3、dev

```shell
#热加载自动生成swagger文档
go install github.com/swaggo/swag/cmd/swag@v1.8.12
cd testpro
./air
```
