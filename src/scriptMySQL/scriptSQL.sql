create database ciclista;

use ciclista;

CREATE TABLE participantes (
    id INT AUTO_INCREMENT PRIMARY KEY,
    nombre VARCHAR(255) NOT NULL,
    apellido_paterno VARCHAR(255) NOT NULL,
    apellido_materno VARCHAR(255),
    email VARCHAR(255) NOT NULL UNIQUE,
    sexo CHAR(1) NOT NULL,
    categoria VARCHAR(100) NOT NULL,
    pago_realizado BOOLEAN DEFAULT false,
    ine_path VARCHAR(255),
    comprobante_pago_path VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

#una ves creado la tabla ejecutar este query para añadir nueva propiedad
ALTER TABLE participantes
ADD COLUMN participant_code VARCHAR(50) NOT NULL UNIQUE,
ADD INDEX (participant_code);