package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "time"

    "andaria-backend/internal/database"
    "andaria-backend/internal/models"
    "andaria-backend/pkg/utils"

    "github.com/gorilla/mux"
)

type agenciaDashboardMetrics struct {
    IngresosMes       float64                  `json:"ingresos_mes"`
    VentasConfirmadas int                      `json:"ventas_confirmadas"`
    TuristasAtendidos int                      `json:"turistas_atendidos"`
    PendientesPago    int                      `json:"pendientes_pago"`
    TopPaquetes       []agenciaDashboardPaquete `json:"top_paquetes"`
}

type agenciaDashboardPaquete struct {
    PaqueteID uint    `json:"paquete_id"`
    Nombre    string  `json:"nombre"`
    Ventas    int     `json:"ventas"`
    Ingresos  float64 `json:"ingresos"`
}

type dashboardMonthAggregate struct {
    Ingresos float64 `gorm:"column:ingresos"`
    Ventas   int64   `gorm:"column:ventas"`
    Turistas int64   `gorm:"column:turistas"`
}

type dashboardVentasRow struct {
    Mes      time.Time `gorm:"column:mes"`
    Ventas   int64     `gorm:"column:ventas"`
    Ingresos float64   `gorm:"column:ingresos"`
}

type dashboardOcupacionRow struct {
    Fecha            time.Time `gorm:"column:fecha"`
    CupoMaximo       int       `gorm:"column:cupo_maximo"`
    CuposConfirmados int       `gorm:"column:cupos_confirmados"`
    CuposReservados  int       `gorm:"column:cupos_reservados"`
}

// GetAgenciaDashboard retorna metricas y series para el dashboard ejecutivo.
func (h *AgenciaHandler) GetAgenciaDashboard(w http.ResponseWriter, r *http.Request) {
    claims, ok := getClaimsOrUnauthorized(w, r)
    if !ok {
        return
    }

    vars := mux.Vars(r)
    agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
    if err != nil {
        utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
        return
    }

    var agencia models.AgenciaTurismo
    if err := database.GetDB().First(&agencia, agenciaID64).Error; err != nil {
        utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
        return
    }

    if !canManageAgencia(claims, &agencia) {
        utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver el dashboard de esta agencia", nil, http.StatusForbidden)
        return
    }

    now := time.Now()
    mes := int(now.Month())
    anio := now.Year()

    if value := r.URL.Query().Get("mes"); value != "" {
        parsed, err := strconv.Atoi(value)
        if err != nil || parsed < 1 || parsed > 12 {
            utils.ErrorResponse(w, "VALIDATION_ERROR", "mes invalido (1-12)", nil, http.StatusBadRequest)
            return
        }
        mes = parsed
    }

    if value := r.URL.Query().Get("anio"); value != "" {
        parsed, err := strconv.Atoi(value)
        if err != nil || parsed < 2000 || parsed > 2100 {
            utils.ErrorResponse(w, "VALIDATION_ERROR", "anio invalido", nil, http.StatusBadRequest)
            return
        }
        anio = parsed
    }

    loc := now.Location()
    start := time.Date(anio, time.Month(mes), 1, 0, 0, 0, 0, loc)
    end := start.AddDate(0, 1, 0)
    prevStart := start.AddDate(0, -1, 0)

    db := database.GetDB()

    var raw string
    if err := db.Raw("SELECT public.get_agencia_metrics(?, ?, ?)::text", agencia.ID, mes, anio).Scan(&raw).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener metricas", err.Error(), http.StatusInternalServerError)
        return
    }

    metrics := agenciaDashboardMetrics{}
    if raw != "" {
        if err := json.Unmarshal([]byte(raw), &metrics); err != nil {
            utils.ErrorResponse(w, "DB_ERROR", "Error al leer metricas", err.Error(), http.StatusInternalServerError)
            return
        }
    }

    if metrics.TopPaquetes == nil {
        metrics.TopPaquetes = []agenciaDashboardPaquete{}
    }

    var pendientesMonto float64
    if err := db.Raw(`
        SELECT COALESCE(SUM(pc.monto), 0)
        FROM pagos_compras pc
        JOIN compras_paquetes cp ON cp.id = pc.compra_id
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE pt.agencia_id = ?
          AND pc.estado = 'pendiente'
    `, agencia.ID).Scan(&pendientesMonto).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener pagos pendientes", err.Error(), http.StatusInternalServerError)
        return
    }

    var paquetesActivos int64
    if err := db.Model(&models.PaqueteTuristico{}).
        Where("agencia_id = ? AND status = ?", agencia.ID, "activo").
        Count(&paquetesActivos).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener paquetes activos", err.Error(), http.StatusInternalServerError)
        return
    }

    var paquetesNuevos int64
    if err := db.Model(&models.PaqueteTuristico{}).
        Where("agencia_id = ? AND status = ? AND created_at >= ? AND created_at < ?", agencia.ID, "activo", start, end).
        Count(&paquetesNuevos).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener paquetes nuevos", err.Error(), http.StatusInternalServerError)
        return
    }

    var prevMetrics dashboardMonthAggregate
    if err := db.Raw(`
        SELECT
            COALESCE(SUM(cp.precio_total), 0) AS ingresos,
            COUNT(*) AS ventas,
            COALESCE(SUM(cp.total_participantes), 0) AS turistas
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE pt.agencia_id = ?
          AND cp.status = 'confirmada'
          AND cp.fecha_confirmacion >= ?
          AND cp.fecha_confirmacion < ?
    `, agencia.ID, prevStart, start).Scan(&prevMetrics).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener metricas previas", err.Error(), http.StatusInternalServerError)
        return
    }

    ingresosPercent := 0.0
    if prevMetrics.Ingresos > 0 {
        ingresosPercent = ((metrics.IngresosMes - prevMetrics.Ingresos) / prevMetrics.Ingresos) * 100
    }

    ventasDiff := int64(metrics.VentasConfirmadas) - prevMetrics.Ventas

    turistasPromedio := 0.0
    if metrics.VentasConfirmadas > 0 {
        turistasPromedio = float64(metrics.TuristasAtendidos) / float64(metrics.VentasConfirmadas)
    }

    seriesStart := start.AddDate(0, -5, 0)
    var ventasRows []dashboardVentasRow
    if err := db.Raw(`
        SELECT date_trunc('month', cp.fecha_confirmacion) AS mes,
               COUNT(*) AS ventas,
               COALESCE(SUM(cp.precio_total), 0) AS ingresos
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE pt.agencia_id = ?
          AND cp.status = 'confirmada'
          AND cp.fecha_confirmacion >= ?
          AND cp.fecha_confirmacion < ?
        GROUP BY mes
        ORDER BY mes ASC
    `, agencia.ID, seriesStart, end).Scan(&ventasRows).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener ventas por mes", err.Error(), http.StatusInternalServerError)
        return
    }

    ventasMap := make(map[string]dashboardVentasRow)
    for _, row := range ventasRows {
        key := row.Mes.Format("2006-01")
        ventasMap[key] = row
    }

    ventasMensuales := make([]map[string]interface{}, 0, 6)
    for i := 0; i < 6; i++ {
        monthDate := seriesStart.AddDate(0, i, 0)
        key := monthDate.Format("2006-01")
        row := ventasMap[key]
        ventasMensuales = append(ventasMensuales, map[string]interface{}{
            "year":     monthDate.Year(),
            "month":    int(monthDate.Month()),
            "ventas":   row.Ventas,
            "ingresos": row.Ingresos,
        })
    }

    var ocupacionRows []dashboardOcupacionRow
    if err := db.Raw(`
        SELECT s.fecha_salida::date AS fecha,
               COALESCE(SUM(s.cupo_maximo), 0) AS cupo_maximo,
               COALESCE(SUM(s.cupos_confirmados), 0) AS cupos_confirmados,
               COALESCE(SUM(s.cupos_reservados), 0) AS cupos_reservados
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos pt ON pt.id = s.paquete_id
        WHERE pt.agencia_id = ?
          AND s.fecha_salida >= ?
          AND s.fecha_salida < ?
        GROUP BY fecha
        ORDER BY fecha
    `, agencia.ID, start, end).Scan(&ocupacionRows).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener ocupacion", err.Error(), http.StatusInternalServerError)
        return
    }

    ocupacionData := make([]map[string]interface{}, 0, len(ocupacionRows))
    for _, row := range ocupacionRows {
        ocupacion := 0.0
        if row.CupoMaximo > 0 {
            ocupacion = float64(row.CuposConfirmados) / float64(row.CupoMaximo)
        }
        ocupacionData = append(ocupacionData, map[string]interface{}{
            "fecha":            row.Fecha.Format("2006-01-02"),
            "cupo_maximo":       row.CupoMaximo,
            "cupos_confirmados": row.CuposConfirmados,
            "cupos_reservados":  row.CuposReservados,
            "ocupacion":         ocupacion,
        })
    }

    dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
    dayEnd := dayStart.AddDate(0, 0, 7)
    var salidasProximas int64
    if err := db.Raw(`
        SELECT COUNT(*)
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos pt ON pt.id = s.paquete_id
        WHERE pt.agencia_id = ?
          AND s.fecha_salida >= ?
          AND s.fecha_salida < ?
          AND s.estado IN ('pendiente', 'activa')
    `, agencia.ID, dayStart, dayEnd).Scan(&salidasProximas).Error; err != nil {
        utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas proximas", err.Error(), http.StatusInternalServerError)
        return
    }

    response := map[string]interface{}{
        "mes":  mes,
        "anio": anio,
        "metrics": map[string]interface{}{
            "ingresos_mes":         metrics.IngresosMes,
            "ventas_confirmadas":   metrics.VentasConfirmadas,
            "turistas_atendidos":   metrics.TuristasAtendidos,
            "turistas_promedio":    turistasPromedio,
            "paquetes_activos":     paquetesActivos,
            "paquetes_nuevos":      paquetesNuevos,
            "pendientes_pago":      metrics.PendientesPago,
            "pendientes_pago_monto": pendientesMonto,
        },
        "comparatives": map[string]interface{}{
            "ingresos_percent": ingresosPercent,
            "ventas_diff":      ventasDiff,
        },
        "series": map[string]interface{}{
            "ventas_mensuales": ventasMensuales,
            "top_paquetes":      metrics.TopPaquetes,
            "ingresos_vs_proyeccion": map[string]interface{}{
                "confirmados": metrics.IngresosMes,
                "pendientes":  pendientesMonto,
            },
        },
        "calendario_ocupacion": ocupacionData,
        "alertas": map[string]interface{}{
            "pagos_pendientes": metrics.PendientesPago,
            "salidas_proximas": salidasProximas,
        },
    }

    utils.SuccessResponse(w, response, "Dashboard generado", http.StatusOK)
}
