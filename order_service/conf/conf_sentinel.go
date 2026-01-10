package conf

type Sentinel struct {
	LimitResourceName       string  `mapstructure:"resource_name" json:"resource_name"`
	LimitFlowID             string  `mapstructure:"limit_flow_id" json:"limit_flow_id"`
	CreateLimitFlowID       string  `mapstructure:"create_limit_flow_id" json:"create_limit_flow_id"`
	CreateLimitResourceName string  `mapstructure:"create_limit_resource_name" json:"create_limit_resource_name"`
	Threshold               float64 `mapstructure:"threshold" json:"threshold"`
	StatIntervalInMs        uint32  `mapstructure:"stat_interval_in_ms" json:"stat_interval_in_ms"`
	FuseResourceName        string  `mapstructure:"fuse_resource_name" json:"fuse_resource_name"`
	FuseFlowID              string  `mapstructure:"fuse_flow_id" json:"fuse_flow_id"`
	FuseErrResourceName     string  `mapstructure:"fuse_err_resource_name" json:"fuse_err_resource_name"`
	FuseErrFlowID           string  `mapstructure:"fuse_err_flow_id" json:"fuse_err_flow_id"`
}
