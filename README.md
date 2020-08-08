# golangLogEvent
# 环境
- mac
- logagent2 需要用到 kafka、etcd
- log_transfer 需要用到 kafka、es、kibana
----------
# 环境安装及配置
1. 使用 brew 安装
2. brew install kafka （可能会缺少zookeeper 再 brew install zookeeper）
   
   https://segmentfault.com/a/1190000022561851
   
   需要更改下配置信息  vim  /usr/local/etc/kafka/server.properties 取消掉 9092端口的注释
   
   启动： zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
   
   日志路径： /usr/local/var/lib/kafka-logs/
   
   启动kafka消费者命令：kafka-console-consumer --bootstrap-server=127.0.0.1:9092 --topic=web_log --from-beginning
   
3. brew install etcd 安装etcd

   https://learnku.com/articles/42515
   
   启动：etcd
   
   存放：etcdctl --endpoints=http://127.0.0.1:2379 put key value
   
   查询：etcdctl --endpoints=http://127.0.0.1:2379 get key
   
4. brew install elasticsearch 安装es

  访问：http://127.0.0.1:9200/
  
5. brew install kibana 安装kibana

   可以更改显示为中文  修改配置文件 kibana.yml 最下面改为 ”cn-ZH“

   访问 http://127.0.0.1:5601
---------
# 启动顺序
```
  启动 kafka
  zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
  启动 etcd 
  etcd
  启动 es
  elasticsearch
  启动 kibana
  kibana
  
```
-----------
# 思路
```
  logagent
  0. 首先要在 etcd 中存储 key 与 json（ logpath topic）的内容 （ 可利用demo包中的 etcd_put 存入数据)
  1. 首先从 ini 读取 etcd kafka 配置信息（地址，etcd key）
  2. 初始化连接，读取 etcd 中 key 对应的日志路径 path，利用 tail 包读取日志内容，异步通信方式将其写入到 kafka 对应的 topic 中
 
```
```
  logtransfer
  0. 首先要保证 kafka 对应的topic中有数据，可供读取
  1. 从 ini 读取，kafka etcd 配置信息（地址， topic）
  2. 初始化连接，读取 kafka 中 topic 对应的内容，异步通信通过 channel 传输es，开启多个go程写入es
```
   
    
