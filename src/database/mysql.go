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

func CreateParticipant(p models.Participant) (int64, error) {
	query := `INSERT INTO participantes (
		nombre, apellido_paterno, apellido_materno, email, sexo, categoria, 
		pago_realizado, ine_path, comprobante_pago_path
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	stmt, err := DB.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("error al preparar la consulta: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		p.Nombre, p.ApellidoPaterno, p.ApellidoMaterno, p.Email, p.Sexo, p.Categoria,
		p.PagoRealizado, p.InePath, p.ComprobantePagoPath,
	)
	if err != nil {
		return 0, fmt.Errorf("error al ejecutar la consulta: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("error al obtener el último ID insertado: %w", err)
	}

	return id, nil
}