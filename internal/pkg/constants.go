package pkg

import "time"

// File upload limits
const (
	MaxFileSize = 10 << 20 // 10MB
)

// JWT token expiry
const (
	JWTExpiry        = 24 * time.Hour
	JWTRefreshExpiry = 168 * time.Hour // 7 days
)

// Database connection pool settings
const (
	DBMaxConnLifetime = 30 * time.Minute
	DBMaxConnIdleTime = 5 * time.Minute
)

// Default pagination
const (
	DefaultPerPage = 10
)

// Default fleet settings
const (
	DefaultMaxUnits = 13
)

// Default attendance machine
const (
	DefaultMachine = "FP-01"
)
