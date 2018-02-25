package main

/**********************

prueba para hacer un APIREST con GO

por el momento todas las funciones se encuentran en el archivo main
al igual que las conexiones.

por el momento solo se puede hacer consultar ya sea por id o
ver todos los datos

falta mover las funciones a otro archivo y cosas así xd

***********************/
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

//estructura de datos para almacenar los datos devueltos por las consultas
type Usuario struct {
	ID     int    `json:"id, omitempty"`
	NOMBRE string `json:"nombre, omitempty"`
	EDAD   int    `json:"edad, omitempty"`
}

type PageVariables struct {
	PageTitle string
}

//funcion para realizar inserciones dentro de la base de datos
func insert(nombre string, edad int) {
	//conexion con la base de datos
	db, err := sql.Open("mysql", "root:gundam000@/api_rest")

	//verifica si no hay un error al conectar
	if err != nil {
		panic(err.Error())
	}

	//cierra la conexion de la BD despues de terminar la funcion
	defer db.Close()

	//prepara una ejecucion de QUERY
	inser, err := db.Prepare("INSERT INTO usuarios(nombre, edad) VALUES (?,?)")

	//verifica si hay algun error en la QUERY
	if err != nil {
		panic(err.Error())
	}

	//cierra la insercion al final
	defer inser.Close()

	//ejecuta la insercion de datos
	//al parecer cualquier ejecucion se hace igualandolo a una variable
	//como no requerimos guardar un resultado usamos una variable blanca _
	_, er := inser.Exec(nombre, edad)

	if er != nil {
		panic(er.Error())
	}
}

//funcion para obtener un usuario
func getUsuario(w http.ResponseWriter, r *http.Request) {
	//se insertan los datos necesarios para la conexion
	db, err := sql.Open("mysql", "root:gundam000@/api_rest")

	//recuperamos los parametros que enviamos al acceder a la direccion
	//en este caso solo el ID
	params := mux.Vars(r)

	//variable de tipo USUARIO para usarla en el JSON
	var user = Usuario{}

	//error en la conexion
	if err != nil {
		panic(err.Error())
	}

	//aquí ejecutamos la sentencia
	rows, err := db.Query("SELECT * FROM usuarios WHERE id = ?", params["id"])

	//error en la sentencia
	if err != nil {
		panic(err.Error())
	}

	//iteramos el resultado de la consulta
	for rows.Next() {
		//construimos la variable USER con el resultado
		err := rows.Scan(&user.ID, &user.NOMBRE, &user.EDAD)
		//verificamos si hay algun error
		if err != nil {
			panic(err.Error())
		} else {
			//convertimos en JSON el resultado y lo enviamos
			json.NewEncoder(w).Encode(user)
			return
		}
	}
}

func getUsuarios(w http.ResponseWriter, r *http.Request) {
	var usuarios []Usuario
	db, err := sql.Open("mysql", "root:gundam000@/api_rest")

	if err != nil {
		panic(err.Error())
	}

	rows, err := db.Query("SELECT * FROM usuarios")

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		//creamos una variable de la estructura USUARIO en cada iteracion
		var user Usuario
		err := rows.Scan(&user.ID, &user.NOMBRE, &user.EDAD)
		if err != nil {
			panic(err.Error())
		} else {
			//agregamos ese usuario a un arreglo
			usuarios = append(usuarios, user)
		}
	}

	//convertimos ese arreglo en JSON y lo enviamos
	json.NewEncoder(w).Encode(usuarios)

}

func createUsuario(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:gundam000@/api_rest")

	if err != nil {
		panic(err.Error())
	}

	defer db.Close()

	inser, err := db.Prepare("INSERT INTO usuarios(nombre, edad) VALUES (?,?)")

	if err != nil {
		panic(err.Error())
	}

	defer inser.Close()

	_, er := inser.Exec(r.FormValue("nombre"), r.FormValue("edad"))

	if er != nil {
		panic(er.Error())
	} else {
		fmt.Println("registro correcto :)")
	}

}

func showForm(w http.ResponseWriter, r *http.Request) {
	//esto se supone que permite usar un html de template pero no me sale xd
	// t, err := template.ParseFiles("index.html")

	// var mensaje = "error"

	// MyPageVariables := PageVariables{
	// 	PageTitle: mensaje,
	// }

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }

	// err = t.Execute(w, MyPageVariables)

	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
}

func main() {
	//iniciamos rourter y creamos las rutas
	router := mux.NewRouter()
	router.HandleFunc("/usuarios", getUsuarios).Methods("GET")
	router.HandleFunc("/usuario/{id}", getUsuario).Methods("GET")
	router.HandleFunc("/usuario/create", createUsuario).Methods("POST")
	router.HandleFunc("/usuario/create", showForm).Methods("GET")
	// router.HandleFunc("/usuario/{id}", updateUsuario).Methods("POST")
	// router.HandleFunc("/usuario/{id}", deleteUsuario).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}
