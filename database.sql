CREATE DATABASE api_prueba;

use api_prueba;

CREATE TABLE usuarios(
    id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
    nombre VARCHAR(20) NOT NULL,
    edad INT NOT NULL
);