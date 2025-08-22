// pkg/grpc/handler/handler.go

package handler

// Handlers menampung semua handler aplikasi.
type Handlers struct {
	Auth      *AuthHandler   // AuthHandler adalah struct
	User      *UserHandler   // UserHandler adalah struct
	Store     *StoreHandler  // StoreHandler adalah struct
	Product   ProductHandler // ProductHandler adalah interface
	Order     OrderHandler   // OrderHandler adalah interface
	Company   CompanyHandler // CompanyHandler adalah interface
	Promotion PromotionHandler
}
