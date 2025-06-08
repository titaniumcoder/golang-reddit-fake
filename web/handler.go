package web

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/titaniumcoder/golang-reddit-fake/goreddit"
)

func NewHandler(store goreddit.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)

	h.Get("/", h.Home())

	h.Route("/threads", func(r chi.Router) {
		r.Get("/", h.ThreadsList())
		r.Get("/new", h.ThreadsCreate())
		r.Post("/", h.ThreadsStore())
		r.Get("/{id}", h.ThreadsShow())
		r.Post("/{id}/delete", h.ThreadsDelete())
		r.Post("/{id}", h.PostsStore())
		r.Get("/{id}/new", h.PostsCreate())
		r.Get("/{threadID}/{postID}", h.PostsShow())
		r.Post("/{threadID}/{postID}", h.CommentsStore())
	})

	return h
}

type Handler struct {
	*chi.Mux

	store goreddit.Store
}

func (h *Handler) Home() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/home.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}

func (h *Handler) ThreadsList() http.HandlerFunc {
	type data struct {
		Threads []goreddit.Thread
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/threads.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tt, err := h.store.Threads()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{tt})
	}
}

func (h *Handler) ThreadsShow() http.HandlerFunc {
	type data struct {
		Thread goreddit.Thread
		Posts  []goreddit.PostWithCount
	}

	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		pp, err := h.store.PostsByThread(t.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, data{Thread: t, Posts: pp})
	}
}

func (h *Handler) ThreadsCreate() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/thread_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		tmpl.Execute(w, nil)
	}
}
func (h *Handler) ThreadsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		description := r.FormValue("description")

		if err := h.store.CreateThread(&goreddit.Thread{
			ID:          uuid.New(),
			Title:       title,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
func (h *Handler) ThreadsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")

		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.store.DeleteThread(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/threads", http.StatusFound)
	}
}
func (h *Handler) PostsShow() http.HandlerFunc {
	type data struct {
		Post     goreddit.Post
		Comments []goreddit.Comment
	}
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")
		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		cc, err := h.store.Comments(p.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, data{Post: p, Comments: cc})
	}
}
func (h *Handler) PostsCreate() http.HandlerFunc {
	tmpl := template.Must(template.ParseFiles("templates/layout.html", "templates/post_create.html"))
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, struct {
			Thread goreddit.Thread
		}{
			Thread: t,
		})
	}
}
func (h *Handler) PostsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		content := r.FormValue("text")

		idStr := chi.URLParam(r, "id")
		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		t, err := h.store.Thread(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		newID := uuid.New()
		if err := h.store.CreatePost(&goreddit.Post{
			ID:       newID,
			Title:    title,
			ThreadID: t.ID,
			Content:  content,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/threads/%s/%s", t.ID.String(), newID.String()), http.StatusFound)
	}
}
func (h *Handler) CommentsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := chi.URLParam(r, "postID")
		id, err := uuid.Parse(idStr)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p, err := h.store.Post(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		comment := r.FormValue("comment")

		if err := h.store.CreateComment(&goreddit.Comment{
			ID:      uuid.New(),
			PostID:  p.ID,
			Content: comment,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, fmt.Sprintf("/threads/%s/%s", p.ThreadID, p.ID), http.StatusFound)
	}
}
