package router

import (
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"

	"universev2-backend/internal/config"
	"universev2-backend/internal/handler"
	"universev2-backend/internal/middleware"
	"universev2-backend/internal/repository"
	"universev2-backend/internal/service"
	"universev2-backend/internal/worker"
)

func SetupRoutes(app *fiber.App, cfg *config.Config, db *gorm.DB, fpWorker *worker.FingerprintWorker) {
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
	attRepo := repository.NewAttendanceRepo(db)

	// Initialize RBAC Middleware
	rbac := middleware.NewRBACMiddleware(roleRepo)

	// Initialize Services
	masterSvc := service.NewMasterService(masterRepo)
	fpRepo := repository.NewFingerprintRepo(db)
	displaySvc := service.NewDisplayService(attRepo, ftwRepo, fleetRepo, empRepo, settingsRepo, fpRepo)
	dashSvc := service.NewDashboardService(attRepo, ftwRepo, fleetRepo, rosterRepo, notifRepo, empRepo)

	// Initialize Handlers
	authH := handler.NewAuthHandler(cfg, userRepo, roleRepo)
	empH := handler.NewEmployeeHandler(empRepo, cfg.UploadDir)
	ftwH := handler.NewFitworkHandler(ftwRepo)
	rosterH := handler.NewRosterHandler(rosterRepo, cfg.UploadDir)
	attH := handler.NewAttendanceHandler(attRepo)
	fleetH := handler.NewFleetHandler(fleetRepo)
	prestasiH := handler.NewPrestasiHandler(prestasiRepo, fleetRepo)
	masterH := handler.NewMasterHandler(masterSvc)
	userH := handler.NewUserHandler(userRepo, roleRepo)
	roleH := handler.NewRoleHandler(roleRepo)
	notifH := handler.NewNotificationHandler(notifRepo)
	settingsH := handler.NewSettingsHandler(settingsRepo)
	fpH := handler.NewFingerprintHandler(cfg, fpRepo, fpWorker)
	profileH := handler.NewProfileHandler(userRepo)
	weatherH := handler.NewWeatherHandler()

	// API group
	api := app.Group("/api")

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
	empGroup := protected.Group("/employees", rbac.RequirePermission("employees", "view"))
	empGroup.Get("/", empH.GetEmployees)
	empGroup.Get("/export", empH.ExportEmployees)
	empGroup.Get("/:nik", empH.GetEmployeeByNIK)
	empGroup.Post("/", rbac.RequirePermission("employees", "manage"), empH.CreateEmployee)
	empGroup.Post("/import", rbac.RequirePermission("employees", "manage"), empH.ImportEmployees)
	empGroup.Put("/:nik", rbac.RequirePermission("employees", "manage"), empH.UpdateEmployee)
	empGroup.Delete("/:nik", rbac.RequirePermission("employees", "manage"), empH.DeleteEmployee)

	// Employee competencies
	empGroup.Get("/:nik/competencies", rbac.RequirePermission("employees", "view"), empH.GetCompetencies)
	empGroup.Put("/:nik/competencies", rbac.RequirePermission("employees", "manage"), empH.UpdateCompetencies)

	// Employee photo
	empGroup.Post("/:nik/photo", rbac.RequirePermission("employees", "manage"), empH.UploadPhoto)

	// Fit To Work module
	ftwGroup := protected.Group("/ftw", rbac.RequirePermission("ftw", "view"))
	ftwGroup.Get("/today", ftwH.GetTodayLog)
	ftwGroup.Get("/history", ftwH.GetHistory)
	ftwGroup.Get("/export", ftwH.ExportFTW)
	ftwGroup.Post("/submit", rbac.RequirePermission("ftw", "manage"), ftwH.SubmitLog)

	// Roster module
	rosterGroup := protected.Group("/rosters", rbac.RequirePermission("roster", "view"))
	rosterGroup.Get("/", rosterH.GetRosters)
	rosterGroup.Get("/:key/export", rosterH.ExportRoster)
	rosterGroup.Get("/:key/detail", rosterH.GetRosterDetail)
	rosterGroup.Post("/upload", rbac.RequirePermission("roster", "manage"), rosterH.UploadRoster)

	rosterGroup.Get("/revisions", rosterH.GetRevisions)
	rosterGroup.Get("/revisions/codes", rosterH.GetRevisionCodes)
	rosterGroup.Post("/revisions/batch", rbac.RequirePermission("roster", "manage"), rosterH.SubmitBatchRevision)
	rosterGroup.Delete("/revisions/:id", rbac.RequirePermission("roster", "manage"), rosterH.DeleteRevision)

	rosterGroup.Put("/approvals/:id/approve", rbac.RequirePermission("roster", "manage"), rosterH.ApproveRevision)
	rosterGroup.Patch("/approvals/:id/note", rbac.RequirePermission("roster", "manage"), rosterH.ApproveRevisionWithNote)
	rosterGroup.Put("/approvals/:id/reject", rbac.RequirePermission("roster", "manage"), rosterH.RejectRevision)

	rosterGroup.Get("/attendance", rosterH.GetAttendance)

	// Standalone Attendance module
	attGroup := protected.Group("/attendance", rbac.RequirePermission("roster", "view"))
	attGroup.Get("/today", attH.GetAttendanceToday)
	attGroup.Get("/date", attH.GetAttendanceByDate)
	attGroup.Get("/range", attH.GetAttendanceRange)
	attGroup.Post("/checkin", rbac.RequirePermission("roster", "manage"), attH.RecordCheckIn)
	attGroup.Post("/checkout", rbac.RequirePermission("roster", "manage"), attH.RecordCheckOut)

	// Assets / Fleet module
	fleetGroup := protected.Group("/fleet", rbac.RequirePermission("asset", "view"))
	fleetGroup.Get("/settings", fleetH.GetFleetSettings)
	fleetGroup.Post("/settings", rbac.RequirePermission("asset", "manage"), fleetH.CreateFleetSetting)
	fleetGroup.Put("/settings/:id", rbac.RequirePermission("asset", "manage"), fleetH.UpdateFleetSetting)
	fleetGroup.Delete("/settings/:id", rbac.RequirePermission("asset", "manage"), fleetH.DeleteFleetSetting)

	fleetGroup.Get("/allocations", fleetH.GetAllocations)
	fleetGroup.Post("/allocations/auto", rbac.RequirePermission("asset", "manage"), fleetH.AutoAllocate)

	protected.Get("/units/status", rbac.RequirePermission("asset", "view"), fleetH.GetUnitStatuses)
	protected.Put("/units/:code/status", rbac.RequirePermission("asset", "manage"), fleetH.UpdateUnitStatus)
	protected.Post("/units/:code/status-report", rbac.RequirePermission("asset", "manage"), fleetH.ReportUnitBreakdown)
	protected.Get("/units/:code/history", rbac.RequirePermission("asset", "view"), fleetH.GetUnitHistory)

	// Unit DB inside Assets
	protected.Get("/units/db", rbac.RequirePermission("asset", "view"), fleetH.GetUnitDB)
	protected.Post("/units/db", rbac.RequirePermission("asset", "manage"), fleetH.CreateUnitDB)
	protected.Post("/units/db/import", rbac.RequirePermission("asset", "manage"), fleetH.ImportUnitDB)
	protected.Put("/units/db", rbac.RequirePermission("asset", "manage"), fleetH.UpdateUnitDB)
	protected.Delete("/units/db", rbac.RequirePermission("asset", "manage"), fleetH.DeleteUnitDB)

	// Prestasi module
	prestasiGroup := protected.Group("/prestasi", rbac.RequirePermission("prestasi", "view"))
	prestasiGroup.Get("/leaderboard", prestasiH.GetLeaderboard)
	prestasiGroup.Get("/:nik/history", prestasiH.GetOperatorHistory)
	prestasiGroup.Post("/recalculate", rbac.RequirePermission("prestasi", "manage"), prestasiH.Recalculate)

	// Master Data module
	masterGroup := protected.Group("/master", rbac.RequirePermission("master", "view"))
	masterGroup.Get("/:category", masterH.GetMasterByCategory)
	masterGroup.Get("/:category/export", masterH.ExportMaster)
	masterGroup.Post("/:category", rbac.RequirePermission("master", "manage"), masterH.CreateMasterEntry)
	masterGroup.Post("/:category/import", rbac.RequirePermission("master", "manage"), masterH.ImportMaster)
	masterGroup.Put("/:category/:id", rbac.RequirePermission("master", "manage"), masterH.UpdateMasterEntry)
	masterGroup.Delete("/:category/:id", rbac.RequirePermission("master", "manage"), masterH.DeleteMasterEntry)

	// User & Role Management module
	userGroup := protected.Group("/users", rbac.RequirePermission("users", "view"))
	userGroup.Get("/", userH.GetUsers)
	userGroup.Get("/export", userH.ExportUsers)
	userGroup.Post("/", rbac.RequirePermission("users", "manage"), userH.CreateUser)
	userGroup.Post("/import", rbac.RequirePermission("users", "manage"), userH.ImportUsers)
	userGroup.Patch("/:id/status", rbac.RequirePermission("users", "manage"), userH.ToggleUserStatus)
	userGroup.Put("/:id", rbac.RequirePermission("users", "manage"), userH.UpdateUser)
	userGroup.Delete("/:id", rbac.RequirePermission("users", "manage"), userH.DeleteUser)

	roleGroup := protected.Group("/roles", rbac.RequirePermission("users", "view"))
	roleGroup.Get("/", roleH.GetRoles)
	roleGroup.Get("/export", roleH.ExportRoles)
	roleGroup.Post("/", rbac.RequirePermission("users", "manage"), roleH.CreateRole)
	roleGroup.Put("/:id", rbac.RequirePermission("users", "manage"), roleH.UpdateRole)
	roleGroup.Delete("/:id", rbac.RequirePermission("users", "manage"), roleH.DeleteRole)

	// Notifications
	notifGroup := protected.Group("/notifications")
	notifGroup.Get("/", notifH.GetNotifications)
	notifGroup.Put("/:id/read", notifH.MarkRead)
	notifGroup.Put("/read-all", notifH.MarkAllRead)

	// Settings & Display TV
	settingsGroup := protected.Group("/settings", rbac.RequirePermission("settings", "view"))
	settingsGroup.Get("/", settingsH.GetSettings)
	settingsGroup.Put("/", rbac.RequirePermission("settings", "manage"), settingsH.UpdateSettings)

	settingsGroup.Get("/audio", settingsH.GetAudioSchedules)
	settingsGroup.Post("/audio", rbac.RequirePermission("settings", "manage"), settingsH.CreateAudioSchedule)
	settingsGroup.Put("/audio/:id", rbac.RequirePermission("settings", "manage"), settingsH.UpdateAudioSchedule)
	settingsGroup.Delete("/audio/:id", rbac.RequirePermission("settings", "manage"), settingsH.DeleteAudioSchedule)

	settingsGroup.Get("/displays", settingsH.GetDisplays)
	settingsGroup.Post("/displays", rbac.RequirePermission("settings", "manage"), settingsH.CreateDisplay)
	settingsGroup.Put("/displays/:id", rbac.RequirePermission("settings", "manage"), settingsH.UpdateDisplay)
	settingsGroup.Delete("/displays/:id", rbac.RequirePermission("settings", "manage"), settingsH.DeleteDisplay)

	// Display Heartbeat
	api.Get("/displays/:code/heartbeat", settingsH.GetDisplayHeartbeat)

	// Display TV endpoints
	displayH := handler.NewDisplayHandler(displaySvc)
	protected.Get("/display/attendance", rbac.RequirePermission("display", "view"), displayH.GetDisplayAttendance)
	protected.Get("/display/ftw", rbac.RequirePermission("display", "view"), displayH.GetDisplayFTW)
	protected.Get("/display/fleet", rbac.RequirePermission("display", "view"), displayH.GetDisplayFleet)
	protected.Get("/display/monitor", rbac.RequirePermission("display", "view"), displayH.GetDisplayMonitor)
	protected.Get("/display/fingerprint", rbac.RequirePermission("display", "view"), displayH.GetDisplayFingerprint)

	// Dashboard summary
	dashH := handler.NewDashboardHandler(dashSvc)
	protected.Get("/dashboard/summary", rbac.RequirePermission("dashboard", "view"), dashH.GetDashboardSummary)

	// Fingerprint Devices
	fpGroup := protected.Group("/fingerprint")
	fpGroup.Get("/devices", rbac.RequirePermission("settings", "view"), fpH.GetDeviceStatus)
	fpGroup.Post("/devices", rbac.RequirePermission("settings", "manage"), fpH.CreateDevice)
	fpGroup.Put("/devices/:id", rbac.RequirePermission("settings", "manage"), fpH.UpdateDevice)
	fpGroup.Delete("/devices/:id", rbac.RequirePermission("settings", "manage"), fpH.DeleteDevice)
	fpGroup.Post("/sync", rbac.RequirePermission("settings", "manage"), fpH.SyncNow)
}
