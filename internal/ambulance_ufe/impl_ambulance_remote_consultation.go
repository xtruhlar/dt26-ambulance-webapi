package ambulance_ufe

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type implAmbulanceRemoteConsultationAPI struct {
}

func NewAmbulanceRemoteConsultationApi() AmbulanceRemoteConsultationAPI {
	return &implAmbulanceRemoteConsultationAPI{}
}

func (o implAmbulanceRemoteConsultationAPI) CreateConsultationEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) DeleteConsultationEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) GetConsultationEntries(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) GetConsultationEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) GetConsultationProtocol(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) UpdateConsultationEntry(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}

func (o implAmbulanceRemoteConsultationAPI) UpdateConsultationProtocol(c *gin.Context) {
	c.AbortWithStatus(http.StatusNotImplemented)
}
