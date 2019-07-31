#### 前言

该智能合约实现签到系统的增删改查功能，运行的fabric网络为first-network中的byfn。

#### 前提条件

- hyperledger fabric网络以及fabric-sample（包含first-network）

- cryptotxgen等工具的环境变量配置

  ```bash
  mv checkin-System/chaincode/ fabric-samples/chaincode/checkinsystem/
  mv checkin-System/scripts/ fabric-samples/first-network/scripts/
  chmod 755 fabric-samples/first-network/scripts/setparas.sh
  ```

  

#### 运行网络

```bash
cd first-network
./byfn.sh -m generate -c mychannel
./byfn.sh -m up -c mychannel -s couchdb
```

查看容器

```
CONTAINER ID        IMAGE                                                                                                           COMMAND                  CREATED             STATUS              PORTS                                        NAMES
84707d906728        dev-peer0.org2.example.com-checkinsystem-1.7-ef8d6c7d2d080f0f333a25b0c6b06784f82adc4bb1ab89da4b3a87869d5bd167   "chaincode -peer.add…"   45 minutes ago      Up 45 minutes                                                    dev-peer0.org2.example.com-checkinsystem-1.7
0fea649c1078        dev-peer0.org1.example.com-checkinsystem-1.7-949d269dca6550746f1da0589a3b1221a7a5ce87d95a7ea6dbbf773580265724   "chaincode -peer.add…"   About an hour ago   Up About an hour                                                 dev-peer0.org1.example.com-checkinsystem-1.7
c125478264b1        dev-peer0.org2.example.com-account-1.0-4595f4dca863875247f680c7459d754f543f0424a80aa6105cb8f30c4b83cd01         "chaincode -peer.add…"   3 hours ago         Up 3 hours                                                       dev-peer0.org2.example.com-Account-1.0
9d8d5adc50d1        dev-peer0.org1.example.com-account-1.0-fcc0f336a119de0ce5eb81415f4a4dcbc3b4b49bcbe3af759e323cd4be7c108f         "chaincode -peer.add…"   3 hours ago         Up 3 hours                                                       dev-peer0.org1.example.com-Account-1.0
4524e47897cf        dev-peer1.org2.example.com-account-1.0-8014823b291eeab1102da81e9a25938d67d38203d92084f0ac9a8481866b255c         "chaincode -peer.add…"   3 hours ago         Up 3 hours                                                       dev-peer1.org2.example.com-Account-1.0
df1a3264ecc1        dev-peer1.org2.example.com-mycc-1.0-26c2ef32838554aac4f7ad6f100aca865e87959c9a126e86d764c8d01f8346ab            "chaincode -peer.add…"   4 hours ago         Up 4 hours                                                       dev-peer1.org2.example.com-mycc-1.0
44e269b73698        dev-peer0.org1.example.com-mycc-1.0-384f11f484b9302df90b453200cfb25174305fce8f53f4e94d45ee3b6cab0ce9            "chaincode -peer.add…"   4 hours ago         Up 4 hours                                                       dev-peer0.org1.example.com-mycc-1.0
09615d2b7ce2        dev-peer0.org2.example.com-mycc-1.0-15b571b3ce849066b7ec74497da3b27e54e0df1345daff3951b94245ce09c42b            "chaincode -peer.add…"   4 hours ago         Up 4 hours                                                       dev-peer0.org2.example.com-mycc-1.0
86e3a050928a        hyperledger/fabric-tools:latest                                                                                 "/bin/bash"              4 hours ago         Up 4 hours                                                       cli
5bf0dc558905        hyperledger/fabric-peer:latest                                                                                  "peer node start"        4 hours ago         Up 4 hours          0.0.0.0:7051->7051/tcp                       peer0.org1.example.com
b64dc621e875        hyperledger/fabric-peer:latest                                                                                  "peer node start"        4 hours ago         Up 4 hours          0.0.0.0:9051->9051/tcp                       peer0.org2.example.com
b3a0af489fee        hyperledger/fabric-peer:latest                                                                                  "peer node start"        4 hours ago         Up 4 hours          0.0.0.0:10051->10051/tcp                     peer1.org2.example.com
a056fb0b2a72        hyperledger/fabric-peer:latest                                                                                  "peer node start"        4 hours ago         Up 4 hours          0.0.0.0:8051->8051/tcp                       peer1.org1.example.com
ec6195d91d9a        hyperledger/fabric-couchdb                                                                                      "tini -- /docker-ent…"   4 hours ago         Up 4 hours          4369/tcp, 9100/tcp, 0.0.0.0:6984->5984/tcp   couchdb1
93da7de1ac6d        hyperledger/fabric-couchdb                                                                                      "tini -- /docker-ent…"   4 hours ago         Up 4 hours          4369/tcp, 9100/tcp, 0.0.0.0:7984->5984/tcp   couchdb2
d94df5dfa26f        hyperledger/fabric-orderer:latest                                                                               "orderer"                4 hours ago         Up 4 hours          0.0.0.0:7050->7050/tcp                       orderer.example.com
f5f35d3163d4        hyperledger/fabric-couchdb                                                                                      "tini -- /docker-ent…"   4 hours ago         Up 4 hours          4369/tcp, 9100/tcp, 0.0.0.0:5984->5984/tcp   couchdb0
77c1178b2ba6        hyperledger/fabric-couchdb                                                                                      "tini -- /docker-ent…"   4 hours ago         Up 4 hours          4369/tcp, 9100/tcp, 0.0.0.0:8984->5984/tcp   couchdb3

```



#### 安装、初始化chaincode

首先进入cli容器

```bash
docker exec -it cli bash
```

配置环境变量

```bash
export CC_SRC_PATH="github.com/chaincode/checkinsystem"
export VERSION="1.0"
export CC_NAME="checkinsystem"
source ./scripts/setparas.sh
```

安装chaincode

```bash
source ./scripts/setparas.sh peerenv 0 1
peer chaincode install -n $CC_NAME -v $VERSION -l $LANGUAGE -p $CC_SRC_PATH
 
source ./scripts/setparas.sh peerenv 0 2
export CORE_PEER_ADDRESS=peer0.org2.example.com:9051
peer chaincode install -n $CC_NAME -v $VERSION -l $LANGUAGE -p $CC_SRC_PATH
 
source ./scripts/setparas.sh peerenv 1 1
export CORE_PEER_ADDRESS=peer1.org1.example.com:8051
peer chaincode install -n $CC_NAME -v $VERSION -l $LANGUAGE -p $CC_SRC_PATH
 
source ./scripts/setparas.sh peerenv 1 2
export CORE_PEER_ADDRESS=peer1.org2.example.com:10051
peer chaincode install -n $CC_NAME -v $VERSION -l $LANGUAGE -p $CC_SRC_PATH
```

初始化chaincode

```bash
peer chaincode instantiate -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME -l $LANGUAGE -v $VERSION -c '{"Args":[""]}' -P "$DEFAULT_POLICY"
```

#### 增改查

节点连接

```bash
export PEER_CONN_PARMS="--peerAddresses peer0.org1.example.com:7051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt --peerAddresses peer0.org2.example.com:9051 --tlsRootCertFiles /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org2.example.com/peers/peer0.org2.example.com/tls/ca.crt"
```

初始化数据（init）

```bash
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"init","Args":[]}'
```

增加数据（create）

```bash
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"create","Args":["ACCOUNT0","0001","0002","0003","0004","420222199804278896","123456","abcdefg","10","11","12","0"]}'
```

修改数据（update），目前只实现修改工作计时(Worktime)，后续会添加其他修改功能

```bash
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"update","Args":["ACCOUNT0","100"]}'
```

查询单个数据（query），用于工人

```bash
peer chaincode query -C $CHANNEL_NAME -n $CC_NAME -c '{"function":"query","Args":["ACCOUNT0"]}'
```

结果如下，转化为json后按首字母对Key进行了排序

```
{"GroupID":"0004","HardwareID":"0002","ID":"0001","IDCard":"420222197812150041","ProjectID":"0003","SalaryPer":50,"Value1":"12bcefdffd4536","Value2":"12fdafe56cb7","Workload":51,"Workprice":20,"Worktime":100}
```

查询所有数据（list），用于管理人员

```bash
peer chaincode invoke -o $ORDERER_ADDRESS --tls --cafile $ORDERER_CA -C $CHANNEL_NAME -n $CC_NAME $PEER_CONN_PARMS -c '{"function":"list","Args":[]}'
```

#### CouchDB

浏览器访问 ip:5984/_utils

![](img\1.JPG)

进入mychannel_checjinsystem即可查看数据

```
mychannel  通道名称
checjinsystem   链码名称
```



#### 未完待续。。。