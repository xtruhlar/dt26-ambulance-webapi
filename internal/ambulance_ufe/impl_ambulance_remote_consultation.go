package ambulance_ufe

import (
	"net/http"
	"slices"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type implAmbulanceRemoteConsultationAPI struct{}

func NewAmbulanceRemoteConsultationApi() AmbulanceRemoteConsultationAPI {
	return &implAmbulanceRemoteConsultationAPI{}
}

func (o *implAmbulanceRemoteConsultationAPI) GetConsultationEntries(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entries := make([]ConsultationEntry, len(ambulance.ConsultationEntries))
		for i, e := range ambulance.ConsultationEntries {
			entries[i] = e.ConsultationEntry
		}
		return nil, entries, http.StatusOK
	})
}

func (o *implAmbulanceRemoteConsultationAPI) GetConsultationEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")
		idx := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entryId
		})
		if idx < 0 {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Entry not found"}, http.StatusNotFound
		}
		return nil, ambulance.ConsultationEntries[idx].ConsultationEntry, http.StatusOK
	})
}

func (o *implAmbulanceRemoteConsultationAPI) CreateConsultationEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		var entry ConsultationEntry
		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		if entry.PatientId == "" {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "patientId is required"}, http.StatusBadRequest
		}
		if entry.Id == "" || entry.Id == "@new" {
			entry.Id = uuid.NewString()
		}
		conflict := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entry.Id || e.PatientId == entry.PatientId
		})
		if conflict >= 0 {
			return nil, gin.H{"status": http.StatusConflict, "message": "Entry already exists"}, http.StatusConflict
		}
		ambulance.ConsultationEntries = append(ambulance.ConsultationEntries, ConsultationEntryFull{ConsultationEntry: entry})
		return ambulance, entry, http.StatusOK
	})
}

func (o *implAmbulanceRemoteConsultationAPI) UpdateConsultationEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")
		var entry ConsultationEntry
		if err := c.ShouldBindJSON(&entry); err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entryId
		})
		if idx < 0 {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Entry not found"}, http.StatusNotFound
		}
		if entry.PatientName != "" {
			ambulance.ConsultationEntries[idx].PatientName = entry.PatientName
		}
		if entry.Condition != "" {
			ambulance.ConsultationEntries[idx].Condition = entry.Condition
		}
		if entry.Status != "" {
			ambulance.ConsultationEntries[idx].Status = entry.Status
		}
		return ambulance, ambulance.ConsultationEntries[idx].ConsultationEntry, http.StatusOK
	})
}

func (o *implAmbulanceRemoteConsultationAPI) DeleteConsultationEntry(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")
		idx := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entryId
		})
		if idx < 0 {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Entry not found"}, http.StatusNotFound
		}
		ambulance.ConsultationEntries = append(ambulance.ConsultationEntries[:idx], ambulance.ConsultationEntries[idx+1:]...)
		return ambulance, nil, http.StatusNoContent
	})
}

func (o *implAmbulanceRemoteConsultationAPI) GetConsultationProtocol(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")
		idx := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entryId
		})
		if idx < 0 {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Entry not found"}, http.StatusNotFound
		}
		if ambulance.ConsultationEntries[idx].Protocol == nil {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Protocol not found"}, http.StatusNotFound
		}
		return nil, ambulance.ConsultationEntries[idx].Protocol, http.StatusOK
	})
}

func (o *implAmbulanceRemoteConsultationAPI) UpdateConsultationProtocol(c *gin.Context) {
	updateAmbulanceFunc(c, func(c *gin.Context, ambulance *Ambulance) (*Ambulance, interface{}, int) {
		entryId := c.Param("entryId")
		var protocol CommunicationProtocol
		if err := c.ShouldBindJSON(&protocol); err != nil {
			return nil, gin.H{"status": http.StatusBadRequest, "message": "Invalid request body", "error": err.Error()}, http.StatusBadRequest
		}
		idx := slices.IndexFunc(ambulance.ConsultationEntries, func(e ConsultationEntryFull) bool {
			return e.Id == entryId
		})
		if idx < 0 {
			return nil, gin.H{"status": http.StatusNotFound, "message": "Entry not found"}, http.StatusNotFound
		}
		if ambulance.ConsultationEntries[idx].Protocol == nil {
			protocol.Id = uuid.NewString()
			protocol.EntryId = entryId
			ambulance.ConsultationEntries[idx].Protocol = &protocol
		} else {
			if protocol.Content != "" {
				ambulance.ConsultationEntries[idx].Protocol.Content = protocol.Content
			}
			if protocol.Status != "" {
				ambulance.ConsultationEntries[idx].Protocol.Status = protocol.Status
			}
		}
		ambulance.ConsultationEntries[idx].Protocol.UpdatedAt = time.Now()
		return ambulance, ambulance.ConsultationEntries[idx].Protocol, http.StatusOK
	})
}
