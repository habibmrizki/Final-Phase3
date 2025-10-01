// handlers/comment_handler.go
package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/models"
	"github.com/habibmrizki/finalphase3/internals/repositories"
)

type CommentHandler struct {
	repo *repositories.CommentRepository
}

func NewCommentHandler(repo *repositories.CommentRepository) *CommentHandler {
	return &CommentHandler{repo: repo}
}

// CreateComment godoc
// @Summary Create a new comment
// @Description Create a comment for a specific post
// @Tags Comments
// @Accept json
// @Produce json
// @Param id path int true "Post ID"
// @Param comment body struct{Content string `json:"content"`} true "Comment content"
// @Success 201 {object} models.Comment "Comment created successfully"
// @Failure 400 {object} gin.H "Invalid post ID or content required"
// @Failure 500 {object} gin.H "Failed to create comment"
// @Security ApiKeyAuth
// @Router /comments/{id} [post]
func (h *CommentHandler) CreateComment(ctx *gin.Context) {
	userIDVal, _ := ctx.Get("user_id")
	userID := userIDVal.(int)

	postID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid post ID"})
		return
	}

	var req struct {
		Content string `form:"content" json:"content" binding:"required"`
	}
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Content required"})
		return
	}

	comment := models.Comment{
		UserID:  userID,
		PostID:  postID,
		Content: req.Content,
	}

	newComment, err := h.repo.CreateComment(ctx, comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to create comment"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Comment created successfully",
		"data":    newComment,
	})
}
