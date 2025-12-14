package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type UsuarioHandler struct {
	validate *validator.Validate
}

func NewUsuarioHandler() *UsuarioHandler {
	return &UsuarioHandler{
		validate: validator.New(),
	}
}

func isOnlyDigits(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}

// CheckUsuarioExiste verifica duplicados de email o CI
func (h *UsuarioHandler) CheckUsuarioExiste(w http.ResponseWriter, r *http.Request) {
	email := strings.ToLower(r.URL.Query().Get("email"))
	ci := strings.ToUpper(r.URL.Query().Get("ci"))

	var emailExists, ciExists bool
	db := database.GetDB()

	if email != "" {
		var count int64
		if err := db.Model(&models.Usuario{}).Where("email = ?", email).Count(&count).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, "DB_ERROR", "Error verificando email", err.Error(), http.StatusInternalServerError)
			return
		}
		emailExists = count > 0
	}

	if ci != "" {
		var count int64
		if err := db.Model(&models.Usuario{}).Where("ci = ?", ci).Count(&count).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			utils.ErrorResponse(w, "DB_ERROR", "Error verificando CI", err.Error(), http.StatusInternalServerError)
			return
		}
		ciExists = count > 0
	}

	utils.SuccessResponse(w, map[string]bool{
		"emailExists": emailExists,
		"ciExists":    ciExists,
	}, "Check duplicados", http.StatusOK)
}

// saveUserPhoto guarda la foto en disco y retorna la ruta relativa
func saveUserPhoto(file multipart.File, header *multipart.FileHeader) (string, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return "", fmt.Errorf("formato de imagen no permitido")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		switch contentType {
		case "image/jpeg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		}
	}

	destDir := filepath.Join("uploads", "fotografias", "usuarios")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("usuario_%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return filepath.ToSlash(destPath), nil
}

func deleteUserPhoto(storedPath string) error {
	if storedPath == "" {
		return nil
	}

	baseDir := filepath.Clean(filepath.Join(".", "uploads", "fotografias", "usuarios"))
	clean := filepath.Clean(filepath.Join(".", storedPath))

	// Evitar borrar archivos fuera de la carpeta esperada
	if !strings.HasPrefix(clean, baseDir) {
		return nil
	}

	if _, err := os.Stat(clean); err == nil {
		return os.Remove(clean)
	}

	return nil
}

// GetUsuarios obtiene lista de usuarios con filtros y paginacion
func (h *UsuarioHandler) GetUsuarios(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	search := r.URL.Query().Get("search")
	rol := r.URL.Query().Get("rol")
	status := r.URL.Query().Get("status")
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	db := database.GetDB()
	query := db.Model(&models.Usuario{})

	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where(
			"LOWER(nombre) LIKE ? OR LOWER(apellido_paterno) LIKE ? OR LOWER(apellido_materno) LIKE ? OR LOWER(email) LIKE ? OR LOWER(ci) LIKE ?",
			searchPattern, searchPattern, searchPattern, searchPattern, searchPattern,
		)
	}

	if rol != "" {
		query = query.Where("rol = ?", rol)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit

	var usuarios []models.Usuario
	orderClause := sortBy + " " + sortOrder
	if err := query.Order(orderClause).Limit(limit).Offset(offset).Find(&usuarios).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener usuarios", err.Error(), http.StatusInternalServerError)
		return
	}

	usuariosPublic := make([]models.UsuarioPublic, len(usuarios))
	for i, u := range usuarios {
		usuariosPublic[i] = models.UsuarioPublic{
			ID:              u.ID,
			Nombre:          u.Nombre,
			ApellidoPaterno: u.ApellidoPaterno,
			ApellidoMaterno: u.ApellidoMaterno,
			Email:           u.Email,
			CI:              u.CI,
			Expedido:        u.Expedido,
			Phone:           u.Phone,
			Rol:             u.Rol,
			Status:          u.Status,
			ProfilePhoto:    u.ProfilePhoto,
			Ciudad:          u.Ciudad,
			Nationality:     u.Nationality,
			EmailVerified:   u.EmailVerified,
			LastLogin:       u.LastLogin,
			CreatedAt:       u.CreatedAt,
		}
	}

	response := map[string]interface{}{
		"usuarios": usuariosPublic,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.SuccessResponse(w, response, "Usuarios obtenidos exitosamente", http.StatusOK)
}

// GetUsuario obtiene un usuario por ID
func (h *UsuarioHandler) GetUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, id).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	usuarioDetalle := models.UsuarioDetalle{
		ID:              usuario.ID,
		Nombre:          usuario.Nombre,
		ApellidoPaterno: usuario.ApellidoPaterno,
		ApellidoMaterno: usuario.ApellidoMaterno,
		Email:           usuario.Email,
		CI:              usuario.CI,
		Expedido:        usuario.Expedido,
		Phone:           usuario.Phone,
		FechaNacimiento: usuario.FechaNacimiento,
		Rol:             usuario.Rol,
		Status:          usuario.Status,
		ProfilePhoto:    usuario.ProfilePhoto,
		Ciudad:          usuario.Ciudad,
		Nationality:     usuario.Nationality,
		EmailVerified:   usuario.EmailVerified,
		TermsAccepted:   usuario.TermsAccepted,
		TermsAcceptedAt: usuario.TermsAcceptedAt,
		LastLogin:       usuario.LastLogin,
		LoginAttempts:   usuario.LoginAttempts,
		LockedUntil:     usuario.LockedUntil,
		CreatedAt:       usuario.CreatedAt,
		UpdatedAt:       usuario.UpdatedAt,
	}

	utils.SuccessResponse(w, usuarioDetalle, "Usuario obtenido exitosamente", http.StatusOK)
}

// CreateUsuario crea un nuevo usuario (solo admin)
func (h *UsuarioHandler) CreateUsuario(w http.ResponseWriter, r *http.Request) {
	var req models.CreateUsuarioRequest

	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.ErrorResponse(w, "INVALID_FORM", "No se pudo procesar el formulario", nil, http.StatusBadRequest)
			return
		}

		req = models.CreateUsuarioRequest{
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
			Nationality:     r.FormValue("nationality"),
			Rol:             r.FormValue("rol"),
		}

		if file, header, err := r.FormFile("profile_photo"); err == nil {
			defer file.Close()
			savedPath, err := saveUserPhoto(file, header)
			if err != nil {
				utils.ErrorResponse(w, "INVALID_FILE", "No se pudo guardar la fotografia", err.Error(), http.StatusBadRequest)
				return
			}
			req.ProfilePhoto = savedPath
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
			return
		}
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	if len(req.Phone) != 8 || !isOnlyDigits(req.Phone) {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "El telefono debe tener 8 digitos", nil, http.StatusBadRequest)
		return
	}

	var existingUser models.Usuario
	if err := database.GetDB().Where("email = ?", strings.ToLower(req.Email)).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(w, "EMAIL_EXISTS", "El email ya esta registrado", nil, http.StatusConflict)
		return
	}

	if err := database.GetDB().Where("ci = ?", strings.ToUpper(req.CI)).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(w, "CI_EXISTS", "El CI ya esta registrado", nil, http.StatusConflict)
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ErrorResponse(w, "SERVER_ERROR", "Error al procesar la contrasena", nil, http.StatusInternalServerError)
		return
	}

	fechaNac, err := time.Parse("2006-01-02", req.FechaNacimiento)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_DATE", "Formato de fecha invalido. Use YYYY-MM-DD", nil, http.StatusBadRequest)
		return
	}

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
		Rol:             req.Rol,
		Status:          "active",
		Nationality:     req.Nationality,
		TermsAccepted:   true,
		ProfilePhoto:    req.ProfilePhoto,
	}

	if req.Nationality == "" {
		usuario.Nationality = "Bolivia"
	}

	now := time.Now()
	usuario.TermsAcceptedAt = &now

	if err := database.GetDB().Create(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear usuario", err.Error(), http.StatusInternalServerError)
		return
	}

	usuarioPublic := models.UsuarioPublic{
		ID:              usuario.ID,
		Nombre:          usuario.Nombre,
		ApellidoPaterno: usuario.ApellidoPaterno,
		ApellidoMaterno: usuario.ApellidoMaterno,
		Email:           usuario.Email,
		CI:              usuario.CI,
		Expedido:        usuario.Expedido,
		Phone:           usuario.Phone,
		Rol:             usuario.Rol,
		Status:          usuario.Status,
		ProfilePhoto:    usuario.ProfilePhoto,
		Ciudad:          usuario.Ciudad,
	}

	utils.SuccessResponse(w, usuarioPublic, "Usuario creado exitosamente", http.StatusCreated)
}

// UpdateUsuario actualiza un usuario
func (h *UsuarioHandler) UpdateUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, id).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	oldPhoto := usuario.ProfilePhoto
	newPhotoPath := ""

	var req models.UpdateUsuarioRequest
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			utils.ErrorResponse(w, "INVALID_FORM", "No se pudo procesar el formulario", nil, http.StatusBadRequest)
			return
		}
		req = models.UpdateUsuarioRequest{
			Nombre:          r.FormValue("nombre"),
			ApellidoPaterno: r.FormValue("apellido_paterno"),
			ApellidoMaterno: r.FormValue("apellido_materno"),
			Phone:           r.FormValue("phone"),
			Ciudad:          r.FormValue("ciudad"),
			ProfilePhoto:    "",
		}
		if file, header, err := r.FormFile("profile_photo"); err == nil {
			defer file.Close()
			savedPath, err := saveUserPhoto(file, header)
			if err != nil {
				utils.ErrorResponse(w, "INVALID_FILE", "No se pudo guardar la fotografia", err.Error(), http.StatusBadRequest)
				return
			}
			req.ProfilePhoto = savedPath
			newPhotoPath = savedPath
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
			return
		}
	}

	if req.Nombre != "" {
		usuario.Nombre = req.Nombre
	}
	if req.ApellidoPaterno != "" {
		usuario.ApellidoPaterno = req.ApellidoPaterno
	}
	if req.ApellidoMaterno != "" {
		usuario.ApellidoMaterno = req.ApellidoMaterno
	}
	if req.Phone != "" {
		if len(req.Phone) != 8 || !isOnlyDigits(req.Phone) {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "El telefono debe tener 8 digitos", nil, http.StatusBadRequest)
			return
		}
		usuario.Phone = req.Phone
	}
	if req.Ciudad != "" {
		usuario.Ciudad = req.Ciudad
	}
	if req.ProfilePhoto != "" {
		usuario.ProfilePhoto = req.ProfilePhoto
	}

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		if newPhotoPath != "" && newPhotoPath != oldPhoto {
			_ = deleteUserPhoto(newPhotoPath)
		}
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar usuario", err.Error(), http.StatusInternalServerError)
		return
	}

	if newPhotoPath != "" && oldPhoto != "" && newPhotoPath != oldPhoto {
		_ = deleteUserPhoto(oldPhoto)
	}

	usuarioPublic := models.UsuarioPublic{
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

	utils.SuccessResponse(w, usuarioPublic, "Usuario actualizado exitosamente", http.StatusOK)
}

// UpdateUsuarioRol actualiza el rol de un usuario (solo admin)
func (h *UsuarioHandler) UpdateUsuarioRol(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, id).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	var req models.UpdateRolRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if req.Rol != "admin" && req.Rol != "turista" && req.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "INVALID_ROLE", "Rol invalido", nil, http.StatusBadRequest)
		return
	}

	usuario.Rol = req.Rol

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar rol", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"id":  usuario.ID,
		"rol": usuario.Rol,
	}, "Rol actualizado exitosamente", http.StatusOK)
}

// UpdateUsuarioStatus actualiza el estado de un usuario (solo admin)
func (h *UsuarioHandler) UpdateUsuarioStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, id).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	var req models.UpdateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if req.Status != "active" && req.Status != "inactive" && req.Status != "suspended" {
		utils.ErrorResponse(w, "INVALID_STATUS", "Status invalido", nil, http.StatusBadRequest)
		return
	}

	usuario.Status = req.Status

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar status", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"id":     usuario.ID,
		"status": usuario.Status,
	}, "Status actualizado exitosamente", http.StatusOK)
}

// DeactivateUsuario desactiva un usuario (soft delete)
func (h *UsuarioHandler) DeactivateUsuario(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	if claims.UserID == uint(id) {
		utils.ErrorResponse(w, "CANNOT_DEACTIVATE_SELF", "No puedes desactivar tu propia cuenta", nil, http.StatusBadRequest)
		return
	}

	var usuario models.Usuario
	if err := database.GetDB().First(&usuario, id).Error; err != nil {
		utils.ErrorResponse(w, "USER_NOT_FOUND", "Usuario no encontrado", nil, http.StatusNotFound)
		return
	}

	usuario.Status = "inactive"

	if err := database.GetDB().Save(&usuario).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al desactivar usuario", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Usuario desactivado exitosamente", http.StatusOK)
}

// GetUsuarioStats obtiene estadisticas de usuarios (solo admin)
func (h *UsuarioHandler) GetUsuarioStats(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		Total               int64 `json:"total"`
		Active              int64 `json:"active"`
		Inactive            int64 `json:"inactive"`
		Suspended           int64 `json:"suspended"`
		Admins              int64 `json:"admins"`
		Turistas            int64 `json:"turistas"`
		Encargados          int64 `json:"encargados"`
		Hoy                 int64 `json:"hoy"`
		Mes                 int64 `json:"mes"`
	}

	db := database.GetDB()

	db.Model(&models.Usuario{}).Count(&stats.Total)
	db.Model(&models.Usuario{}).Where("status = ?", "active").Count(&stats.Active)
	db.Model(&models.Usuario{}).Where("status = ?", "inactive").Count(&stats.Inactive)
	db.Model(&models.Usuario{}).Where("status = ?", "suspended").Count(&stats.Suspended)
	db.Model(&models.Usuario{}).Where("rol = ?", "admin").Count(&stats.Admins)
	db.Model(&models.Usuario{}).Where("rol = ?", "turista").Count(&stats.Turistas)
	db.Model(&models.Usuario{}).Where("rol = ?", "encargado_agencia").Count(&stats.Encargados)

	today := time.Now().Truncate(24 * time.Hour)
	db.Model(&models.Usuario{}).Where("created_at >= ?", today).Count(&stats.Hoy)

	firstDayOfMonth := time.Date(time.Now().Year(), time.Now().Month(), 1, 0, 0, 0, 0, time.Now().Location())
	db.Model(&models.Usuario{}).Where("created_at >= ?", firstDayOfMonth).Count(&stats.Mes)

	utils.SuccessResponse(w, stats, "Estadisticas obtenidas exitosamente", http.StatusOK)
}
