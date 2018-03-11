//IMPLEMENTACION DE UNA API REST SENCILLA
package main
import (
	"time"
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"strconv"
)

type Note struct {
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedOn time.Time `json:"created_on"`
}

//creamos un seudo base de datos con un map de tipo Note
var noteStore = make(map[string]Note)
var id int

//Función manejadora GET
func GetNoteHandler(w http.ResponseWriter, r *http.Request){
	//Crear un slice de Note
	var notes []Note
	for _,v := range noteStore {
		notes = append(notes,v)
	}
	//Crear la cabecera
	w.Header().Set("Content-Type","application/json")
	//Convertir Go a JSON
	j,err := json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	//Escribir la cabecera
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}
//Función manejadora POST
func PostNoteHandler(w http.ResponseWriter, r *http.Request){
	var note Note
	//Convertir JSON a Go
	err := json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	//Agregar fecha de creación
	note.CreatedOn = time.Now()
	id++
	k := strconv.Itoa(id)
	noteStore[k] = note
	//Devolver JSON con fecha de creación
	w.Header().Set("Content-Type","application/json")
	j,err := json.Marshal(note)
	if err != nil {
		panic(err)
	}
	//Escribir la cabecera
	w.WriteHeader(http.StatusCreated)
	w.Write(j)
}
//Función manejadora PUT
func PutNoteHandler(w http.ResponseWriter, r *http.Request){
	//Extraer variables del request
	vars := mux.Vars(r)
	k := vars["id"]
	var noteUpdate Note
	//Convertir JSON a estructura Go
	err := json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		panic(err)
	}
	//Actualizar la seudo base de datos
	if note,ok := noteStore[k]; ok {
		noteUpdate.CreatedOn = note.CreatedOn
		delete(noteStore,k)
		noteStore[k] = noteUpdate
	} else {
		log.Printf("No encontramos el id %s", k)
	}
	//Devolver código de que no hay contenido
	w.WriteHeader(http.StatusNoContent)
}
//Función manejadora DELETE
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	k := vars["id"]
	//Borrar el elemento
	if _,ok := noteStore[k]; ok {
		delete(noteStore,k)
	} else {
		log.Printf("No encontramos el id %s", k)
	}
	//Devolver código de que no hay contenido
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	//Creamos el enrutador con Gorilla Mux
	r := mux.NewRouter().StrictSlash(false)

	//Implementamos las funciones como manejador con HandleFunc
	r.HandleFunc("/api/notes",GetNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes",PostNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}",PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}",DeleteNoteHandler).Methods("DELETE")

	//Creamos el servidor
	server := &http.Server{
		Addr:		":8080",
		Handler:	r,
		ReadTimeout:	10*time.Second,
		WriteTimeout:	10*time.Second,
		MaxHeaderBytes:	1 << 20,
	}
	log.Println("Listening http://localhost:8080...")
	server.ListenAndServe()
}

