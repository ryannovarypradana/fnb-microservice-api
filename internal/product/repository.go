// internal/product/repository.go
package product

import (
	"context"
	"encoding/json"
	"fnb-system/pkg/model"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

const (
	// Definisikan kunci cache sebagai konstanta agar tidak ada typo.
	cacheKeyAllMenus = "menus:all"
)

// ProductRepository mendefinisikan kontrak untuk operasi data produk.
type ProductRepository interface {
	// Category methods
	CreateCategory(category *model.Category) (*model.Category, error)
	FindAllCategories() (*[]model.Category, error)

	// Menu methods
	CreateMenu(menu *model.Menu) (*model.Menu, error)
	FindAllMenus() (*[]model.Menu, error)
	FindMenuByID(id uint) (*model.Menu, error)
}

// productRepository adalah implementasi konkret dengan koneksi DB dan Redis.
type productRepository struct {
	db  *gorm.DB
	rdb *redis.Client
	ctx context.Context
}

// NewProductRepository membuat instance baru dari productRepository.
func NewProductRepository(db *gorm.DB, rdb *redis.Client) ProductRepository {
	return &productRepository{
		db:  db,
		rdb: rdb,
		ctx: context.Background(),
	}
}

// --- Category Implementations ---

func (r *productRepository) CreateCategory(category *model.Category) (*model.Category, error) {
	if err := r.db.Create(category).Error; err != nil {
		return nil, err
	}
	return category, nil
}

func (r *productRepository) FindAllCategories() (*[]model.Category, error) {
	var categories []model.Category
	if err := r.db.Find(&categories).Error; err != nil {
		return nil, err
	}
	return &categories, nil
}

// --- Menu Implementations ---

// CreateMenu menyimpan menu baru dan menghapus cache yang sudah ada.
func (r *productRepository) CreateMenu(menu *model.Menu) (*model.Menu, error) {
	if err := r.db.Create(menu).Error; err != nil {
		return nil, err
	}
	// Invalidate cache karena data telah berubah.
	r.rdb.Del(r.ctx, cacheKeyAllMenus)
	return menu, nil
}

// FindAllMenus menerapkan pola cache-aside.
func (r *productRepository) FindAllMenus() (*[]model.Menu, error) {
	// 1. Coba ambil dari Redis terlebih dahulu.
	val, err := r.rdb.Get(r.ctx, cacheKeyAllMenus).Result()
	if err == nil {
		// Cache hit! Deserialisasi JSON dari Redis ke struct.
		var menus []model.Menu
		if json.Unmarshal([]byte(val), &menus) == nil {
			return &menus, nil
		}
	}

	// 2. Cache miss. Ambil dari database.
	var menus []model.Menu
	if err := r.db.Preload("Category").Find(&menus).Error; err != nil {
		return nil, err
	}

	// 3. Simpan hasil ke Redis untuk permintaan selanjutnya.
	jsonData, err := json.Marshal(menus)
	if err == nil {
		// Set cache dengan waktu kadaluarsa 10 menit.
		r.rdb.Set(r.ctx, cacheKeyAllMenus, jsonData, 10*time.Minute)
	}

	return &menus, nil
}

// FindMenuByID mencari satu menu. Caching bisa ditambahkan di sini juga jika diperlukan.
func (r *productRepository) FindMenuByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	if err := r.db.Preload("Category").First(&menu, id).Error; err != nil {
		return nil, err
	}
	return &menu, nil
}
