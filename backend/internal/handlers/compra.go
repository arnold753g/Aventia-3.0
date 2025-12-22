package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/internal/services"
	"andaria-backend/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type CompraHandler struct {
	validate      *validator.Validate
	compraService *services.CompraService
}

func NewCompraHandler() *CompraHandler {
	return &CompraHandler{
		validate:      validator.New(),
		compraService: services.NewCompraService(database.GetDB()),
	}
}

// CrearCompra crea una compra de paquete (solo turista).
func (h *CompraHandler) CrearCompra(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "turista" {
		utils.ErrorResponse(w, "FORBIDDEN", "Solo turistas pueden realizar compras", nil, http.StatusForbidden)
		return
	}

	var req models.CrearCompraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.compraService.CrearCompra(claims.UserID, &req)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"compra_id":    result.CompraID,
		"salida_id":    result.SalidaID,
		"precio_total": result.PrecioTotal,
	}, result.Mensaje, http.StatusCreated)
}

// ObtenerDetalleCompra retorna el detalle de una compra del turista.
func (h *CompraHandler) ObtenerDetalleCompra(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "turista" {
		utils.ErrorResponse(w, "FORBIDDEN", "Solo turistas pueden ver sus compras", nil, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	resp, err := h.compraService.ObtenerDetalleCompra(uint(id64), claims.UserID)
	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", err.Error(), nil, http.StatusNotFound)
		return
	}

	utils.SuccessResponse(w, resp, "Compra obtenida exitosamente", http.StatusOK)
}

// ListarMisCompras lista compras del turista con paginación.
func (h *CompraHandler) ListarMisCompras(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "turista" {
		utils.ErrorResponse(w, "FORBIDDEN", "Solo turistas pueden listar sus compras", nil, http.StatusForbidden)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	compras, total, err := h.compraService.ListarComprasTurista(claims.UserID, page, pageSize)
	if err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener compras", err.Error(), http.StatusInternalServerError)
		return
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"compras": compras,
		"pagination": map[string]interface{}{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	}, "Compras obtenidas exitosamente", http.StatusOK)
}

// CancelarCompra cancela una compra pendiente del turista
func (h *CompraHandler) CancelarCompra(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "turista" {
		utils.ErrorResponse(w, "FORBIDDEN", "Solo turistas pueden cancelar sus compras", nil, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	// Obtener razón del body (opcional)
	var req struct {
		Razon string `json:"razon"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	if err := h.compraService.CancelarCompra(uint(id64), claims.UserID, req.Razon); err != nil {
		utils.ErrorResponse(w, "CANCELATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Compra cancelada exitosamente", http.StatusOK)
}
