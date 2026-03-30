package amplitude

// Event type constants for the base (auth) service.
const (
	// Auth
	EventUserSignedUp  = "user_signed_up"
	EventUserLoggedIn  = "user_logged_in"
	EventUserLoggedOut = "user_logged_out"

	// Password management
	EventPasswordResetCompleted = "password_reset_completed"

	// Email verification
	EventEmailVerified = "email_verified"

	// Contacts / Customers (form submissions)
	EventContactCreated = "contact_created"
)

// Property key constants — use these instead of raw strings to avoid typos.
const (
	PropEmail      = "email"
	PropAuthMethod = "auth_method" // "email_password" | "google" | "email_code"
)
