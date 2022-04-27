package constants

const (
	JWTKeyUserID     = "user_id"
	JWTKeyExpiration = "exp"

	// Format: "1s", "13h" (max - h, d - days - not allowed)
	// See "time.ParseDuration"
	ViperJWTTTLKey    = "service.jwt_ttl"
	ViperJWTSecretKey = "service.jwt_secret"

	ViperCSRFTTLKey    = "service.csrf_ttl"
	ViperCSRFSecretKey = "service.csrf_secret"
)
