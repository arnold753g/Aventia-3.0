package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/internal/services"
	"andaria-backend/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	validate *validator.Validate
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		validate: validator.New(),
	}
}

// Register maneja el registro de nuevos usuarios
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	var profilePhotoPath string

	contentType := r.Header.Get("Content-Type")

	// Manejar multipart/form-data (con foto) o JSON (sin foto)
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.ErrorResponse(w, "INVALID_FORM", "No se pudo procesar el formulario", nil, http.StatusBadRequest)
			return
		}

		req = models.RegisterRequest{
			Nombre:          r.FormValue("nombre"),
			ApellidoPaterno: r.FormValue("apellido_paterno"),
			ApellidoMaterno: r.FormValue("apellido_materno"),
			Email:           r.FormValue("email"),
			Password:        r.FormValue("password"),
			CI:              r.FormValue("ci"),
			Expedido:        r.FormValue("expedido"),
			Phone:           r.FormValue("phone"),
			FechaNacimiento: r.FormValue("fecha_nacimiento"),
			Ciudad:          r.FormValue("ciudad"),
			Rol:             r.FormValue("rol"),
		}

		// Manejar archivo de foto si existe
		if file, header, err := r.FormFile("profile_photo"); err == nil {
			defer file.Close()
			savedPath, err := saveUserPhoto(file, header)
			if err != nil {
				utils.ErrorResponse(w, "INVALID_FILE", "No se pudo guardar la fotografía", err.Error(), http.StatusBadRequest)
				return
			}
			profilePhotoPath = savedPath
		}
	} else {
		// Decodificar JSON (sin foto)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
			return
		}
	}

	// Validar datos
	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	normalizedPhone, ok := validatePhone(req.Phone)
	if !ok {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El telefono es invalido", nil, http.StatusBadRequest)
		return
	}
	req.Phone = normalizedPhone

	// Verificar si el email ya existe
	var existingUser models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(w, "EMAIL_EXISTS", "El email ya está registrado", nil, http.StatusConflict)
		return
	}

	// Verificar si el CI ya existe
	if err := database.GetDB().Where("ci = ?", strings.ToUpper(req.CI)).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(w, "CI_EXISTS", "El CI ya está registrado", nil, http.StatusConflict)
		return
	}

	// Hashear contraseña
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(w, "SERVER_ERROR", "Error al procesar la contraseña", nil, http.StatusInternalServerError)
		return
	}

	// Parsear fecha de nacimiento
	fechaNac, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_DATE", "Formato de fecha inválido. Use YYYY-MM-DD", nil, http.StatusBadRequest)
		return
	}

	// SEGURIDAD: Forzar rol 'turista' para registros públicos
	// No se permite especificar el rol desde el frontend para prevenir
	// que usuarios maliciosos se registren como admin o encargado_agencia
	rol := "turista"

	// Generar OTP de verificacion
	code, err := utils.GenerateOTP6()
	if err != nil {
		utils.ErrorResponse(w, "OTP_ERROR", "Error al generar el codigo de verificacion", nil, http.StatusInternalServerError)
		return
	}

	codeHash := utils.HashOTP("email_verify", strings.ToLower(req.Email), code)
	expiresAt := time.Now().Add(otpExpiryDuration())

	// Crear usuario
	usuario := models.Usuario{
		Nombre:          req.Nombre,
		ApellidoPaterno: req.ApellidoPaterno,
		ApellidoMaterno: req.ApellidoMaterno,
		Email:           strings.ToLower(req.Email),
		PasswordHash:    hashedPassword,
		CI:              strings.ToUpper(req.CI),
		Expedido:        req.Expedido,
		Phone:           req.Phone,
		FechaNacimiento: fechaNac,
		Ciudad:          req.Ciudad,
		Rol:             rol,
		Status:          "active",
		Nationality:     "Bolivia",
		TermsAccepted:   true,
		ProfilePhoto:    profilePhotoPath,
		EmailVerified:   false,
		EmailVerifyCodeHash:   &codeHash,
		EmailVerifyExpiresAt:  &expiresAt,
		EmailVerifyAttempts:   0,
		EmailVerifyLastSentAt: timePtr(time.Now()),
		EmailVerifySentCount:  1,
	}

	now := time.Now()
	usuario.TermsAcceptedAt = &now

	// Guardar en base de datos
	if err := database.GetDB().Create(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear usuario", err.Error(), http.StatusInternalServerError)
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendVerificationCode(usuario.Email, code); err != nil {
		log.Printf("Error enviando email de verificacion: %v", err)
	}

	utils.SuccessResponse(w, nil, "Cuenta creada exitosamente. Te enviamos un codigo de verificacion a tu correo.", http.StatusCreated)
}

// Login maneja el inicio de sesión
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest

	// Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	// Validar datos
	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	// Buscar usuario por email
	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "INVALID_CREDENTIALS", "Credenciales inválidas", nil, http.StatusUnauthorized)
		return
	}

	// Verificar si la cuenta está bloqueada
	if usuario.LockedUntil != nil && usuario.LockedUntil.After(time.Now()) {
		utils.ErrorResponse(w, "ACCOUNT_LOCKED", "Cuenta bloqueada. Intente más tarde", nil, http.StatusForbidden)
		return
	}

	// Verificar si la cuenta está suspendida
	if usuario.Status == "suspended" || usuario.Status == "inactive" {
		utils.ErrorResponse(w, "ACCOUNT_SUSPENDED", "Cuenta suspendida o inactiva", nil, http.StatusForbidden)
		return
	}

	// Verificar email (admins no requieren verificacion)
	if usuario.Rol != "admin" && !usuario.EmailVerified {
		utils.ErrorResponse(w, "EMAIL_NOT_VERIFIED", "Debes verificar tu correo antes de iniciar sesion. Revisa tu bandeja de entrada.", nil, http.StatusForbidden)
		return
	}

	// Encargados deben establecer contrasena inicial
	if usuario.Rol == "encargado_agencia" && usuario.PasswordHash == "" {
		utils.ErrorResponse(w, "PASSWORD_NOT_SET", "Debes establecer tu contrasena inicial primero", nil, http.StatusForbidden)
		return
	}

	// Verificar contraseña
	if !utils.CheckPassword(req.Password, usuario.PasswordHash) {
		// Incrementar intentos fallidos
		usuario.LoginAttempts++

		if usuario.LoginAttempts >= 5 {
			lockedUntil := time.Now().Add(30 * time.Minute)
			usuario.LockedUntil = &lockedUntil
			usuario.Status = "suspended"
		}

		database.GetDB().Save(&usuario)

		utils.ErrorResponse(w, "INVALID_CREDENTIALS", "Credenciales inválidas", nil, http.StatusUnauthorized)
		return
	}

	// Login exitoso - resetear intentos fallidos
	usuario.LoginAttempts = 0
	usuario.LockedUntil = nil
	if usuario.Status == "suspended" {
		usuario.Status = "active"
	}
	now := time.Now()
	usuario.LastLogin = &now
	database.GetDB().Save(&usuario)

	// Generar tokens
	token, err := utils.GenerateToken(usuario.ID, usuario.Email, usuario.Rol, 24*time.Hour)
	if err != nil {
		utils.ErrorResponse(w, "TOKEN_ERROR", "Error al generar token", nil, http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateToken(usuario.ID, usuario.Email, usuario.Rol, 168*time.Hour)
	if err != nil {
		utils.ErrorResponse(w, "TOKEN_ERROR", "Error al generar refresh token", nil, http.StatusInternalServerError)
		return
	}

	// Respuesta
	response := models.LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		User: models.UsuarioPublic{
			ID:              usuario.ID,
			Nombre:          usuario.Nombre,
			ApellidoPaterno: usuario.ApellidoPaterno,
			ApellidoMaterno: usuario.ApellidoMaterno,
			Email:           usuario.Email,
			Rol:             usuario.Rol,
			Status:          usuario.Status,
			ProfilePhoto:    usuario.ProfilePhoto,
			Ciudad:          usuario.Ciudad,
		},
	}

	utils.SuccessResponse(w, response, "Login exitoso", http.StatusOK)
}

// RefreshToken maneja la renovación de tokens
func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshTokenRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	// Validar el refresh token
	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_TOKEN", "Token inválido o expirado", nil, http.StatusUnauthorized)
		return
	}

	// Generar nuevos tokens
	newToken, err := utils.GenerateToken(claims.UserID, claims.Email, claims.Rol, 24*time.Hour)
	if err != nil {
		utils.ErrorResponse(w, "TOKEN_ERROR", "Error al generar token", nil, http.StatusInternalServerError)
		return
	}

	newRefreshToken, err := utils.GenerateToken(claims.UserID, claims.Email, claims.Rol, 168*time.Hour)
	if err != nil {
		utils.ErrorResponse(w, "TOKEN_ERROR", "Error al generar refresh token", nil, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"token":         newToken,
		"refresh_token": newRefreshToken,
	}

	utils.SuccessResponse(w, response, "Token renovado exitosamente", http.StatusOK)
}

// GetProfile obtiene el perfil del usuario autenticado
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Obtener claims del contexto (puestos por el middleware)
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	// Buscar usuario
	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, claims.UserID).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	response := models.UsuarioPublic{
		ID:              usuario.ID,
		Nombre:          usuario.Nombre,
		ApellidoPaterno: usuario.ApellidoPaterno,
		ApellidoMaterno: usuario.ApellidoMaterno,
		Email:           usuario.Email,
		Rol:             usuario.Rol,
		Status:          usuario.Status,
		ProfilePhoto:    usuario.ProfilePhoto,
		Ciudad:          usuario.Ciudad,
	}

	utils.SuccessResponse(w, response, "Perfil obtenido exitosamente", http.StatusOK)
}

// Logout maneja el cierre de sesión
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Por ahora solo respondemos OK
	// En producción, aquí podrías invalidar el token en una lista negra
	utils.SuccessResponse(w, nil, "Logout exitoso", http.StatusOK)
}

// VerifyEmail valida el codigo OTP para verificar el correo.
func (h *AuthHandler) VerifyEmail(w http.ResponseWriter, r *http.Request) {
	var req models.VerifyEmailRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "INVALID_CODE", "Codigo invalido o expirado", nil, http.StatusUnauthorized)
		return
	}

	if usuario.EmailVerified {
		utils.SuccessResponse(w, nil, "El correo ya esta verificado", http.StatusOK)
		return
	}

	if usuario.EmailVerifyCodeHash == nil || usuario.EmailVerifyExpiresAt == nil {
		utils.ErrorResponse(w, "NO_PENDING_CODE", "No hay codigo de verificacion pendiente", nil, http.StatusBadRequest)
		return
	}

	if time.Now().After(*usuario.EmailVerifyExpiresAt) {
		utils.ErrorResponse(w, "CODE_EXPIRED", "El codigo ha expirado. Solicita uno nuevo.", nil, http.StatusUnauthorized)
		return
	}

	maxAttempts := otpMaxAttempts()
	if usuario.EmailVerifyAttempts >= maxAttempts {
		utils.ErrorResponse(w, "TOO_MANY_ATTEMPTS", "Demasiados intentos fallidos. Solicita un nuevo codigo.", nil, http.StatusTooManyRequests)
		return
	}

	if !utils.VerifyOTP("email_verify", usuario.Email, req.Code, *usuario.EmailVerifyCodeHash) {
		usuario.EmailVerifyAttempts++
		_ = database.GetDB().Save(&usuario)

		remaining := maxAttempts - usuario.EmailVerifyAttempts
		utils.ErrorResponse(w, "INVALID_CODE", fmt.Sprintf("Codigo incorrecto. Te quedan %d intentos.", remaining), nil, http.StatusUnauthorized)
		return
	}

	usuario.EmailVerified = true
	usuario.EmailVerifyCodeHash = nil
	usuario.EmailVerifyExpiresAt = nil
	usuario.EmailVerifyAttempts = 0
	usuario.EmailVerifyLastSentAt = nil
	usuario.EmailVerifySentCount = 0

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al verificar email", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "¡Correo verificado exitosamente! Ya puedes iniciar sesion.", http.StatusOK)
}

// ResendEmailCode reenvia el codigo de verificacion con cooldown y limites.
func (h *AuthHandler) ResendEmailCode(w http.ResponseWriter, r *http.Request) {
	var req models.ResendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	genericMessage := "Si el correo existe y no esta verificado, recibiras un nuevo codigo."

	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
		return
	}

	if usuario.EmailVerified {
		utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
		return
	}

	cooldown := otpResendCooldown()
	if usuario.EmailVerifyLastSentAt != nil {
		elapsed := time.Since(*usuario.EmailVerifyLastSentAt)
		if elapsed < cooldown {
			utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
			return
		}
	}

	resetDailyCount(usuario.EmailVerifyLastSentAt, &usuario.EmailVerifySentCount)
	if usuario.EmailVerifySentCount >= otpMaxResends() {
		utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
		return
	}

	code, err := utils.GenerateOTP6()
	if err != nil {
		utils.ErrorResponse(w, "OTP_ERROR", "Error al generar el codigo", nil, http.StatusInternalServerError)
		return
	}

	codeHash := utils.HashOTP("email_verify", usuario.Email, code)
	expiresAt := time.Now().Add(otpExpiryDuration())

	usuario.EmailVerifyCodeHash = &codeHash
	usuario.EmailVerifyExpiresAt = &expiresAt
	usuario.EmailVerifyAttempts = 0
	usuario.EmailVerifyLastSentAt = timePtr(time.Now())
	usuario.EmailVerifySentCount++

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al procesar solicitud", err.Error(), http.StatusInternalServerError)
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendVerificationCode(usuario.Email, code); err != nil {
		log.Printf("Error enviando email de verificacion: %v", err)
	}

	utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
}

// ForgotPassword envia el codigo de recuperacion con respuesta generica.
func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	genericMessage := "Si el correo existe, recibiras un codigo para restablecer tu contrasena."

	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
		return
	}

	cooldown := otpResendCooldown()
	if usuario.PasswordResetLastSentAt != nil {
		elapsed := time.Since(*usuario.PasswordResetLastSentAt)
		if elapsed < cooldown {
			utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
			return
		}
	}

	resetDailyCount(usuario.PasswordResetLastSentAt, &usuario.PasswordResetSentCount)
	if usuario.PasswordResetSentCount >= otpMaxResends() {
		utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
		return
	}

	code, err := utils.GenerateOTP6()
	if err != nil {
		utils.ErrorResponse(w, "OTP_ERROR", "Error al generar el codigo", nil, http.StatusInternalServerError)
		return
	}

	tokenHash := utils.HashOTP("password_reset", usuario.Email, code)
	expiresAt := time.Now().Add(passwordResetExpiryDuration())

	usuario.PasswordResetTokenHash = &tokenHash
	usuario.PasswordResetExpiresAt = &expiresAt
	usuario.PasswordResetAttempts = 0
	usuario.PasswordResetLastSentAt = timePtr(time.Now())
	usuario.PasswordResetSentCount++

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al procesar solicitud", err.Error(), http.StatusInternalServerError)
		return
	}

	emailService := services.NewEmailService()
	if err := emailService.SendPasswordResetCode(usuario.Email, code); err != nil {
		log.Printf("Error enviando email de recuperacion: %v", err)
	}

	utils.SuccessResponse(w, nil, genericMessage, http.StatusOK)
}

// ResetPassword valida el codigo y actualiza la contrasena.
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req models.ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "INVALID_CODE", "Codigo invalido o expirado", nil, http.StatusUnauthorized)
		return
	}

	if usuario.PasswordResetTokenHash == nil || usuario.PasswordResetExpiresAt == nil {
		utils.ErrorResponse(w, "NO_PENDING_RESET", "No hay solicitud de restablecimiento pendiente", nil, http.StatusBadRequest)
		return
	}

	if time.Now().After(*usuario.PasswordResetExpiresAt) {
		utils.ErrorResponse(w, "CODE_EXPIRED", "El codigo ha expirado. Solicita uno nuevo.", nil, http.StatusUnauthorized)
		return
	}

	maxAttempts := passwordResetMaxAttempts()
	if usuario.PasswordResetAttempts >= maxAttempts {
		utils.ErrorResponse(w, "TOO_MANY_ATTEMPTS", "Demasiados intentos fallidos. Solicita un nuevo codigo.", nil, http.StatusTooManyRequests)
		return
	}

	if !utils.VerifyOTP("password_reset", usuario.Email, req.Code, *usuario.PasswordResetTokenHash) {
		usuario.PasswordResetAttempts++
		_ = database.GetDB().Save(&usuario)

		remaining := maxAttempts - usuario.PasswordResetAttempts
		utils.ErrorResponse(w, "INVALID_CODE", fmt.Sprintf("Codigo incorrecto. Te quedan %d intentos.", remaining), nil, http.StatusUnauthorized)
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(w, "SERVER_ERROR", "Error al procesar la contrasena", nil, http.StatusInternalServerError)
		return
	}

	usuario.PasswordHash = hashedPassword
	usuario.PasswordResetTokenHash = nil
	usuario.PasswordResetExpiresAt = nil
	usuario.PasswordResetAttempts = 0
	usuario.PasswordResetLastSentAt = nil
	usuario.PasswordResetSentCount = 0

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar contrasena", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Contrasena restablecida exitosamente. Ya puedes iniciar sesion.", http.StatusOK)
}

// ChangePassword permite cambiar la contrasena para usuarios autenticados.
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	var req models.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, claims.UserID).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	if !utils.CheckPassword(req.CurrentPassword, usuario.PasswordHash) {
		utils.ErrorResponse(w, "INVALID_CREDENTIALS", "La contrasena actual es incorrecta", nil, http.StatusUnauthorized)
		return
	}

	if req.CurrentPassword == req.NewPassword {
		utils.ErrorResponse(w, "INVALID_PASSWORD", "La nueva contrasena debe ser diferente a la actual", nil, http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(w, "SERVER_ERROR", "Error al procesar la contrasena", nil, http.StatusInternalServerError)
		return
	}

	usuario.PasswordHash = hashedPassword

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar contrasena", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Contrasena actualizada exitosamente. Las sesiones activas en otros dispositivos seguiran funcionando.", http.StatusOK)
}

// SetInitialPassword establece la primera contrasena para encargados verificados.
func (h *AuthHandler) SetInitialPassword(w http.ResponseWriter, r *http.Request) {
	var req models.SetInitialPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	if !usuario.EmailVerified {
		utils.ErrorResponse(w, "EMAIL_NOT_VERIFIED", "Debes verificar tu correo primero", nil, http.StatusForbidden)
		return
	}

	if usuario.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "FORBIDDEN", "Esta operacion solo esta disponible para encargados de agencia", nil, http.StatusForbidden)
		return
	}

	if usuario.PasswordHash != "" {
		utils.ErrorResponse(w, "PASSWORD_ALREADY_SET", "Ya tienes una contrasena establecida. Usa cambiar contrasena u olvide mi contrasena.", nil, http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ErrorResponse(w, "SERVER_ERROR", "Error al procesar la contrasena", nil, http.StatusInternalServerError)
		return
	}

	usuario.PasswordHash = hashedPassword

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al establecer contrasena", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Contrasena establecida exitosamente. Ya puedes iniciar sesion.", http.StatusOK)
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func otpExpiryDuration() time.Duration {
	return time.Duration(utils.GetEnvInt("OTP_EXPIRY_MINUTES", 10)) * time.Minute
}

func otpMaxAttempts() int {
	return utils.GetEnvInt("OTP_MAX_ATTEMPTS", 5)
}

func otpResendCooldown() time.Duration {
	return time.Duration(utils.GetEnvInt("OTP_RESEND_COOLDOWN_SECONDS", 60)) * time.Second
}

func otpMaxResends() int {
	return utils.GetEnvInt("OTP_MAX_RESENDS_PER_DAY", 3)
}

func passwordResetExpiryDuration() time.Duration {
	return time.Duration(utils.GetEnvInt("PASSWORD_RESET_EXPIRY_MINUTES", 20)) * time.Minute
}

func passwordResetMaxAttempts() int {
	return utils.GetEnvInt("PASSWORD_RESET_MAX_ATTEMPTS", 5)
}

func resetDailyCount(lastSent *time.Time, count *int) {
	if lastSent == nil {
		return
	}
	if time.Since(*lastSent) > 24*time.Hour {
		*count = 0
	}
}
