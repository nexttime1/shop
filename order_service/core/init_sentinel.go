package core

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
	"order_service/global"
)

func InitSentinel() {
	// 初始化Sentinel
	err := api.InitDefault()
	if err != nil {
		zap.S().Fatalf("Sentinel初始化失败: %v", err)
	}

	// 1. 加载服务层限流规则
	_, _ = flow.LoadRules([]*flow.Rule{ // order list
		{
			ID:                     global.Config.Sentinel.LimitFlowID,
			Resource:               global.Config.Sentinel.LimitResourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              global.Config.Sentinel.Threshold, // 服务层QPS阈值，比Web层低
			StatIntervalInMs:       global.Config.Sentinel.StatIntervalInMs,
		},
		{ // orderCreate
			ID:                     global.Config.Sentinel.CreateLimitFlowID,
			Resource:               global.Config.Sentinel.CreateLimitResourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              400, // 服务层QPS阈值，比Web层低
			StatIntervalInMs:       global.Config.Sentinel.StatIntervalInMs,
		},
	})

	// 加载服务层熔断规则
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
