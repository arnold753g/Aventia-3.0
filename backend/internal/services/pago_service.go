package services

import (
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"andaria-backend/internal/models"

	"gorm.io/gorm"
)

type PagoService struct {
	db *gorm.DB
}

func NewPagoService(db *gorm.DB) *PagoService {
	return &PagoService{db: db}
}

func (s *PagoService) CrearPago(turistaID uint, req *models.CrearPagoRequest) (*models.PagoResponse, error) {
	var compra models.CompraPaquete
	if err := s.db.Where("id = ? AND turista_id = ?", req.CompraID, turistaID).First(&compra).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("compra no encontrada")
		}
		return nil, err
	}

	if compra.Status != "pendiente_confirmacion" {
		return nil, errors.New("la compra no está pendiente de confirmación")
	}

	if math.Abs(req.Monto-compra.PrecioTotal) > 0.01 {
		return nil, fmt.Errorf("el monto debe ser %.2f Bs", compra.PrecioTotal)
	}

	// Evitar múltiples pagos pendientes para la misma compra (previene inconsistencias con triggers).
	var existing models.PagoCompra
	if err := s.db.Where("compra_id = ? AND estado = ?", req.CompraID, "pendiente").First(&existing).Error; err == nil {
		return nil, errors.New("ya existe un pago pendiente para esta compra")
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	needsComprobante := req.MetodoPago == "qr" || req.MetodoPago == "transferencia"
	if needsComprobante && req.Comprobante == nil {
		return nil, errors.New("debe adjuntar comprobante para pagos QR o transferencia")
	}

	var comprobantePath *string
	if req.Comprobante != nil {
		path, err := saveComprobante(req.Comprobante, req.CompraID)
		if err != nil {
			return nil, err
		}
		comprobantePath = &path
	}

	pago := models.PagoCompra{
		CompraID:        req.CompraID,
		MetodoPago:      req.MetodoPago,
		Monto:           req.Monto,
		ComprobanteFoto: comprobantePath,
		Estado:          "pendiente",
	}

	if err := s.db.Create(&pago).Error; err != nil {
		return nil, err
	}

	return &models.PagoResponse{
		ID:              pago.ID,
		CompraID:        pago.CompraID,
		MetodoPago:      pago.MetodoPago,
		Monto:           pago.Monto,
		Estado:          pago.Estado,
		ComprobanteFoto: pago.ComprobanteFoto,
		Mensaje:         "Pago registrado. Esperando confirmación del encargado.",
	}, nil
}

func (s *PagoService) ObtenerPagoConContexto(pagoID uint) (*models.PagoCompra, error) {
	var pago models.PagoCompra
	if err := s.db.
		Preload("Compra.Paquete.Agencia").
		First(&pago, pagoID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("pago no encontrado")
		}
		return nil, err
	}
	return &pago, nil
}

func (s *PagoService) ConfirmarPago(pagoID uint, confirmadoPor uint, notas *string) error {
	now := time.Now()
	res := s.db.Model(&models.PagoCompra{}).
		Where("id = ? AND estado = ?", pagoID, "pendiente").
		Updates(map[string]interface{}{
			"estado":             "confirmado",
			"confirmado_por":     confirmadoPor,
			"fecha_confirmacion": now,
			"notas_encargado":    notas,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("pago no encontrado o ya fue procesado")
	}
	return nil
}

func (s *PagoService) RechazarPago(pagoID uint, confirmadoPor uint, razon string, notas *string) error {
	res := s.db.Model(&models.PagoCompra{}).
		Where("id = ? AND estado = ?", pagoID, "pendiente").
		Updates(map[string]interface{}{
			"estado":          "rechazado",
			"confirmado_por":  confirmadoPor,
			"razon_rechazo":   razon,
			"notas_encargado": notas,
		})

	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("pago no encontrado o ya fue procesado")
	}
	return nil
}

func saveComprobante(fileHeader *multipart.FileHeader, compraID uint) (string, error) {
	allowedTypes := map[string]string{
		"image/jpeg": ".jpg",
		"image/png":  ".png",
		"image/webp": ".webp",
	}

	contentType := fileHeader.Header.Get("Content-Type")
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		if guessed, ok := allowedTypes[contentType]; ok {
			ext = guessed
		}
	}

	if contentType != "" {
		if _, ok := allowedTypes[contentType]; !ok {
			return "", fmt.Errorf("formato de archivo no permitido: %s", contentType)
		}
	}

	if ext == "" {
		return "", fmt.Errorf("no se pudo determinar la extensión del archivo")
	}

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	destDir := filepath.Join("uploads", "comprobantes")
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return "", err
	}

	filename := fmt.Sprintf("comprobante_%d_%d%s", compraID, time.Now().UnixNano(), ext)
	destPath := filepath.Join(destDir, filename)

	dst, err := os.Create(destPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", err
	}

	return filepath.ToSlash(destPath), nil
}
