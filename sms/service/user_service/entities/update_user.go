package entities

type UpdateUserRequest struct {
	Username string `json:"username"`
	Role     string `json:"role"`
}

type UpdateUserResponse struct {
	Message string `json:"message"`
}

type UpdateUserBadRequestResponse struct {
	Error string `json:"error"`
}

type UpdateUserNotFound struct {
	Error string `json:"error"`
}

type UpdateUserInternalServerErrorResponse struct {
	Error string `json:"error"`
}
