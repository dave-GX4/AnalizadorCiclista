package database

import (
	"database/sql"
	"fmt"
	"compilerciclista/src/models"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB(dataSourceName string) error {
	var err error
	DB, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}

	return DB.Ping()
}

func CountParticipants() (int64, error) {
	var count int64
	query := "SELECT COUNT(id) FROM participantes"
	err := DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func CreateParticipant(p models.Participant) (int64, error) {
	// Se añade 'participant_code' al query y a los valores.
	query := `INSERT INTO participantes (
		participant_code, nombre, apellido_paterno, apellido_materno, email, sexo, categoria, 
		pago_realizado, ine_path, comprobante_pago_path
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		p.ParticipantCode, // <-- Se añade el nuevo valor aquí
		p.Nombre, p.ApellidoPaterno, p.ApellidoMaterno, p.Email, p.Sexo, p.Categoria,
		p.PagoRealizado, p.InePath, p.ComprobantePagoPath,
	)
	if err != nil {
		// El error de "Duplicate entry" ahora podría ser por el email o por el participant_code
		return 0, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el último ID insertado: %w", err)
	}

	return id, nil
}