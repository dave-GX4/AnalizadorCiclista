package semantic

import (
	"compilerciclista/src/parser"
	"fmt"
	"strings"
)

type ParticipantData map[string]interface{}

func Analyze(data parser.ParticipantData) []string {
	var errors []string

	// 1. Validar campos requeridos
	requiredFields := []string{"nombre", "apellido_paterno", "email", "sexo", "categoria"}
	for _, field := range requiredFields {
		if _, ok := data[field]; !ok {
			errors = append(errors, fmt.Sprintf("Error semántico: el campo requerido '%s' no fue encontrado.", field))
		}
	}
	
	if len(errors) > 0 {
		return errors
	}

	// 2. Validar formato de email (debe ser de Gmail)
	email, ok := data["email"].(string)
	if !ok {
		errors = append(errors, "Error semántico: el campo 'email' debe ser una cadena de texto.")
	} else if !strings.HasSuffix(email, "@gmail.com") {
		errors = append(errors, "Error semántico: el correo electrónico debe ser de Gmail.")
	}

	// 3. Validar el valor de 'sexo'
	sexo, ok := data["sexo"].(string)
	if !ok {
		errors = append(errors, "Error semántico: el campo 'sexo' debe ser una cadena de texto.")
	} else if sexo != "M" && sexo != "F" {
		errors = append(errors, "Error semántico: el valor de 'sexo' debe ser 'M' o 'F'.")
	}

	// 4. Validar la categoría
	validCategories := map[string]bool{"Elite": true, "Aficionado": true, "Juvenil": true}
	categoria, ok := data["categoria"].(string)
	if !ok {
		errors = append(errors, "Error semántico: el campo 'categoria' debe ser una cadena de texto.")
	} else if !validCategories[categoria] {
		errors = append(errors, fmt.Sprintf("Error semántico: la categoría '%s' no es válida.", categoria))
	}
	
	// 5. Validar consistencia de pago
    pago, pagoExists := data["pago_realizado"].(bool)
    comprobante, comprobanteExists := data["comprobante_pago_path"].(string)
    
    if pagoExists && pago && (!comprobanteExists || comprobante == "") {
        errors = append(errors, "Error semántico: si 'pago_realizado' es true, 'comprobante_pago_path' no puede estar vacío.")
    }

	return errors
}