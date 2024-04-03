package routes

import (
	"fmt"
	"net/http"

	"github.com/asdutoit/go_backend_template/models"
	"github.com/gin-gonic/gin"
)

func getDeployments(ctx *gin.Context) {
	deployments, err := models.GetAllDeployments()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch deployments", "error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, deployments)
}

func GetDeploymentByQuery(ctx *gin.Context) {
	organization := ctx.Query("organization")
	product := ctx.Query("product")
	systemLayer := ctx.Query("systemLayer")
	environment := ctx.Query("environment")

	fmt.Println("organization", organization)

	deployment, err := models.GetDeploymentByQuery(organization, product, systemLayer, environment)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch deployment", "error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, deployment)
}
