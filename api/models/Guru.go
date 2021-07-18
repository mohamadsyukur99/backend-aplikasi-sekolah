package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Guru ...
type Guru struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nip           string    `gorm:"size:255;not null;unique" json:"nip"`
	Nama          string    `gorm:"size:255;not null;unique" json:"nama"`
	Jenis_kelamin string    `gorm:"size:255;not null;unique" json:"jenis_kelamin"`
	Tempat_lahir  string    `gorm:"size:100;not null;unique" json:"tempat_lahir"`
	Tanggal_lahir string    `gorm:"size:100;not null;" json:"tanggal_lahir"`
	Alamat        string    `gorm:"size:100;not null;" json:"alamat"`
	Agama         string    `gorm:"size:100;not null;" json:"agama"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare ...
func (g *Guru) Prepare() {
	g.Nip = html.EscapeString(strings.TrimSpace(g.Nip))
	g.Nama = html.EscapeString(strings.TrimSpace(g.Nama))
	g.Tempat_lahir = html.EscapeString(strings.TrimSpace(g.Tempat_lahir))
	g.Tanggal_lahir = html.EscapeString(strings.TrimSpace(g.Tanggal_lahir))
	g.Alamat = html.EscapeString(strings.TrimSpace(g.Alamat))
	g.Agama = html.EscapeString(strings.TrimSpace(g.Alamat))
	g.CreatedAt = time.Now()
	g.UpdatedAt = time.Now()
}

// SaveGuru ...
func (g *Guru) SaveGuru(db *gorm.DB) (*Guru, error) {
	var err error
	err = db.Debug().Create(&g).Error
	if err != nil {
		return &Guru{}, err
	}
	return g, nil
}

// FindAllGuru ...
func (g *Guru) FindAllGuru(db *gorm.DB) (*[]Guru, error) {
	var err error
	dataguru := []Guru{}
	err = db.Debug().Model(&Guru{}).Where("status=1").Limit(100).Find(&dataguru).Error
	if err != nil {
		return &[]Guru{}, err
	}
	return &dataguru, err
}

// FindGuruByID ...
func (g *Guru) FindGuruByID(db *gorm.DB, uid uint32) (*Guru, error) {
	var err error
	err = db.Debug().Model(Guru{}).Where("id = ?", uid).Take(&g).Error
	if err != nil {
		return &Guru{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Guru{}, errors.New("User Not Found")
	}
	return g, err
}

// GetGuruByNames ...
func (g *Guru) GetGuruByNames(db *gorm.DB, nama string) (*[]Guru, error) {
	var err error
	dataguru := []Guru{}
	err = db.Debug().Model(&Guru{}).Where("nama LIKE ? OR nip LIKE ? AND STATUS = ?", `%`+nama+`%`, `%`+nama+`%`, "1").Take(&dataguru).Error
	if err != nil {
		return &[]Guru{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]Guru{}, errors.New("Guru Not Found")
	}
	return &dataguru, err
}

// GetGuruNip ...
func (g *Guru) GetGuruNip(db *gorm.DB, nip string) (*[]Guru, error) {
	var err error
	dataguru := []Guru{}
	err = db.Debug().Model(&Guru{}).Where("nip = ? ", nip).Take(&dataguru).Error
	if err != nil {
		return &[]Guru{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]Guru{}, errors.New("Siswa Not Found")
	}
	return &dataguru, err
}

// UpdateAGuru ...
func (g *Guru) UpdateAGuru(db *gorm.DB, uid uint32) (*Guru, error) {

	db = db.Debug().Model(&Guru{}).Where("id = ?", uid).Take(&Guru{}).UpdateColumn(
		map[string]interface{}{
			"nama":          g.Nama,
			"tempat_lahir":  g.Tempat_lahir,
			"Jenis_kelamin": g.Jenis_kelamin,
			"tanggal_lahir": g.Tanggal_lahir,
			"alamat":        g.Alamat,
			"agama":         g.Agama,
			"updated_at":    time.Now(),
		},
	)
	if db.Error != nil {
		return &Guru{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Guru{}).Where("id = ?", uid).Take(&g).Error
	if err != nil {
		return &Guru{}, err
	}
	return g, nil
}

// DeleteAGuru ...
func (g *Guru) DeleteAGuru(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Guru{}).Where("id = ?", uid).Take(&Guru{}).Delete(&Guru{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
