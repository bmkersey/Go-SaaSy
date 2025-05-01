package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/bmkersey/Go-SaaSy/internal/admin"
	"github.com/bmkersey/Go-SaaSy/internal/auth"
	"github.com/bmkersey/Go-SaaSy/internal/billing"
	"github.com/bmkersey/Go-SaaSy/internal/config"
	"github.com/bmkersey/Go-SaaSy/internal/db"
	"github.com/bmkersey/Go-SaaSy/internal/email"
	"github.com/bmkersey/Go-SaaSy/internal/orgs"
	stripeclient "github.com/bmkersey/Go-SaaSy/internal/stripe"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

func main() {
	cfg := config.LoadConfig()

	stripeclient.Init()

	conn, err := sql.Open("postgres", cfg.DB)
	if err != nil {
		log.Fatal("Could not establish connection to DB", err)
	}
	defer conn.Close()
	queries := db.New(conn)
	store := db.NewStore(queries)
	sender := email.NewEmailSender()

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir("./public/assets"))))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("SaaSy is running 🚀 (via Chi)"))
	})

	r.Post("/api/register", auth.RegisterHandler(store))
	r.Post("/api/login", auth.LoginHandler(store, cfg.JwtSecret))
	r.Post("/api/billing/webhook", billing.WebhookHandler(store, sender))
	r.Post("/api/forgot-password", auth.ForgotPasswordHandler(store, sender))
	r.Post("/api/reset-password", auth.ResetPasswordHandler(store))
	r.With(
		auth.AuthMiddleware(cfg.JwtSecret),
		orgs.OrgMiddleware(store),
		orgs.RequireOwner(store),
	).Get("/api/org/dashboard", orgs.OrgDashboardHandler(store))
	r.With(
		auth.AuthMiddleware(cfg.JwtSecret),
		auth.RequireAdmin(store),
	).Get("/api/admin/overview", admin.AdminOverviewHandler(store))

	r.Route("/api", func(r chi.Router) {

		r.Use(auth.AuthMiddleware(cfg.JwtSecret))

		r.Get("/me", auth.MeHandler(store))
		r.Post("/createorg", orgs.CreateOrganizationHandler(store))
		r.Route("/billing", func(r chi.Router) {
			r.Use(orgs.OrgMiddleware(store))
			r.Post("/checkout", billing.CreateCheckoutHandler(cfg))
		})

		r.Route("/premium", func(r chi.Router) {
			r.Use(orgs.OrgMiddleware(store))
			r.Use(orgs.RequirePaidPlan(store))

			// r.Get("/exclusive", premiumHandler)
			// r.Get("/analytics", analyticsHandler)
		})

		r.Route("/orgs", func(r chi.Router) {
			r.Use(orgs.OrgMiddleware(store))
			r.Get("/{id}/members", orgs.GetOrgMembers(store))
		})
	})

	log.Printf("Starting server on port %s...\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
