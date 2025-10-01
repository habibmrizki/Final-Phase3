package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/habibmrizki/finalphase3/pkg"
	"github.com/jackc/pgx/v5"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepository
	hashCfg  *pkg.HashConfig
}

func NewAuthHandler(authRepo *repositories.AuthRepository, hashCfg *pkg.HashConfig) *AuthHandler {
	return &AuthHandler{
		authRepo: authRepo,
		hashCfg:  hashCfg,
	}
}

func respondWithError(ctx *gin.Context, status int, message string) {
	ctx.JSON(status, models.Response{
		Message: message,
		Status:  http.StatusText(status),
	})
}

// Register godoc
// @Summary Register user
// @Description Buat akun baru (name, email, password)
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.RegisterRequest true "Register Request"
// @Success 201 {object} models.RegisterResponse
// @Failure 400 {object} models.Response "Invalid request body or validation failed"
// @Failure 409 {object} models.Response "User already exists"
// @Failure 500 {object} models.Response "Database or hashing error"
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterRequest
	if err := ctx.ShouldBind(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request body or validation failed")
		return
	}

	_, err := h.authRepo.FindUserByEmail(ctx, req.Email)
	if err == nil {
		respondWithError(ctx, http.StatusConflict, "User with this email already exists")
		return
	}
	if err != pgx.ErrNoRows {
		log.Println("DB error:", err)
		respondWithError(ctx, http.StatusInternalServerError, "Database error while checking user")
		return
	}

	h.hashCfg.UseRecommended()
	hashedPassword, err := h.hashCfg.GenHash(req.Password)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to hash password")
		return
	}

	newUser, err := h.authRepo.CreateUser(ctx, req, hashedPassword)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to create user")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"id":      newUser.ID,
		"email":   newUser.Email,
		"name":    newUser.Name,
	})
}

// Login godoc
// @Summary Login user
// @Description Login dengan email & password, hasilkan JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body models.LoginRequest true "Login Request"
// @Success 200 {object} models.LoginResponse
// @Failure 400 {object} models.Response "Invalid request body"
// @Failure 401 {object} models.Response "Invalid credentials"
// @Failure 500 {object} models.Response "Server error"
// @Router /auth/login [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest
	if err := ctx.ShouldBind(&req); err != nil {
		respondWithError(ctx, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.authRepo.FindUserByEmail(ctx, req.Email)
	if err == pgx.ErrNoRows {
		respondWithError(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	} else if err != nil {
		log.Println("DB error:", err)
		respondWithError(ctx, http.StatusInternalServerError, "Database error")
		return
	}

	h.hashCfg.UseRecommended()
	isMatch, err := h.hashCfg.CompareHashAndPassword(req.Password, user.Password)
	if err != nil || !isMatch {
		respondWithError(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	claims := pkg.NewJWTClaims(user.ID, "user")
	tokenString, err := claims.GenToken()
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login berhasil",
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"name":  user.Name,
		},
		"token": tokenString,
	})

}
