package ambulance_ufe

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xtruhlar/dt26-ambulance-webapi/internal/db_service"
)

type implAmbulancesAPI struct{}

func NewAmbulancesApi() *implAmbulancesAPI {
	return &implAmbulancesAPI{}
}

func (o *implAmbulancesAPI) CreateAmbulance(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service not found"})
		return
	}
	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service context is not of required type"})
		return
	}

	var ambulance Ambulance
	if err := c.BindJSON(&ambulance); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "Bad Request", "message": "Invalid request body", "error": err.Error()})
		return
	}
	if ambulance.Id == "" {
		ambulance.Id = uuid.NewString()
	}
	if ambulance.ConsultationEntries == nil {
		ambulance.ConsultationEntries = []ConsultationEntryFull{}
	}

	err := db.CreateDocument(c.Request.Context(), ambulance.Id, &ambulance)
	switch err {
	case nil:
		c.JSON(http.StatusCreated, ambulance)
	case db_service.ErrConflict:
		c.JSON(http.StatusConflict, gin.H{"status": "Conflict", "message": "Ambulance already exists", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to create ambulance", "error": err.Error()})
	}
}

func (o *implAmbulancesAPI) DeleteAmbulance(c *gin.Context) {
	value, exists := c.Get("db_service")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service not found"})
		return
	}
	db, ok := value.(db_service.DbService[Ambulance])
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "Internal Server Error", "message": "db_service context is not of required type"})
		return
	}

	ambulanceId := c.Param("ambulanceId")
	err := db.DeleteDocument(c.Request.Context(), ambulanceId)
	switch err {
	case nil:
		c.AbortWithStatus(http.StatusNoContent)
	case db_service.ErrNotFound:
		c.JSON(http.StatusNotFound, gin.H{"status": "Not Found", "message": "Ambulance not found", "error": err.Error()})
	default:
		c.JSON(http.StatusBadGateway, gin.H{"status": "Bad Gateway", "message": "Failed to delete ambulance", "error": err.Error()})
	}
}
