package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type AtraccionHandler struct {
	validate *validator.Validate
}

func NewAtraccionHandler() *AtraccionHandler {
	return &AtraccionHandler{
		validate: validator.New(),
	}
}

// GetAtracciones lista atracciones con filtros
func (h *AtraccionHandler) GetAtracciones(w http.ResponseWriter, r *http.Request) {
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
	provinciaID := r.URL.Query().Get("provincia_id")
	departamentoID := r.URL.Query().Get("departamento_id")
	categoriaID := r.URL.Query().Get("categoria_id")
	subcategoriaID := r.URL.Query().Get("subcategoria_id")
	nivelDificultad := r.URL.Query().Get("nivel_dificultad")
	status := r.URL.Query().Get("status")
	requiereAgencia := r.URL.Query().Get("requiere_agencia")
	visiblePublico := r.URL.Query().Get("visible_publico")
	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "created_at"
	}
	sortOrder := r.URL.Query().Get("sort_order")
	if sortOrder == "" {
		sortOrder = "desc"
	}

	db := database.GetDB()
	query := db.Model(&models.AtraccionTuristica{}).
		Preload("Provincia.Departamento").
		Preload("MesInicio").
		Preload("MesFin").
		Preload("Subcategorias.Subcategoria.Categoria").
		Preload("Fotos").
		Preload("Dias")

	// Filtros
	if search != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(nombre) LIKE ? OR LOWER(descripcion) LIKE ?", searchPattern, searchPattern)
	}

	if provinciaID != "" {
		query = query.Where("provincia_id = ?", provinciaID)
	}

	if departamentoID != "" {
		query = query.Joins("JOIN provincias ON provincias.id = atracciones_turisticas.provincia_id").
			Where("provincias.departamento_id = ?", departamentoID)
	}

	if categoriaID != "" {
		query = query.Joins("JOIN atraccion_subcategorias ON atraccion_subcategorias.atraccion_id = atracciones_turisticas.id").
			Joins("JOIN subcategorias_atracciones ON subcategorias_atracciones.id = atraccion_subcategorias.subcategoria_id").
			Where("subcategorias_atracciones.categoria_id = ?", categoriaID).
			Distinct()
	}

	if subcategoriaID != "" {
		query = query.Joins("JOIN atraccion_subcategorias ON atraccion_subcategorias.atraccion_id = atracciones_turisticas.id").
			Where("atraccion_subcategorias.subcategoria_id = ?", subcategoriaID).
			Distinct()
	}

	if nivelDificultad != "" {
		query = query.Where("nivel_dificultad = ?", nivelDificultad)
	}

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if requiereAgencia != "" {
		query = query.Where("requiere_agencia = ?", requiereAgencia == "true")
	}

	if visiblePublico != "" {
		query = query.Where("visible_publico = ?", visiblePublico == "true")
	}

	var total int64
	query.Count(&total)

	offset := (page - 1) * limit

	var atracciones []models.AtraccionTuristica
	orderClause := sortBy + " " + sortOrder
	if err := query.Order(orderClause).Limit(limit).Offset(offset).Find(&atracciones).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener atracciones", err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"atracciones": atracciones,
		"pagination": map[string]interface{}{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	}

	utils.SuccessResponse(w, response, "Atracciones obtenidas exitosamente", http.StatusOK)
}

// GetAtraccion obtiene una atraccion por ID
func (h *AtraccionHandler) GetAtraccion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var atraccion models.AtraccionTuristica
	if err := database.GetDB().
		Preload("Provincia.Departamento").
		Preload("MesInicio").
		Preload("MesFin").
		Preload("Subcategorias.Subcategoria.Categoria").
		Preload("Fotos").
		Preload("Dias").
		First(&atraccion, id).Error; err != nil {
		utils.ErrorResponse(w, "ATRACCION_NOT_FOUND", "Atraccion no encontrada", nil, http.StatusNotFound)
		return
	}

	utils.SuccessResponse(w, atraccion, "Atraccion obtenida exitosamente", http.StatusOK)
}

// CreateAtraccion crea una nueva atraccion
func (h *AtraccionHandler) CreateAtraccion(w http.ResponseWriter, r *http.Request) {
	var req models.CreateAtraccionRequest
	var savedPhotoPaths []string

	contentType := r.Header.Get("Content-Type")

	// Manejar multipart/form-data (con fotos) o JSON (sin fotos)
	if strings.Contains(contentType, "multipart/form-data") {
		if err := r.ParseMultipartForm(50 << 20); err != nil { // 50MB max
			utils.ErrorResponse(w, "INVALID_FORM", "No se pudo procesar el formulario", err.Error(), http.StatusBadRequest)
			return
		}

		// Parsear datos del formulario
		req = parseCreateAtraccionFromForm(r)

		// Procesar fotos subidas (máximo 10)
		form := r.MultipartForm
		files := form.File["fotos"]
		if len(files) > 10 {
			utils.ErrorResponse(w, "TOO_MANY_FILES", "Máximo 10 fotos permitidas", nil, http.StatusBadRequest)
			return
		}

		for _, fileHeader := range files {
			file, err := fileHeader.Open()
			if err != nil {
				utils.ErrorResponse(w, "FILE_ERROR", "Error al abrir archivo", err.Error(), http.StatusBadRequest)
				return
			}
			defer file.Close()

			savedPath, err := saveAtraccionPhoto(file, fileHeader)
			if err != nil {
				// Limpiar fotos ya guardadas si hay error
				for _, path := range savedPhotoPaths {
					deleteAtraccionPhoto(path)
				}
				utils.ErrorResponse(w, "INVALID_FILE", "Error al guardar fotografía", err.Error(), http.StatusBadRequest)
				return
			}
			savedPhotoPaths = append(savedPhotoPaths, savedPath)
		}
	} else {
		// Modo JSON (sin fotos)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
			return
		}
	}

	if err := h.validate.Struct(req); err != nil {
		// Limpiar fotos guardadas si hay error de validación
		for _, path := range savedPhotoPaths {
			deleteAtraccionPhoto(path)
		}
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	// Obtener usuario autenticado
	claims, ok := r.Context().Value("claims").(*utils.JWTClaims)
	if !ok {
		utils.ErrorResponse(w, "UNAUTHORIZED", "No autorizado", nil, http.StatusUnauthorized)
		return
	}

	// Crear atraccion
	atraccion := models.AtraccionTuristica{
		Nombre:           req.Nombre,
		Descripcion:      req.Descripcion,
		ProvinciaID:      req.ProvinciaID,
		Direccion:        req.Direccion,
		Latitud:          req.Latitud,
		Longitud:         req.Longitud,
		PrecioEntrada:    req.PrecioEntrada,
		NivelDificultad:  req.NivelDificultad,
		RequiereAgencia:  req.RequiereAgencia,
		AccesoParticular: req.AccesoParticular,
		MesInicioID:      req.MesInicioID,
		MesFinID:         req.MesFinID,
		Status:           req.Status,
		VisiblePublico:   req.VisiblePublico,
		Telefono:         req.Telefono,
		Email:            strings.ToLower(req.Email),
		SitioWeb:         req.SitioWeb,
		Facebook:         req.Facebook,
		Instagram:        req.Instagram,
		CreatedBy:        claims.UserID,
	}

	// Horarios
	if req.HorarioApertura != "" {
		atraccion.HorarioApertura = &req.HorarioApertura
	}

	if req.HorarioCierre != "" {
		atraccion.HorarioCierre = &req.HorarioCierre
	}

	// Iniciar transaccion
	tx := database.GetDB().Begin()

	if err := tx.Create(&atraccion).Error; err != nil {
		tx.Rollback()
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear atraccion", err.Error(), http.StatusInternalServerError)
		return
	}

	// Agregar subcategorias
	if len(req.SubcategoriasIDs) > 0 {
		for i, subcatID := range req.SubcategoriasIDs {
			atraccionSubcat := models.AtraccionSubcategoria{
				AtraccionID:    atraccion.ID,
				SubcategoriaID: subcatID,
				EsPrincipal:    i == 0, // La primera es principal
			}
			if err := tx.Create(&atraccionSubcat).Error; err != nil {
				tx.Rollback()
				utils.ErrorResponse(w, "DB_ERROR", "Error al agregar subcategorias", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Agregar dias
	if len(req.DiasIDs) > 0 {
		for _, diaID := range req.DiasIDs {
			if err := tx.Exec("INSERT INTO atraccion_dias (atraccion_id, dia_id) VALUES (?, ?)", atraccion.ID, diaID).Error; err != nil {
				tx.Rollback()
				utils.ErrorResponse(w, "DB_ERROR", "Error al agregar dias", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Agregar fotos guardadas en el sistema de archivos
	if len(savedPhotoPaths) > 0 {
		for i, fotoPath := range savedPhotoPaths {
			foto := models.AtraccionFoto{
				AtraccionID: atraccion.ID,
				Foto:        fotoPath,
				EsPrincipal: i == 0, // La primera es principal
				Orden:       i,
			}
			if err := tx.Create(&foto).Error; err != nil {
				tx.Rollback()
				// Limpiar fotos guardadas si falla la transacción
				for _, path := range savedPhotoPaths {
					deleteAtraccionPhoto(path)
				}
				utils.ErrorResponse(w, "DB_ERROR", "Error al agregar fotos", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Si hay error en commit, limpiar fotos
	if err := tx.Commit().Error; err != nil {
		for _, path := range savedPhotoPaths {
			deleteAtraccionPhoto(path)
		}
		utils.ErrorResponse(w, "DB_ERROR", "Error al crear atracción", err.Error(), http.StatusInternalServerError)
		return
	}

	// Recargar con relaciones
	database.GetDB().
		Preload("Provincia.Departamento").
		Preload("MesInicio").
		Preload("MesFin").
		Preload("Subcategorias.Subcategoria.Categoria").
		Preload("Fotos").
		Preload("Dias").
		First(&atraccion, atraccion.ID)

	utils.SuccessResponse(w, atraccion, "Atraccion creada exitosamente", http.StatusCreated)
}

// UpdateAtraccion actualiza una atraccion
func (h *AtraccionHandler) UpdateAtraccion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var atraccion models.AtraccionTuristica
	if err := database.GetDB().First(&atraccion, id).Error; err != nil {
		utils.ErrorResponse(w, "ATRACCION_NOT_FOUND", "Atraccion no encontrada", nil, http.StatusNotFound)
		return
	}

	var req models.UpdateAtraccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validacion", err.Error(), http.StatusBadRequest)
		return
	}

	// Actualizar campos proporcionados
	if req.Nombre != "" {
		atraccion.Nombre = req.Nombre
	}
	if req.Descripcion != "" {
		atraccion.Descripcion = req.Descripcion
	}
	if req.ProvinciaID > 0 {
		atraccion.ProvinciaID = req.ProvinciaID
	}
	if req.Direccion != "" {
		atraccion.Direccion = req.Direccion
	}
	if req.Latitud != nil {
		atraccion.Latitud = req.Latitud
	}
	if req.Longitud != nil {
		atraccion.Longitud = req.Longitud
	}
	if req.PrecioEntrada != nil {
		atraccion.PrecioEntrada = *req.PrecioEntrada
	}
	if req.NivelDificultad != "" {
		atraccion.NivelDificultad = req.NivelDificultad
	}
	if req.RequiereAgencia != nil {
		atraccion.RequiereAgencia = *req.RequiereAgencia
	}
	if req.AccesoParticular != nil {
		atraccion.AccesoParticular = *req.AccesoParticular
	}
	if req.MesInicioID != nil {
		atraccion.MesInicioID = req.MesInicioID
	}
	if req.MesFinID != nil {
		atraccion.MesFinID = req.MesFinID
	}
	if req.Status != "" {
		atraccion.Status = req.Status
	}
	if req.VisiblePublico != nil {
		atraccion.VisiblePublico = *req.VisiblePublico
	}
	if req.Telefono != "" {
		atraccion.Telefono = req.Telefono
	}
	if req.Email != "" {
		atraccion.Email = strings.ToLower(req.Email)
	}
	if req.SitioWeb != "" {
		atraccion.SitioWeb = req.SitioWeb
	}
	if req.Facebook != "" {
		atraccion.Facebook = req.Facebook
	}
	if req.Instagram != "" {
		atraccion.Instagram = req.Instagram
	}

	// Horarios
	if req.HorarioApertura != "" {
		atraccion.HorarioApertura = &req.HorarioApertura
	}

	if req.HorarioCierre != "" {
		atraccion.HorarioCierre = &req.HorarioCierre
	}

	if err := database.GetDB().Save(&atraccion).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al actualizar atraccion", err.Error(), http.StatusInternalServerError)
		return
	}

	// Recargar con relaciones
	database.GetDB().
		Preload("Provincia.Departamento").
		Preload("MesInicio").
		Preload("MesFin").
		Preload("Subcategorias.Subcategoria.Categoria").
		Preload("Fotos").
		Preload("Dias").
		First(&atraccion, atraccion.ID)

	utils.SuccessResponse(w, atraccion, "Atraccion actualizada exitosamente", http.StatusOK)
}

// DeleteAtraccion desactiva una atraccion (soft delete)
func (h *AtraccionHandler) DeleteAtraccion(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var atraccion models.AtraccionTuristica
	if err := database.GetDB().First(&atraccion, id).Error; err != nil {
		utils.ErrorResponse(w, "ATRACCION_NOT_FOUND", "Atraccion no encontrada", nil, http.StatusNotFound)
		return
	}

	atraccion.Status = "inactiva"
	atraccion.VisiblePublico = false

	if err := database.GetDB().Save(&atraccion).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al desactivar atraccion", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Atraccion desactivada exitosamente", http.StatusOK)
}

// AddSubcategoria agrega una subcategoria a la atraccion
func (h *AtraccionHandler) AddSubcategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var req models.AddSubcategoriaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	// Verificar limite de 4 subcategorias
	var count int64
	database.GetDB().Model(&models.AtraccionSubcategoria{}).
		Where("atraccion_id = ?", id).
		Count(&count)

	if count >= 4 {
		utils.ErrorResponse(w, "MAX_SUBCATEGORIAS", "Una atraccion no puede tener mas de 4 subcategorias", nil, http.StatusBadRequest)
		return
	}

	atraccionSubcat := models.AtraccionSubcategoria{
		AtraccionID:    uint(id),
		SubcategoriaID: req.SubcategoriaID,
		EsPrincipal:    req.EsPrincipal,
	}

	if err := database.GetDB().Create(&atraccionSubcat).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al agregar subcategoria", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, atraccionSubcat, "Subcategoria agregada exitosamente", http.StatusCreated)
}

// RemoveSubcategoria elimina una subcategoria de la atraccion
func (h *AtraccionHandler) RemoveSubcategoria(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	atraccionID, _ := strconv.ParseUint(vars["id"], 10, 32)
	subcategoriaID, _ := strconv.ParseUint(vars["subcategoria_id"], 10, 32)

	if err := database.GetDB().
		Where("atraccion_id = ? AND subcategoria_id = ?", atraccionID, subcategoriaID).
		Delete(&models.AtraccionSubcategoria{}).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar subcategoria", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, nil, "Subcategoria eliminada exitosamente", http.StatusOK)
}

// AddFoto agrega una foto a la atraccion
func (h *AtraccionHandler) AddFoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return
	}

	var req models.AddFotoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON invalido", nil, http.StatusBadRequest)
		return
	}

	// Verificar limite de 10 fotos
	var count int64
	database.GetDB().Model(&models.AtraccionFoto{}).
		Where("atraccion_id = ?", id).
		Count(&count)

	if count >= 10 {
		utils.ErrorResponse(w, "MAX_FOTOS", "Una atraccion no puede tener mas de 10 fotos", nil, http.StatusBadRequest)
		return
	}

	foto := models.AtraccionFoto{
		AtraccionID: uint(id),
		Foto:        req.Foto,
		EsPrincipal: req.EsPrincipal,
		Orden:       req.Orden,
	}

	if err := database.GetDB().Create(&foto).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al agregar foto", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, foto, "Foto agregada exitosamente", http.StatusCreated)
}

// RemoveFoto elimina una foto de la atraccion y el archivo físico
func (h *AtraccionHandler) RemoveFoto(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	atraccionID, _ := strconv.ParseUint(vars["id"], 10, 32)
	fotoID, _ := strconv.ParseUint(vars["foto_id"], 10, 32)

	// Buscar la foto para obtener la ruta del archivo
	var foto models.AtraccionFoto
	if err := database.GetDB().
		Where("id = ? AND atraccion_id = ?", fotoID, atraccionID).
		First(&foto).Error; err != nil {
		utils.ErrorResponse(w, "FOTO_NOT_FOUND", "Foto no encontrada", nil, http.StatusNotFound)
		return
	}

	// Eliminar de la base de datos
	if err := database.GetDB().Delete(&foto).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al eliminar foto", err.Error(), http.StatusInternalServerError)
		return
	}

	// Eliminar archivo físico
	deleteAtraccionPhoto(foto.Foto)

	utils.SuccessResponse(w, nil, "Foto eliminada exitosamente", http.StatusOK)
}

// GetStats obtiene estadisticas de atracciones
func (h *AtraccionHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var stats struct {
		TotalAtracciones        int64            `json:"total_atracciones"`
		AtraccionesActivas      int64            `json:"atracciones_activas"`
		AtraccionesInactivas    int64            `json:"atracciones_inactivas"`
		AtraccionesMantenimiento int64           `json:"atracciones_mantenimiento"`
		PorDepartamento         map[string]int64 `json:"por_departamento"`
		PorCategoria            map[string]int64 `json:"por_categoria"`
		RequierenAgencia        int64            `json:"requieren_agencia"`
		AccesoParticular        int64            `json:"acceso_particular"`
	}

	db := database.GetDB()

	db.Model(&models.AtraccionTuristica{}).Count(&stats.TotalAtracciones)
	db.Model(&models.AtraccionTuristica{}).Where("status = ?", "activa").Count(&stats.AtraccionesActivas)
	db.Model(&models.AtraccionTuristica{}).Where("status = ?", "inactiva").Count(&stats.AtraccionesInactivas)
	db.Model(&models.AtraccionTuristica{}).Where("status = ?", "mantenimiento").Count(&stats.AtraccionesMantenimiento)
	db.Model(&models.AtraccionTuristica{}).Where("requiere_agencia = ?", true).Count(&stats.RequierenAgencia)
	db.Model(&models.AtraccionTuristica{}).Where("acceso_particular = ?", true).Count(&stats.AccesoParticular)

	// Por departamento
	stats.PorDepartamento = make(map[string]int64)
	var deptoStats []struct {
		Nombre string
		Total  int64
	}
	db.Raw(`
        SELECT d.nombre, COUNT(a.id) as total
        FROM departamentos d
        LEFT JOIN provincias p ON p.departamento_id = d.id
        LEFT JOIN atracciones_turisticas a ON a.provincia_id = p.id
        GROUP BY d.nombre
    `).Scan(&deptoStats)

	for _, stat := range deptoStats {
		stats.PorDepartamento[stat.Nombre] = stat.Total
	}

	// Por categoria
	stats.PorCategoria = make(map[string]int64)
	var catStats []struct {
		Nombre string
		Total  int64
	}
	db.Raw(`
        SELECT c.nombre, COUNT(DISTINCT a.id) as total
        FROM categorias_atracciones c
        LEFT JOIN subcategorias_atracciones s ON s.categoria_id = c.id
        LEFT JOIN atraccion_subcategorias acs ON acs.subcategoria_id = s.id
        LEFT JOIN atracciones_turisticas a ON a.id = acs.atraccion_id
        GROUP BY c.nombre
    `).Scan(&catStats)

	for _, stat := range catStats {
		stats.PorCategoria[stat.Nombre] = stat.Total
	}

	utils.SuccessResponse(w, stats, "Estadisticas obtenidas exitosamente", http.StatusOK)
}

// Metodos auxiliares para categorias y provincias

func (h *AtraccionHandler) GetCategorias(w http.ResponseWriter, r *http.Request) {
	var categorias []models.CategoriaAtraccion
	if err := database.GetDB().Order("orden").Find(&categorias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener categorias", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, categorias, "Categorias obtenidas exitosamente", http.StatusOK)
}

func (h *AtraccionHandler) GetSubcategorias(w http.ResponseWriter, r *http.Request) {
	categoriaID := r.URL.Query().Get("categoria_id")

	query := database.GetDB().Preload("Categoria")

	if categoriaID != "" {
		query = query.Where("categoria_id = ?", categoriaID)
	}

	var subcategorias []models.SubcategoriaAtraccion
	if err := query.Order("orden").Find(&subcategorias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener subcategorias", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, subcategorias, "Subcategorias obtenidas exitosamente", http.StatusOK)
}

func (h *AtraccionHandler) GetDepartamentos(w http.ResponseWriter, r *http.Request) {
	var departamentos []models.Departamento
	if err := database.GetDB().Find(&departamentos).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener departamentos", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, departamentos, "Departamentos obtenidos exitosamente", http.StatusOK)
}

func (h *AtraccionHandler) GetProvincias(w http.ResponseWriter, r *http.Request) {
	departamentoID := r.URL.Query().Get("departamento_id")

	query := database.GetDB().Preload("Departamento")

	if departamentoID != "" {
		query = query.Where("departamento_id = ?", departamentoID)
	}

	var provincias []models.Provincia
	if err := query.Find(&provincias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener provincias", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, provincias, "Provincias obtenidas exitosamente", http.StatusOK)
}

func (h *AtraccionHandler) GetDias(w http.ResponseWriter, r *http.Request) {
	var dias []models.Dia
	if err := database.GetDB().Find(&dias).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener dias", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, dias, "Dias obtenidos exitosamente", http.StatusOK)
}

func (h *AtraccionHandler) GetMeses(w http.ResponseWriter, r *http.Request) {
	var meses []models.Mes
	if err := database.GetDB().Find(&meses).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener meses", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, meses, "Meses obtenidos exitosamente", http.StatusOK)
}
