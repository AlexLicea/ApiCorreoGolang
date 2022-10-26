package main

import (
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slices"
)

//Se declara la estructura a utilizar con los campos de usuario, telefono, correo y contrasena
type datos struct {
	Usuario  string `json:"user"`
	Telefono string `json:"phone"`
	Correo   string `json:"mail"`
	Password string `json:"password"`
}

//Se crea un slice con datos por default para validar mas adelante
var datosDefault = []datos{
	{Usuario: "alicea", Telefono: "4422602875", Correo: "alexislicea@gmail.com", Password: "123Abc+"},
	{Usuario: "rleon", Telefono: "4422602873", Correo: "alexisleon@gmail.com", Password: "123Abc+w"},
}

//Retorna todos los datos que se encuentran en datosDefault para mostrarlos en un json al recibir una peticion GET
func getDatos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, datosDefault)
}

//Agrega un dato a datosDefault si se cumplen las validacions, retorna los datos recibidos si son exitosas las validaciones, caso contrario retorna un message de error
func addDatos(context *gin.Context) {
	var newDatos datos
	if err := context.BindJSON(&newDatos); err != nil {
		return
	}

	//Realiza la busqueda en el slice datosDefault para verificar si ya se encuentra registrado el telefono o el correo
	indiceTelefono := slices.IndexFunc(datosDefault, func(c datos) bool { return c.Telefono == newDatos.Telefono })
	indiceCorreo := slices.IndexFunc(datosDefault, func(c datos) bool { return c.Correo == newDatos.Correo })
	if indiceTelefono >= 0 {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "El telefono ya se encuentra registado"})
		return
	}
	if indiceCorreo >= 0 {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "El correo ya se encuentra registado"})
		return
	}

	//Validaciones con expresiones regulares para validar que el mail y el telefono son validos
	reMail := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	reTelefono := regexp.MustCompile("^(\\d{10})$")
	if reMail.MatchString(newDatos.Correo) != true {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "El correo no tiene el formato valido"})
		return
	}
	if reTelefono.MatchString(newDatos.Telefono) != true {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "El telefono no tiene el formato valido"})
		return
	}

	//Validaciones para verificar que los campos no vienen vacios
	if newDatos.Usuario == "" {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Falta el campo de usuario"})
		return
	}
	if newDatos.Telefono == "" {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Falta el campo de telefono"})
		return
	}
	if newDatos.Correo == "" {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Falta el campo de correo"})
		return
	}
	if newDatos.Password == "" {
		context.IndentedJSON(http.StatusConflict, gin.H{"message": "Falta el campo de password"})
		return
	}

	context.IndentedJSON(http.StatusCreated, newDatos)
}

func main() {
	router := gin.Default()
	//Metodo get que retorna todos los valores en registrados en datosDefault
	router.GET("/datos", getDatos)
	//Metodo post para ingresar datos de un nuevo usuario a datosDefault
	router.POST("/datos", addDatos)
	//ruta en la que se ejecuta nuestra api
	router.Run("localhost:9090")
}
