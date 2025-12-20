package services

import (
	"errors"
	"fmt"
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
