package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"andaria-backend/internal/config"
	"andaria-backend/internal/database"
	"andaria-backend/internal/handlers"
	"andaria-backend/internal/middleware"
	"andaria-backend/internal/services"
	"andaria-backend/internal/websocket"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuracion
	cfg := config.LoadConfig()

	validateEnv()

	// Inicializar JWT
	utils.InitJWT(cfg.JWTSecret)
	log.Println("OK. JWT initialized")

	// Conectar a la base de datos
	if err := database.Connect(cfg); err != nil {
		log.Fatal("Error connecting to database:", err)
	}

	// Iniciar worker de expiración de compras
	minutosExpiracion, _ := strconv.Atoi(cfg.CompraExpiracionMinutos)
	if minutosExpiracion < 1 {
		minutosExpiracion = 30
	}
	services.StartExpirationWorker(database.GetDB(), minutosExpiracion, 5)
	log.Printf("OK. Worker de expiración de compras iniciado (%d minutos)", minutosExpiracion)

	// Iniciar WebSocket Hub
	hub := websocket.NewHub()
	go hub.Run()
	log.Println("OK. WebSocket Hub iniciado")

	// Iniciar PostgreSQL Listener para notificaciones
	listener := services.NewNotificationListener(database.GetConnPool(), hub)
	go func() {
		if err := listener.Start(); err != nil {
			log.Printf("Error en PostgreSQL Listener: %v", err)
		}
	}()
	log.Println("OK. PostgreSQL Listener iniciado")

	// Crear router
	router := mux.NewRouter()

	// Servir archivos estaticos (fotografias)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Setup CORS
	corsHandler := middleware.SetupCORS()

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Handlers
	authHandler := handlers.NewAuthHandler()
	usuarioHandler := handlers.NewUsuarioHandler()
	agenciaHandler := handlers.NewAgenciaHandler()
	agenciaVisitasHandler := handlers.NewAgenciaVisitasHandler()
	compraHandler := handlers.NewCompraHandler()
	pagoHandler := handlers.NewPagoHandler()
	atraccionHandler := handlers.NewAtraccionHandler()
	notificacionHandler := handlers.NewNotificacionHandler()
	wsHandler := handlers.NewWebSocketHandler(hub)
	salidaHandler := handlers.NewSalidaHandler()

	// ========== RUTAS PÚBLICAS (sin autenticación) ==========
	// Aplicar rate limiting (100 requests/minuto) y caché (5 minutos)
	publicAPI := api.PathPrefix("/public").Subrouter()
	publicAPI.Use(middleware.RateLimitMiddleware(100))
	publicAPI.Use(middleware.CacheMiddleware(5 * time.Minute))

	// Paquetes turísticos públicos
	publicAPI.HandleFunc("/paquetes", agenciaHandler.GetPaquetesPublicos).Methods("GET")
	publicAPI.HandleFunc("/paquetes/{id:[0-9]+}", agenciaHandler.GetPaquetePublico).Methods("GET")
	publicAPI.HandleFunc("/paquetes/{id:[0-9]+}/salidas", agenciaHandler.GetPaqueteSalidasPublicas).Methods("GET")
	publicAPI.HandleFunc("/salidas-confirmadas", agenciaHandler.GetSalidasConfirmadasPublicas).Methods("GET")
	publicAPI.HandleFunc("/salidas-disponibles", salidaHandler.ObtenerSalidasPublicas).Methods("GET")

	// Agencias públicas
	publicAPI.HandleFunc("/agencias", agenciaHandler.GetAgencias).Methods("GET")
	publicAPI.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.GetAgencia).Methods("GET")
	// Agencias por slug (público, con caché) - Excluye palabras reservadas
	publicAPI.HandleFunc("/agencias/{id:[a-zA-Z0-9-]+}", agenciaHandler.GetAgencia).Methods("GET")
	// Registrar visitas (público, sin caché porque es POST)
	publicAPI.HandleFunc("/agencias/{id:[0-9]+}/visitas", agenciaVisitasHandler.RegistrarVisita).Methods("POST")
	publicAPI.HandleFunc("/agencias/{id:[a-zA-Z0-9-]+}/visitas", agenciaVisitasHandler.RegistrarVisita).Methods("POST")

	// Rutas públicas sin caché para agencias (soporta ID o slug) - DEPRECADAS, usa /public/
	// IMPORTANTE: Usa regex para excluir palabras reservadas como "me", "rapida", "completa", "data"
	api.HandleFunc("/agencias/{id:(?!me$|rapida$|completa$|data$)[a-zA-Z0-9-]+}", agenciaHandler.GetAgencia).Methods("GET")
	api.HandleFunc("/agencias/{id:(?!me$|rapida$|completa$|data$)[a-zA-Z0-9-]+}/visitas", agenciaVisitasHandler.RegistrarVisita).Methods("POST")

	// Atracciones públicas
	publicAPI.HandleFunc("/atracciones", atraccionHandler.GetAtracciones).Methods("GET")
	publicAPI.HandleFunc("/atracciones/{id:[0-9]+}", atraccionHandler.GetAtraccion).Methods("GET")

	// Datos auxiliares públicos (sin caché ya que rara vez cambian)
	dataAPI := api.PathPrefix("/data").Subrouter()
	dataAPI.HandleFunc("/departamentos", agenciaHandler.GetDepartamentos).Methods("GET")
	dataAPI.HandleFunc("/provincias", atraccionHandler.GetProvincias).Methods("GET")
	dataAPI.HandleFunc("/categorias", atraccionHandler.GetCategorias).Methods("GET")
	dataAPI.HandleFunc("/subcategorias", atraccionHandler.GetSubcategorias).Methods("GET")
	dataAPI.HandleFunc("/dias", atraccionHandler.GetDias).Methods("GET")
	dataAPI.HandleFunc("/meses", atraccionHandler.GetMeses).Methods("GET")

	// Auth routes (públicas pero sin caché)
	auth := api.PathPrefix("/auth").Subrouter()
	auth.Use(middleware.RateLimitMiddleware(20)) // Más restrictivo: 20 req/min
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST")
	auth.HandleFunc("/verify-email", authHandler.VerifyEmail).Methods("POST")
	auth.HandleFunc("/resend-email-code", authHandler.ResendEmailCode).Methods("POST")
	auth.HandleFunc("/forgot-password", authHandler.ForgotPassword).Methods("POST")
	auth.HandleFunc("/reset-password", authHandler.ResetPassword).Methods("POST")
	auth.HandleFunc("/set-initial-password", authHandler.SetInitialPassword).Methods("POST")

	// Check duplicados (público sin caché)
	api.HandleFunc("/usuarios/check", usuarioHandler.CheckUsuarioExiste).Methods("GET")

	// ========== WEBSOCKET (requiere autenticación vía token en query) ==========
	api.HandleFunc("/ws", wsHandler.HandleWebSocket)

	// ========== RUTAS PROTEGIDAS (requieren autenticación) ==========
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Profile
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")

	authProtected := protected.PathPrefix("/auth").Subrouter()
	authProtected.HandleFunc("/change-password", authHandler.ChangePassword).Methods("POST")

	// ========== NOTIFICACIONES ==========
	protected.HandleFunc("/notificaciones", notificacionHandler.GetNotificaciones).Methods("GET")
	protected.HandleFunc("/notificaciones/no-leidas/count", notificacionHandler.GetContadorNoLeidas).Methods("GET")
	protected.HandleFunc("/notificaciones/{id:[0-9]+}/marcar-leida", notificacionHandler.MarcarComoLeida).Methods("PUT")
	protected.HandleFunc("/notificaciones/marcar-todas-leidas", notificacionHandler.MarcarTodasLeidas).Methods("PUT")
	protected.HandleFunc("/notificaciones/{id:[0-9]+}", notificacionHandler.EliminarNotificacion).Methods("DELETE")

	// ========== RUTAS DE USUARIOS ==========
	protected.HandleFunc("/usuarios/{id:[0-9]+}", usuarioHandler.GetUsuario).Methods("GET")
	protected.HandleFunc("/usuarios/{id:[0-9]+}", usuarioHandler.UpdateUsuario).Methods("PUT")

	// ========== RUTAS DE AGENCIAS TURISTICAS (Protegidas) ==========
	// Datos auxiliares protegidos
	protected.HandleFunc("/agencias/data/encargados", agenciaHandler.GetEncargados).Methods("GET")

	// Rutas protegidas (requieren autenticación)
	protected.HandleFunc("/agencias/rapida", agenciaHandler.CreateAgenciaRapida).Methods("POST")
	protected.HandleFunc("/agencias/completa", agenciaHandler.CreateAgenciaCompleta).Methods("POST")
	protected.HandleFunc("/agencias/me", agenciaHandler.GetMiAgencia).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/dashboard", agenciaHandler.GetAgenciaDashboard).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/reportes/ventas", agenciaHandler.GetAgenciaReporteVentas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/reportes/ocupacion", agenciaHandler.GetAgenciaReporteOcupacion).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/reportes/financiero", agenciaHandler.GetAgenciaReporteFinanciero).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/reportes/turistas", agenciaHandler.GetAgenciaReporteTuristas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.UpdateAgencia).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/fotos/upload", agenciaHandler.UploadAgenciaFoto).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/fotos/{foto_id:[0-9]+}", agenciaHandler.RemoveFotoWithFile).Methods("DELETE")
	protected.HandleFunc("/agencias/{id:[0-9]+}/especialidades", agenciaHandler.AddEspecialidad).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/especialidades/{especialidad_id:[0-9]+}", agenciaHandler.RemoveEspecialidad).Methods("DELETE")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquete-politicas", agenciaHandler.GetPaquetePoliticas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquete-politicas", agenciaHandler.UpdatePaquetePoliticas).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/datos-pago", agenciaHandler.GetAgenciaDatosPago).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/datos-pago", agenciaHandler.UpdateAgenciaDatosPago).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/datos-pago/qr/upload", agenciaHandler.UploadAgenciaDatosPagoQrFoto).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/capacidad", agenciaHandler.GetAgenciaCapacidad).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/capacidad", agenciaHandler.UpdateAgenciaCapacidad).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/ventas/pagos", agenciaHandler.GetAgenciaVentasPagos).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/ventas/salidas", agenciaHandler.GetAgenciaVentasSalidas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/ventas/salidas/{salida_id:[0-9]+}/compras", agenciaHandler.GetAgenciaVentasSalidaCompras).Methods("GET")

	// Estadísticas de visitas (solo para encargado de la agencia o admin)
	protected.HandleFunc("/agencias/{id:[0-9]+}/estadisticas-visitas", agenciaVisitasHandler.GetEstadisticasVisitas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/visitas-detalle", agenciaVisitasHandler.GetVisitasDetalle).Methods("GET")

	// ========== PAQUETES TURISTICOS (Encargado/Admin) ==========
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes", agenciaHandler.GetAgenciaPaquetes).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes", agenciaHandler.CreateAgenciaPaquete).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}", agenciaHandler.GetAgenciaPaquete).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}", agenciaHandler.UpdateAgenciaPaquete).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}", agenciaHandler.DeleteAgenciaPaquete).Methods("DELETE")

	// Fotos del paquete (máximo 6)
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/fotos/upload", agenciaHandler.UploadPaqueteFoto).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/fotos/{foto_id:[0-9]+}", agenciaHandler.RemovePaqueteFoto).Methods("DELETE")

	// Itinerario
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/itinerario", agenciaHandler.GetPaqueteItinerario).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/itinerario", agenciaHandler.CreatePaqueteItinerario).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/itinerario/{itinerario_id:[0-9]+}", agenciaHandler.UpdatePaqueteItinerario).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/itinerario/{itinerario_id:[0-9]+}", agenciaHandler.DeletePaqueteItinerario).Methods("DELETE")

	// Atracciones del paquete
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/atracciones", agenciaHandler.GetPaqueteAtracciones).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/atracciones", agenciaHandler.AddPaqueteAtraccion).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/atracciones/{paquete_atraccion_id:[0-9]+}", agenciaHandler.UpdatePaqueteAtraccion).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/atracciones/{paquete_atraccion_id:[0-9]+}", agenciaHandler.RemovePaqueteAtraccion).Methods("DELETE")

	// Salidas habilitadas (edición logística/estado)
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas", agenciaHandler.GetPaqueteSalidas).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas", agenciaHandler.CreatePaqueteSalida).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas/{salida_id:[0-9]+}", agenciaHandler.UpdatePaqueteSalida).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas/{salida_id:[0-9]+}/activar", agenciaHandler.ActivarSalida).Methods("POST")

	// ========== GESTIÓN DE SALIDAS MANUALES (Encargado Agencia) ==========
	// Nuevas rutas para crear y gestionar salidas manualmente
	protected.HandleFunc("/agencias/paquetes/{paquete_id:[0-9]+}/salidas-manuales", salidaHandler.CrearSalidaManual).Methods("POST")
	protected.HandleFunc("/agencias/paquetes/{paquete_id:[0-9]+}/salidas-manuales", salidaHandler.ObtenerSalidasPorPaquete).Methods("GET")
	protected.HandleFunc("/agencias/salidas/{salida_id:[0-9]+}", salidaHandler.ActualizarSalida).Methods("PUT")
	protected.HandleFunc("/agencias/salidas/{salida_id:[0-9]+}/cancelar", salidaHandler.CancelarSalida).Methods("POST")

	// ========== COMPRAS DE PAQUETES (Turista) ==========
	protected.HandleFunc("/compras", compraHandler.CrearCompra).Methods("POST")
	protected.HandleFunc("/compras/{id:[0-9]+}", compraHandler.ObtenerDetalleCompra).Methods("GET")
	protected.HandleFunc("/compras/{id:[0-9]+}/cancelar", compraHandler.CancelarCompra).Methods("POST")
	protected.HandleFunc("/mis-compras", compraHandler.ListarMisCompras).Methods("GET")

	// ========== PAGOS DE COMPRAS ==========
	protected.HandleFunc("/pagos", pagoHandler.CrearPago).Methods("POST")

	pagosManager := protected.PathPrefix("").Subrouter()
	pagosManager.Use(middleware.RoleMiddleware("admin", "encargado_agencia"))
	pagosManager.HandleFunc("/pagos/{id:[0-9]+}/confirmar", pagoHandler.ConfirmarPago).Methods("PUT")
	pagosManager.HandleFunc("/pagos/{id:[0-9]+}/rechazar", pagoHandler.RechazarPago).Methods("PUT")

	// ========== RUTAS DE ATRACCIONES TURISTICAS (Protegidas) ==========
	// Rutas protegidas (requieren autenticacion)
	protected.HandleFunc("/atracciones", atraccionHandler.CreateAtraccion).Methods("POST")
	protected.HandleFunc("/atracciones/{id}", atraccionHandler.UpdateAtraccion).Methods("PUT")
	protected.HandleFunc("/atracciones/{id}/subcategorias", atraccionHandler.AddSubcategoria).Methods("POST")
	protected.HandleFunc("/atracciones/{id}/subcategorias/{subcategoria_id}", atraccionHandler.RemoveSubcategoria).Methods("DELETE")
	protected.HandleFunc("/atracciones/{id}/fotos", atraccionHandler.AddFoto).Methods("POST")
	protected.HandleFunc("/atracciones/{id}/fotos/{foto_id}", atraccionHandler.RemoveFoto).Methods("DELETE")

	// Rutas solo para admin
	adminRouter := protected.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.RoleMiddleware("admin"))

	adminRouter.HandleFunc("/usuarios", usuarioHandler.GetUsuarios).Methods("GET")
	adminRouter.HandleFunc("/usuarios", usuarioHandler.CreateUsuario).Methods("POST")
	adminRouter.HandleFunc("/agency-managers", usuarioHandler.CreateAgencyManager).Methods("POST")
	adminRouter.HandleFunc("/usuarios/{id}/rol", usuarioHandler.UpdateUsuarioRol).Methods("PATCH")
	adminRouter.HandleFunc("/usuarios/{id}/status", usuarioHandler.UpdateUsuarioStatus).Methods("PATCH")
	adminRouter.HandleFunc("/usuarios/{id}/deactivate", usuarioHandler.DeactivateUsuario).Methods("POST")
	adminRouter.HandleFunc("/usuarios/stats", usuarioHandler.GetUsuarioStats).Methods("GET")
	adminRouter.HandleFunc("/atracciones/{id}", atraccionHandler.DeleteAtraccion).Methods("DELETE")
	adminRouter.HandleFunc("/atracciones/stats", atraccionHandler.GetStats).Methods("GET")
	adminRouter.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.DeleteAgencia).Methods("DELETE")
	adminRouter.HandleFunc("/agencias/{id:[0-9]+}/status", agenciaHandler.UpdateAgenciaStatus).Methods("PATCH")
	adminRouter.HandleFunc("/agencias/stats", agenciaHandler.GetStats).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok", "message": "ANDARIA - Sistema de Gestion Turistica", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
	}).Methods("GET")

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)
	log.Printf("Server starting on http://%s\n", addr)
	log.Printf("Health check: http://%s/health\n", addr)
	log.Printf("API v1: http://%s/api/v1\n", addr)

	log.Printf("\n ===== ENDPOINTS PÚBLICOS (sin login) =====\n")
	log.Printf("   GET    http://%s/api/v1/public/paquetes (Listar paquetes - CACHEADO)\n", addr)
	log.Printf("   GET    http://%s/api/v1/public/paquetes/{id} (Ver paquete - CACHEADO)\n", addr)
	log.Printf("   GET    http://%s/api/v1/public/agencias (Listar agencias - CACHEADO)\n", addr)
	log.Printf("   GET    http://%s/api/v1/public/atracciones (Listar atracciones - CACHEADO)\n", addr)
	log.Printf("   GET    http://%s/api/v1/data/departamentos (Datos auxiliares)\n", addr)
	log.Printf("   GET    http://%s/api/v1/data/categorias (Categorías)\n", addr)
	log.Printf("\n ===== AUTENTICACIÓN =====\n")
	log.Printf("   POST   http://%s/api/v1/auth/register (Registro)\n", addr)
	log.Printf("   POST   http://%s/api/v1/auth/login (Login)\n", addr)
	log.Printf("\n ===== ENDPOINTS PROTEGIDOS (requieren login) =====\n")
	log.Printf("   POST   http://%s/api/v1/compras (Crear compra - Turista)\n", addr)
	log.Printf("   POST   http://%s/api/v1/pagos (Crear pago - Turista)\n", addr)
	log.Printf("   GET    http://%s/api/v1/profile (Mi perfil)\n", addr)

	handler := corsHandler.Handler(router)
	log.Fatal(http.ListenAndServe(addr, handler))
}

func validateEnv() {
	otpSecret := os.Getenv("OTP_SECRET")
	if otpSecret == "" {
		log.Fatal("ERROR CRITICO: OTP_SECRET no esta configurado en .env")
	}

	forbiddenSecrets := []string{
		"change-this",
		"changethis",
		"change_this",
		"your-secret",
		"secret-key",
		"my-secret",
	}

	lowerSecret := strings.ToLower(otpSecret)
	for _, forbidden := range forbiddenSecrets {
		if strings.Contains(lowerSecret, forbidden) {
			log.Fatal("ERROR CRITICO: OTP_SECRET parece ser un placeholder. Usa: openssl rand -hex 32")
		}
	}

	if len(otpSecret) < 32 {
		log.Fatal("ERROR CRITICO: OTP_SECRET debe tener al menos 32 caracteres")
	}

	log.Println("OK. OTP_SECRET validado correctamente")

	appEnv := strings.ToLower(strings.TrimSpace(os.Getenv("APP_ENV")))
	if appEnv == "production" {
		jwtSecret := os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("ERROR CRITICO: JWT_SECRET no esta configurado en .env")
		}

		forbiddenJWT := []string{
			"secret",
			"change-this",
			"changethis",
			"change_this",
			"your-secret",
			"secret-key",
			"my-secret",
		}

		lowerJWT := strings.ToLower(jwtSecret)
		for _, forbidden := range forbiddenJWT {
			if strings.Contains(lowerJWT, forbidden) {
				log.Fatal("ERROR CRITICO: JWT_SECRET parece ser un placeholder. Usa: openssl rand -hex 32")
			}
		}

		if len(jwtSecret) < 32 {
			log.Fatal("ERROR CRITICO: JWT_SECRET debe tener al menos 32 caracteres")
		}

		log.Println("OK. JWT_SECRET validado correctamente")
	}
}
