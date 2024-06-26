# fgc(Fabric Generate Config)

[![Qodana](https://github.com/chaunsin/fgc/actions/workflows/qodana_code_quality.yml/badge.svg?branch=master)](https://github.com/chaunsin/fgc/actions/workflows/qodana_code_quality.yml)

由于fabric证书配置复杂编写配置容易搞错,本人想着便捷于是诞生了此工具,用于生成sdk链接配置文件.

# 注意！！！

目前此库还不完善不过可以生成golang可用的配置

# 安装

方式一

```shell
go install github.com/chaunsin/fgc@latest 
```

此方式生成的`fgc`二进制可执行程序会放到$GOBIN目录下

方式二

```shell
git clone https://github.com/chaunsin/fgc.git
cd fgc
make install
```

执行之后会在fgc目录下生成`fgc`可执行文件,推荐放入到全局可引导的目录位置比如`/usr/local/bin`目录下

# 使用

生成fabric-sdk-go sdk链接配置证书

```shell
# tips: crypto-config为fabric证书目录
fgc go -i ./crypto-config
```

帮助

```shell
fgc -h
```

# 功能

- [ ] 支持生成普通配置文件生成
    - [x] 支持golang普通配置文件生成
    - [ ] 支持java普通配置文件生成
    - [ ] 支持nodejs普通配置文件生成
- [ ] 配置文件格式
    - [x] 支持生成yaml配置文件
    - [ ] 支持生成json配置文件(目前能生成但是配置文件未必能使用)
- [ ] 支持生成gateway连接配置文件
    - [ ] golang网关钱包配置生成
    - [ ] java网关钱包配置生成
    - [ ] nodejs网关钱包配置生成
- [ ] 支持sftp读取配置文件
- [ ] 支持ftp读取配置文件

细节功能：

- [x] 可控制生成双tls认证连接方式
- [ ] 生成 Metrics Operations CA模块配置
- [ ] 可控生成文件是硬编码方式还是路径方式,以及golang环境魔法变量${FABRIC_SDK_GO_PROJECT_PATH}/${CRYPTOCONFIG_FIXTURES_PATH}
- [ ] 支持魔法变量导入路径或者参数例如:$(pwd)或者${pwd}
- [ ] 增加配置注释内容

# 问题

由于fabric组织节点服务较多,关系配置复杂,天生自带分布式属性多机部署,在实际环境中会更加恶劣,因此此工具也面临着一些配置文件需要二次修改的问题,目前碰到的痛点有如下

1. mspid 不太容易获取
    1. docker命令方式获取?
    2. 配置区块中获取?
    3. configtx.yaml?
    4. 进入容器读取环境变量 CORE_PEER_LOCALMSPID?
    5. 使用Discover服务来获取相关配置信息,但也面临着二次配置证书公私钥等信息?
2. 获取组织服务的真实ip、域名或端口问题
    1. 使用docker命令获取
       `docker ps --format "table{{.Image}}\t{{.Names}}\t{{.Ports}}" | grep "hyperledger/fabric-peer\|hyperledger/fabric-orderer" | awk '{print $2,$3}'`
3. peer下面有两个组织每个组织有两个节点,但是每个组织只生成一个节点需要排查修改(貌似没问题)
