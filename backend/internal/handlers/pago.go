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

type PagoHandler struct {
	validate    *validator.Validate
	pagoService *services.PagoService
}

func NewPagoHandler() *PagoHandler {
	return &PagoHandler{
		validate:    validator.New(),
		pagoService: services.NewPagoService(database.GetDB()),
	}
}

// CrearPago registra un pago para una compra (solo turista).
func (h *PagoHandler) CrearPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "turista" {
		utils.ErrorResponse(w, "FORBIDDEN", "Solo turistas pueden registrar pagos", nil, http.StatusForbidden)
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		utils.ErrorResponse(w, "INVALID_FORM", "No se pudo procesar el formulario", nil, http.StatusBadRequest)
		return
	}

	compraID64, err := strconv.ParseUint(r.FormValue("compra_id"), 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "compra_id inválido", nil, http.StatusBadRequest)
		return
	}

	monto, err := strconv.ParseFloat(r.FormValue("monto"), 64)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "monto inválido", nil, http.StatusBadRequest)
		return
	}

	req := models.CrearPagoRequest{
		CompraID:   uint(compraID64),
		MetodoPago: r.FormValue("metodo_pago"),
		Monto:      monto,
	}

	if file, header, err := r.FormFile("comprobante"); err == nil {
		file.Close()
		req.Comprobante = header
	}

	if err := h.validate.Struct(req); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.pagoService.CrearPago(claims.UserID, &req)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, result, result.Mensaje, http.StatusCreated)
}

// ConfirmarPago confirma un pago (admin o encargado_agencia).
func (h *PagoHandler) ConfirmarPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "admin" && claims.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para confirmar pagos", nil, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	var body models.ConfirmarPagoRequest
	_ = json.NewDecoder(r.Body).Decode(&body)

	pago, err := h.pagoService.ObtenerPagoConContexto(uint(id64))
	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", err.Error(), nil, http.StatusNotFound)
		return
	}

	if claims.Rol != "admin" {
		compra := pago.Compra
		if compra == nil || compra.Paquete == nil || compra.Paquete.Agencia == nil {
			utils.ErrorResponse(w, "DB_ERROR", "No se pudo resolver la agencia asociada al pago", nil, http.StatusInternalServerError)
			return
		}
		if !canManageAgencia(claims, compra.Paquete.Agencia) {
			utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para confirmar pagos de esta agencia", nil, http.StatusForbidden)
			return
		}
	}

	if err := h.pagoService.ConfirmarPago(uint(id64), claims.UserID, body.NotasEncargado); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Pago confirmado exitosamente", http.StatusOK)
}

// RechazarPago rechaza un pago (admin o encargado_agencia).
func (h *PagoHandler) RechazarPago(w http.ResponseWriter, r *http.Request) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return
	}

	if claims.Rol != "admin" && claims.Rol != "encargado_agencia" {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para rechazar pagos", nil, http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}

	var body models.RechazarPagoRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		utils.ErrorResponse(w, "INVALID_JSON", "JSON inválido", nil, http.StatusBadRequest)
		return
	}

	if err := h.validate.Struct(body); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "Error de validación", err.Error(), http.StatusBadRequest)
		return
	}

	pago, err := h.pagoService.ObtenerPagoConContexto(uint(id64))
	if err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", err.Error(), nil, http.StatusNotFound)
		return
	}

	if claims.Rol != "admin" {
		compra := pago.Compra
		if compra == nil || compra.Paquete == nil || compra.Paquete.Agencia == nil {
			utils.ErrorResponse(w, "DB_ERROR", "No se pudo resolver la agencia asociada al pago", nil, http.StatusInternalServerError)
			return
		}
		if !canManageAgencia(claims, compra.Paquete.Agencia) {
			utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para rechazar pagos de esta agencia", nil, http.StatusForbidden)
			return
		}
	}

	if err := h.pagoService.RechazarPago(uint(id64), claims.UserID, body.RazonRechazo, nil); err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.SuccessResponse(w, nil, "Pago rechazado exitosamente", http.StatusOK)
}
