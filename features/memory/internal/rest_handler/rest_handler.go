package rest_handler

import (
    "net/http"
    "github.com/brnsampson/echopilot/features/memory/templates"
    "github.com/brnsampson/echopilot/features/memory/records"
    "github.com/brnsampson/echopilot/pkg/option"
    "github.com/go-chi/chi/v5"
    "github.com/charmbracelet/log"
)

func NewRestHandler(store *records.MemoryStore) chi.Router {
    rh := MemoryResourceHandler { store }
    router := chi.NewRouter()
    router.Get("/", rh.listMemories)
    router.Post("/", rh.postMemory)
    return router
}

type MemoryResourceHandler struct {
    store *records.MemoryStore
}

func (rh *MemoryResourceHandler) postMemory(w http.ResponseWriter, r *http.Request) {
    r.ParseForm()
    if r.Form.Has("content") {
        content := r.Form.Get("content")
        m := records.NewMemory(content)
        log.Infof("created memories %v", m)
        rh.store.Create(m)
    }
    http.Redirect(w, r, "/memory", http.StatusSeeOther)
}

func (rh *MemoryResourceHandler) listMemories(w http.ResponseWriter, r *http.Request) {
    memories := rh.store.List(option.None[*records.Memory](), 0, 0)
    log.Infof("current memories %v", memories)

    templates.Page(memories).Render(r.Context(), w)
}
