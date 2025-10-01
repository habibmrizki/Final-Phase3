// handlers/like_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/repositories"
)

type LikeHandler struct {
	repo *repositories.LikeRepository
}

func NewLikeHandler(repo *repositories.LikeRepository) *LikeHandler {
	return &LikeHandler{repo: repo}
}

// LikePost godoc
// @Summary Like a post
// @Description Allows the authenticated user to like a specific post by ID
// @Tags Likes
// @Accept json
// @Produce json
// @Param id path int true "ID of the post to like"
// @Success 200 {object} gin.H "Post liked successfully"
// @Failure 400 {object} gin.H "Invalid post ID"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Failed to like post"
// @Security ApiKeyAuth
// @Router /likes/{id} [post]
func (h *LikeHandler) LikePost(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("user_id")
	userID := userIDVal.(int)

	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	if err := h.repo.LikePost(ctx, userID, postID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to like post"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Post liked"})
}
