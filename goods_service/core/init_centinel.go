package core

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
	"goods_service/global"
)

func InitSentinel() {
	// 初始化Sentinel
	err := api.InitDefault()
	if err != nil {
		zap.S().Fatalf("Sentinel初始化失败: %v", err)
	}

	// 1. 加载服务层限流规则
	_, _ = flow.LoadRules([]*flow.Rule{
		{
			ID:                     global.Config.Sentinel.LimitFlowID,
			Resource:               global.Config.Sentinel.LimitResourceName,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              global.Config.Sentinel.Threshold, // 服务层QPS阈值，比Web层低
			StatIntervalInMs:       global.Config.Sentinel.StatIntervalInMs,
		},
	})

	// 加载服务层熔断规则
	_, _ = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Id:                           global.Config.Sentinel.FuseFlowID,
			Resource:                     global.Config.Sentinel.FuseResourceName,
			Strategy:                     circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             100,
			StatIntervalMs:               1000,
			StatSlidingWindowBucketCount: 1,
			MaxAllowedRtMs:               200,
			Threshold:                    0.5,
			ProbeNum:                     10,
		},
		{
			Id:                           global.Config.Sentinel.FuseErrFlowID,
			Resource:                     global.Config.Sentinel.FuseErrResourceName,
			Strategy:                     circuitbreaker.ErrorRatio,
			RetryTimeoutMs:               3000,
			MinRequestAmount:             100,
			StatIntervalMs:               1000,
			StatSlidingWindowBucketCount: 1,
			MaxAllowedRtMs:               0,
			Threshold:                    0.3,
			ProbeNum:                     10,
		},
	})
}
