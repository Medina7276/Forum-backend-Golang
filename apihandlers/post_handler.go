package apihandlers

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/qjawko/forum/dto"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/jwt"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	uuid "github.com/satori/go.uuid"
)

type PostHandler struct {
	SubforumService     *service.SubforumService
	SubforumRoleService *service.SubforumRoleService
	PostService         *service.PostService
	CommentService      *service.CommentService
	Endpoint            string
}

func NewPostHandler(endpoint string, postService *service.PostService, commentService *service.CommentService) *PostHandler {
	return &PostHandler{
		PostService:    postService,
		CommentService: commentService,
		Endpoint:       endpoint,
	}
}

func contains(roles []model.SubforumRole, id uuid.UUID) bool {
	for _, r := range roles {
		if r.ID == id {
			return true
		}
	}

	return false
}

func (ph *PostHandler) checkForPostDelete(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var user model.User
		if err := jwt.Unmarshal(cookie.Value, "supersecret", user); err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		id := r.URL.Path[len(ph.Endpoint):]
		post, err := ph.PostService.GetPostById(uuid.FromStringOrNil(id))
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		subforum, err := ph.SubforumService.GetSubforumById(post.SubofrumID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		admins, err := ph.SubforumRoleService.GetBySubforumId(subforum.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		if user.ID == post.UserID || contains(admins, user.ID) {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "У вас нет прав", http.StatusBadRequest)
			return
		}
	})
}

func (ph *PostHandler) Route(w http.ResponseWriter, r *http.Request) {
	var funcToCall func(http.ResponseWriter, *http.Request)

	switch r.Method {
	case http.MethodGet:
		endpoint := r.URL.Path[len(ph.Endpoint):]
		if len(endpoint) == 0 {
			funcToCall = ph.GetAll
		} else {
			funcToCall = ph.GetPostByID
		}

	case http.MethodPost:
		funcToCall = ph.CreatePost
	case http.MethodPut:
		funcToCall = ph.checkForPostUpdate(http.HandlerFunc(ph.Update)).ServeHTTP
	case http.MethodDelete:
		funcToCall = ph.checkForPostDelete(http.HandlerFunc(ph.Delete)).ServeHTTP
	default:
		http.Error(w, "Route Not found", http.StatusNotFound)
		return
	}

	funcToCall(w, r)
}

//CreatePost q
func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := ph.PostService.CreatePost(&post)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

//GetAll qwe
func (ph *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {

	userID := r.FormValue("userid")
	subforumID := r.FormValue("subforumid")

	if userID != "" {
		posts, err := ph.PostService.GetUserPosts(uuid.FromStringOrNil(userID))
		if err != nil {
			httpErr := err.(*http_errors.HttpError)
			http.Error(w, httpErr.Error(), httpErr.Code)
			return
		}
		json.NewEncoder(w).Encode(posts)
		return
	}

	if subforumID != "" {
		posts, err := ph.PostService.GetAllPostsBySubforumId(uuid.FromStringOrNil(subforumID))
		if err != nil {
			httpErr := err.(*http_errors.HttpError)
			http.Error(w, httpErr.Error(), httpErr.Code)
			return
		}
		json.NewEncoder(w).Encode(posts)
		return
	}

	posts, err := ph.PostService.GetAllPosts()
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(posts)
}

//GetPostByID возвращает postDto, который хранит пост с указанным ID и его комментарии
func (ph *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len(ph.Endpoint):]

	post, err := ph.PostService.GetPostById(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	comments, err := ph.CommentService.GetAllCommentsByPostID(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	postDto := &dto.PostDto{
		ID:           post.ID,
		Title:        post.Title,
		Content:      post.Content,
		CreationDate: post.CreationDate,
		SubofrumID:   post.SubofrumID,
		UserID:       post.UserID,
		Comments:     comments,
	}

	json.NewEncoder(w).Encode(postDto)
}

func (ph *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	var post model.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	idFromURL := r.URL.Path[len(ph.Endpoint):]
	post.ID = uuid.FromStringOrNil(idFromURL)

	updated, err := ph.PostService.UpdatePost(&post)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (ph *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len(ph.Endpoint):]
	if err := ph.PostService.DeletePost(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
	}
}
