package controllers

import (
	"backbu/internal/data"
	"backbu/pkg/mailing"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func MailSender(c *gin.Context) {
	var mensaje mailing.Message

	err := json.NewDecoder(c.Request.Body).Decode(&mensaje)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, data.JsonError{Message: err.Error()})
		return
	}

	err = mailing.SendMail(mensaje)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, data.JsonError{Message: err.Error()})
		return
	}

	c.Status(http.StatusAccepted)
}
