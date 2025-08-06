package services

import (
	"fmt"
	"log"
	"path/filepath"
)

func UploadFile(localPath string) (string, error) {
	if localPath == "" {
		return "", nil // No hay nada que subir
	}

	fileName := filepath.Base(localPath)
	drivePath := fmt.Sprintf("/gdrive/uploads/%s", fileName)

	// Simulación
	log.Printf("[GDRIVE SIMULADO] Subiendo archivo '%s' a Google Drive...", fileName)
	log.Printf("[GDRIVE SIMULADO] Archivo subido exitosamente. Ruta en Drive: %s", drivePath)

	return drivePath, nil
}