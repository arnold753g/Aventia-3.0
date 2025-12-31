package models

type RegisterRequest struct {
	Nombre          string `json:"nombre" validate:"required,min=2"`
	ApellidoPaterno string `json:"apellido_paterno" validate:"required,min=2"`
	ApellidoMaterno string `json:"apellido_materno" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	CI              string `json:"ci" validate:"required,min=5"`
	Expedido        string `json:"expedido" validate:"required"`
	Phone           string `json:"phone" validate:"required"`
	FechaNacimiento string `json:"fecha_nacimiento" validate:"required"`
	Ciudad          string `json:"ciudad"`
	Rol             string `json:"rol" validate:"required,oneof=turista encargado_agencia"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	User         UsuarioPublic `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type VerifyEmailRequest struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6,numeric"`
}

type ResendCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6,numeric"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
}

type SetInitialPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}
