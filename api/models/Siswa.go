package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Siswa ...
type Siswa struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	No_induk      string    `gorm:"size:255;not null;unique" json:"no_induk"`
	Nama          string    `gorm:"size:255;not null;unique" json:"nama"`
	Jenis_kelamin string    `gorm:"size:255;not null;unique" json:"jenis_kelamin"`
	Tempat_lahir  string    `gorm:"size:100;not null;unique" json:"tempat_lahir"`
	Tanggal_lahir string    `gorm:"size:100;not null;" json:"tanggal_lahir"`
	Nama_wali     string    `gorm:"size:100;not null;" json:"nama_wali"`
	Alamat        string    `gorm:"size:100;not null;" json:"alamat"`
	Agama         string    `gorm:"size:100;not null;" json:"agama"`
	Kelas         string    `gorm:"size:100;not null;" json:"kelas"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare ...
func (s *Siswa) Prepare() {
	s.No_induk = html.EscapeString(strings.TrimSpace(s.No_induk))
	s.Nama = html.EscapeString(strings.TrimSpace(s.Nama))
	s.Tempat_lahir = html.EscapeString(strings.TrimSpace(s.Tempat_lahir))
	s.Tanggal_lahir = html.EscapeString(strings.TrimSpace(s.Tanggal_lahir))
	s.Nama_wali = html.EscapeString(strings.TrimSpace(s.Nama_wali))
	s.Alamat = html.EscapeString(strings.TrimSpace(s.Alamat))
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
}

// SaveSiswa ...
func (s *Siswa) SaveSiswa(db *gorm.DB) (*Siswa, error) {
	var err error
	err = db.Debug().Create(&s).Error
	if err != nil {
		return &Siswa{}, err
	}
	return s, nil
}

// FindAllSiswa ...
func (s *Siswa) FindAllSiswa(db *gorm.DB) (*[]Siswa, error) {
	var err error
	datasiswa := []Siswa{}
	err = db.Debug().Model(&Siswa{}).Where("status=1").Limit(100).Find(&datasiswa).Error
	if err != nil {
		return &[]Siswa{}, err
	}
	return &datasiswa, err
}

// FindSiswaByID ...
func (s *Siswa) FindSiswaByID(db *gorm.DB, uid uint32) (*Siswa, error) {
	var err error
	err = db.Debug().Model(Siswa{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Siswa{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &Siswa{}, errors.New("User Not Found")
	}
	return s, err
}

// GetSiswaByNames ...
func (s *Siswa) GetSiswaByNames(db *gorm.DB, nama string) (*[]Siswa, error) {
	var err error
	datasiswa := []Siswa{}
	err = db.Debug().Model(&Siswa{}).Where("nama LIKE ? OR no_induk LIKE ? AND STATUS = ?", `%`+nama+`%`, `%`+nama+`%`, "1").Take(&datasiswa).Error
	if err != nil {
		return &[]Siswa{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]Siswa{}, errors.New("Siswa Not Found")
	}
	return &datasiswa, err
}

// GetSiswaNoInduk ...
func (s *Siswa) GetSiswaNoInduk(db *gorm.DB, no_induk string) (*[]Siswa, error) {
	var err error
	datasiswa := []Siswa{}
	err = db.Debug().Model(&Siswa{}).Where("no_induk = ? ", no_induk).Take(&datasiswa).Error
	if err != nil {
		return &[]Siswa{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &[]Siswa{}, errors.New("Siswa Not Found")
	}
	return &datasiswa, err
}

// UpdateASiswa ...
func (s *Siswa) UpdateASiswa(db *gorm.DB, uid uint32) (*Siswa, error) {

	db = db.Debug().Model(&Siswa{}).Where("id = ?", uid).Take(&Siswa{}).UpdateColumn(
		map[string]interface{}{
			"nama":          s.Nama,
			"tempat_lahir":  s.Tempat_lahir,
			"Jenis_kelamin": s.Jenis_kelamin,
			"tanggal_lahir": s.Tanggal_lahir,
			"nama_wali":     s.Nama_wali,
			"alamat":        s.Alamat,
			"agama":         s.Agama,
			"kelas":         s.Kelas,
			"updated_at":    time.Now(),
		},
	)
	if db.Error != nil {
		return &Siswa{}, db.Error
	}
	// This is the display the updated user
	err := db.Debug().Model(&Siswa{}).Where("id = ?", uid).Take(&s).Error
	if err != nil {
		return &Siswa{}, err
	}
	return s, nil
}

// DeleteASiswa ...
func (s *Siswa) DeleteASiswa(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&Siswa{}).Where("id = ?", uid).Take(&Siswa{}).Delete(&Siswa{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
