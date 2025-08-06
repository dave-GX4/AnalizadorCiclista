package services

import (
	"fmt"
	"os"
	"time"
	"compilerciclista/src/models"

	"github.com/golang-jwt/jwt/v4"
)

// GenerateToken crea un nuevo token JWT para un participante.
func GenerateToken(participant models.Participant) (string, error) {
	// CAMBIO CLAVE: Leemos la variable de entorno DENTRO de la función.
	var jwtSecretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

	// Ahora esta validación funcionará correctamente.
	if len(jwtSecretKey) == 0 {
		return "", fmt.Errorf("la clave secreta JWT_SECRET_KEY no está configurada o no se pudo leer del .env")
	}

	// El resto de la función es exactamente igual.
	claims := jwt.MapClaims{
		"participant_id": participant.ID,
		"email":          participant.Email,
		"nombre":         participant.Nombre,
		"exp":            time.Now().Add(time.Hour * 72).Unix(),
		"iat":            time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("error al firmar el token: %w", err)
	}

	return tokenString, nil
}