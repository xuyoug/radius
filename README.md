# radius
rfc2865协议（radius协议）的部分实现


radius              //radius报文解包封包基本操作，属性厂商列表定义域
radius/radiuscli    //radius的客户端实现
radius/radiusserv   //radius服务端的实现
radius/radiusfast   //radius报文的快速解包实现  不进行某些rfc2865规定的校验
radius/radiusiptv   //基于radius协议的iptv协议实现
