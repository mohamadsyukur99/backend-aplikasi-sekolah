package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"github.com/mohamadsyukur99/fullstack/api/security"
)

// User ...
type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null;unique" json:"username"`
	Name      string    `gorm:"size:255;not null;unique" json:"name"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	Level     string    `gorm:"size:100;not null;" json:"level"`
	Status    string    `gorm:"size:100;not null;" json:"status"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeSave ...
func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// Prepare ...
func (u *User) Prepare() {
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Name = html.EscapeString(strings.TrimSpace(u.Name))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.Level = html.EscapeString(strings.TrimSpace(u.Level))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// Validate ...
func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		// if u.Nickname == "" {
		// 	err = errors.New("Required Nickname")
		// 	errorMessages["Required_nickname"] = err.Error()
		// }
		// if u.Password == "" {
		// 	err = errors.New("Required Password")
		// 	errorMessages["Required_password"] = err.Error()
		// }
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}

	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Username == "" {
			err = errors.New("Required Username")
			errorMessages["Required_username"] = err.Error()

		}
	default:

		if u.Username == "" {
			err = errors.New("Required Username")
			errorMessages["Required_username"] = err.Error()
		}

		if u.Name == "" {
			err = errors.New("Required Nickname")
			errorMessages["Required_nickname"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages
}

// SaveUser ...
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// FindAllUsers ...
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Where("status=1").Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

// FindUserByID ...
func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// GetUserByName ...
func (u *User) GetUserByNames(db *gorm.DB, nama string) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Where("name LIKE ? OR username LIKE ? AND STATUS = ?", `%`+nama+`%`, `%`+nama+`%`, "1").Take(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]User{}, errors.New("User Not Found")
	}
	return &users, err
}

// GetUserByName ...
func (u *User) GetUserUsername(db *gorm.DB, nama string) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Where("username = ? ", nama).Take(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]User{}, errors.New("User Not Found")
	}
	return &users, err
}

// UpdateAUser ...
func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {
	// to has the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumn(
		map[string]interface{}{
			"name":       u.Name,
			"email":      u.Email,
			"level":      u.Level,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// UpdateAUserPassword ...
func (u *User) UpdateAUserPassword(db *gorm.DB, uid uint32) (*User, error) {
	// to has the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumn(
		map[string]interface{}{
			"password":   u.Password,
			"name":       u.Name,
			"email":      u.Email,
			"level":      u.Level,
			"updated_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// DeleteAUser ...
func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// UpdatePassword ...
func (u *User) UpdatePassword(db *gorm.DB) error {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
