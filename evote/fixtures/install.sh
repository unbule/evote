#!/bin/bash

CONFIG_ROOT=/opt/gopath/src/github.com/hyperledger/fabric/peer
ORG1_MSPCONFIGPATH=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp
ORG1_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/peerOrganizations/org1.example.com/peers/peer0.org1.example.com/tls/ca.crt
ORDERER_TLS_ROOTCERT_FILE=${CONFIG_ROOT}/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

#创建通道
echo "通道开始创建"
docker exec cli \
  peer channel create \
    -o orderer.example.com:7050 \
    -c mychannel \
    -f ${CONFIG_ROOT}/channel-artifacts/channel.tx \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE}
sleep 10

#加入通道
echo "节点加入通道"
docker exec \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
  -e CORE_PEER_LOCALMSPID="Org1MSP" \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
  cli \
  peer channel join \
    -b mychannel.block

#更新锚节点
echo "开始更新锚节点"
docker exec \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
  -e CORE_PEER_LOCALMSPID="Org1MSP" \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
  cli \
  peer channel update \
    -o orderer.example.com:7050 \
    -c mychannel \
    -f ${CONFIG_ROOT}/channel-artifacts/Org1MSPanchors.tx \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE}

#安装链码
echo "安装链码"
docker exec \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
  -e CORE_PEER_LOCALMSPID="Org1MSP" \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode install \
    -n mycc \
    -v 1.0 \
    -p github.com/chaincode

#初始化链码
echo "初始化链码"
docker exec \
  -e CORE_PEER_MSPCONFIGPATH=${ORG1_MSPCONFIGPATH} \
  -e CORE_PEER_ADDRESS=peer0.org1.example.com:7051 \
  -e CORE_PEER_LOCALMSPID="Org1MSP" \
  -e CORE_PEER_TLS_ROOTCERT_FILE=${ORG1_TLS_ROOTCERT_FILE} \
  cli \
  peer chaincode instantiate \
    -o orderer.example.com:7050 \
    --tls \
    --cafile ${ORDERER_TLS_ROOTCERT_FILE} \
    -C mychannel \
    -n mycc \
    -v 1.0 \
    -c '{"Args":["init","a", "100", "b","200"]}'\
    -P "AND ('Org1MSP.peer')"