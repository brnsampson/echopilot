package appserver

import (
	"net/http"
    "time"

    "github.com/brnsampson/echopilot/internal/templates"
    "github.com/brnsampson/echopilot/pkg/option"
    "github.com/brnsampson/echopilot/rpc/echo"
    "github.com/go-chi/chi/v5"
)

func routeRoot(r chi.Router) {
    r.Get("/", getIndex)
    r.Get("/echo/{content}", getEcho)
    r.Post("/echo", postEcho)
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		filePath := "web/404.html"
		http.ServeFile(w, r, filePath)
	}
}

func getIndex(w http.ResponseWriter, r *http.Request) {
    templates.Index().Render(r.Context(), w)
}

func getEcho(w http.ResponseWriter, r *http.Request) {
    content := chi.URLParam(r, "content")
    templates.Echo(content).Render(r.Context(), w)
}

func postEcho(w http.ResponseWriter, r *http.Request) {
    var content string
    r.ParseForm()
    if r.Form.Has("content") {
        content = r.Form.Get("content")
    } else {
        content = ""
    }
	timeout := option.NewOption(time.Duration(10) * time.Second)
	skipVerify := option.NewOption(true)
    client, err := echo.NewRemoteEchoClient("https://127.0.0.1:1443", timeout, skipVerify)
    if err != nil {
        errorHandler(w, r, 500)
        return
    }
    req := echo.NewStringRequest(content)
    res, err := client.EchoString(req)
    if err != nil {
        errorHandler(w, r, 500)
        return
    }

    echoed := echo.ReadStringResult(res)
    http.Redirect(w, r, "/echo/" + echoed, http.StatusSeeOther)
}
