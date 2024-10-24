// domains/response.go
package domains

// BaseResponse is the general structure for all API responses
type BaseResponse struct {
    Code      string      `json:"code"`                 // HTTP response code
    Message   string      `json:"message"`              // Response message
    Data      interface{} `json:"data,omitempty"`       // Data payload (optional)
    Error     string      `json:"error,omitempty"`      // Error details (optional)
    Parameter string      `json:"parameter,omitempty"`  // Related parameter (optional)
}

// FormatError sets Error to "null" if it's an empty string
func (r *BaseResponse) FormatError() {
    if r.Error == "" {
        r.Error = "null"
    }
}

// TokenResponse represents a response with a token
type TokenResponse struct {
    Token string `json:"token"`  // JWT token string
}

// UserResponse represents the user details in the response
type UserResponse struct {
    UserID   string `json:"user_id"`      // Unique user ID
    Username string `json:"username"`     // User's username
    Email    string `json:"email"`        // User's email
    Password string `json:"-"`            // Password is omitted in the response
}

// DeleteResponse represents the response after a user is deleted
type DeleteResponse struct {
    UserID string `json:"user_id"`  // ID of the deleted user
}

// RegisterResponse represents the response after user registration
type RegisterResponse struct {
    Username  string `json:"username"`   // Registered username
    Email     string `json:"email"`      // Registered email
    Password1 string `json:"-"`          // Password is hidden in the response
    Password2 string `json:"-"`          // Password is hidden in the response
}

// ErrorResponse is used to format error messages with extra details
type ErrorResponse struct {
    Code      string            `json:"code"`                 // HTTP response code
    Message   string            `json:"message"`              // Error message
    Errors    map[string]string `json:"errors,omitempty"`     // Map of field errors (optional)
    Parameter string            `json:"parameter,omitempty"`  // Related parameter (optional)
}
