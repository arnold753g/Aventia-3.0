package handlers

import (
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
)

func savePaqueteFoto(file multipart.File, header *multipart.FileHeader) (string, error) {
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

	destDir := filepath.Join("uploads", "fotografias", "paquetes")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("paquete_%d%s", time.Now().UnixNano(), ext)
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

func deletePaqueteFotoFile(storedPath string) error {
	if storedPath == "" {
		return nil
	}

	baseDir := filepath.Clean(filepath.Join(".", "uploads", "fotografias", "paquetes"))
	clean := filepath.Clean(filepath.Join(".", storedPath))
	if !strings.HasPrefix(clean, baseDir) {
		return nil
	}

	if _, err := os.Stat(clean); err == nil {
		return os.Remove(clean)
	}
	return nil
}

// UploadPaqueteFoto sube una foto para un paquete (máximo 6).
func (h *AgenciaHandler) UploadPaqueteFoto(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
		return
	}
	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	db := database.GetDB()
	var count int64
	db.Model(&models.PaqueteFoto{}).Where("paquete_id = ?", paqueteID).Count(&count)
	if count >= 6 {
		utils.ErrorResponse(w, "MAX_FOTOS_EXCEEDED", "El paquete ya tiene el máximo de 6 fotos", nil, http.StatusBadRequest)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, "PARSE_ERROR", "Error al procesar el formulario", err.Error(), http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("foto")
	if err != nil {
		utils.ErrorResponse(w, "NO_FILE", "No se proporcionó ningún archivo", err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 5<<20 {
		utils.ErrorResponse(w, "FILE_TOO_LARGE", "El archivo no debe superar 5MB", nil, http.StatusBadRequest)
		return
	}

	fotoPath, err := savePaqueteFoto(file, header)
	if err != nil {
		utils.ErrorResponse(w, "SAVE_ERROR", "Error al guardar la foto", err.Error(), http.StatusInternalServerError)
		return
	}

	esPrincipal := r.FormValue("es_principal") == "true"
	orden, _ := strconv.Atoi(r.FormValue("orden"))

	// Asegurar 1 principal (también en INSERT: el trigger puede no cubrir NEW.id NULL).
	if esPrincipal {
		db.Model(&models.PaqueteFoto{}).
			Where("paquete_id = ?", paqueteID).
			Update("es_principal", false)
	}

	foto := models.PaqueteFoto{
		PaqueteID:   paqueteID,
		Foto:        fotoPath,
		EsPrincipal: esPrincipal,
		Orden:       orden,
	}

	if err := db.Create(&foto).Error; err != nil {
		_ = deletePaqueteFotoFile(fotoPath)
		utils.ErrorResponse(w, "DB_ERROR", "Error al guardar en base de datos", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, foto, "Foto subida exitosamente", http.StatusCreated)
}

// RemovePaqueteFoto elimina una foto del paquete y el archivo en disco.
func (h *AgenciaHandler) RemovePaqueteFoto(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID64, err := strconv.ParseUint(vars["paquete_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de paquete inválido", nil, http.StatusBadRequest)
		return
	}
	fotoID64, err := strconv.ParseUint(vars["foto_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de foto inválido", nil, http.StatusBadRequest)
		return
	}

	agenciaID := uint(agenciaID64)
	paqueteID := uint(paqueteID64)
	fotoID := uint(fotoID64)

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}
	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var paquete models.PaqueteTuristico
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", paqueteID, agenciaID).First(&paquete).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	db := database.GetDB()
	var foto models.PaqueteFoto
	if err := db.Where("id = ? AND paquete_id = ?", fotoID, paqueteID).First(&foto).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Foto no encontrada", nil, http.StatusNotFound)
		return
	}

	_ = deletePaqueteFotoFile(foto.Foto)

	if err := db.Delete(&foto).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar foto", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Foto eliminada exitosamente", http.StatusOK)
}

