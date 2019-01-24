# flwoer

* web服务基于[gin](https://github.com/gin-gonic/gin)
* 异步分布式任务队列基于[machinery](https://github.com/RichardKnop/machinery)
* 任务调度基于[robfig/cron](https://github.com/robfig/cron)


---
### Build

1. Web Service 接口服务
    ```
    sh web.sh
    ```
2. Worker Service 任务执行服务，基于redis or rabbitmq 分布式执行任务

3. cronJob Service 任务分发服务。


### 文件夹约定
1. web 
    ```
    存放web相关
    ```
2. jobs
    ```
    自定义任务
    ``` 
