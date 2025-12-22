package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"andaria-backend/internal/database"
	"andaria-backend/internal/models"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type CompraService struct {
	db *gorm.DB
}

func NewCompraService(db *gorm.DB) *CompraService {
	return &CompraService{db: db}
}

func isUndefinedFunctionError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "42883"
}

func (s *CompraService) CrearCompra(turistaID uint, req *models.CrearCompraRequest) (*models.ProcesarCompraPaqueteResult, error) {
	fecha, err := time.Parse("2006-01-02", req.FechaSeleccionada)
	if err != nil {
		return nil, fmt.Errorf("fecha_seleccionada inválida (use YYYY-MM-DD)")
	}

	ninosPagan := req.CantidadNinosPagan
	if ninosPagan < 0 {
		ninosPagan = 0
	}
	ninosGratis := req.CantidadNinosGratis
	if ninosGratis < 0 {
		ninosGratis = 0
	}
	totalParticipantes := req.CantidadAdultos + ninosPagan + ninosGratis

	// Si no existe una salida habilitada para la fecha, exigir cupo mínimo para crear la primera salida compartida.
	if req.TipoCompra == "compartido" {
		fechaStr := fecha.Format("2006-01-02")

		var salidaExiste bool
		if err := s.db.Raw(`
			SELECT EXISTS (
				SELECT 1
				FROM paquete_salidas_habilitadas
				WHERE paquete_id = ?
				  AND fecha_salida = ?
				  AND tipo_salida = 'compartido'
				  AND estado IN ('pendiente', 'activa')
			)
		`, req.PaqueteID, fechaStr).Scan(&salidaExiste).Error; err != nil {
			return nil, err
		}

		if !salidaExiste {
			var paquete models.PaqueteTuristico
			if err := s.db.Select("cupo_minimo").First(&paquete, req.PaqueteID).Error; err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return nil, err
				}
			} else {
				min := paquete.CupoMinimo
				if min < 1 {
					min = 1
				}
				if totalParticipantes < min {
					return nil, fmt.Errorf("para habilitar la primera salida en esta fecha debe registrar al menos %d participantes", min)
				}
			}
		}
	}

	query := `SELECT * FROM public.procesar_compra_paquete(?::int, ?::int, ?::date, ?::text, ?::boolean, ?::int, ?::int, ?::int, ?::boolean, ?::text, ?::text)`
	args := []interface{}{
		turistaID,
		req.PaqueteID,
		fecha,
		req.TipoCompra,
		req.Extranjero,
		req.CantidadAdultos,
		req.CantidadNinosPagan,
		req.CantidadNinosGratis,
		req.TieneDiscapacidad,
		req.DescripcionDiscapacidad,
		req.NotasTurista,
	}

	var result models.ProcesarCompraPaqueteResult
	if err := s.db.Raw(query, args...).Scan(&result).Error; err != nil {
		if isUndefinedFunctionError(err) {
			if bootstrapErr := database.ApplySQLBootstrap(s.db); bootstrapErr != nil {
				return nil, fmt.Errorf("la base de datos no est\u00e1 preparada (procesar_compra_paquete faltante): %w", bootstrapErr)
			}

			if retryErr := s.db.Raw(query, args...).Scan(&result).Error; retryErr != nil {
				return nil, retryErr
			}
		} else {
			return nil, err
		}
	}

	if !result.Success {
		if result.Mensaje == "" {
			return nil, errors.New("no se pudo procesar la compra")
		}
		return nil, errors.New(result.Mensaje)
	}

	if result.CompraID == 0 {
		return nil, errors.New("no se pudo crear la compra (respuesta vacía)")
	}

	return &result, nil
}

func (s *CompraService) ObtenerDetalleCompra(compraID uint, turistaID uint) (*models.CompraDetalleResponse, error) {
	var compra models.CompraPaquete

	err := s.db.
		Preload("Paquete").
		Preload("Salida").
		Preload("Pagos", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Order("id DESC").Limit(1)
		}).
		Where("id = ? AND turista_id = ?", compraID, turistaID).
		First(&compra).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("compra no encontrada")
		}
		return nil, err
	}

	resp := &models.CompraDetalleResponse{
		ID:                 compra.ID,
		FechaCompra:        compra.FechaCompra,
		FechaSeleccionada:  compra.FechaSeleccionada,
		TipoCompra:         compra.TipoCompra,
		TotalParticipantes: compra.TotalParticipantes,
		PrecioTotal:        compra.PrecioTotal,
		Status:             compra.Status,
		Paquete: models.PaqueteSimpleResponse{
			ID:           compra.PaqueteID,
			Nombre:       "",
			Frecuencia:   "",
			DuracionDias: nil,
			Horario:      nil,
		},
		Salida:     nil,
		UltimoPago: nil,
	}

	if compra.Paquete != nil {
		resp.Paquete = models.PaqueteSimpleResponse{
			ID:           compra.Paquete.ID,
			Nombre:       compra.Paquete.Nombre,
			Frecuencia:   compra.Paquete.Frecuencia,
			DuracionDias: compra.Paquete.DuracionDias,
			Horario:      compra.Paquete.Horario,
		}
	}

	if compra.Salida != nil && compra.SalidaID != nil {
		resp.Salida = &models.SalidaSimpleResponse{
			ID:               compra.Salida.ID,
			FechaSalida:      compra.Salida.FechaSalida,
			TipoSalida:       compra.Salida.TipoSalida,
			Estado:           compra.Salida.Estado,
			CuposDisponibles: compra.Salida.CuposDisponibles(),
		}
	}

	if len(compra.Pagos) > 0 {
		p := compra.Pagos[0]
		resp.UltimoPago = &models.PagoSimpleResponse{
			ID:                p.ID,
			MetodoPago:        p.MetodoPago,
			Monto:             p.Monto,
			Estado:            p.Estado,
			ComprobanteFoto:   p.ComprobanteFoto,
			FechaConfirmacion: p.FechaConfirmacion,
		}
	}

	return resp, nil
}

func (s *CompraService) ListarComprasTurista(turistaID uint, page int, pageSize int) ([]models.CompraDetalleResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var total int64
	if err := s.db.Model(&models.CompraPaquete{}).
		Where("turista_id = ?", turistaID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	var compras []models.CompraPaquete
	if err := s.db.
		Preload("Paquete").
		Preload("Salida").
		Where("turista_id = ?", turistaID).
		Order("fecha_compra DESC").
		Order("id DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&compras).Error; err != nil {
		return nil, 0, err
	}

	out := make([]models.CompraDetalleResponse, 0, len(compras))
	for _, compra := range compras {
		item := models.CompraDetalleResponse{
			ID:                 compra.ID,
			FechaCompra:        compra.FechaCompra,
			FechaSeleccionada:  compra.FechaSeleccionada,
			TipoCompra:         compra.TipoCompra,
			TotalParticipantes: compra.TotalParticipantes,
			PrecioTotal:        compra.PrecioTotal,
			Status:             compra.Status,
			Paquete: models.PaqueteSimpleResponse{
				ID:           compra.PaqueteID,
				Nombre:       "",
				Frecuencia:   "",
				DuracionDias: nil,
				Horario:      nil,
			},
			Salida:     nil,
			UltimoPago: nil,
		}

		if compra.Paquete != nil {
			item.Paquete = models.PaqueteSimpleResponse{
				ID:           compra.Paquete.ID,
				Nombre:       compra.Paquete.Nombre,
				Frecuencia:   compra.Paquete.Frecuencia,
				DuracionDias: compra.Paquete.DuracionDias,
				Horario:      compra.Paquete.Horario,
			}
		}

		if compra.Salida != nil && compra.SalidaID != nil {
			item.Salida = &models.SalidaSimpleResponse{
				ID:               compra.Salida.ID,
				FechaSalida:      compra.Salida.FechaSalida,
				TipoSalida:       compra.Salida.TipoSalida,
				Estado:           compra.Salida.Estado,
				CuposDisponibles: compra.Salida.CuposDisponibles(),
			}
		}

		// Último pago (N+1, acotado por paginación)
		var pago models.PagoCompra
		if err := s.db.
			Where("compra_id = ?", compra.ID).
			Order("created_at DESC").
			Order("id DESC").
			First(&pago).Error; err == nil {
			item.UltimoPago = &models.PagoSimpleResponse{
				ID:                pago.ID,
				MetodoPago:        pago.MetodoPago,
				Monto:             pago.Monto,
				Estado:            pago.Estado,
				FechaConfirmacion: pago.FechaConfirmacion,
			}
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, err
		}

		out = append(out, item)
	}

	return out, total, nil
}

// ExpirarComprasPendientes expira compras sin pago después de X minutos y libera cupos
func (s *CompraService) ExpirarComprasPendientes(minutosLimite int) (int64, error) {
	if minutosLimite < 1 {
		minutosLimite = 30 // default 30 minutos
	}

	fechaLimite := time.Now().Add(-time.Duration(minutosLimite) * time.Minute)

	// Buscar compras pendientes que no tienen ningún pago pendiente/confirmado
	// y fueron creadas hace más de X minutos
	var comprasExpirar []models.CompraPaquete
	err := s.db.
		Where("status = ?", "pendiente_confirmacion").
		Where("created_at < ?", fechaLimite).
		Where("id NOT IN (?)",
			s.db.Table("pagos_compras").
				Select("compra_id").
				Where("estado IN ('pendiente', 'confirmado')"),
		).
		Find(&comprasExpirar).Error

	if err != nil {
		return 0, fmt.Errorf("error buscando compras a expirar: %w", err)
	}

	if len(comprasExpirar) == 0 {
		return 0, nil
	}

	var expiradas int64 = 0

	for _, compra := range comprasExpirar {
		err := s.db.Transaction(func(tx *gorm.DB) error {
			// Liberar cupos reservados
			if compra.SalidaID != nil {
				if err := tx.Exec(`
					UPDATE paquete_salidas_habilitadas
					SET cupos_reservados = GREATEST(0, cupos_reservados - ?),
					    updated_at = NOW()
					WHERE id = ?
				`, compra.TotalParticipantes, *compra.SalidaID).Error; err != nil {
					return err
				}
			}

			// Marcar compra como expirada
			if err := tx.Model(&compra).Updates(map[string]interface{}{
				"status":       "expirada",
				"razon_rechazo": "Compra expirada por falta de pago",
				"updated_at":  time.Now(),
			}).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Printf("Error expirando compra %d: %v", compra.ID, err)
			continue
		}

		expiradas++
	}

	return expiradas, nil
}

// CancelarCompra cancela una compra y libera los cupos
func (s *CompraService) CancelarCompra(compraID uint, turistaID uint, razon string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		var compra models.CompraPaquete
		if err := tx.Where("id = ? AND turista_id = ?", compraID, turistaID).First(&compra).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("compra no encontrada")
			}
			return err
		}

		// Validar que se puede cancelar
		if compra.Status == "cancelada" || compra.Status == "expirada" {
			return errors.New("la compra ya está cancelada o expirada")
		}

		if compra.Status == "confirmada" {
			return errors.New("no se puede cancelar una compra confirmada, contacte a la agencia")
		}

		// Solo se puede cancelar si está pendiente_confirmacion
		if compra.Status != "pendiente_confirmacion" {
			return fmt.Errorf("no se puede cancelar una compra con estado: %s", compra.Status)
		}

		// Verificar si hay pagos pendientes
		var pagosPendientes int64
		if err := tx.Model(&models.PagoCompra{}).
			Where("compra_id = ? AND estado = ?", compraID, "pendiente").
			Count(&pagosPendientes).Error; err != nil {
			return err
		}

		if pagosPendientes > 0 {
			return errors.New("no se puede cancelar, hay pagos pendientes de revisión")
		}

		// Liberar cupos reservados
		if compra.SalidaID != nil {
			if err := tx.Exec(`
				UPDATE paquete_salidas_habilitadas
				SET cupos_reservados = GREATEST(0, cupos_reservados - ?),
				    updated_at = NOW()
				WHERE id = ?
			`, compra.TotalParticipantes, *compra.SalidaID).Error; err != nil {
				return err
			}
		}

		// Marcar compra como cancelada
		now := time.Now()
		razonFinal := "Cancelada por el turista"
		if razon != "" {
			razonFinal = razon
		}

		if err := tx.Model(&compra).Updates(map[string]interface{}{
			"status":        "cancelada",
			"razon_rechazo": razonFinal,
			"fecha_rechazo": now,
			"updated_at":    now,
		}).Error; err != nil {
			return err
		}

		return nil
	})
}

// StartExpirationWorker inicia un worker que expira compras periódicamente
func StartExpirationWorker(db *gorm.DB, minutosExpiracion int, intervaloChequeoMinutos int) {
	if intervaloChequeoMinutos < 1 {
		intervaloChequeoMinutos = 5 // chequear cada 5 minutos por defecto
	}

	service := NewCompraService(db)

	go func() {
		ticker := time.NewTicker(time.Duration(intervaloChequeoMinutos) * time.Minute)
		defer ticker.Stop()

		log.Printf("Worker de expiración iniciado: expira compras después de %d minutos, chequea cada %d minutos",
			minutosExpiracion, intervaloChequeoMinutos)

		// Ejecutar inmediatamente al iniciar
		expiradas, err := service.ExpirarComprasPendientes(minutosExpiracion)
		if err != nil {
			log.Printf("Error en expiración inicial: %v", err)
		} else if expiradas > 0 {
			log.Printf("Expiración inicial: %d compras expiradas", expiradas)
		}

		for range ticker.C {
			expiradas, err := service.ExpirarComprasPendientes(minutosExpiracion)
			if err != nil {
				log.Printf("Error en worker de expiración: %v", err)
				continue
			}
			if expiradas > 0 {
				log.Printf("Worker de expiración: %d compras expiradas", expiradas)
			}
		}
	}()
}
