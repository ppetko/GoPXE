package routes

import (
	"flag"
	"net/http"

	h "github.com/ppetko/gopxe/handlers"

	//External dependencies
	"github.com/gorilla/mux"
)

func New() http.Handler {

	const localrepo string = "/opt/localrepo"
	tftpPath := flag.Lookup("tftpPath").Value.(flag.Getter).Get().(string)

	router := mux.NewRouter()
	router = mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/localrepo").Handler(http.StripPrefix("/localrepo/", http.FileServer(http.Dir(localrepo))))
	router.PathPrefix("/pxelinux").Handler(http.StripPrefix("/pxelinux/", http.FileServer(http.Dir(tftpPath))))
	router.HandleFunc("/", h.Index)
	router.HandleFunc("/viewbootaction", h.BootactionHandler)
	router.HandleFunc("/viewpxeboot", h.PxebootHandler)
	router.HandleFunc("/health", h.StatusHandler)
	router.HandleFunc("/bootaction/{key}", h.GetBA).Methods("GET")
	router.HandleFunc("/bootaction/{key}", h.PutBA).Methods("POST")
	router.HandleFunc("/bootaction", h.GetAllBA).Methods("GET")
	router.HandleFunc("/kickstart/", h.KsGenerate)
	router.HandleFunc("/pxeboot", h.PXEBOOT).Methods("POST")
	h.LoadTemplates()
	return router
}
