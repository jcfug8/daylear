package model

// User -
type User struct {
	UserId     int64   `gorm:"primaryKey;bigint;not null;<-:false"`
	Email      string  `gorm:"string;not null;uniqueIndex"`
	Username   string  `gorm:"string;not null;uniqueIndex"`
	AmazonId   *string `gorm:"string;->:false;<-:create;uniqueIndex"`
	FacebookId *string `gorm:"string;->:false;<-:create;uniqueIndex"`
	GoogleId   *string `gorm:"string;->:false;<-:create;uniqueIndex"`
}

// TableName -
func (User) TableName() string {
	return "public.user"
}
