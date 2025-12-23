package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"andaria-backend/internal/config"
	"andaria-backend/internal/database"
	"andaria-backend/internal/handlers"
	"andaria-backend/internal/middleware"
	"andaria-backend/internal/services"
	"andaria-backend/pkg/utils"

	"github.com/gorilla/mux"
)

func main() {
	// Cargar configuracion
	cfg := config.LoadConfig()

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
	compraHandler := handlers.NewCompraHandler()
	pagoHandler := handlers.NewPagoHandler()
	atraccionHandler := handlers.NewAtraccionHandler()

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

	// Agencias públicas
	publicAPI.HandleFunc("/agencias", agenciaHandler.GetAgencias).Methods("GET")
	publicAPI.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.GetAgencia).Methods("GET")

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

	// Check duplicados (público sin caché)
	api.HandleFunc("/usuarios/check", usuarioHandler.CheckUsuarioExiste).Methods("GET")

	// ========== RUTAS PROTEGIDAS (requieren autenticación) ==========
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Profile
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")

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
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas/{salida_id:[0-9]+}", agenciaHandler.UpdatePaqueteSalida).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/paquetes/{paquete_id:[0-9]+}/salidas/{salida_id:[0-9]+}/activar", agenciaHandler.ActivarSalida).Methods("POST")

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
