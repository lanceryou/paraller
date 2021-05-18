# paraller
一个简单的异步库，支持用同步方式写异步代码

##  安装方式

```go
go get -u github.com/lanceryou/paraller
```

## paraller api

+ Async 返回paraller对象，异步的去调用函数
+ Await 只能在异步代码中使用，会阻塞等待结果

## 开始使用

```go

s1 := paraller.Async(func()error{
	// do something
})

s2 := paraller.Async(func()error{
// do something
})

s3 := paraller.Async(func() error{
    if err := s1.Await(); err != nil{
		return err
    }
    
    if err := s2.Await(); err != nil{
        return err
    }
})

err := s3.Await()
// do something
```