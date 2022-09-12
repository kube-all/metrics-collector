# 集群指标收集器
根据定时策略收集集群的指标数据，目前支持的维度：namespace，node，pod，service(deployment)，用于分析资源使用情况以及资源预估，
该项目只用于数据收集，数据分析和利用需另开发
支持的数据库：
- [x] mysql
- [ ] clickhouse
## mysql
```
docker run   -p 3307:3306 --name kube_cloud_metrics -e MYSQL_DATABASE=kube_cloud_metrics -e MYSQL_ROOT_PASSWORD=123456 -dit mysql:5.7 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci --character-set-client-handshake=false 
```
## 新增维度
参考namespace维度指标收集
## 维度新增指标
- 在对应维度的model中增加指标字段(如Metrics1)
- 在对应维度的*.yaml(如namespace.yaml)中增加对应的query语句即可，其中的code需与增加的字段名相同(如上Metrics1)
- 编译打包制作镜像运行即可
## prometheus查询语句修改
提供两种方式
- 修改yaml文件，最终同步到数据库(只支持新增)，采集以数据库的为准
- 通过api修改数据库记录


## 应用配置
```
databaseConfig: # 数据库信息
  name: kube_cloud_metrics
  user: root
  password: 123456
  address: localhost:3307
```
