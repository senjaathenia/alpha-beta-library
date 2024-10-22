package domains

// BaseResponse is the general structure for all API responses
type BaseResponse struct {
    Code      string      `json:"code"`
    Message   string      `json:"message"`
    Data      interface{} `json:"data,omitempty"`
    Error     interface{} `json:"error,omitempty"`
    Parameter string      `json:"parameter,omitempty"` // Parameter added
}

// RegisterResponse defines the structure for register response
type RegisterResponse struct {
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password1 string `json:"password_1"`
    Password2 string `json:"password_2"`
}

// UserResponse defines the structure for user-related responses
type UserResponse struct {
    UserID    string `json:"user_id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    Password  string `json:"password"`
}

// DeleteResponse defines the structure for delete response
type DeleteResponse struct {
    UserID string `json:"user_id"`
}
