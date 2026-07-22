package router

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"universev2-backend/internal/config"
	"universev2-backend/internal/handler"
	"universev2-backend/internal/middleware"
	"universev2-backend/internal/repository"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB) {
	// Initialize GORM Repositories
	empRepo := repository.NewEmployeeRepo(db)
	ftwRepo := repository.NewFTWRepo(db)
	rosterRepo := repository.NewRosterRepo(db)
	fleetRepo := repository.NewFleetRepo(db)
	userRepo := repository.NewUserRepo(db)
	roleRepo := repository.NewRoleRepo(db)
	notifRepo := repository.NewNotificationRepo(db)
	settingsRepo := repository.NewSettingsRepo(db)
	prestasiRepo := repository.NewPrestasiRepo(db)
	masterRepo := repository.NewMasterRepo(db)

	// Initialize Handlers with Repositories
	authH := handler.NewAuthHandler(cfg, userRepo, roleRepo)
	empH := handler.NewEmployeeHandler(empRepo)
	ftwH := handler.NewFitworkHandler(ftwRepo)
	rosterH := handler.NewRosterHandler(rosterRepo)
	attH := handler.NewAttendanceHandler()
	fleetH := handler.NewFleetHandler(fleetRepo)
	prestasiH := handler.NewPrestasiHandler(prestasiRepo)
	masterH := handler.NewMasterHandler(masterRepo)
	userH := handler.NewUserHandler(userRepo)
	roleH := handler.NewRoleHandler(roleRepo)
	notifH := handler.NewNotificationHandler(notifRepo)
	settingsH := handler.NewSettingsHandler(settingsRepo)
	fpH := handler.NewFingerprintHandler()
	profileH := handler.NewProfileHandler()
	weatherH := handler.NewWeatherHandler()

	// API v1 group
	api := app.Group("/api/v1")

	// Public Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", authH.Login)
	auth.Post("/register", authH.Register)
	auth.Post("/refresh", authH.RefreshToken)

	// Protected routes (JWT middleware)
	protected := api.Group("", middleware.AuthMiddleware(cfg))
	protected.Post("/auth/logout", authH.Logout)

	// Profile module
	profileGroup := protected.Group("/profile")
	profileGroup.Get("/", profileH.GetProfile)
	profileGroup.Put("/", profileH.UpdateProfile)
	profileGroup.Put("/password", profileH.UpdatePassword)

	// Weather Proxy
	protected.Get("/weather/current", weatherH.GetCurrentWeather)

	// Employees module
	empGroup := protected.Group("/employees", middleware.RequirePermission("employees", "view"))
	empGroup.Get("/", empH.GetEmployees)
	empGroup.Get("/:nik", empH.GetEmployeeByNIK)
	empGroup.Post("/", middleware.RequirePermission("employees", "manage"), empH.CreateEmployee)
	empGroup.Put("/:nik", middleware.RequirePermission("employees", "manage"), empH.UpdateEmployee)
	empGroup.Delete("/:nik", middleware.RequirePermission("employees", "manage"), empH.DeleteEmployee)

	// Fit To Work module
	ftwGroup := protected.Group("/ftw", middleware.RequirePermission("ftw", "view"))
	ftwGroup.Get("/today", ftwH.GetTodayLog)
	ftwGroup.Get("/history", ftwH.GetHistory)
	ftwGroup.Post("/submit", middleware.RequirePermission("ftw", "manage"), ftwH.SubmitLog)

	// Roster module
	rosterGroup := protected.Group("/rosters", middleware.RequirePermission("roster", "view"))
	rosterGroup.Get("/", rosterH.GetRosters)
	rosterGroup.Get("/:key/export", rosterH.ExportRoster)
	rosterGroup.Post("/upload", middleware.RequirePermission("roster", "manage"), rosterH.UploadRoster)

	rosterGroup.Get("/revisions", rosterH.GetRevisions)
	rosterGroup.Get("/revisions/codes", rosterH.GetRevisionCodes)
	rosterGroup.Post("/revisions/batch", middleware.RequirePermission("roster", "manage"), rosterH.SubmitBatchRevision)
	rosterGroup.Delete("/revisions/:id", middleware.RequirePermission("roster", "manage"), rosterH.DeleteRevision)

	rosterGroup.Put("/approvals/:id/approve", middleware.RequirePermission("roster", "manage"), rosterH.ApproveRevision)
	rosterGroup.Patch("/approvals/:id/note", middleware.RequirePermission("roster", "manage"), rosterH.ApproveRevisionWithNote)
	rosterGroup.Put("/approvals/:id/reject", middleware.RequirePermission("roster", "manage"), rosterH.RejectRevision)

	rosterGroup.Get("/attendance", rosterH.GetAttendance)

	// Standalone Attendance module
	attGroup := protected.Group("/attendance", middleware.RequirePermission("roster", "view"))
	attGroup.Get("/today", attH.GetAttendanceToday)
	attGroup.Get("/date", attH.GetAttendanceByDate)
	attGroup.Get("/range", attH.GetAttendanceRange)
	attGroup.Post("/checkin", middleware.RequirePermission("roster", "manage"), attH.RecordCheckIn)
	attGroup.Post("/checkout", middleware.RequirePermission("roster", "manage"), attH.RecordCheckOut)

	// Assets / Fleet module
	fleetGroup := protected.Group("/fleet", middleware.RequirePermission("asset", "view"))
	fleetGroup.Get("/settings", fleetH.GetFleetSettings)
	fleetGroup.Post("/settings", middleware.RequirePermission("asset", "manage"), fleetH.CreateFleetSetting)
	fleetGroup.Put("/settings/:id", middleware.RequirePermission("asset", "manage"), fleetH.UpdateFleetSetting)
	fleetGroup.Delete("/settings/:id", middleware.RequirePermission("asset", "manage"), fleetH.DeleteFleetSetting)

	fleetGroup.Get("/allocations", fleetH.GetAllocations)
	fleetGroup.Post("/allocations/auto", middleware.RequirePermission("asset", "manage"), fleetH.AutoAllocate)

	protected.Get("/units/status", middleware.RequirePermission("asset", "view"), fleetH.GetUnitStatuses)
	protected.Put("/units/:code/status", middleware.RequirePermission("asset", "manage"), fleetH.UpdateUnitStatus)
	protected.Post("/units/:code/status-report", middleware.RequirePermission("asset", "manage"), fleetH.ReportUnitBreakdown)
	protected.Get("/units/:code/history", middleware.RequirePermission("asset", "view"), fleetH.GetUnitHistory)

	// Unit DB inside Assets
	protected.Get("/units/db", middleware.RequirePermission("asset", "view"), fleetH.GetUnitDB)
	protected.Post("/units/db", middleware.RequirePermission("asset", "manage"), fleetH.CreateUnitDB)
	protected.Put("/units/db", middleware.RequirePermission("asset", "manage"), fleetH.UpdateUnitDB)
	protected.Delete("/units/db", middleware.RequirePermission("asset", "manage"), fleetH.DeleteUnitDB)

	// Prestasi module
	prestasiGroup := protected.Group("/prestasi", middleware.RequirePermission("prestasi", "view"))
	prestasiGroup.Get("/leaderboard", prestasiH.GetLeaderboard)

	// Master Data module
	masterGroup := protected.Group("/master", middleware.RequirePermission("master", "view"))
	masterGroup.Get("/:category", masterH.GetMasterByCategory)
	masterGroup.Post("/:category", middleware.RequirePermission("master", "manage"), masterH.CreateMasterEntry)
	masterGroup.Put("/:category/:id", middleware.RequirePermission("master", "manage"), masterH.UpdateMasterEntry)
	masterGroup.Delete("/:category/:id", middleware.RequirePermission("master", "manage"), masterH.DeleteMasterEntry)

	// User & Role Management module
	userGroup := protected.Group("/users", middleware.RequirePermission("users", "view"))
	userGroup.Get("/", userH.GetUsers)
	userGroup.Post("/", middleware.RequirePermission("users", "manage"), userH.CreateUser)
	userGroup.Put("/:id", middleware.RequirePermission("users", "manage"), userH.UpdateUser)
	userGroup.Delete("/:id", middleware.RequirePermission("users", "manage"), userH.DeleteUser)

	roleGroup := protected.Group("/roles", middleware.RequirePermission("users", "view"))
	roleGroup.Get("/", roleH.GetRoles)
	roleGroup.Post("/", middleware.RequirePermission("users", "manage"), roleH.CreateRole)
	roleGroup.Put("/:id", middleware.RequirePermission("users", "manage"), roleH.UpdateRole)
	roleGroup.Delete("/:id", middleware.RequirePermission("users", "manage"), roleH.DeleteRole)

	// Notifications
	notifGroup := protected.Group("/notifications")
	notifGroup.Get("/", notifH.GetNotifications)
	notifGroup.Put("/:id/read", notifH.MarkRead)
	notifGroup.Put("/read-all", notifH.MarkAllRead)

	// Settings & Display TV
	settingsGroup := protected.Group("/settings", middleware.RequirePermission("settings", "view"))
	settingsGroup.Get("/", settingsH.GetSettings)
	settingsGroup.Put("/", middleware.RequirePermission("settings", "manage"), settingsH.UpdateSettings)

	settingsGroup.Get("/audio", settingsH.GetAudioSchedules)
	settingsGroup.Post("/audio", middleware.RequirePermission("settings", "manage"), settingsH.CreateAudioSchedule)
	settingsGroup.Put("/audio/:id", middleware.RequirePermission("settings", "manage"), settingsH.UpdateAudioSchedule)
	settingsGroup.Delete("/audio/:id", middleware.RequirePermission("settings", "manage"), settingsH.DeleteAudioSchedule)

	settingsGroup.Get("/displays", settingsH.GetDisplays)
	settingsGroup.Post("/displays", middleware.RequirePermission("settings", "manage"), settingsH.CreateDisplay)
	settingsGroup.Put("/displays/:id", middleware.RequirePermission("settings", "manage"), settingsH.UpdateDisplay)
	settingsGroup.Delete("/displays/:id", middleware.RequirePermission("settings", "manage"), settingsH.DeleteDisplay)

	// Display Heartbeat
	api.Get("/displays/:id/heartbeat", settingsH.GetDisplayHeartbeat)

	// Fingerprint
	protected.Get("/fingerprint/devices", fpH.GetDeviceStatus)
}
