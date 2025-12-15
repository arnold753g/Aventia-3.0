package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"andaria-backend/internal/config"
	"andaria-backend/internal/database"
	"andaria-backend/internal/handlers"
	"andaria-backend/internal/middleware"
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

	// Crear router
	router := mux.NewRouter()

	// Servir archivos estaticos (fotografias)
	router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Setup CORS
	corsHandler := middleware.SetupCORS()

	// API v1 routes
	api := router.PathPrefix("/api/v1").Subrouter()

	// Auth routes (publicas)
	authHandler := handlers.NewAuthHandler()
	usuarioHandler := handlers.NewUsuarioHandler()
	agenciaHandler := handlers.NewAgenciaHandler()

	auth := api.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/register", authHandler.Register).Methods("POST")
	auth.HandleFunc("/login", authHandler.Login).Methods("POST")
	auth.HandleFunc("/refresh", authHandler.RefreshToken).Methods("POST")
	auth.HandleFunc("/logout", authHandler.Logout).Methods("POST")

	// Check duplicados (publico)
	api.HandleFunc("/usuarios/check", usuarioHandler.CheckUsuarioExiste).Methods("GET")

	// Protected routes (requieren autenticacion)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// ========== RUTAS DE USUARIOS ==========
	// Rutas autenticadas para todos los usuarios
	protected.HandleFunc("/usuarios/{id}", usuarioHandler.GetUsuario).Methods("GET")
	protected.HandleFunc("/usuarios/{id}", usuarioHandler.UpdateUsuario).Methods("PUT")

	// ========== RUTAS DE AGENCIAS TURISTICAS ==========
	// Rutas públicas
	api.HandleFunc("/agencias", agenciaHandler.GetAgencias).Methods("GET")
	api.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.GetAgencia).Methods("GET")
	api.HandleFunc("/agencias/data/departamentos", agenciaHandler.GetDepartamentos).Methods("GET")
	api.HandleFunc("/agencias/data/categorias", agenciaHandler.GetCategorias).Methods("GET")
	api.HandleFunc("/agencias/data/dias", agenciaHandler.GetDias).Methods("GET")
	api.HandleFunc("/agencias/data/encargados", agenciaHandler.GetEncargados).Methods("GET")

	// Rutas protegidas (requieren autenticación)
	protected.HandleFunc("/agencias/rapida", agenciaHandler.CreateAgenciaRapida).Methods("POST")
	protected.HandleFunc("/agencias/completa", agenciaHandler.CreateAgenciaCompleta).Methods("POST")
	protected.HandleFunc("/agencias/me", agenciaHandler.GetMiAgencia).Methods("GET")
	protected.HandleFunc("/agencias/{id:[0-9]+}", agenciaHandler.UpdateAgencia).Methods("PUT")
	protected.HandleFunc("/agencias/{id:[0-9]+}/fotos/upload", agenciaHandler.UploadAgenciaFoto).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/fotos/{foto_id:[0-9]+}", agenciaHandler.RemoveFotoWithFile).Methods("DELETE")
	protected.HandleFunc("/agencias/{id:[0-9]+}/especialidades", agenciaHandler.AddEspecialidad).Methods("POST")
	protected.HandleFunc("/agencias/{id:[0-9]+}/especialidades/{especialidad_id:[0-9]+}", agenciaHandler.RemoveEspecialidad).Methods("DELETE")

	// ========== RUTAS DE ATRACCIONES TURISTICAS ==========
	atraccionHandler := handlers.NewAtraccionHandler()

	// Rutas publicas
	api.HandleFunc("/atracciones", atraccionHandler.GetAtracciones).Methods("GET")
	api.HandleFunc("/atracciones/{id}", atraccionHandler.GetAtraccion).Methods("GET")

	// Rutas de datos auxiliares (publicas)
	api.HandleFunc("/categorias", atraccionHandler.GetCategorias).Methods("GET")
	api.HandleFunc("/subcategorias", atraccionHandler.GetSubcategorias).Methods("GET")
	api.HandleFunc("/departamentos", atraccionHandler.GetDepartamentos).Methods("GET")
	api.HandleFunc("/provincias", atraccionHandler.GetProvincias).Methods("GET")
	api.HandleFunc("/dias", atraccionHandler.GetDias).Methods("GET")
	api.HandleFunc("/meses", atraccionHandler.GetMeses).Methods("GET")

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

	// Profile
	protected.HandleFunc("/profile", authHandler.GetProfile).Methods("GET")

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

	log.Printf("\n ENDPOINTS DE USUARIOS:\n")
	log.Printf("   GET    http://%s/api/v1/admin/usuarios (Lista - Admin)\n", addr)
	log.Printf("   POST   http://%s/api/v1/admin/usuarios (Crear - Admin)\n", addr)
	log.Printf("   GET    http://%s/api/v1/usuarios/{id} (Ver detalle)\n", addr)
	log.Printf("   PUT    http://%s/api/v1/usuarios/{id} (Actualizar)\n", addr)
	log.Printf("   PATCH  http://%s/api/v1/admin/usuarios/{id}/rol (Cambiar rol - Admin)\n", addr)
	log.Printf("   PATCH  http://%s/api/v1/admin/usuarios/{id}/status (Cambiar status - Admin)\n", addr)
	log.Printf("   POST   http://%s/api/v1/admin/usuarios/{id}/deactivate (Desactivar - Admin)\n", addr)
	log.Printf("   GET    http://%s/api/v1/admin/usuarios/stats (Estadisticas - Admin)\n", addr)
	log.Printf("\n ENDPOINTS DE ATRACCIONES:\n")
	log.Printf("   GET    http://%s/api/v1/atracciones (Lista publica)\n", addr)
	log.Printf("   GET    http://%s/api/v1/atracciones/{id} (Ver detalle)\n", addr)
	log.Printf("   POST   http://%s/api/v1/atracciones (Crear)\n", addr)
	log.Printf("   PUT    http://%s/api/v1/atracciones/{id} (Actualizar)\n", addr)
	log.Printf("   DELETE http://%s/api/v1/admin/atracciones/{id} (Desactivar - Admin)\n", addr)
	log.Printf("   GET    http://%s/api/v1/admin/atracciones/stats (Estadisticas - Admin)\n", addr)
	log.Printf("   GET    http://%s/api/v1/categorias (Categorias)\n", addr)
	log.Printf("   GET    http://%s/api/v1/subcategorias (Subcategorias)\n", addr)
	log.Printf("   GET    http://%s/api/v1/departamentos (Departamentos)\n", addr)
	log.Printf("   GET    http://%s/api/v1/provincias (Provincias)\n", addr)

	handler := corsHandler.Handler(router)
	log.Fatal(http.ListenAndServe(addr, handler))
}
