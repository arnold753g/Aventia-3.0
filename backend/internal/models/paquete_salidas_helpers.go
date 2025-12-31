package models

// CuposDisponibles retorna cupos libres considerando reservados y confirmados.
func (s PaqueteSalidaHabilitada) CuposDisponibles() int {
	return s.CupoMaximo - s.CuposReservados - s.CuposConfirmados
}
