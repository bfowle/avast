package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "text/template"

    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
)

var (
    router     *mux.Router
    addr       string
    apiVersion string
    indexTpl   *template.Template = template.Must(template.ParseFiles("index.html"))
)

func setHeaders(w http.ResponseWriter, headers map[string]string) {
    for k, v := range headers {
        w.Header().Set(http.CanonicalHeaderKey(k), v)
    }
}

func processEnv() {
    if addr = os.Getenv("AVAST_ADDR"); addr == "" {
        addr = ":8080"
    }
    if apiVersion = os.Getenv("AVAST_API_VERSION"); apiVersion == "" {
        apiVersion = "v1"
    }
}

func startWebserver() {
    processEnv()

    router = mux.NewRouter()
    router.HandleFunc("/ws", wrap(wsHandler))
    router.HandleFunc("/", wrap(indexHandler))

    dockerRouter := router.PathPrefix(fmt.Sprintf("/api/%v/docker", apiVersion)).Subrouter()
    dockerRouter.HandleFunc("/containers",         wrap(dockerClient.ContainersHandler))
    dockerRouter.HandleFunc("/container/{name}",   wrap(dockerClient.ContainerHandler))
    dockerRouter.HandleFunc("/images",             wrap(dockerClient.ImagesHandler))
    dockerRouter.HandleFunc("/image/history/{id}", wrap(dockerClient.HistoryHandler))
    dockerRouter.HandleFunc("/info",               wrap(dockerClient.InfoHandler))

    consulRouter := router.PathPrefix(fmt.Sprintf("/api/%v/consul", apiVersion)).Subrouter()
    consulRouter.HandleFunc("/datacenters",        wrap(consulRegistry.DatacentersHandler))
    consulRouter.HandleFunc("/nodes",              wrap(consulRegistry.NodesHandler))
    consulRouter.HandleFunc("/nodes/{dc}",         wrap(consulRegistry.NodesHandler))
    consulRouter.HandleFunc("/node/{name}",        wrap(consulRegistry.NodeHandler))
    consulRouter.HandleFunc("/health/{name}",      wrap(consulRegistry.HealthHandler))
    consulRouter.HandleFunc("/health/{name}/{dc}", wrap(consulRegistry.HealthHandler))

    http.Handle("/", router)
    loggedRouter := handlers.CombinedLoggingHandler(os.Stdout, router)
    panic(http.ListenAndServe(addr, handlers.CompressHandler(loggedRouter)))
}

func indexHandler(w http.ResponseWriter, r *http.Request) (interface{}, error) {
    indexTpl.Execute(w, r.Host)
    return nil, nil
}

func wrap(handler func(w http.ResponseWriter, r *http.Request) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) {
    f := func(w http.ResponseWriter, r *http.Request) {
        var headers = make(map[string]string)
        headers["Content-Type"] = "text/html; charset=utf8"

        if origin := r.Header.Get("Origin"); origin != "" {
            headers["Access-Control-Allow-Origin"] = origin
            headers["Access-Control-Allow-Credentials"] = "true"
            headers["Access-Control-Allow-Methods"] = "POST, GET, OPTIONS, PUT, DELETE"
            headers["Access-Control-Allow-Headers"] = "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization"
        }
        setHeaders(w, headers)

        obj, err := handler(w, r)

    HAS_ERR:
        if err != nil {
            code := http.StatusInternalServerError
            errMsg := err.Error()
            if strings.Contains(errMsg, "Permission denied") {
                code = http.StatusForbidden
            }

            w.WriteHeader(code)
            w.Write([]byte(err.Error()))
            return
        }

        if obj != nil {
            buf, err := json.Marshal(obj)
            if err != nil {
                goto HAS_ERR
            }

            w.Header().Set("Content-Type", "application/json")
            w.Write(buf)
        }
    }

    return f
}
