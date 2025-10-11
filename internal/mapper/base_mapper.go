package mapper

// BaseMapper 基础映射器接口
type BaseMapper interface {
	// Save 保存实体
	Save(entity interface{}) error

	// DeleteByID 根据ID删除实体
	DeleteByID(id interface{}) error

	// GetByID 根据ID获取实体
	GetByID(id interface{}) (interface{}, error)

	// GetAll 获取所有实体
	GetAll() ([]interface{}, error)

	// Update 更新实体
	Update(entity interface{}) error
}