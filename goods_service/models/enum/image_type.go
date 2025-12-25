package enum

import (
	"database/sql/driver"
	"errors"
	"go.uber.org/zap"
)

type ImageType int

const (
	MainImageType   ImageType = 1
	DetailImageType ImageType = 2
	OtherImageType  ImageType = 3
)

// Value 实现 driver.Valuer 接口：将 Go 类型转数据库类型（tinyint）
func (t ImageType) Value() (driver.Value, error) {
	return int(t), nil
}

// Scan 实现 sql.Scanner 接口：将数据库类型（tinyint）转 Go 类型
func (t *ImageType) Scan(value interface{}) error {
	val, ok := value.(int64)
	if !ok {
		zap.S().Error("转换失败")
		return errors.New("invalid image type value")
	}
	*t = ImageType(val)
	return nil
}

// IsValid 校验图片类型是否合法（oneof 1/2/3）
func (t ImageType) IsValid() bool {
	return t == MainImageType || t == DetailImageType || t == OtherImageType
}

// String 转字符串描述（方便日志打印/返回给前端）
func (t ImageType) String() string {
	switch t {
	case MainImageType:
		return "main"
	case DetailImageType:
		return "detail"
	case OtherImageType:
		return "other"
	default:
		return "unknown"
	}
}
