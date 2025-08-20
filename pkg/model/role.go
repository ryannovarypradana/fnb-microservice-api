package model

// CORRECTED: Added the type definition for Role.
type Role string

const (
	RoleSuperAdmin Role = "SUPER_ADMIN"
	RoleCompanyRep Role = "COMPANY_REP"
	RoleStoreAdmin Role = "STORE_ADMIN"
	RoleCashier    Role = "CASHIER"
	RoleCustomer   Role = "CUSTOMER"
	RoleAdmin      Role = "ADMIN"
)
