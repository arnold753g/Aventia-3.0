package models

import (
    "time"
)

type Usuario struct {
    ID                      uint       `gorm:"primaryKey" json:"id"`
    Nombre                  string     `gorm:"size:100;not null" json:"nombre" validate:"required"`
    ApellidoPaterno         string     `gorm:"size:100;not null" json:"apellido_paterno" validate:"required"`
    ApellidoMaterno         string     `gorm:"size:100;not null" json:"apellido_materno" validate:"required"`
    FechaNacimiento         time.Time  `gorm:"type:date;not null" json:"fecha_nacimiento" validate:"required"`
    Phone                   string     `gorm:"size:20" json:"phone"`
    CI                      string     `gorm:"size:50;unique;not null" json:"ci" validate:"required"`
    Expedido                string     `gorm:"size:5" json:"expedido"`
    Email                   string     `gorm:"size:255;unique;not null" json:"email" validate:"required,email"`
    PasswordHash            string     `gorm:"size:255;not null" json:"-"`
    ProfilePhoto            string     `gorm:"type:text" json:"profile_photo"`
    Rol                     string     `gorm:"size:50;not null" json:"rol" validate:"required,oneof=admin turista encargado_agencia"`
    Status                  string     `gorm:"size:20;default:active" json:"status"`
    Nationality             string     `gorm:"size:100;default:Bolivia" json:"nationality"`
    Ciudad                  string     `gorm:"size:100" json:"ciudad"`
    EmailVerified           bool       `gorm:"default:false" json:"email_verified"`
    EmailVerificationToken  string     `gorm:"size:255" json:"-"`
    PasswordResetToken      string     `gorm:"size:255" json:"-"`
    PasswordResetExpires    *time.Time `json:"-"`
    LastLogin               *time.Time `json:"last_login"`
    LoginAttempts           int        `gorm:"default:0" json:"-"`
    LockedUntil             *time.Time `json:"-"`
    TermsAccepted           bool       `gorm:"default:false" json:"terms_accepted"`
    TermsAcceptedAt         *time.Time `json:"terms_accepted_at"`
    CreatedAt               time.Time  `json:"created_at"`
    UpdatedAt               time.Time  `json:"updated_at"`
}

func (Usuario) TableName() string {
    return "usuarios"
}