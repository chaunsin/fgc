# fgc(fabric-gen-config)

由于fabric证书配置复杂编写配置容易搞错,本人想着便捷于是诞生了此工具,生成sdk链接配置文件

# 使用

```shell
go run main go
```

# 功能

- [ ] 支持生成普通配置文件生成
    - [x] 支持golang普通配置文件生成
    - [ ] 支持java普通配置文件生成
    - [ ] 支持nodejs普通配置文件生成
- [x] 支持生成 yaml json配置文件
- [ ] 支持生成gateway链接配置文件
    - [ ] golang网关钱包配置生成
    - [ ] java网关钱包配置生成
    - [ ] nodejs网关钱包配置生成
- [ ] 支持sftp读取配置文件
- [ ] 支持ftp读取配置文件

细节功能：

- [x] 可控制生成双tls认证方式
- [ ] 可控制生成 Metrics Operations CA配置作用于模块配置
- [ ] 可控生成文件是硬编码方式还是路径方式,以及golang环境魔法变量${FABRIC_SDK_GO_PROJECT_PATH}/${CRYPTOCONFIG_FIXTURES_PATH}
- [ ] 可以支持魔法变量导入路径或者参数例如$(pwd)或者${pwd}
- [ ] 增加配置注释内容

# 问题

由于fabric组件服务较多,配置复杂,天生自带分布式属性多机部署,在实际生成环境中会更加恶劣,因此此工具也面临着一些配置文件需要二次修改的问题,目前碰到的痛点有如下

1. mspid 不太容易获取
    1. docker命令方式获取
    2. 配置区块中获取
    3. configtx.yaml
    4. 进入容器读取环境变量 CORE_PEER_LOCALMSPID
    5. 使用Discover服务来获取相关配置信息,但也面临着二次配置证书公私钥等信息
2. 组织服务的真实ip、域名或端口获取问题
    1. 使用docker命令获取
       `docker ps --format "table{{.Image}}\t{{.Names}}\t{{.Ports}}" | grep "hyperledger/fabric-peer\|hyperledger/fabric-orderer" | awk '{print $2,$3}'`
3. 