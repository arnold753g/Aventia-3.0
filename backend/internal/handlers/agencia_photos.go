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

// saveAgenciaFoto guarda la foto en disco y retorna la ruta relativa
func saveAgenciaFoto(file multipart.File, header *multipart.FileHeader) (string, error) {
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

	// Obtener extensión
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

	// Crear directorio si no existe
	destDir := filepath.Join("uploads", "fotografias", "agencias")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	// Generar nombre único
	filename := fmt.Sprintf("agencia_%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	// Crear archivo
	out, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Copiar contenido
	if _, err := io.Copy(out, file); err != nil {
		return "", err
	}

	return filepath.ToSlash(destPath), nil
}

// deleteAgenciaFoto elimina una foto del disco
func deleteAgenciaFoto(storedPath string) error {
	if storedPath == "" {
		return nil
	}

	baseDir := filepath.Clean(filepath.Join(".", "uploads", "fotografias", "agencias"))
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

// UploadAgenciaFoto sube una foto de agencia
func (h *AgenciaHandler) UploadAgenciaFoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agenciaID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	// Verificar que la agencia exista
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

	// Verificar límite de 10 fotos
	var fotosCount int64
	database.GetDB().Model(&models.AgenciaFoto{}).Where("agencia_id = ?", agenciaID).Count(&fotosCount)
	if fotosCount >= 10 {
		utils.ErrorResponse(w, "MAX_FOTOS_EXCEEDED", "La agencia ya tiene el máximo de 10 fotos", nil, http.StatusBadRequest)
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, "PARSE_ERROR", "Error al procesar el formulario", err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener archivo
	file, header, err := r.FormFile("foto")
	if err != nil {
		utils.ErrorResponse(w, "NO_FILE", "No se proporcionó ningún archivo", err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validar tamaño (max 5MB)
	if header.Size > 5<<20 {
		utils.ErrorResponse(w, "FILE_TOO_LARGE", "El archivo no debe superar 5MB", nil, http.StatusBadRequest)
		return
	}

	// Guardar foto
	fotoPath, err := saveAgenciaFoto(file, header)
	if err != nil {
		utils.ErrorResponse(w, "SAVE_ERROR", "Error al guardar la foto", err.Error(), http.StatusInternalServerError)
		return
	}

	// Obtener datos adicionales del form
	titulo := r.FormValue("titulo")
	descripcion := r.FormValue("descripcion")
	esPrincipal := r.FormValue("es_principal") == "true"
	orden, _ := strconv.Atoi(r.FormValue("orden"))

	// Si es principal, quitar flag de otras fotos
	if esPrincipal {
		database.GetDB().Model(&models.AgenciaFoto{}).
			Where("agencia_id = ?", agenciaID).
			Update("es_principal", false)
	}

	// Crear registro en BD
	foto := models.AgenciaFoto{
		AgenciaID:   uint(agenciaID),
		FotoURL:     fotoPath,
		Titulo:      titulo,
		Descripcion: descripcion,
		EsPrincipal: esPrincipal,
		Orden:       orden,
	}

	if err := database.GetDB().Create(&foto).Error; err != nil {
		// Si falla, eliminar archivo subido
		deleteAgenciaFoto(fotoPath)
		utils.ErrorResponse(w, "DB_ERROR", "Error al guardar en base de datos", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, foto, "Foto subida exitosamente", http.StatusCreated)
}

// RemoveFotoWithFile elimina una foto de la agencia y el archivo del disco
func (h *AgenciaHandler) RemoveFotoWithFile(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	vars := mux.Vars(r)
	agenciaID, _ := strconv.ParseUint(vars["id"], 10, 32)
	fotoID, err := strconv.ParseUint(vars["foto_id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID de foto inválido", nil, http.StatusBadRequest)
		return
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para gestionar esta agencia", nil, http.StatusForbidden)
		return
	}

	var foto models.AgenciaFoto
	if err := database.GetDB().Where("id = ? AND agencia_id = ?", fotoID, agenciaID).First(&foto).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Foto no encontrada", nil, http.StatusNotFound)
		return
	}

	// Eliminar archivo del disco
	if err := deleteAgenciaFoto(foto.FotoURL); err != nil {
		// Registrar error pero continuar
		fmt.Printf("Error al eliminar archivo: %v\n", err)
	}

	// Eliminar de BD
	if err := database.GetDB().Delete(&foto).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar foto", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Foto eliminada exitosamente", http.StatusOK)
}
