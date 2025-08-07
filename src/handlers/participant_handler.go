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

// RegisterParticipantHandler procesa la solicitud completa para registrar un nuevo participante.
func RegisterParticipantHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Leer y procesar el DSL de entrada
	body, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "No se pudo leer el cuerpo de la solicitud")
		return
	}
	input := string(body)

	// 2. Pipeline del Compilador (Léxico, Sintáctico, Semántico)
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

	// 3. Poblar el modelo de Go con los datos validados
	participantModel, err := populateModel(participantData)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// 4. Generar el CÓDIGO DE PARTICIPANTE único
	participantCode, err := services.GenerateParticipantCode(participantModel.Categoria)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "No se pudo generar el código de participante.")
		return
	}
	// Asignamos el código generado a nuestro modelo antes de guardarlo.
	participantModel.ParticipantCode = participantCode

	// 5. (Simulado) Subir archivos y actualizar rutas
	if participantModel.InePath != "" {
		driveInePath, _ := services.UploadFile(participantModel.InePath)
		participantModel.InePath = driveInePath
	}
	if participantModel.ComprobantePagoPath != "" {
		driveComprobantePath, _ := services.UploadFile(participantModel.ComprobantePagoPath)
		participantModel.ComprobantePagoPath = driveComprobantePath
	}

	// 6. Guardar el participante en la Base de Datos
	// El modelo ahora contiene el código de participante generado.
	id, err := database.CreateParticipant(participantModel)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate entry") {
			// Este error ahora puede ser por un email o un código de participante duplicado
			respondWithError(w, http.StatusConflict, fmt.Sprintf("Conflicto de datos: El email '%s' o el código generado ya existe.", participantModel.Email))
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error al guardar el participante en la base de datos")
		return
	}
	participantModel.ID = id // Asignamos el ID autoincremental de la DB

	// 7. Preparar la respuesta JSON final
	responsePayload := map[string]interface{}{
		"message":     "Participante registrado exitosamente.",
		"participant": participantModel, // El modelo completo con ID y Código de Participante
	}

	// 8. Generar el TOKEN JWT (para la seguridad de la API)
	jwtToken, err := services.GenerateToken(participantModel)
	if err != nil {
		log.Printf("ADVERTENCIA: No se pudo generar el token JWT: %v", err)
		responsePayload["token_warning"] = "No se pudo generar el token de acceso JWT."
	} else {
		responsePayload["access_token"] = jwtToken
	}

	// 9. Enviar el CORREO DE CONFIRMACIÓN (pasando el código de participante)
	if err_email := services.SendConfirmationEmail(participantModel.Email, participantModel.Nombre, participantModel.ParticipantCode); err_email != nil {
		log.Printf("ADVERTENCIA: El correo para el participante %d no se pudo enviar: %v", id, err_email)
		responsePayload["email_warning"] = "El servicio de correo falló: " + err_email.Error()
	}

	// 10. Enviar la respuesta final completa al cliente
	respondWithJSON(w, http.StatusCreated, responsePayload)
}

// populateModel convierte el mapa de datos del parser a un struct de Participant.
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

// respondWithError es una función helper para enviar respuestas de error en JSON.
func respondWithError(w http.ResponseWriter, code int, message interface{}) {
	respondWithJSON(w, code, map[string]interface{}{"error": message})
}

// respondWithJSON es una función helper para enviar respuestas exitosas en JSON.
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}