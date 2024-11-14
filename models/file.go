package models

type FileModel struct {
	BaseModel      `gorm:"embedded"`
	FileName       string            `gorm:"column:file_name;type:varchar(255);comment:文件名" json:"fileName"`
	SourceFileName string            `gorm:"column:source_file_name;type:varchar(255);comment:原始文件名" json:"sourceFileName"`
	FilePath       string            `gorm:"column:file_path;type:varchar(255);comment:文件路径" json:"filePath"`
	Type           string            `gorm:"column:type;type:varchar(255);comment:文件类型" json:"type"`
	Size           int64             `gorm:"column:size;type:bigint;comment:文件大小" json:"size"`
	SliceCount     int               `gorm:"column:slice_count;type:int;comment:文件切片数量" json:"sliceCount"`
	UserId         string            `gorm:"column:user_id;type:varchar(64);comment:上传文件用户id" json:"userId"`
	FileChunkModel []*FileChunkModel `gorm:"foreignKey:FileModelId;constraint:OnDelete:CASCADE;" json:"fileChunkModel"`
}

func (f *FileModel) TableName() string {
	return "file_model"
}

type FileChunkModel struct {
	BaseModel   `gorm:"embedded"`
	FileName    string `gorm:"column:file_name;type:varchar(255);comment:文件切片名" json:"fileName"`
	Index       int    `gorm:"column:index;type:int;comment:文件切片索引" json:"index"`
	FilePath    string `gorm:"column:file_path;type:varchar(255);comment:文件切片路径" json:"filePath"`
	FileModelId string `gorm:"column:file_model_id;type:varchar(64);comment:上传文件id" json:"fileModelId"`
}

func (f *FileChunkModel) TableName() string {
	return "file_chunk_model"
}
