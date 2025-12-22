package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

type paqueteSalidaPublica struct {
	ID               uint   `json:"id"`
	FechaSalida      string `json:"fecha_salida"`
	TipoSalida       string `json:"tipo_salida"`
	Estado           string `json:"estado"`
	CupoMinimo       int    `json:"cupo_minimo"`
	CupoMaximo       int    `json:"cupo_maximo"`
	CuposReservados  int    `json:"cupos_reservados"`
	CuposConfirmados int    `json:"cupos_confirmados"`
	CuposDisponibles int    `json:"cupos_disponibles"`
}

// GetPaqueteSalidasPublicas lista salidas habilitadas (pendiente/activa) de un paquete visible al público.
// Uso principal: validar si existe salida para una fecha/tipo antes de permitir la primera compra (cupo mínimo).
func (h *AgenciaHandler) GetPaqueteSalidasPublicas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID inválido", nil, http.StatusBadRequest)
		return
	}
	paqueteID := uint(id64)

	fecha := strings.TrimSpace(r.URL.Query().Get("fecha")) // YYYY-MM-DD
	tipo := strings.TrimSpace(r.URL.Query().Get("tipo"))   // compartido | privado

	if fecha != "" {
		if _, err := time.Parse("2006-01-02", fecha); err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "fecha inválida (use YYYY-MM-DD)", nil, http.StatusBadRequest)
			return
		}
	}

	if tipo != "" && tipo != "compartido" && tipo != "privado" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "tipo inválido (use compartido o privado)", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()

	// Validar que el paquete sea visible al público y la agencia esté activa/visible.
	var paqueteVisible bool
	if err := db.Raw(`
		SELECT EXISTS (
			SELECT 1
			FROM paquetes_turisticos p
			JOIN agencias_turismo a ON a.id = p.agencia_id
			WHERE p.id = ?
			  AND p.status = 'activo'
			  AND p.visible_publico = TRUE
			  AND a.status = 'activa'
			  AND a.visible_publico = TRUE
		)
	`, paqueteID).Scan(&paqueteVisible).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al validar paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	if !paqueteVisible {
		utils.ErrorResponse(w, "NOT_FOUND", "Paquete no encontrado", nil, http.StatusNotFound)
		return
	}

	query := db.Model(&models.PaqueteSalidaHabilitada{}).
		Where("paquete_id = ?", paqueteID).
		Where("estado IN ('pendiente', 'activa')")

	if fecha != "" {
		query = query.Where("fecha_salida = ?", fecha)
	}
	if tipo != "" {
		query = query.Where("tipo_salida = ?", tipo)
	}

	var salidas []models.PaqueteSalidaHabilitada
	if err := query.
		Order("fecha_salida ASC").
		Order("tipo_salida ASC").
		Order("id ASC").
		Find(&salidas).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	out := make([]paqueteSalidaPublica, 0, len(salidas))
	for _, s := range salidas {
		out = append(out, paqueteSalidaPublica{
			ID:               s.ID,
			FechaSalida:      s.FechaSalida,
			TipoSalida:       s.TipoSalida,
			Estado:           s.Estado,
			CupoMinimo:       s.CupoMinimo,
			CupoMaximo:       s.CupoMaximo,
			CuposReservados:  s.CuposReservados,
			CuposConfirmados: s.CuposConfirmados,
			CuposDisponibles: s.CuposDisponibles(),
		})
	}

	utils.SuccessResponse(w, map[string]interface{}{
		"existe": len(out) > 0,
		"salidas": out,
	}, "Salidas obtenidas exitosamente", http.StatusOK)
}

