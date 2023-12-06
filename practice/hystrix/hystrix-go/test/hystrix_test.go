package test

import (
	"errors"
	"github.com/afex/hystrix-go/hystrix"
	"log"
	"net/http"
	"testing"
	"time"
)

func Test_Func1(t *testing.T) {
	//1.创建hystrix流量统计服务器
	//(实际可以理解为启动一个专门的hystrix服务用来管理熔断器,
	//通过熔断器监控指定的业务接口的流量情况,当启动此处的hystrix流控服务后可以访问这个地址,
	//例如此处访问:127.0.0.1:8074 就可以拿到监控的访问情况)
	hystrixStreamHandler := hystrix.NewStreamHandler()
	hystrixStreamHandler.Start()
	go http.ListenAndServe(":8074", hystrixStreamHandler)

	//2.配置限流规则hystrix.CommandConfig
	commandConfig := hystrix.CommandConfig{
		Timeout:                1000, //单次请求超时时间,默认时间是1000毫秒
		MaxConcurrentRequests:  1,    // 最大并发量,默认值是10(注意此处并不是设置为1就是1)
		SleepWindow:            5000, // 熔断后多久去尝试服务是否可用,默认值是5000毫秒(熔断器打开到半打开的时间)
		RequestVolumeThreshold: 1,    //一个统计窗口10秒内请求数量。达到这个请求数量后才去判断是否要开启熔断,默认值是20
		ErrorPercentThreshold:  1,    //错误百分比,默认值是50(当错误百分比超过这个限制时则进行熔断)
	}

	//3.设置熔断器
	//第一个参数:当前创建的熔断器名称
	//第二个参数: hystrix.CommandConfig配置的限流规则
	hystrix.ConfigureCommand("aaa", commandConfig)

	//4.hystrix.Do()同步: hystrix-go与业务方法的同步整合
	//假设当前有100个请求依次进来
	for i := 0; i < 100; i++ {

		//5.通过hystrix.Do()对业务方法实现流控
		//hystrix.Do需要三个参数
		//参数一: 指定使用哪个熔断器,设置对应的熔断器名称,例如上面创建了"aaa"熔断器,当前则可以设置"aaa"
		//参数二: 需要被流控的业务方法
		//参数三: 降级方法,当代码返回一个错误时,或者当它基于各种健康检查无法完成时,就会触发此事件
		//对应Do()方法hystrix还提供了一个Doc(),Do()内部实际调用的就是Doc()
		//Doc()中多了一个Context,思考一下能否使用Context传递参数与接收响应
		err := hystrix.Do("aaa", func() error {
			//test case 1 并发测试
			if i == 0 {
				return errors.New("service error")
			}
			//test case 2 超时测试
			//time.Sleep(2 * time.Second)
			log.Println("do services")
			return nil
		}, nil)

		//6.hystrix.Do()执行完毕后返回的error,该error会影响熔断
		//返回的error会记录到熔断器的错误百分比中,当超过设置阈值则触发熔断,
		//熔断后当到达熔断器设置的SleepWindow时间后,开始尝试恢复
		if err != nil {
			log.Println("hystrix err:" + err.Error())
			time.Sleep(1 * time.Second)
			log.Println("sleep 1 second")
		}
	}
	time.Sleep(100 * time.Second)
}
