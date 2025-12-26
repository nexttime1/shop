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

func (t ImageType) Value() (driver.Value, error) {
	return int64(t), nil // 改为返回int64，符合driver.Value规范
}

// Scan 实现 sql.Scanner 接口：将数据库类型（tinyint）转 Go 类型
func (t *ImageType) Scan(value interface{}) error {
	// 兼容更多类型（比如int/uint），避免转换失败
	var val int64
	switch v := value.(type) {
	case int64:
		val = v
	case int:
		val = int64(v)
	case uint:
		val = int64(v)
	case uint64:
		val = int64(v)
	default:
		zap.S().Errorf("图片类型转换失败，不支持的类型: %T, 值: %v", value, value)
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
