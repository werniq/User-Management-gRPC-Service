package controllers

import (
	"github.com/Chained/auth-service/cmd/utils"
	"github.com/Chained/auth-service/internal/models"
	"github.com/gin-gonic/gin"
)

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

type ResetPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

func ForgetPassword(ctx *gin.Context) {
	var userCred *ForgotPasswordInput

	if err := ctx.ShouldBindJSON(&userCred); err != nil {
		ctx.JSON(200, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})

		return
	}

	msg := "You will receive a reset email if user with this email exists"

	//get user by email

	i := &models.DatabaseModel{DB: nil}

	user, err := i.GetUserByEmail(userCred.Email)
	if err != nil {
		ctx.JSON(200, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})

		return
	}

	resetToken := generateRandomString(20)

	// update user in database
	//query := bson.D{{Key: "email", Value: strings.ToLower(user.Email)}}
	//update := bson.D{{Key: "$set", Value: bson.D{{Key: "passwordResetToken", Value: passwordResetToken}, {Key: "passwordResetAt", Value: time.Now().Add(time.Minute * 15)}}}}

	// update user in mongoDB collection

	emailData := &utils.EmailData{
		URL:       "http://localhost:8080/resetPassword/" + resetToken,
		Firstname: user.GetFirstName(),
		Subject:   "Password reset",
	}

	err = utils.SendEmail(user.Email, emailData, "resetPassword.gohtml")
	if err != nil {
		ctx.JSON(200, gin.H{
			"status":  "failed",
			"message": err.Error(),
		})

		return
	}

	ctx.JSON(200, gin.H{
		"status":  "success",
		"message": msg,
	})
}
