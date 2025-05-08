package model

// User -
type User struct {
	UserId     int64   `gorm:"primaryKey;bigint;not null;<-:false"`
	Email      string  `gorm:"size:255;not null;uniqueIndex"`
	Username   string  `gorm:"size:50;not null;uniqueIndex"`
	GivenName  string  `gorm:"size:100"`
	FamilyName string  `gorm:"size:100"`
	AmazonId   *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`
	FacebookId *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`
	GoogleId   *string `gorm:"size:255;->:false;<-:create;uniqueIndex"`
}

// TableName -
func (User) TableName() string {
	return "daylear_user"
}
