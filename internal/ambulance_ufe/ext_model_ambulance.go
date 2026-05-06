package ambulance_ufe

// Ambulance is the top-level MongoDB document. Protocol is embedded inside each entry.
type Ambulance struct {
	Id                  string                  `json:"id" bson:"id"`
	ConsultationEntries []ConsultationEntryFull `json:"consultationEntries" bson:"consultationEntries"`
}

// ConsultationEntryFull extends ConsultationEntry with an embedded protocol for storage.
type ConsultationEntryFull struct {
	ConsultationEntry `bson:",inline"`
	Protocol          *CommunicationProtocol `json:"protocol,omitempty" bson:"protocol,omitempty"`
}
