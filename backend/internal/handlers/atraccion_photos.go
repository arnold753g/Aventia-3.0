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

	"andaria-backend/internal/models"
)

// saveAtraccionPhoto guarda una foto de atracción en disco y retorna la ruta relativa
func saveAtraccionPhoto(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Validar tipo de contenido
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/webp": true,
		"image/jpg":  true,
	}

	contentType := header.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return "", fmt.Errorf("formato de imagen no permitido: %s", contentType)
	}

	// Determinar extensión
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
	destDir := filepath.Join("uploads", "fotografias", "atracciones")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", fmt.Errorf("error al crear directorio: %w", err)
	}

	// Generar nombre de archivo único
	filename := fmt.Sprintf("atraccion_%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	// Crear archivo y guardar contenido
	out, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("error al crear archivo: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("error al guardar archivo: %w", err)
	}

	// Retornar ruta relativa normalizada
	return filepath.ToSlash(destPath), nil
}

// deleteAtraccionPhoto elimina una foto de atracción del sistema de archivos
func deleteAtraccionPhoto(storedPath string) error {
	if storedPath == "" {
		return nil
	}

	// Directorio base esperado
	baseDir := filepath.Clean(filepath.Join(".", "uploads", "fotografias", "atracciones"))
	clean := filepath.Clean(filepath.Join(".", storedPath))

	// Evitar borrar archivos fuera de la carpeta esperada (seguridad)
	if !strings.HasPrefix(clean, baseDir) {
		return nil
	}

	// Verificar si el archivo existe antes de intentar borrarlo
	if _, err := os.Stat(clean); err == nil {
		return os.Remove(clean)
	}

	return nil
}

// parseCreateAtraccionFromForm parsea los datos del formulario multipart a CreateAtraccionRequest
func parseCreateAtraccionFromForm(r *http.Request) models.CreateAtraccionRequest {
	req := models.CreateAtraccionRequest{
		Nombre:           r.FormValue("nombre"),
		Descripcion:      r.FormValue("descripcion"),
		Direccion:        r.FormValue("direccion"),
		HorarioApertura:  r.FormValue("horario_apertura"),
		HorarioCierre:    r.FormValue("horario_cierre"),
		NivelDificultad:  r.FormValue("nivel_dificultad"),
		Status:           r.FormValue("status"),
		Telefono:         r.FormValue("telefono"),
		Email:            r.FormValue("email"),
		SitioWeb:         r.FormValue("sitio_web"),
		Facebook:         r.FormValue("facebook"),
		Instagram:        r.FormValue("instagram"),
	}

	// Parsear campos numéricos
	if provinciaID, err := strconv.ParseUint(r.FormValue("provincia_id"), 10, 32); err == nil {
		req.ProvinciaID = uint(provinciaID)
	}

	if precioEntrada, err := strconv.ParseFloat(r.FormValue("precio_entrada"), 64); err == nil {
		req.PrecioEntrada = precioEntrada
	}

	// Parsear coordenadas (pueden ser nil)
	if latStr := r.FormValue("latitud"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			req.Latitud = &lat
		}
	}

	if lngStr := r.FormValue("longitud"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			req.Longitud = &lng
		}
	}

	// Parsear booleanos
	req.RequiereAgencia = r.FormValue("requiere_agencia") == "true"
	req.AccesoParticular = r.FormValue("acceso_particular") == "true"
	req.VisiblePublico = r.FormValue("visible_publico") == "true"

	// Parsear meses (pueden ser nil)
	if mesInicioStr := r.FormValue("mes_inicio_id"); mesInicioStr != "" {
		if mesInicio, err := strconv.ParseUint(mesInicioStr, 10, 32); err == nil {
			mesInicioUint := uint(mesInicio)
			req.MesInicioID = &mesInicioUint
		}
	}

	if mesFinStr := r.FormValue("mes_fin_id"); mesFinStr != "" {
		if mesFin, err := strconv.ParseUint(mesFinStr, 10, 32); err == nil {
			mesFinUint := uint(mesFin)
			req.MesFinID = &mesFinUint
		}
	}

	// Parsear arrays (subcategorias y dias)
	if subcatsStr := r.FormValue("subcategorias_ids"); subcatsStr != "" {
		var subcats []uint
		for _, idStr := range strings.Split(subcatsStr, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				subcats = append(subcats, uint(id))
			}
		}
		req.SubcategoriasIDs = subcats
	}

	if diasStr := r.FormValue("dias_ids"); diasStr != "" {
		var dias []uint
		for _, idStr := range strings.Split(diasStr, ",") {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				dias = append(dias, uint(id))
			}
		}
		req.DiasIDs = dias
	}

	return req
}