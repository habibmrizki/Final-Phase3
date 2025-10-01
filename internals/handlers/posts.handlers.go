package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/habibmrizki/finalphase3/internals/repositories"
	"github.com/redis/go-redis/v9"
)

type PostHandler struct {
	postRepo *repositories.PostRepository
	rdb      *redis.Client
}

func NewPostHandler(postRepo *repositories.PostRepository, rdb *redis.Client) *PostHandler {
	return &PostHandler{
		postRepo: postRepo,
		rdb:      rdb,
	}
}

// CreatePost godoc
// @Summary      Membuat Postingan Baru
// @Description  Membuat postingan baru dengan teks dan/atau gambar. Harus login.
// @Tags         posts
// @Accept       multipart/form-data
// @Produce      json
// @Param        content  formData  string                false  "Konten teks postingan"
// @Param        image    formData  file                  false  "File gambar (jpg, png, dll)"
// @Success      201 {object} map[string]interface{}    "Postingan berhasil dibuat"
// @Failure      400 {object} map[string]interface{}    "Format data tidak valid / konten kosong"
// @Failure      401 {object} map[string]interface{}    "User tidak terautentikasi"
// @Failure      500 {object} map[string]interface{}    "Server error / gagal simpan file / DB"
// @Router       /posts/ [post]
// @Security     ApiKeyAuth
func (h *PostHandler) CreatePost(ctx *gin.Context) {
	var req models.CreatePostRequest

	if err := ctx.ShouldBind(&req); err != nil {
		log.Printf("[ERROR CreatePost][Handler][Binding]: %v\n", err)
		respondWithError(ctx, http.StatusBadRequest, "Format data tidak valid atau binding gagal")
		return
	}

	// Ambil UserID dari Context (ASUMSI: JWT Middleware sudah menyimpan ini)
	userIDValue, exists := ctx.Get("user_id")
	if !exists {
		respondWithError(ctx, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}
	userID, ok := userIDValue.(int)
	if !ok {
		respondWithError(ctx, http.StatusInternalServerError, "Gagal mendapatkan User ID")
		return
	}

	// Validasi Konten: Harus ada teks atau gambar
	if req.Content == "" && req.Image == nil {
		respondWithError(ctx, http.StatusBadRequest, "Postingan harus memiliki konten teks atau gambar")
		return
	}

	// Proses Upload Gambar
	var imagePath string = ""
	if req.Image != nil {
		imageFile := req.Image

		// Membuat nama file unik
		imageExt := filepath.Ext(imageFile.Filename)
		imageFilename := fmt.Sprintf("post_%d_%s%s", userID, strconv.FormatInt(time.Now().Unix(), 10), imageExt)

		// Lokasi penyimpanan file
		uploadDir := filepath.Join("public", "uploads", "posts")
		imageLocation := filepath.Join(uploadDir, imageFilename)

		// Membuat direktori jika belum ada
		if err := os.MkdirAll(uploadDir, 0755); err != nil {
			log.Printf("[ERROR CreatePost][Handler][MakeDir]: %v\n", err)
			respondWithError(ctx, http.StatusInternalServerError, "Gagal membuat direktori upload")
			return
		}

		// Simpan file
		if err := ctx.SaveUploadedFile(imageFile, imageLocation); err != nil {
			log.Printf("[ERROR CreatePost][Handler][Upload]: %v\n", err)
			respondWithError(ctx, http.StatusInternalServerError, "Gagal menyimpan file gambar")
			return
		}

		// Path yang akan disimpan ke DB (akses via web)
		imagePath = "/uploads/posts/" + imageFilename
	}

	// Insert ke Database
	newPostData := models.Post{
		UserID:  userID,
		Content: req.Content,
		Image:   imagePath,
	}

	postID, err := h.postRepo.CreatePost(ctx.Request.Context(), newPostData)
	if err != nil {
		log.Printf("[ERROR CreatePost][Handler][Repo]: %v\n", err)
		respondWithError(ctx, http.StatusInternalServerError, "Gagal membuat postingan di database")
		return
	}

	// Response Sukses
	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Postingan berhasil dibuat",
		"post_id": postID,
		"content": req.Content,
		"image":   imagePath,
	})
}

// GetFeed godoc
// @Summary Get user feed
// @Description Retrieves the feed (posts) for the authenticated user, with caching support
// @Tags Posts
// @Accept json
// @Produce json
// @Success 200 {array} models.Post "List of posts in feed"
// @Failure 401 {object} gin.H "User not authenticated"
// @Failure 500 {object} gin.H "Failed to retrieve feed"
// @Security ApiKeyAuth
// @Router /posts/feed [get]
func (h *PostHandler) GetFeed(ctx *gin.Context) {
	userIDVal, exists := ctx.Get("user_id")
	if !exists {
		respondWithError(ctx, http.StatusUnauthorized, "User tidak terautentikasi")
		return
	}
	userID, ok := userIDVal.(int)
	if !ok {
		respondWithError(ctx, http.StatusInternalServerError, "Gagal mendapatkan User ID")
		return
	}

	cacheKey := fmt.Sprintf("feed:%d", userID)
	cached, err := h.rdb.Get(context.Background(), cacheKey).Result()
	if err == nil {

		ctx.Data(http.StatusOK, "application/json", []byte(cached))
		return
	}

	posts, err := h.postRepo.GetFeed(ctx.Request.Context(), userID)
	if err != nil {
		respondWithError(ctx, http.StatusInternalServerError, "Gagal mengambil feed dari DB")
		return
	}

	jsonData, _ := json.Marshal(posts)
	h.rdb.Set(context.Background(), cacheKey, jsonData, 60*time.Second)

	ctx.JSON(http.StatusOK, posts)
}
