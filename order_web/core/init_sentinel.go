package core

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"go.uber.org/zap"
	"order_web/global"
)

func InitSentinel() {
	// 初始化Sentinel
	err := api.InitDefault()
	if err != nil {
		zap.S().Fatalf("Sentinel初始化失败: %v", err)
	}

	//  QPS限流规则
	_, _ = flow.LoadRules([]*flow.Rule{ //订单List
		{
			ID:                     global.Config.Sentinel.LimitFlowID,       // 规则唯一ID，自定义
			Resource:               global.Config.Sentinel.LimitResourceName, // 资源名：必选，对应商品列表接口的标识，要和中间件里一致
			TokenCalculateStrategy: flow.Direct,                              // 固定值：直接限流，按QPS/并发数计算令牌 WarmUp（预热限流，适合秒杀）
			ControlBehavior:        flow.Reject,                              // 限流触发后的行为：快速失败（最适合查询接口）
			Threshold:              global.Config.Sentinel.Threshold,         // 单机QPS阈值，1500
			StatIntervalInMs:       global.Config.Sentinel.StatIntervalInMs,  // 统计间隔：1000ms=1s，固定值
		},
		// orderCreate 限流规则
		{
			ID:                     global.Config.Sentinel.CreateLimitFlowID,
			Resource:               global.Config.Sentinel.CreateLimitResourceName, // 商品详情资源名
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              520, // 详情接口单机QPS阈值
			StatIntervalInMs:       1000,
		},
	})

	// 热点参数限流规则- 防商品接口参数滥用  订单List
	_, _ = hotspot.LoadRules([]*hotspot.Rule{
		{
			ID:                global.Config.Sentinel.LimitFlowID,       // 规则唯一ID
			Resource:          global.Config.Sentinel.LimitResourceName, // 必须和限流规则的Resource一致
			MetricType:        hotspot.QPS,                              // 限流维度：按QPS限流
			ControlBehavior:   hotspot.Reject,                           // 触发后行为：快速失败
			ParamIndex:        0,                                        // 拿第一个参数
			Threshold:         500,                                      // 单个页码的QPS阈值  page = 1 一秒可以500个  page = 2 也是  互不影响
			BurstCount:        20,                                       // 突发流量容忍数，20-50 允许瞬间超过阈值20个请求
			MaxQueueingTimeMs: 0,                                        // 排队超时时间，Reject模式下填0
			DurationInSec:     1,                                        // 统计周期，1s 和限流规则一致
			ParamsMaxCapacity: 100,                                      // 缓存的参数数量，100足够
			SpecificItems: map[interface{}]int64{
				9999: 10, // 极端分页限流
			},
		},
	})

	// 熔断规则

	_, _ = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Id:                           global.Config.Sentinel.FuseFlowID,       //  唯一标识
			Resource:                     global.Config.Sentinel.FuseResourceName, // 服务层调用MySQL的标识，必须和service层一致
			Strategy:                     circuitbreaker.SlowRequestRatio,         // 熔断策略：慢调用比例
			RetryTimeoutMs:               5000,                                    // 熔断5秒（下单接口熔断时长略长）
			MinRequestAmount:             100,                                     // 最小请求数，累计100次调用才触发熔断，防误触
			StatIntervalMs:               1000,                                    // 统计间隔1秒（固定值）
			StatSlidingWindowBucketCount: 1,                                       // 滑动窗口桶数
			MaxAllowedRtMs:               500,                                     // 慢调用阈值，超过200ms的请求算慢调用
			Threshold:                    0.5,                                     // 慢调用比例阈值，50%的请求是慢调用就熔断
			ProbeNum:                     10,                                      // 半开状态试探请求数，放行10个请求判断是否恢复
		},
		{
			Id:                           global.Config.Sentinel.FuseErrFlowID, //  唯一标识
			Resource:                     global.Config.Sentinel.FuseErrResourceName,
			Strategy:                     circuitbreaker.ErrorRatio, // 熔断策略 错误率
			RetryTimeoutMs:               3000,                      // =熔断时长3秒
			MinRequestAmount:             80,                        // 最小请求数80
			StatIntervalMs:               1000,                      // 统计间隔1秒
			StatSlidingWindowBucketCount: 1,                         // 默认1
			MaxAllowedRtMs:               0,                         // 错误率策略下，该字段无效，填0
			Threshold:                    0.3,                       // 错误率阈值，30%的请求报错就熔断
			ProbeNum:                     10,                        // 半开状态试探10个请求
		},
	})

}
