package models

type Participant struct {
	ID                  int64  `json:"id"`
	ParticipantCode     string `json:"participant_code"`
	Nombre              string `json:"nombre"`
	ApellidoPaterno     string `json:"apellido_paterno"`
	ApellidoMaterno     string `json:"apellido_materno,omitempty"`
	Email               string `json:"email"`
	Sexo                string `json:"sexo"`
	Categoria           string `json:"categoria"`
	PagoRealizado       bool   `json:"pago_realizado"`
	InePath             string `json:"ine_path,omitempty"`
	ComprobantePagoPath string `json:"comprobante_pago_path,omitempty"`
}