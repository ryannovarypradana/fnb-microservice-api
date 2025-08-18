// pkg/dto/user.go
package dto

// UserUpdateRequest adalah struct untuk request body saat memperbarui data user.
// Kita hanya mengizinkan nama dan email yang diubah melalui endpoint ini.
// Password memiliki alur perubahannya sendiri (misal: "lupa password").
type UserUpdateRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// UserResponse adalah struct untuk data user yang aman untuk dikirim kembali ke klien.
// Ini mencegah kita mengirim data sensitif seperti `PasswordHash`.
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
