package base

import (
	"time"

	"gorm.io/gorm"
)

// UuidModel базовая модель для сущностей которые не отправляются на МУ как справочники
type UuidModel struct {
	Uuid      string         `gorm:"TYPE:uuid;primarykey;default:uuid_generate_v4();uniqueIndex"` // Uuid записи
	CreatedAt time.Time      // Дата создания
	UpdatedAt time.Time      // Дата обновления
	DeletedAt gorm.DeletedAt `gorm:"index"` // Дата Удаления
}

// DictionaryUUIDModel базовая модель для справочников которые отправляются на МУ и потенциально могут иметь размер более 100к элементов
type DictionaryUUIDModel struct {
	Uuid      string         `gorm:"TYPE:uuid;primarykey;default:uuid_generate_v4();uniqueIndex"` // Uuid записи
	CreatedAt time.Time      `gorm:"index"`                                                       // Дата создания
	UpdatedAt time.Time      `gorm:"index"`                                                       // Дата обновления
	DeletedAt gorm.DeletedAt `gorm:"index"`                                                       // Дата Удаления
}

// DictionaryIDModel базовая модель для справочников которые отправляются на МУ и потенциально могут иметь размер более 100к элементов
type DictionaryIDModel struct {
	ID        uint           `gorm:"primarykey"`                      // Id записи
	CreatedAt time.Time      `gorm:"index;default:current_timestamp"` // Дата создания
	UpdatedAt time.Time      `gorm:"index"`                           // Дата обновления
	DeletedAt gorm.DeletedAt `gorm:"index"`                           // Дата Удаления
}
