package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/habibmrizki/finalphase3/internals/repositories"
)

type FollowHandler struct {
	followRepo *repositories.FollowRepository
}

func NewFollowHandler(followRepo *repositories.FollowRepository) *FollowHandler {
	return &FollowHandler{followRepo: followRepo}
}

// Follow godoc
// @Summary Follow a user
// @Description Allows the authenticated user to follow another user by ID
// @Tags Follow
// @Accept json
// @Produce json
// @Param id path int true "ID of the user to follow"
// @Success 200 {object} gin.H "Successfully followed user"
// @Failure 400 {object} gin.H "Invalid user ID or cannot follow yourself"
// @Failure 401 {object} gin.H "Unauthorized"
// @Failure 500 {object} gin.H "Failed to follow user"
// @Security ApiKeyAuth
// @Router /follow/{id} [post]

func (h *FollowHandler) Follow(ctx *gin.Context) {
	followerID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	}

	followingIDStr := ctx.Param("id")
	followingID, err := strconv.Atoi(followingIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}

	if followerID.(int) == followingID {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Cannot follow yourself"})
		return
	}

	err = h.followRepo.Follow(ctx, followerID.(int), followingID)
	if err != nil {
		log.Println("Follow error:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to follow user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":      "Successfully followed user",
		"follower_id":  followerID,
		"following_id": followingID,
	})
}
