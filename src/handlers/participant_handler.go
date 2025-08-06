package handlers

import (
	"compilerciclista/src/database"
	"compilerciclista/src/lexer"
	"compilerciclista/src/models"
	"compilerciclista/src/parser"
	"compilerciclista/src/semantic"
	"compilerciclista/src/services"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// =================================================================================
// REEMPLAZA TU FUNCIÓN ENTERA CON ESTA
// =================================================================================
func RegisterParticipantHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No se pudo leer el cuerpo de la solicitud")
		return
	}
	input := string(body)

	l := lexer.New(input)
	p := parser.New(l)
	participantData, parsingErrors := p.ParseProgram()

	if len(parsingErrors) > 0 {
		respondWithError(w, http.StatusBadRequest, parsingErrors)
		return
	}

	semanticErrors := semantic.Analyze(participantData)
	if len(semanticErrors) > 0 {
		respondWithError(w, http.StatusBadRequest, semanticErrors)
		return
	}

	participantModel, err := populateModel(participantData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	
	// Simulación de subida de archivos (opcional, pero lo dejamos por consistencia)
	if participantModel.InePath != "" {
		driveInePath, _ := services.UploadFile(participantModel.InePath)
		participantModel.InePath = driveInePath
	}
	if participantModel.ComprobantePagoPath != "" {
		driveComprobantePath, _ := services.UploadFile(participantModel.ComprobantePagoPath)
		participantModel.ComprobantePagoPath = driveComprobantePath
	}

	// Guardar en la Base de Datos
	id, err := database.CreateParticipant(participantModel)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			respondWithError(w, http.StatusConflict, fmt.Sprintf("El email '%s' ya está registrado.", participantModel.Email))
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al guardar el participante en la base de datos")
		return
	}
	participantModel.ID = id
	
	// <-- ¡ESTE ES EL CAMBIO MÁS IMPORTANTE! -->
	// De aquí en adelante, el código es nuevo.
	
	// Crear un mapa para la respuesta final, en lugar de enviar solo el modelo.
	responsePayload := map[string]interface{}{
		"message":     "Participante registrado exitosamente.",
		"participant": participantModel,
	}

	// Generar Token
	token, err := services.GenerateToken(participantModel)
	if err != nil {
		log.Printf("ADVERTENCIA: No se pudo generar el token para el participante %d: %v", id, err)
		responsePayload["token_warning"] = "No se pudo generar el token."
	} else {
		responsePayload["access_token"] = token
	}

	// Enviar Correo de Confirmación
	if err_email := services.SendConfirmationEmail(participantModel.Email, participantModel.Nombre, token); err_email != nil {
		log.Printf("ADVERTENCIA: El correo para el participante %d no se pudo enviar: %v", id, err_email)
		responsePayload["email_warning"] = "El servicio de correo falló: " + err_email.Error()
	}

	// Responder con el mapa completo que contiene toda la información.
	respondWithJSON(w, http.StatusCreated, responsePayload)
}


// El resto del archivo (populateModel, respondWithError, respondWithJSON)
// puede quedar exactamente como está.
// ... (resto de las funciones sin cambios) ...

func populateModel(data parser.ParticipantData) (models.Participant, error) {
	var p models.Participant
	var ok bool

	if p.Nombre, ok = data["nombre"].(string); !ok {
		return models.Participant{}, fmt.Errorf("campo 'nombre' inválido o ausente")
	}
	if p.ApellidoPaterno, ok = data["apellido_paterno"].(string); !ok {
		return models.Participant{}, fmt.Errorf("campo 'apellido_paterno' inválido o ausente")
	}
	if p.Email, ok = data["email"].(string); !ok {
		return models.Participant{}, fmt.Errorf("campo 'email' inválido o ausente")
	}
    if p.Sexo, ok = data["sexo"].(string); !ok {
		return models.Participant{}, fmt.Errorf("campo 'sexo' inválido o ausente")
	}
    if p.Categoria, ok = data["categoria"].(string); !ok {
		return models.Participant{}, fmt.Errorf("campo 'categoria' inválido o ausente")
	}

    // Campos opcionales
	p.ApellidoMaterno, _ = data["apellido_materno"].(string)
	p.PagoRealizado, _ = data["pago_realizado"].(bool)
	p.InePath, _ = data["ine_path"].(string)
	p.ComprobantePagoPath, _ = data["comprobante_pago_path"].(string)

	return p, nil
}

func respondWithError(w http.ResponseWriter, code int, message interface{}) {
	respondWithJSON(w, code, map[string]interface{}{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}