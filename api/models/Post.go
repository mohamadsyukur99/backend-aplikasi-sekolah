package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Post ...
type Post struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	Content   string    `gorm:"size:255;not null;" json:"content"`
	Author    User      `json:"author"`
	AuthorID  uint32    `gorm:"not null" json:"author_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare ...
func (p *Post) Prepare() {
	p.ID = 0
	p.Title = html.EscapeString(strings.TrimSpace(p.Title))
	p.Content = html.EscapeString(strings.TrimSpace(p.Content))
	p.Author = User{}
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()
}

// Validate ...
func (p *Post) Validate() map[string]string {
	var err error

	var errorMessages = make(map[string]string)

	if p.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()

	}
	if p.Content == "" {
		err = errors.New("Required Content")
		errorMessages["Required_content"] = err.Error()

	}

	if p.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}
	return errorMessages
}

// SavePost ...
func (p *Post) SavePost(db *gorm.DB) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Create(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// FindAllPosts ...
func (p *Post) FindAllPosts(db *gorm.DB) (*[]Post, error) {
	var err error
	posts := []Post{}
	err = db.Debug().Model(&Post{}).Limit(100).Order("created_at desc").Find(&posts).Error
	if err != nil {
		return &[]Post{}, err
	}
	if len(posts) > 0 {
		for i := range posts {
			err := db.Debug().Model(&User{}).Where("id = ?", posts[i].AuthorID).Take(&posts[i].Author).Error
			if err != nil {
				return &[]Post{}, err
			}
		}
	}
	return &posts, nil
}

// FindPostByID ...
func (p *Post) FindPostByID(db *gorm.DB, pid uint64) (*Post, error) {
	var err error
	err = db.Debug().Model(&Post{}).Where("id = ?", pid).Take(&p).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// UpdateAPost ...
func (p *Post) UpdateAPost(db *gorm.DB) (*Post, error) {
	var err error

	err = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Updates(Post{Title: p.Title, Content: p.Content, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Post{}, err
	}
	if p.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", p.AuthorID).Take(&p.Author).Error
		if err != nil {
			return &Post{}, err
		}
	}
	return p, nil
}

// DeleteAPost ...
func (p *Post) DeleteAPost(db *gorm.DB) (int64, error) {
	db = db.Debug().Model(&Post{}).Where("id = ?", p.ID).Take(&Post{}).Delete(&Post{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

//DeleteUserPosts When a user is deleted, we also delete the post that the user had
func (p *Post) DeleteUserPosts(db *gorm.DB, uid uint32) (int64, error) {
	posts := []Post{}
	db = db.Debug().Model(&Post{}).Where("author_id = ?", uid).Find(&posts).Delete(&posts)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
