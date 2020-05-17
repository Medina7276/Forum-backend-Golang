package apihandlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.01.alem.school/qjawko/forum/dto"
	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	uuid "github.com/satori/go.uuid"
)

type SubforumHandler struct {
	SubforumRoleService *service.SubforumRoleService
	PostService         *service.PostService
	SubforumService     *service.SubforumService
	Endpoint            string //
}

func NewSubforumHandler(endpoint string, service *service.SubforumService, post *service.PostService, roleService *service.SubforumRoleService) *SubforumHandler {
	return &SubforumHandler{
		PostService:     post,
		SubforumService: service,
		SubforumRoleService: roleService,
		Endpoint:        endpoint,
	}
}

func (sh *SubforumHandler) Route(w http.ResponseWriter, r *http.Request) {
	var funcToCall func(http.ResponseWriter, *http.Request)

	switch r.Method {
	case http.MethodGet:
		endpoint := r.URL.Path[len(sh.Endpoint):]
		_, ok := r.URL.Query()["name"] //йога
		if ok {
			fmt.Println("name is reached")
			funcToCall = sh.GetSubforumByName
		} else {
			if len(endpoint) == 0 {
				funcToCall = sh.GetAllSubforums
			} else {
				funcToCall = sh.GetSubforumById
			}
		}
	case http.MethodPost:
		funcToCall = sh.CreateSubforum
	case http.MethodPut:
		funcToCall = sh.Update
	case http.MethodDelete:
		funcToCall = sh.Delete
	default:
		http.Error(w, "Route Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	funcToCall(w, r)
}

func (sh *SubforumHandler) CreateSubforum(w http.ResponseWriter, r *http.Request) {
	var subforum model.Subforum

	if err := json.NewDecoder(r.Body).Decode(&subforum); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := sh.SubforumService.CreateSubforum(&subforum)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (sh *SubforumHandler) GetAllSubforums(w http.ResponseWriter, r *http.Request) {
	subforums, err := sh.SubforumService.GetAllSubforums()
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
	json.NewEncoder(w).Encode(subforums)
}

func (sh *SubforumHandler) GetSubforumById(w http.ResponseWriter, r *http.Request) {

	//subforum/459864056/
	id := r.URL.Path[len(sh.Endpoint):]

	subforum, err := sh.SubforumService.GetSubforumById(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	children, err := sh.SubforumService.GetSubforumsByParentId(subforum.ID)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	posts, err := sh.PostService.GetAllPostsBySubforumId(subforum.ID)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
	//data transfer object
	subforumDto := &dto.SubforumDto{
		ID:          subforum.ID,
		Name:        subforum.Name,
		Description: subforum.Description,
		ParentID:    subforum.ParentID,
		Children:    children,
		Posts:       posts,
	}
	json.NewEncoder(w).Encode(subforumDto)
}

func (sh *SubforumHandler) GetSubforumByName(w http.ResponseWriter, r *http.Request) {
	name, _ := r.URL.Query()["name"]
	subforum, err := sh.SubforumService.GetSubforumByName(name[0])
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(subforum)
}

func (sh *SubforumHandler) Update(w http.ResponseWriter, r *http.Request) {

	var subforum model.Subforum
	if err := json.NewDecoder(r.Body).Decode(&subforum); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	updated, err := sh.SubforumService.UpdateSubforum(&subforum)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(updated)
}

func (sh *SubforumHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len(sh.Endpoint):]
	if err := sh.SubforumService.DeleteSubforum(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
	}
}
