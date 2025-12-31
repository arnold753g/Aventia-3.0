package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
	"github.com/jung-kurt/gofpdf"
)

type reportRange struct {
	Start      time.Time
	End        time.Time
	StartLabel string
	EndLabel   string
	Month      int
	Year       int
}

type reporteVentasResumen struct {
	Ingresos      float64 `gorm:"column:ingresos"`
	Ventas        int64   `gorm:"column:ventas"`
	Participantes int64   `gorm:"column:participantes"`
}

type reporteMetodoPagoRow struct {
	MetodoPago string  `json:"metodo_pago" gorm:"column:metodo_pago"`
	Cantidad   int64   `json:"cantidad" gorm:"column:cantidad"`
	Monto      float64 `json:"monto" gorm:"column:monto"`
}

type reporteVentaRow struct {
	CompraID           uint      `json:"compra_id" gorm:"column:compra_id"`
	FechaConfirmacion  time.Time `json:"fecha_confirmacion" gorm:"column:fecha_confirmacion"`
	PaqueteID          uint      `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre      string    `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	TipoCompra         string    `json:"tipo_compra" gorm:"column:tipo_compra"`
	TotalParticipantes int       `json:"total_participantes" gorm:"column:total_participantes"`
	PrecioTotal        float64   `json:"precio_total" gorm:"column:precio_total"`
	MetodoPago         *string   `json:"metodo_pago,omitempty" gorm:"column:metodo_pago"`
	EstadoPago         *string   `json:"estado_pago,omitempty" gorm:"column:estado_pago"`
}

type reporteOcupacionSalidaRow struct {
	SalidaID         uint      `json:"salida_id" gorm:"column:salida_id"`
	PaqueteID        uint      `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre    string    `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	FechaSalida      time.Time `json:"fecha_salida" gorm:"column:fecha_salida"`
	TipoSalida       string    `json:"tipo_salida" gorm:"column:tipo_salida"`
	CupoMaximo       int       `json:"cupo_maximo" gorm:"column:cupo_maximo"`
	CuposConfirmados int       `json:"cupos_confirmados" gorm:"column:cupos_confirmados"`
	CuposReservados  int       `json:"cupos_reservados" gorm:"column:cupos_reservados"`
	Estado           string    `json:"estado" gorm:"column:estado"`
}

type reporteOcupacionPaqueteRow struct {
	PaqueteID        uint    `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre    string  `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	Salidas          int64   `json:"salidas" gorm:"column:salidas"`
	CupoMaximo       int64   `json:"cupo_maximo" gorm:"column:cupo_maximo"`
	CuposConfirmados int64   `json:"cupos_confirmados" gorm:"column:cupos_confirmados"`
	CuposReservados  int64   `json:"cupos_reservados" gorm:"column:cupos_reservados"`
	Ocupacion        float64 `json:"ocupacion"`
}

type reporteFinancieroPaqueteRow struct {
	PaqueteID     uint    `json:"paquete_id" gorm:"column:paquete_id"`
	PaqueteNombre string  `json:"paquete_nombre" gorm:"column:paquete_nombre"`
	Ventas        int64   `json:"ventas" gorm:"column:ventas"`
	Ingresos      float64 `json:"ingresos" gorm:"column:ingresos"`
	Participantes int64   `json:"participantes" gorm:"column:participantes"`
}

type reporteFinancieroPendiente struct {
	Pendientes int64   `json:"pendientes" gorm:"column:pendientes"`
	Monto      float64 `json:"monto" gorm:"column:monto"`
}

type reporteTuristaRow struct {
	TuristaID       uint    `json:"turista_id" gorm:"column:turista_id"`
	Nombre          string  `json:"nombre" gorm:"column:nombre"`
	ApellidoPaterno string  `json:"apellido_paterno" gorm:"column:apellido_paterno"`
	ApellidoMaterno string  `json:"apellido_materno" gorm:"column:apellido_materno"`
	Nacionalidad    string  `json:"nacionalidad" gorm:"column:nacionalidad"`
	Compras         int64   `json:"compras" gorm:"column:compras"`
	Participantes   int64   `json:"participantes" gorm:"column:participantes"`
	TotalGastado    float64 `json:"total_gastado" gorm:"column:total_gastado"`
}

type reporteTuristasResumen struct {
	TotalTuristas int64   `json:"total_turistas"`
	Nacionales    int64   `json:"nacionales"`
	Extranjeros   int64   `json:"extranjeros"`
	Repetidores   int64   `json:"repetidores"`
	Nuevos        int64   `json:"nuevos"`
	EdadPromedio  float64 `json:"edad_promedio"`
}

func parseReportRange(r *http.Request) (reportRange, error) {
	now := time.Now()
	loc := now.Location()

	if startStr := strings.TrimSpace(r.URL.Query().Get("fecha_inicio")); startStr != "" {
		start, err := time.ParseInLocation("2006-01-02", startStr, loc)
		if err != nil {
			return reportRange{}, fmt.Errorf("fecha_inicio invalida")
		}

		endStr := strings.TrimSpace(r.URL.Query().Get("fecha_fin"))
		if endStr == "" {
			endStr = startStr
		}
		end, err := time.ParseInLocation("2006-01-02", endStr, loc)
		if err != nil {
			return reportRange{}, fmt.Errorf("fecha_fin invalida")
		}
		if end.Before(start) {
			return reportRange{}, fmt.Errorf("rango de fechas invalido")
		}

		return reportRange{
			Start:      start,
			End:        end.AddDate(0, 0, 1),
			StartLabel: start.Format("2006-01-02"),
			EndLabel:   end.Format("2006-01-02"),
			Month:      int(start.Month()),
			Year:       start.Year(),
		}, nil
	}

	mes := int(now.Month())
	anio := now.Year()

	if value := strings.TrimSpace(r.URL.Query().Get("mes")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 1 || parsed > 12 {
			return reportRange{}, fmt.Errorf("mes invalido")
		}
		mes = parsed
	}

	if value := strings.TrimSpace(r.URL.Query().Get("anio")); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil || parsed < 2000 || parsed > 2100 {
			return reportRange{}, fmt.Errorf("anio invalido")
		}
		anio = parsed
	}

	start := time.Date(anio, time.Month(mes), 1, 0, 0, 0, 0, loc)
	end := start.AddDate(0, 1, 0)
	endLabel := end.AddDate(0, 0, -1).Format("2006-01-02")

	return reportRange{
		Start:      start,
		End:        end,
		StartLabel: start.Format("2006-01-02"),
		EndLabel:   endLabel,
		Month:      mes,
		Year:       anio,
	}, nil
}

func parseReportFormat(r *http.Request) string {
	format := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("formato")))
	if format == "" {
		return "json"
	}
	if format == "excel" {
		return "csv"
	}
	return format
}

func reportFilename(prefix string, r reportRange, format string) string {
	base := fmt.Sprintf("%s_%s_%s", prefix, r.StartLabel, r.EndLabel)
	return fmt.Sprintf("%s.%s", base, format)
}

func writeCSV(w http.ResponseWriter, filename string, rows [][]string) error {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	writer := csv.NewWriter(w)
	for _, row := range rows {
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()
	return writer.Error()
}

func writePDF(w http.ResponseWriter, filename string, pdf *gofpdf.Fpdf) error {
	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(buf.Bytes())
	return err
}

func newReportPDF(title string, agencia *models.AgenciaTurismo, r reportRange) *gofpdf.Fpdf {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(10, 12, 10)
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.Cell(0, 10, title)
	pdf.Ln(8)
	pdf.SetFont("Helvetica", "", 10)
	pdf.Cell(0, 6, fmt.Sprintf("Agencia: %s", agencia.NombreComercial))
	pdf.Ln(5)
	pdf.Cell(0, 6, fmt.Sprintf("Rango: %s a %s", r.StartLabel, r.EndLabel))
	pdf.Ln(8)
	return pdf
}

func pdfKeyValue(pdf *gofpdf.Fpdf, label string, value string) {
	pdf.SetFont("Helvetica", "B", 10)
	pdf.CellFormat(45, 6, label, "", 0, "", false, 0, "")
	pdf.SetFont("Helvetica", "", 10)
	pdf.CellFormat(0, 6, value, "", 1, "", false, 0, "")
}

func pdfTable(pdf *gofpdf.Fpdf, headers []string, widths []float64, rows [][]string) {
	lineHeight := 6.0
	header := func() {
		pdf.SetFont("Helvetica", "B", 9)
		for i, title := range headers {
			pdf.CellFormat(widths[i], lineHeight, title, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
		pdf.SetFont("Helvetica", "", 9)
	}

	header()
	for _, row := range rows {
		if pdf.GetY()+lineHeight > 280 {
			pdf.AddPage()
			header()
		}
		for i, value := range row {
			pdf.CellFormat(widths[i], lineHeight, value, "1", 0, "L", false, 0, "")
		}
		pdf.Ln(-1)
	}
}

func truncateText(value string, max int) string {
	if max <= 0 {
		return value
	}
	runes := []rune(value)
	if len(runes) <= max {
		return value
	}
	if max < 3 {
		return string(runes[:max])
	}
	return string(runes[:max-3]) + "..."
}

func (h *AgenciaHandler) loadAgenciaForReport(w http.ResponseWriter, r *http.Request) (*models.AgenciaTurismo, bool) {
	claims, ok := getClaimsOrUnauthorized(w, r)
	if !ok {
		return nil, false
	}

	vars := mux.Vars(r)
	agenciaID64, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		utils.ErrorResponse(w, "INVALID_ID", "ID invalido", nil, http.StatusBadRequest)
		return nil, false
	}

	var agencia models.AgenciaTurismo
	if err := database.GetDB().First(&agencia, agenciaID64).Error; err != nil {
		utils.ErrorResponse(w, "NOT_FOUND", "Agencia no encontrada", nil, http.StatusNotFound)
		return nil, false
	}

	if !canManageAgencia(claims, &agencia) {
		utils.ErrorResponse(w, "FORBIDDEN", "No tiene permisos para ver reportes de esta agencia", nil, http.StatusForbidden)
		return nil, false
	}

	return &agencia, true
}

// GetAgenciaReporteVentas genera el reporte de ventas.
func (h *AgenciaHandler) GetAgenciaReporteVentas(w http.ResponseWriter, r *http.Request) {
	agencia, ok := h.loadAgenciaForReport(w, r)
	if !ok {
		return
	}

	rango, err := parseReportRange(r)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	format := parseReportFormat(r)
	if format != "json" && format != "csv" && format != "pdf" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "formato invalido (json|csv|pdf)", nil, http.StatusBadRequest)
		return
	}

	var paqueteID uint64
	if value := strings.TrimSpace(r.URL.Query().Get("paquete_id")); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "paquete_id invalido", nil, http.StatusBadRequest)
			return
		}
		paqueteID = parsed
	}

	tipoCompra := strings.ToLower(strings.TrimSpace(r.URL.Query().Get("tipo_compra")))
	if tipoCompra != "" && tipoCompra != "compartido" && tipoCompra != "privado" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "tipo_compra invalido (compartido|privado)", nil, http.StatusBadRequest)
		return
	}

	conditions := []string{
		"pt.agencia_id = ?",
		"cp.status = 'confirmada'",
		"cp.fecha_confirmacion >= ?",
		"cp.fecha_confirmacion < ?",
	}
	args := []interface{}{agencia.ID, rango.Start, rango.End}

	if paqueteID > 0 {
		conditions = append(conditions, "pt.id = ?")
		args = append(args, paqueteID)
	}
	if tipoCompra != "" {
		conditions = append(conditions, "cp.tipo_compra = ?")
		args = append(args, tipoCompra)
	}

	whereClause := strings.Join(conditions, " AND ")
	db := database.GetDB()

	var resumen reporteVentasResumen
	resumenQuery := fmt.Sprintf(`
        SELECT COALESCE(SUM(cp.precio_total), 0) AS ingresos,
               COUNT(*) AS ventas,
               COALESCE(SUM(cp.total_participantes), 0) AS participantes
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE %s
    `, whereClause)
	if err := db.Raw(resumenQuery, args...).Scan(&resumen).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener resumen", err.Error(), http.StatusInternalServerError)
		return
	}

	var metodos []reporteMetodoPagoRow
	metodosQuery := fmt.Sprintf(`
        SELECT COALESCE(pc.metodo_pago, 'sin_pago') AS metodo_pago,
               COUNT(*) AS cantidad,
               COALESCE(SUM(pc.monto), 0) AS monto
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        LEFT JOIN LATERAL (
            SELECT metodo_pago, monto
            FROM pagos_compras
            WHERE compra_id = cp.id
            ORDER BY created_at DESC, id DESC
            LIMIT 1
        ) pc ON TRUE
        WHERE %s
        GROUP BY metodo_pago
        ORDER BY cantidad DESC
    `, whereClause)
	if err := db.Raw(metodosQuery, args...).Scan(&metodos).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener metodos de pago", err.Error(), http.StatusInternalServerError)
		return
	}

	var rows []reporteVentaRow
	rowsQuery := fmt.Sprintf(`
        SELECT
            cp.id AS compra_id,
            cp.fecha_confirmacion,
            pt.id AS paquete_id,
            pt.nombre AS paquete_nombre,
            cp.tipo_compra,
            cp.total_participantes,
            cp.precio_total,
            pc.metodo_pago,
            pc.estado AS estado_pago
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        LEFT JOIN LATERAL (
            SELECT metodo_pago, estado
            FROM pagos_compras
            WHERE compra_id = cp.id
            ORDER BY created_at DESC, id DESC
            LIMIT 1
        ) pc ON TRUE
        WHERE %s
        ORDER BY cp.fecha_confirmacion DESC, cp.id DESC
    `, whereClause)
	if err := db.Raw(rowsQuery, args...).Scan(&rows).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener ventas", err.Error(), http.StatusInternalServerError)
		return
	}

	promedio := 0.0
	if resumen.Ventas > 0 {
		promedio = resumen.Ingresos / float64(resumen.Ventas)
	}

	if format == "json" {
		utils.SuccessResponse(w, map[string]interface{}{
			"rango": map[string]interface{}{
				"inicio": rango.StartLabel,
				"fin":    rango.EndLabel,
			},
			"resumen": map[string]interface{}{
				"ingresos":        resumen.Ingresos,
				"ventas":          resumen.Ventas,
				"participantes":   resumen.Participantes,
				"promedio_ticket": promedio,
			},
			"metodos_pago": metodos,
			"compras":      rows,
		}, "Reporte generado", http.StatusOK)
		return
	}

	filename := reportFilename("reporte_ventas", rango, format)

	if format == "csv" {
		csvRows := [][]string{
			{"Fecha", "Paquete", "Tipo compra", "Participantes", "Monto", "Metodo pago", "Estado pago"},
		}
		for _, row := range rows {
			fecha := row.FechaConfirmacion.Format("2006-01-02")
			metodo := ""
			if row.MetodoPago != nil {
				metodo = *row.MetodoPago
			}
			estado := ""
			if row.EstadoPago != nil {
				estado = *row.EstadoPago
			}
			csvRows = append(csvRows, []string{
				fecha,
				row.PaqueteNombre,
				row.TipoCompra,
				strconv.Itoa(row.TotalParticipantes),
				fmt.Sprintf("%.2f", row.PrecioTotal),
				metodo,
				estado,
			})
		}
		if err := writeCSV(w, filename, csvRows); err != nil {
			utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar CSV", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	pdf := newReportPDF("Reporte de ventas", agencia, rango)
	pdfKeyValue(pdf, "Ingresos", fmt.Sprintf("Bs %.2f", resumen.Ingresos))
	pdfKeyValue(pdf, "Ventas", fmt.Sprintf("%d", resumen.Ventas))
	pdfKeyValue(pdf, "Participantes", fmt.Sprintf("%d", resumen.Participantes))
	pdfKeyValue(pdf, "Promedio ticket", fmt.Sprintf("Bs %.2f", promedio))
	pdf.Ln(4)

	headers := []string{"Fecha", "Paquete", "Tipo", "Pax", "Monto", "Pago"}
	widths := []float64{22, 60, 18, 12, 22, 40}
	pdfRows := make([][]string, 0, len(rows))
	for _, row := range rows {
		metodo := "sin"
		if row.MetodoPago != nil {
			metodo = *row.MetodoPago
		}
		estado := ""
		if row.EstadoPago != nil {
			estado = *row.EstadoPago
		}
		pdfRows = append(pdfRows, []string{
			row.FechaConfirmacion.Format("2006-01-02"),
			truncateText(row.PaqueteNombre, 32),
			truncateText(row.TipoCompra, 10),
			strconv.Itoa(row.TotalParticipantes),
			fmt.Sprintf("%.0f", row.PrecioTotal),
			truncateText(fmt.Sprintf("%s %s", metodo, estado), 18),
		})
	}
	pdfTable(pdf, headers, widths, pdfRows)

	if err := writePDF(w, filename, pdf); err != nil {
		utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar PDF", err.Error(), http.StatusInternalServerError)
	}
}

// GetAgenciaReporteOcupacion genera el reporte de ocupacion.
func (h *AgenciaHandler) GetAgenciaReporteOcupacion(w http.ResponseWriter, r *http.Request) {
	agencia, ok := h.loadAgenciaForReport(w, r)
	if !ok {
		return
	}

	rango, err := parseReportRange(r)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	format := parseReportFormat(r)
	if format != "json" && format != "csv" && format != "pdf" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "formato invalido (json|csv|pdf)", nil, http.StatusBadRequest)
		return
	}

	var paqueteID uint64
	if value := strings.TrimSpace(r.URL.Query().Get("paquete_id")); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "paquete_id invalido", nil, http.StatusBadRequest)
			return
		}
		paqueteID = parsed
	}

	conditions := []string{
		"pt.agencia_id = ?",
		"s.fecha_salida >= ?",
		"s.fecha_salida < ?",
		"s.estado <> 'cancelada'",
	}
	args := []interface{}{agencia.ID, rango.Start, rango.End}
	if paqueteID > 0 {
		conditions = append(conditions, "pt.id = ?")
		args = append(args, paqueteID)
	}
	whereClause := strings.Join(conditions, " AND ")
	db := database.GetDB()

	var salidas []reporteOcupacionSalidaRow
	salidasQuery := fmt.Sprintf(`
        SELECT
            s.id AS salida_id,
            pt.id AS paquete_id,
            pt.nombre AS paquete_nombre,
            s.fecha_salida,
            s.tipo_salida,
            s.cupo_maximo,
            s.cupos_confirmados,
            s.cupos_reservados,
            s.estado
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos pt ON pt.id = s.paquete_id
        WHERE %s
        ORDER BY s.fecha_salida ASC, pt.nombre ASC
    `, whereClause)
	if err := db.Raw(salidasQuery, args...).Scan(&salidas).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener salidas", err.Error(), http.StatusInternalServerError)
		return
	}

	var paquetes []reporteOcupacionPaqueteRow
	paquetesQuery := fmt.Sprintf(`
        SELECT
            pt.id AS paquete_id,
            pt.nombre AS paquete_nombre,
            COUNT(*) AS salidas,
            COALESCE(SUM(s.cupo_maximo), 0) AS cupo_maximo,
            COALESCE(SUM(s.cupos_confirmados), 0) AS cupos_confirmados,
            COALESCE(SUM(s.cupos_reservados), 0) AS cupos_reservados
        FROM paquete_salidas_habilitadas s
        JOIN paquetes_turisticos pt ON pt.id = s.paquete_id
        WHERE %s
        GROUP BY pt.id, pt.nombre
        ORDER BY pt.nombre ASC
    `, whereClause)
	if err := db.Raw(paquetesQuery, args...).Scan(&paquetes).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener ocupacion por paquete", err.Error(), http.StatusInternalServerError)
		return
	}

	var totalSalidas int64
	var totalCupo int64
	var totalConfirmados int64
	for i := range paquetes {
		if paquetes[i].CupoMaximo > 0 {
			paquetes[i].Ocupacion = float64(paquetes[i].CuposConfirmados) / float64(paquetes[i].CupoMaximo)
		}
		totalSalidas += paquetes[i].Salidas
		totalCupo += paquetes[i].CupoMaximo
		totalConfirmados += paquetes[i].CuposConfirmados
	}

	ocupacionPromedio := 0.0
	if totalCupo > 0 {
		ocupacionPromedio = float64(totalConfirmados) / float64(totalCupo)
	}

	if format == "json" {
		utils.SuccessResponse(w, map[string]interface{}{
			"rango": map[string]interface{}{
				"inicio": rango.StartLabel,
				"fin":    rango.EndLabel,
			},
			"resumen": map[string]interface{}{
				"salidas":            totalSalidas,
				"cupo_maximo":        totalCupo,
				"confirmados":        totalConfirmados,
				"ocupacion_promedio": ocupacionPromedio,
			},
			"por_paquete": paquetes,
			"por_salida":  salidas,
		}, "Reporte generado", http.StatusOK)
		return
	}

	filename := reportFilename("reporte_ocupacion", rango, format)
	if format == "csv" {
		csvRows := [][]string{
			{"Fecha", "Paquete", "Tipo", "Cupo maximo", "Confirmados", "Reservados", "Ocupacion", "Estado"},
		}
		for _, row := range salidas {
			ocupacion := 0.0
			if row.CupoMaximo > 0 {
				ocupacion = float64(row.CuposConfirmados) / float64(row.CupoMaximo)
			}
			csvRows = append(csvRows, []string{
				row.FechaSalida.Format("2006-01-02"),
				row.PaqueteNombre,
				row.TipoSalida,
				strconv.Itoa(row.CupoMaximo),
				strconv.Itoa(row.CuposConfirmados),
				strconv.Itoa(row.CuposReservados),
				fmt.Sprintf("%.2f", ocupacion*100),
				row.Estado,
			})
		}
		if err := writeCSV(w, filename, csvRows); err != nil {
			utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar CSV", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	pdf := newReportPDF("Reporte de ocupacion", agencia, rango)
	pdfKeyValue(pdf, "Salidas", fmt.Sprintf("%d", totalSalidas))
	pdfKeyValue(pdf, "Cupo maximo", fmt.Sprintf("%d", totalCupo))
	pdfKeyValue(pdf, "Confirmados", fmt.Sprintf("%d", totalConfirmados))
	pdfKeyValue(pdf, "Ocupacion promedio", fmt.Sprintf("%.1f%%", ocupacionPromedio*100))
	pdf.Ln(4)

	headers := []string{"Fecha", "Paquete", "Tipo", "Cupo", "Conf", "Res", "%", "Estado"}
	widths := []float64{20, 55, 14, 14, 12, 12, 10, 22}
	pdfRows := make([][]string, 0, len(salidas))
	for _, row := range salidas {
		ocupacion := 0.0
		if row.CupoMaximo > 0 {
			ocupacion = float64(row.CuposConfirmados) / float64(row.CupoMaximo)
		}
		pdfRows = append(pdfRows, []string{
			row.FechaSalida.Format("2006-01-02"),
			truncateText(row.PaqueteNombre, 30),
			truncateText(row.TipoSalida, 8),
			strconv.Itoa(row.CupoMaximo),
			strconv.Itoa(row.CuposConfirmados),
			strconv.Itoa(row.CuposReservados),
			fmt.Sprintf("%.0f", ocupacion*100),
			truncateText(row.Estado, 10),
		})
	}
	pdfTable(pdf, headers, widths, pdfRows)

	if err := writePDF(w, filename, pdf); err != nil {
		utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar PDF", err.Error(), http.StatusInternalServerError)
	}
}

// GetAgenciaReporteFinanciero genera el reporte financiero.
func (h *AgenciaHandler) GetAgenciaReporteFinanciero(w http.ResponseWriter, r *http.Request) {
	agencia, ok := h.loadAgenciaForReport(w, r)
	if !ok {
		return
	}

	rango, err := parseReportRange(r)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	format := parseReportFormat(r)
	if format != "json" && format != "csv" && format != "pdf" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "formato invalido (json|csv|pdf)", nil, http.StatusBadRequest)
		return
	}

	var paqueteID uint64
	if value := strings.TrimSpace(r.URL.Query().Get("paquete_id")); value != "" {
		parsed, err := strconv.ParseUint(value, 10, 32)
		if err != nil {
			utils.ErrorResponse(w, "VALIDATION_ERROR", "paquete_id invalido", nil, http.StatusBadRequest)
			return
		}
		paqueteID = parsed
	}

	conditions := []string{
		"pt.agencia_id = ?",
		"cp.status = 'confirmada'",
		"cp.fecha_confirmacion >= ?",
		"cp.fecha_confirmacion < ?",
	}
	args := []interface{}{agencia.ID, rango.Start, rango.End}
	if paqueteID > 0 {
		conditions = append(conditions, "pt.id = ?")
		args = append(args, paqueteID)
	}
	whereClause := strings.Join(conditions, " AND ")
	db := database.GetDB()

	var resumen reporteVentasResumen
	resumenQuery := fmt.Sprintf(`
        SELECT COALESCE(SUM(cp.precio_total), 0) AS ingresos,
               COUNT(*) AS ventas,
               COALESCE(SUM(cp.total_participantes), 0) AS participantes
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE %s
    `, whereClause)
	if err := db.Raw(resumenQuery, args...).Scan(&resumen).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener ingresos", err.Error(), http.StatusInternalServerError)
		return
	}

	var paquetes []reporteFinancieroPaqueteRow
	paquetesQuery := fmt.Sprintf(`
        SELECT
            pt.id AS paquete_id,
            pt.nombre AS paquete_nombre,
            COUNT(*) AS ventas,
            COALESCE(SUM(cp.precio_total), 0) AS ingresos,
            COALESCE(SUM(cp.total_participantes), 0) AS participantes
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE %s
        GROUP BY pt.id, pt.nombre
        ORDER BY ingresos DESC
    `, whereClause)
	if err := db.Raw(paquetesQuery, args...).Scan(&paquetes).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener desglose", err.Error(), http.StatusInternalServerError)
		return
	}

	var pendientes reporteFinancieroPendiente
	if err := db.Raw(`
        SELECT COUNT(*) AS pendientes, COALESCE(SUM(pc.monto), 0) AS monto
        FROM pagos_compras pc
        JOIN compras_paquetes cp ON cp.id = pc.compra_id
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE pt.agencia_id = ?
          AND pc.estado = 'pendiente'
          AND pc.created_at >= ?
          AND pc.created_at < ?
    `, agencia.ID, rango.Start, rango.End).Scan(&pendientes).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener pagos pendientes", err.Error(), http.StatusInternalServerError)
		return
	}

	var ingresosFuturos float64
	if err := db.Raw(`
        SELECT COALESCE(SUM(cp.precio_total), 0)
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        WHERE pt.agencia_id = ?
          AND cp.status = 'confirmada'
          AND cp.fecha_seleccionada >= CURRENT_DATE
    `, agencia.ID).Scan(&ingresosFuturos).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener proyeccion", err.Error(), http.StatusInternalServerError)
		return
	}

	if format == "json" {
		utils.SuccessResponse(w, map[string]interface{}{
			"rango": map[string]interface{}{
				"inicio": rango.StartLabel,
				"fin":    rango.EndLabel,
			},
			"resumen": map[string]interface{}{
				"ingresos_totales": resumen.Ingresos,
				"ventas":           resumen.Ventas,
				"participantes":    resumen.Participantes,
			},
			"ingresos_por_paquete": paquetes,
			"pagos_pendientes":     pendientes,
			"ingresos_futuros":     ingresosFuturos,
		}, "Reporte generado", http.StatusOK)
		return
	}

	filename := reportFilename("reporte_financiero", rango, format)
	if format == "csv" {
		csvRows := [][]string{
			{"Paquete", "Ventas", "Participantes", "Ingresos"},
		}
		for _, row := range paquetes {
			csvRows = append(csvRows, []string{
				row.PaqueteNombre,
				fmt.Sprintf("%d", row.Ventas),
				fmt.Sprintf("%d", row.Participantes),
				fmt.Sprintf("%.2f", row.Ingresos),
			})
		}
		if err := writeCSV(w, filename, csvRows); err != nil {
			utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar CSV", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	pdf := newReportPDF("Reporte financiero", agencia, rango)
	pdfKeyValue(pdf, "Ingresos totales", fmt.Sprintf("Bs %.2f", resumen.Ingresos))
	pdfKeyValue(pdf, "Ventas", fmt.Sprintf("%d", resumen.Ventas))
	pdfKeyValue(pdf, "Participantes", fmt.Sprintf("%d", resumen.Participantes))
	pdfKeyValue(pdf, "Pagos pendientes", fmt.Sprintf("%d (Bs %.2f)", pendientes.Pendientes, pendientes.Monto))
	pdfKeyValue(pdf, "Ingresos futuros", fmt.Sprintf("Bs %.2f", ingresosFuturos))
	pdf.Ln(4)

	headers := []string{"Paquete", "Ventas", "Pax", "Ingresos"}
	widths := []float64{80, 20, 20, 30}
	pdfRows := make([][]string, 0, len(paquetes))
	for _, row := range paquetes {
		pdfRows = append(pdfRows, []string{
			truncateText(row.PaqueteNombre, 40),
			fmt.Sprintf("%d", row.Ventas),
			fmt.Sprintf("%d", row.Participantes),
			fmt.Sprintf("%.0f", row.Ingresos),
		})
	}
	pdfTable(pdf, headers, widths, pdfRows)

	if err := writePDF(w, filename, pdf); err != nil {
		utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar PDF", err.Error(), http.StatusInternalServerError)
	}
}

// GetAgenciaReporteTuristas genera el reporte de turistas.
func (h *AgenciaHandler) GetAgenciaReporteTuristas(w http.ResponseWriter, r *http.Request) {
	agencia, ok := h.loadAgenciaForReport(w, r)
	if !ok {
		return
	}

	rango, err := parseReportRange(r)
	if err != nil {
		utils.ErrorResponse(w, "VALIDATION_ERROR", err.Error(), nil, http.StatusBadRequest)
		return
	}

	format := parseReportFormat(r)
	if format != "json" && format != "csv" && format != "pdf" {
		utils.ErrorResponse(w, "VALIDATION_ERROR", "formato invalido (json|csv|pdf)", nil, http.StatusBadRequest)
		return
	}

	db := database.GetDB()
	conditions := []string{
		"pt.agencia_id = ?",
		"cp.status = 'confirmada'",
		"cp.fecha_confirmacion >= ?",
		"cp.fecha_confirmacion < ?",
	}
	args := []interface{}{agencia.ID, rango.Start, rango.End}
	whereClause := strings.Join(conditions, " AND ")

	var turistas []reporteTuristaRow
	turistasQuery := fmt.Sprintf(`
        SELECT
            u.id AS turista_id,
            u.nombre,
            u.apellido_paterno,
            u.apellido_materno,
            COALESCE(u.nationality, 'Bolivia') AS nacionalidad,
            COUNT(*) AS compras,
            COALESCE(SUM(cp.total_participantes), 0) AS participantes,
            COALESCE(SUM(cp.precio_total), 0) AS total_gastado
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        JOIN usuarios u ON u.id = cp.turista_id
        WHERE %s
        GROUP BY u.id, u.nombre, u.apellido_paterno, u.apellido_materno, u.nationality
        ORDER BY total_gastado DESC
    `, whereClause)
	if err := db.Raw(turistasQuery, args...).Scan(&turistas).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener turistas", err.Error(), http.StatusInternalServerError)
		return
	}

	var resumen reporteTuristasResumen
	if err := db.Raw(fmt.Sprintf(`
        SELECT
            COUNT(DISTINCT u.id) AS total_turistas,
            COUNT(DISTINCT CASE WHEN LOWER(COALESCE(u.nationality, 'bolivia')) = 'bolivia' THEN u.id END) AS nacionales,
            COUNT(DISTINCT CASE WHEN LOWER(COALESCE(u.nationality, 'bolivia')) <> 'bolivia' THEN u.id END) AS extranjeros
        FROM compras_paquetes cp
        JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
        JOIN usuarios u ON u.id = cp.turista_id
        WHERE %s
    `, whereClause), args...).Scan(&resumen).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener resumen", err.Error(), http.StatusInternalServerError)
		return
	}

	var repetidores struct {
		Repetidores int64 `gorm:"column:repetidores"`
		Nuevos      int64 `gorm:"column:nuevos"`
	}
	if err := db.Raw(`
        SELECT
            SUM(CASE WHEN t.es_repetidor THEN 1 ELSE 0 END) AS repetidores,
            SUM(CASE WHEN NOT t.es_repetidor THEN 1 ELSE 0 END) AS nuevos
        FROM (
            SELECT DISTINCT cp.turista_id,
                EXISTS (
                    SELECT 1
                    FROM compras_paquetes cp2
                    JOIN paquetes_turisticos pt2 ON cp2.paquete_id = pt2.id
                    WHERE cp2.turista_id = cp.turista_id
                      AND pt2.agencia_id = ?
                      AND cp2.status = 'confirmada'
                      AND cp2.fecha_confirmacion < ?
                ) AS es_repetidor
            FROM compras_paquetes cp
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            WHERE pt.agencia_id = ?
              AND cp.status = 'confirmada'
              AND cp.fecha_confirmacion >= ?
              AND cp.fecha_confirmacion < ?
        ) t
    `, agencia.ID, rango.Start, agencia.ID, rango.Start, rango.End).Scan(&repetidores).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener repetidores", err.Error(), http.StatusInternalServerError)
		return
	}

	var edadPromedio float64
	if err := db.Raw(fmt.Sprintf(`
        SELECT COALESCE(AVG(EXTRACT(YEAR FROM age(CURRENT_DATE, t.fecha_nacimiento))), 0)
        FROM (
            SELECT DISTINCT u.id, u.fecha_nacimiento
            FROM compras_paquetes cp
            JOIN paquetes_turisticos pt ON cp.paquete_id = pt.id
            JOIN usuarios u ON u.id = cp.turista_id
            WHERE %s
        ) t
    `, whereClause), args...).Scan(&edadPromedio).Error; err != nil {
		utils.ErrorResponse(w, "DB_ERROR", "Error al obtener edad promedio", err.Error(), http.StatusInternalServerError)
		return
	}

	resumen.Repetidores = repetidores.Repetidores
	resumen.Nuevos = repetidores.Nuevos
	resumen.EdadPromedio = edadPromedio

	if format == "json" {
		utils.SuccessResponse(w, map[string]interface{}{
			"rango": map[string]interface{}{
				"inicio": rango.StartLabel,
				"fin":    rango.EndLabel,
			},
			"resumen":  resumen,
			"turistas": turistas,
		}, "Reporte generado", http.StatusOK)
		return
	}

	filename := reportFilename("reporte_turistas", rango, format)
	if format == "csv" {
		csvRows := [][]string{
			{"Turista", "Nacionalidad", "Compras", "Participantes", "Total gastado"},
		}
		for _, row := range turistas {
			nombre := strings.TrimSpace(fmt.Sprintf("%s %s %s", row.Nombre, row.ApellidoPaterno, row.ApellidoMaterno))
			csvRows = append(csvRows, []string{
				nombre,
				row.Nacionalidad,
				fmt.Sprintf("%d", row.Compras),
				fmt.Sprintf("%d", row.Participantes),
				fmt.Sprintf("%.2f", row.TotalGastado),
			})
		}
		if err := writeCSV(w, filename, csvRows); err != nil {
			utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar CSV", err.Error(), http.StatusInternalServerError)
		}
		return
	}

	pdf := newReportPDF("Reporte de turistas", agencia, rango)
	pdfKeyValue(pdf, "Total turistas", fmt.Sprintf("%d", resumen.TotalTuristas))
	pdfKeyValue(pdf, "Nacionales", fmt.Sprintf("%d", resumen.Nacionales))
	pdfKeyValue(pdf, "Extranjeros", fmt.Sprintf("%d", resumen.Extranjeros))
	pdfKeyValue(pdf, "Repetidores", fmt.Sprintf("%d", resumen.Repetidores))
	pdfKeyValue(pdf, "Nuevos", fmt.Sprintf("%d", resumen.Nuevos))
	pdfKeyValue(pdf, "Edad promedio", fmt.Sprintf("%.1f", resumen.EdadPromedio))
	pdf.Ln(4)

	headers := []string{"Turista", "Nacionalidad", "Compras", "Pax", "Total"}
	widths := []float64{70, 28, 18, 14, 24}
	pdfRows := make([][]string, 0, len(turistas))
	for _, row := range turistas {
		nombre := strings.TrimSpace(fmt.Sprintf("%s %s %s", row.Nombre, row.ApellidoPaterno, row.ApellidoMaterno))
		pdfRows = append(pdfRows, []string{
			truncateText(nombre, 40),
			truncateText(row.Nacionalidad, 14),
			fmt.Sprintf("%d", row.Compras),
			fmt.Sprintf("%d", row.Participantes),
			fmt.Sprintf("%.0f", row.TotalGastado),
		})
	}
	pdfTable(pdf, headers, widths, pdfRows)

	if err := writePDF(w, filename, pdf); err != nil {
		utils.ErrorResponse(w, "EXPORT_ERROR", "Error al generar PDF", err.Error(), http.StatusInternalServerError)
	}
}
