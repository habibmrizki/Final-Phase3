package models

type FollowRequest struct {
	FollowingID int `json:"following_id" form:"following_id" binding:"required"`
}

type FollowResponse struct {
	Message     string `json:"message"`
	FollowerID  int    `json:"follower_id"`
	FollowingID int    `json:"following_id"`
}
