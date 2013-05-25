## Sdfs
=======

Sdfs是一个 Golang轻量级分布式文件存储分发系统


##licensed

Sdfs is licensed under the BSD Licence


## Install
============
    go get github.com/insionng/sdfs


## 设计思路
============
在web前端做负载均衡，通过302状态对图片路径进行转发到不同的img server上；
让浏览器访问图片的时候都做一次跳转（或永久重定向）。
这样可以配合后端的服务器压力状态进行动态分配。
