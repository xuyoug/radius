# radius简介
RFC 2865 协议（radius协议）的部分实现

Author:xuyoug
xuyoug@yeah.net

### 说明
本项目__还在开发过程中，还远没有达到生产可用的水平__,因为完成度太低，注释文档也还没有跟进，现在一片混乱
本人初学者，请多指教，勿喷
注重效率还是便捷性，这是个问题。。。

### 包构架
- radius              
radius报文解包封包基本操作，属性厂商列表定义域
- radius/radiuscli    
radius的客户端实现
- radius/radiusserv   
radius服务端的实现
- radius/radiusfast   
radius报文的快速解包实现  不进行某些rfc2865规定的校验
- radius/radiusiptv   
基于radius协议的iptv协议实现

### 完成进度
- 20150602  
- 20150427  radius 设计、基础定义、构架建设
- 20150512  重新设计 定义 版本升至0.2
- 20150512  radius报文的基本解包封包
- 20150525  能够正常运行，但是还有很多地方需要优化调整
- 20150526  上传代码到github
- 20150526  做一些变更和优化



# radius

### 说明

### 使用方式



# 尾记
学开发，有点痛苦……