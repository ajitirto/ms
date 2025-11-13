package usecase

import (
	"context"
	"log"

	"github.com/ajitirto/ms/product-service/internal/domain"
	"github.com/ajitirto/ms/product-service/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
	productCache repository.ProductCache
}

const allProductsCacheKey = "all_products"

func NewProductService(repo repository.ProductRepository, cache repository.ProductCache) ProductUsecase {
	return &productService{productRepo: repo, productCache: cache}
}


func (s *productService) GetProduct(ctx context.Context, id int) (*domain.Product, error) {
    product, err := s.productCache.Get(ctx, id)
    
    if err == nil && product != nil {
        log.Printf("[CACHE HIT] Product ID %d diambil dari Redis.", id)
        return product, nil 
    }
    
    if err != nil && err.Error() != "redis: nil" {
         log.Printf("[CACHE ERROR] Gagal akses Redis untuk ID %d: %v. Lanjut ke DB.", id, err)
    }

    log.Printf("[DB ACCESS] Product ID %d diambil dari PostgreSQL.", id)
    product, err = s.productRepo.GetByID(ctx, id)
    if err != nil {
        return nil, err // Gagal dari DB
    }

    if err := s.productCache.Set(ctx, product); err != nil {
        log.Printf("[CACHE FAIL] Gagal menyimpan Product ID %d ke Redis: %v", id, err)
    } else {
        log.Printf("[CACHE SET] Product ID %d berhasil disimpan ke Redis.", id)
    }

    return product, nil
}


// func (s *productService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
// 	return s.productRepo.GetAll(ctx)
// }

func (s *productService) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	// 1. COBA AMBIL DARI CACHE (REDIS)
	// Kita perlu asumsi bahwa interface ProductCache memiliki metode GetAll dan SetAll
	products, err := s.productCache.GetAll(ctx, allProductsCacheKey)

	if err == nil && products != nil {
		log.Printf("[CACHE HIT] Semua produk diambil dari Redis.")
		return products, nil // Success: Data dari Cache
	}

	// Log error jika Redis bermasalah (tapi tetap lanjut ke DB)
	if err != nil && err.Error() != "redis: nil" {
		log.Printf("[CACHE ERROR] Gagal akses Redis untuk semua produk: %v. Lanjut ke DB.", err)
	}

	// 2. CACHE MISS: AMBIL DARI REPOSITORY (POSTGRES)
	log.Printf("[DB ACCESS] Semua produk diambil dari PostgreSQL.")
	products, err = s.productRepo.GetAll(ctx)
	if err != nil {
		return nil, err // Gagal dari DB
	}

	// 3. SET KE CACHE (REDIS) UNTUK REQUEST BERIKUTNYA
	if err := s.productCache.SetAll(ctx, allProductsCacheKey, products); err != nil {
		log.Printf("[CACHE FAIL] Gagal menyimpan semua produk ke Redis: %v", err)
	} else {
		log.Printf("[CACHE SET] Semua produk berhasil disimpan ke Redis.")
	}

	return products, nil
}