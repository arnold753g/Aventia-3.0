package models

import "time"

type CreateUsuarioRequest struct {
	Nombre          string `json:"nombre" validate:"required,min=2"`
	ApellidoPaterno string `json:"apellido_paterno" validate:"required,min=2"`
	ApellidoMaterno string `json:"apellido_materno" validate:"required,min=2"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	CI              string `json:"ci" validate:"required,min=5"`
	Expedido        string `json:"expedido" validate:"required,oneof=LP CB SC PT OR TJ CH BE BN PD"`
	Phone           string `json:"phone" validate:"required"`
	FechaNacimiento string `json:"fecha_nacimiento" validate:"required"`
	Ciudad          string `json:"ciudad"`
	Nationality     string `json:"nationality"`
	Rol             string `json:"rol" validate:"required,oneof=admin turista encargado_agencia"`
	ProfilePhoto    string `json:"profile_photo"`
}

type UpdateUsuarioRequest struct {
	Nombre          string `json:"nombre"`
	ApellidoPaterno string `json:"apellido_paterno"`
	ApellidoMaterno string `json:"apellido_materno"`
	Phone           string `json:"phone"`
	Ciudad          string `json:"ciudad"`
	ProfilePhoto    string `json:"profile_photo"`
}

type UpdateRolRequest struct {
	Rol string `json:"rol" validate:"required,oneof=admin turista encargado_agencia"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=active inactive suspended"`
}

type UsuarioPublic struct {
	ID              uint       `json:"id"`
	Nombre          string     `json:"nombre"`
	ApellidoPaterno string     `json:"apellido_paterno"`
	ApellidoMaterno string     `json:"apellido_materno"`
	Email           string     `json:"email"`
	CI              string     `json:"ci"`
	Expedido        string     `json:"expedido"`
	Phone           string     `json:"phone"`
	Rol             string     `json:"rol"`
	Status          string     `json:"status"`
	ProfilePhoto    string     `json:"profile_photo"`
	Ciudad          string     `json:"ciudad"`
	Nationality     string     `json:"nationality"`
	EmailVerified   bool       `json:"email_verified"`
	LastLogin       *time.Time `json:"last_login"`
	CreatedAt       time.Time  `json:"created_at"`
}

type UsuarioDetalle struct {
	ID              uint       `json:"id"`
	Nombre          string     `json:"nombre"`
	ApellidoPaterno string     `json:"apellido_paterno"`
	ApellidoMaterno string     `json:"apellido_materno"`
	Email           string     `json:"email"`
	CI              string     `json:"ci"`
	Expedido        string     `json:"expedido"`
	Phone           string     `json:"phone"`
	FechaNacimiento time.Time  `json:"fecha_nacimiento"`
	Rol             string     `json:"rol"`
	Status          string     `json:"status"`
	ProfilePhoto    string     `json:"profile_photo"`
	Ciudad          string     `json:"ciudad"`
	Nationality     string     `json:"nationality"`
	EmailVerified   bool       `json:"email_verified"`
	TermsAccepted   bool       `json:"terms_accepted"`
	TermsAcceptedAt *time.Time `json:"terms_accepted_at"`
	LastLogin       *time.Time `json:"last_login"`
	LoginAttempts   int        `json:"login_attempts"`
	LockedUntil     *time.Time `json:"locked_until"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateAgencyManagerRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Nombre          string `json:"nombre" validate:"required"`
	ApellidoPaterno string `json:"apellido_paterno" validate:"required"`
	ApellidoMaterno string `json:"apellido_materno"`
	CI              string `json:"ci" validate:"required"`
	Expedido        string `json:"expedido" validate:"required,oneof=LP CB SC PT OR TJ CH BE BN PD"`
	Telefono        string `json:"telefono" validate:"required"`
	FechaNacimiento string `json:"fecha_nacimiento" validate:"required"`
	Ciudad          string `json:"ciudad"`
	AgenciaID       uint   `json:"agencia_id" validate:"required"`
}
