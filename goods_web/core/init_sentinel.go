package core

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/core/hotspot"
	"go.uber.org/zap"
	"goods_web/global"
)

func InitSentinel() {
	// 初始化Sentinel
	err := api.InitDefault()
	if err != nil {
		zap.S().Fatalf("Sentinel初始化失败: %v", err)
	}

	//  QPS限流规则
	_, _ = flow.LoadRules([]*flow.Rule{
		{
			ID:                     global.Config.Sentinel.LimitFlowID,       // 规则唯一ID，自定义
			Resource:               global.Config.Sentinel.LimitResourceName, // 资源名：必选，对应商品列表接口的标识，要和中间件里一致
			TokenCalculateStrategy: flow.Direct,                              // 固定值：直接限流，按QPS/并发数计算令牌 WarmUp（预热限流，适合秒杀）
			ControlBehavior:        flow.Reject,                              // 限流触发后的行为：快速失败（最适合查询接口）
			Threshold:              global.Config.Sentinel.Threshold,         // 核心参数：单机QPS阈值，8核16G服务器建议1500-3000
			StatIntervalInMs:       global.Config.Sentinel.StatIntervalInMs,  // 统计间隔：1000ms=1s，固定值
		},
		// goodsDetail限流规则（详情接口QPS可略低，比如1000）
		{
			ID:                     "good_Detail_flow_rule",
			Resource:               "api:good:Detail", // 商品详情资源名
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              1000, // 详情接口单机QPS阈值
			StatIntervalInMs:       1000,
		},
	})

	// 热点参数限流规则- 防商品接口参数滥用
	_, _ = hotspot.LoadRules([]*hotspot.Rule{
		{
			ID:                global.Config.Sentinel.LimitFlowID,       // 规则唯一ID
			Resource:          global.Config.Sentinel.LimitResourceName, // 必须和限流规则的Resource一致
			MetricType:        hotspot.QPS,                              // 限流维度：按QPS限流（商品列表必选）
			ControlBehavior:   hotspot.Reject,                           // 触发后行为：快速失败
			ParamIndex:        0,                                        // 拦截第0个入参（categoryId 分类ID）
			Threshold:         800,                                      // 单个分类ID的QPS阈值
			BurstCount:        20,                                       // 突发流量容忍数，20-50 允许瞬间超过阈值20个请求
			MaxQueueingTimeMs: 0,                                        // 排队超时时间，Reject模式下填0
			DurationInSec:     1,                                        // 统计周期，1s 和限流规则一致
			ParamsMaxCapacity: 200,                                      // 缓存的参数数量，200足够
			// 重点：爆款分类单独限流（比如女装分类ID=1，单独设置阈值更低）
			SpecificItems: map[interface{}]int64{
				1: 300, // 分类ID=1的QPS阈值只有300，防止爆款分类压垮数据库 只能手机
				5: 400, // 分类ID=2的QPS阈值400   耳机音响
			},
		},
		// goodsDetail热点规则（入参是goodsId，重点防爆款商品）
		{
			ID:                "good_Detail_flow_rule",
			Resource:          "api:good:Detail",
			MetricType:        hotspot.QPS,
			ControlBehavior:   hotspot.Reject,
			ParamIndex:        0,   // 第0个入参：goodsI
			Threshold:         500, // 单个商品ID的QPS阈值
			BurstCount:        10,
			MaxQueueingTimeMs: 0,
			DurationInSec:     1,
			ParamsMaxCapacity: 500, // 缓存更多商品ID
			SpecificItems: map[interface{}]int64{ // 爆款商品单独限流
				1001: 100, // 爆款商品ID=1001，QPS仅允许100
				2002: 150, // 爆款耳机ID=2002，QPS允许150
			},
		},
	})

	// 熔断规则

	_, _ = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Id:                           global.Config.Sentinel.FuseFlowID,       //  唯一标识
			Resource:                     global.Config.Sentinel.FuseResourceName, // 服务层调用MySQL的标识，必须和service层一致
			Strategy:                     circuitbreaker.SlowRequestRatio,         // 熔断策略：慢调用比例
			RetryTimeoutMs:               3000,                                    // 熔断时长3秒，熔断后3秒内不调用MySQL
			MinRequestAmount:             100,                                     // 最小请求数，累计100次调用才触发熔断，防误触
			StatIntervalMs:               1000,                                    // 统计间隔1秒（固定值）
			StatSlidingWindowBucketCount: 1,                                       // 滑动窗口桶数
			MaxAllowedRtMs:               200,                                     // 慢调用阈值，超过200ms的请求算慢调用
			Threshold:                    0.5,                                     // 慢调用比例阈值，50%的请求是慢调用就熔断
			ProbeNum:                     10,                                      // 半开状态试探请求数，放行10个请求判断是否恢复
		},
		{
			Id:                           global.Config.Sentinel.FuseErrFlowID, //  唯一标识
			Resource:                     global.Config.Sentinel.FuseErrResourceName,
			Strategy:                     circuitbreaker.ErrorRatio, // 熔断策略 错误率
			RetryTimeoutMs:               3000,                      // =熔断时长3秒
			MinRequestAmount:             100,                       // 最小请求数100
			StatIntervalMs:               1000,                      // 统计间隔1秒
			StatSlidingWindowBucketCount: 1,                         // 默认1
			MaxAllowedRtMs:               0,                         // 错误率策略下，该字段无效，填0
			Threshold:                    0.3,                       // 错误率阈值，30%的请求报错就熔断
			ProbeNum:                     10,                        // 半开状态试探10个请求
		},
	})

}
