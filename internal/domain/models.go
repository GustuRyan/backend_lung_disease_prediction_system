package domain

import "gorm.io/gorm"

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"size:120;not null"`
	Email    string `json:"email" gorm:"size:255;uniqueIndex;not null"`
	Password string `json:"-" gorm:"size:255;not null"`
	Role     string `json:"role" gorm:"size:32;not null;default:user"`
}

type Disease struct {
	gorm.Model
	DiseaseName string `json:"disease_name" gorm:"size:160;uniqueIndex;not null"`
	Description string `json:"description" gorm:"type:text"`
}

type Dataset struct {
	gorm.Model
	UserID      uint   `json:"user_id" gorm:"not null;index"`
	User        User   `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DatasetName string `json:"dataset_name" gorm:"size:160;not null"`
	FilePath    string `json:"file_path" gorm:"size:1024;not null"`
	Description string `json:"description" gorm:"type:text"`
}

type Model struct {
	gorm.Model
	ModelName string  `json:"model_name" gorm:"size:160;not null"`
	FilePath  string  `json:"file_path" gorm:"size:1024;not null"`
	Accuracy  float64 `json:"accuracy" gorm:"type:decimal(7,6);not null"`
	Type      string  `json:"type" gorm:"size:80;not null"`
	Active    bool    `json:"active" gorm:"not null;default:false;index"`
}

type Recommendation struct {
	gorm.Model
	DiseaseID          uint    `json:"disease_id" gorm:"not null;index"`
	Disease            Disease `json:"disease,omitempty" gorm:"foreignKey:DiseaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	Type               string  `json:"type" gorm:"size:80;not null"`
	RecommendationText string  `json:"recommendation_text" gorm:"type:text;not null"`
}

type History struct {
	gorm.Model
	UserID           uint    `json:"user_id" gorm:"not null;index"`
	User             User    `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	DiseaseID        uint    `json:"disease_id" gorm:"not null;index"`
	Disease          Disease `json:"disease,omitempty" gorm:"foreignKey:DiseaseID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
	ImagePath        string  `json:"image_path" gorm:"size:1024;not null"`
	PredictionResult string  `json:"prediction_result" gorm:"size:160;not null"`
	ConfidenceResult float64 `json:"confidence_result" gorm:"type:decimal(7,6);not null"`
}

type DatasetModel struct {
	gorm.Model
	DatasetID uint    `json:"dataset_id" gorm:"not null;uniqueIndex:idx_dataset_model"`
	Dataset   Dataset `json:"dataset,omitempty" gorm:"foreignKey:DatasetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ModelID   uint    `json:"model_id" gorm:"not null;uniqueIndex:idx_dataset_model"`
	MLModel   Model   `json:"model,omitempty" gorm:"foreignKey:ModelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &Disease{}, &Dataset{}, &Model{}, &Recommendation{}, &History{}, &DatasetModel{})
}
