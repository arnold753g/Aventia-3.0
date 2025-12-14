package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
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
	}

	now := time.Now()
	usuario.TermsAcceptedAt = &now

	// Guardar en base de datos
	if err := database.GetDB().Create(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear usuario", err.Error(), http.StatusInternalServerError)
		return
	}

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

	utils.SuccessResponse(w, response, "Usuario registrado exitosamente", http.StatusCreated)
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
