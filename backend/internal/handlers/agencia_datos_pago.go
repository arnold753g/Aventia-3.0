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

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var allowedBancos = map[string]bool{
	"Banco Nacional de Bolivia S.A.":       true,
	"Banco de Crédito de Bolivia S.A.":     true,
	"Banco Mercantil Santa Cruz S.A.":      true,
	"Banco Ganadero S.A.":                  true,
	"Banco Económico S.A.":                 true,
	"Banco Unión S.A.":                     true,
	"Banco BISA S.A.":                      true,
	"Banco FIE S.A.":                       true,
	"BancoSol S.A.":                        true,
	"Banco Ecofuturo S.A.":                 true,
	"Banco Prodem S.A.":                    true,
	"Banco Fortaleza S.A.":                 true,
}

type agenciaDatosPagoUpdateRequest struct {
	NombreBanco   *string `json:"nombre_banco"`
	NumeroCuenta  *string `json:"numero_cuenta"`
	NombreTitular *string `json:"nombre_titular"`
	Activo        *bool   `json:"activo"`
}

func ensureAgenciaDatosPagoRow(db *gorm.DB, agenciaID uint) (*models.AgenciaDatosPago, error) {
	var datos models.AgenciaDatosPago
	err := db.Where("agencia_id = ?", agenciaID).First(&datos).Error
	if err == nil {
		return &datos, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	datos = models.AgenciaDatosPago{
		AgenciaID: agenciaID,
		Activo:   true,
	}

	if err := db.Create(&datos).Error; err != nil {
		return nil, err
	}

	return &datos, nil
}

// GetAgenciaDatosPago obtiene (y crea si no existe) los datos de pago configurados por la agencia.
func (h *AgenciaHandler) GetAgenciaDatosPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	datos, err := ensureAgenciaDatosPagoRow(database.GetDB(), uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener datos de pago", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, datos, "Datos de pago obtenidos exitosamente", http.StatusOK)
}

// UpdateAgenciaDatosPago actualiza los datos de pago (transferencia) configurados por la agencia.
func (h *AgenciaHandler) UpdateAgenciaDatosPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, id).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var req agenciaDatosPagoUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if req.NombreBanco != nil && strings.TrimSpace(*req.NombreBanco) != "" {
		n := strings.TrimSpace(*req.NombreBanco)
		if !allowedBancos[n] {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "Banco no valido", nil, http.StatusBadRequest)
			return
		}
	}

	if req.NumeroCuenta != nil && len(strings.TrimSpace(*req.NumeroCuenta)) > 50 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Numero de cuenta demasiado largo", nil, http.StatusBadRequest)
		return
	}

	if req.NombreTitular != nil && len(strings.TrimSpace(*req.NombreTitular)) > 255 {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Nombre del titular demasiado largo", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	datos, err := ensureAgenciaDatosPagoRow(db, uint(id))
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al preparar datos de pago", err.Error(), http.StatusInternalServerError)
		return
	}

	if req.NombreBanco != nil {
		trimmed := strings.TrimSpace(*req.NombreBanco)
		if trimmed == "" {
			datos.NombreBanco = nil
		} else {
			datos.NombreBanco = &trimmed
		}
	}

	if req.NumeroCuenta != nil {
		trimmed := strings.TrimSpace(*req.NumeroCuenta)
		if trimmed == "" {
			datos.NumeroCuenta = nil
		} else {
			datos.NumeroCuenta = &trimmed
		}
	}

	if req.NombreTitular != nil {
		trimmed := strings.TrimSpace(*req.NombreTitular)
		if trimmed == "" {
			datos.NombreTitular = nil
		} else {
			datos.NombreTitular = &trimmed
		}
	}

	if req.Activo != nil {
		datos.Activo = *req.Activo
	}

	if err := db.Save(datos).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar datos de pago", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, datos, "Datos de pago actualizados exitosamente", http.StatusOK)
}

// saveAgenciaDatosPagoQrFoto guarda la foto del QR en disco y retorna la ruta relativa
func saveAgenciaDatosPagoQrFoto(file multipart.File, header *multipart.FileHeader) (string, error) {
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/jpg":  true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return "", fmt.Errorf("formato de imagen no permitido. Use JPG, PNG o WEBP")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext == "" {
		switch contentType {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/webp":
			ext = ".webp"
		}
	}

	destDir := filepath.Join("uploads", "fotografias", "agencias_datos_pago")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("qr_pago_%d%s", time.Now().UnixNano(), ext)
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

func deleteAgenciaDatosPagoQrFoto(storedPath string) error {
	if storedPath == "" {
		return nil
	}

	baseDir := filepath.Clean(filepath.Join(".", "uploads", "fotografias", "agencias_datos_pago"))
	clean := filepath.Clean(filepath.Join(".", storedPath))

	if !strings.HasPrefix(clean, baseDir) {
		return nil
	}

	if _, err := os.Stat(clean); err == nil {
		return os.Remove(clean)
	}

	return nil
}

// UploadAgenciaDatosPagoQrFoto sube/reemplaza la foto QR de pagos de la agencia
func (h *AgenciaHandler) UploadAgenciaDatosPagoQrFoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agenciaID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, "PARSE_ERROR", "Error al procesar el formulario", err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("qr_pago_foto")
	if err != nil {
		file, header, err = r.FormFile("foto")
	}
	if err != nil {
		utils.ErrorResponse(w, "NO_FILE", "No se proporciono ningun archivo", err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 5<<20 {
		utils.ErrorResponse(w, "FILE_TOO_LARGE", "El archivo no debe superar 5MB", nil, http.StatusBadRequest)
		return
	}

	newPath, err := saveAgenciaDatosPagoQrFoto(file, header)
	if err != nil {
		utils.ErrorResponse(w, "SAVE_ERROR", "Error al guardar la foto", err.Error(), http.StatusInternalServerError)
		return
	}

	db := database.GetDB()
	datos, err := ensureAgenciaDatosPagoRow(db, uint(agenciaID))
	if err != nil {
		_ = deleteAgenciaDatosPagoQrFoto(newPath)
		utils.ErrorResponse(w, "DB_ERROR", "Error al preparar datos de pago", err.Error(), http.StatusInternalServerError)
		return
	}

	oldPath := ""
	if datos.QrPagoFoto != nil {
		oldPath = *datos.QrPagoFoto
	}

	datos.QrPagoFoto = &newPath

	if err := db.Save(datos).Error; err != nil {
		_ = deleteAgenciaDatosPagoQrFoto(newPath)
		utils.ErrorResponse(w, "DB_ERROR", "Error al guardar en base de datos", err.Error(), http.StatusInternalServerError)
		return
	}

	if oldPath != "" && oldPath != newPath {
		_ = deleteAgenciaDatosPagoQrFoto(oldPath)
	}

	utils.SuccessResponse(w, datos, "QR actualizado exitosamente", http.StatusOK)
}

