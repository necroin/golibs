package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/necroin/golibs/libs/winappstream"
	"github.com/necroin/golibs/utils/winapi"
	"github.com/necroin/golibs/utils/winutils"
)

func RecursiveFind(allProcesses []*winutils.Process, parentPid winapi.ProcessId, level int) {
	processes := winutils.FindProcessesByParentPid(allProcesses, parentPid)
	for _, process := range processes {
		fmt.Printf("%s%s\n", strings.Repeat("\t", level), process)
		RecursiveFind(allProcesses, process.Pid, level+1)
	}
}

func main() {
	// allProcesses, err := winutils.GetAllProcesses()
	// if err != nil {
	// 	panic(err)
	// }

	// RecursiveFind(allProcesses, 19668, 0)

	app, err := winappstream.NewApp(20548)
	if err != nil {
		panic(err)
	}
	defer app.Destroy()
	app.LaunchStream()

	router := mux.NewRouter()

	router.HandleFunc("/stream", func(responseWriter http.ResponseWriter, r *http.Request) {
		data, _ := os.ReadFile("winappstream/example/page.html")
		responseWriter.Write(data)
	}).Methods("GET")
	router.Handle("/image", app.HttpImageCaptureHandler())

	server := http.Server{
		Addr:      "localhost:3301",
		Handler:   router,
		TLSConfig: &tls.Config{},
	}
	server.ListenAndServe()
}
