package routes

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tadeasf/eve-ran/src/db"
	"github.com/tadeasf/eve-ran/src/db/models"
)

// GetKillComments retrieves all comments for a specific kill
// @Summary Get kill comments
// @Description Fetch all comments for a specific killmail
// @Tags comments
// @Accept json
// @Produce json
// @Param killmail_id path int true "Killmail ID"
// @Security ApiKeyAuth
// @Success 200 {array} models.KillComment
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /kills/{killmail_id}/comments [get]
func GetKillComments(c *gin.Context) {
	killmailID, err := strconv.ParseInt(c.Param("killmail_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid killmail ID"})
		return
	}

	var comments []models.KillComment
	if err := db.DB.Where("killmail_id = ?", killmailID).Order("created_at DESC").Find(&comments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, comments)
}

// CreateKillComment creates a new comment on a kill
// @Summary Create kill comment
// @Description Add a new comment to a killmail
// @Tags comments
// @Accept json
// @Produce json
// @Param killmail_id path int true "Killmail ID"
// @Param comment body models.KillComment true "Comment data"
// @Security ApiKeyAuth
// @Success 201 {object} models.KillComment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /kills/{killmail_id}/comments [post]
func CreateKillComment(c *gin.Context) {
	killmailID, err := strconv.ParseInt(c.Param("killmail_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid killmail ID"})
		return
	}

	// Verify kill exists
	var kill models.Kill
	if err := db.DB.Where("killmail_id = ?", killmailID).First(&kill).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kill not found"})
		return
	}

	var comment models.KillComment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set killmail_id and kill_id
	comment.KillmailID = killmailID
	comment.KillID = kill.ID

	// Validate required fields
	if comment.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username is required"})
		return
	}
	if comment.Comment == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Comment text is required"})
		return
	}

	// Create comment
	if err := db.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, comment)
}

// UpdateKillComment updates an existing comment
// @Summary Update kill comment
// @Description Update an existing comment on a killmail
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Param comment body models.KillComment true "Updated comment data"
// @Security ApiKeyAuth
// @Success 200 {object} models.KillComment
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [put]
func UpdateKillComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Check if comment exists
	var existingComment models.KillComment
	if err := db.DB.First(&existingComment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	var updateData models.KillComment
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Only update the comment text
	if updateData.Comment != "" {
		existingComment.Comment = updateData.Comment
	}

	if err := db.DB.Save(&existingComment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, existingComment)
}

// DeleteKillComment deletes a comment
// @Summary Delete kill comment
// @Description Delete a comment from a killmail
// @Tags comments
// @Accept json
// @Produce json
// @Param id path int true "Comment ID"
// @Security ApiKeyAuth
// @Success 204 "No Content"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /comments/{id} [delete]
func DeleteKillComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment ID"})
		return
	}

	// Check if comment exists
	var comment models.KillComment
	if err := db.DB.First(&comment, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found"})
		return
	}

	if err := db.DB.Delete(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
