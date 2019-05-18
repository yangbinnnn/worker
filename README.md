## worker

简单工作任务池

## 接口
- NewWorkPool 新建工作队列
- Start 开启工作队列
- Push 添加新任务
- WaitAndClose 等待任务完成, 并关闭工作队列
- IsRunning 查看是否运行中

## 如何使用
实现Worker 接口即可, 参考[Demo](./worker_test.go)
