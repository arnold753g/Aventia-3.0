package models

import "mime/multipart"

type CrearPagoRequest struct {
	CompraID    uint                  `validate:"required"`
	MetodoPago  string                `validate:"required,oneof=efectivo qr transferencia"`
	Monto       float64               `validate:"required,gt=0"`
	Comprobante *multipart.FileHeader `validate:"-"`
}

type ConfirmarPagoRequest struct {
	NotasEncargado *string `json:"notas_encargado"`
}

type RechazarPagoRequest struct {
	RazonRechazo string `json:"razon_rechazo" validate:"required"`
}

type PagoResponse struct {
	ID              uint    `json:"id"`
	CompraID        uint    `json:"compra_id"`
	MetodoPago      string  `json:"metodo_pago"`
	Monto           float64 `json:"monto"`
	Estado          string  `json:"estado"`
	ComprobanteFoto *string `json:"comprobante_foto"`
	Mensaje         string  `json:"mensaje"`
}
