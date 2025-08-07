package services

import (
	"compilerciclista/src/database"
	"fmt"
	"strings"
)

// GenerateParticipantCode crea un código único para un nuevo participante.
func GenerateParticipantCode(categoria string) (string, error) {
	// 1. Obtener el prefijo de la categoría (ej: "JUV", "ELI", "AFI")
	prefix := "GEN" // Prefijo genérico por si acaso
	if len(categoria) >= 3 {
		prefix = strings.ToUpper(categoria[0:3])
	}

	// 2. Contar los participantes existentes para obtener el siguiente número
	count, err := database.CountParticipants()
	if err != nil {
		return "", fmt.Errorf("no se pudo obtener el conteo de participantes: %w", err)
	}
	nextNumber := count + 1

	// 3. Formatear el código final (ej: "JUV-001")
	// %03d significa: formatea como un entero, con 3 dígitos, rellenando con ceros a la izquierda.
	participantCode := fmt.Sprintf("%s-%03d", prefix, nextNumber)

	return participantCode, nil
}