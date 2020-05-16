package apihandlers

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/qjawko/forum/http_errors"
	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	uuid "github.com/satori/go.uuid"
)

type SubforumRoleHandler struct {
	SubforumRoleService *service.SubforumRoleService
	Endpoint            string
}

func NewSubforumRoleHandler(endpoint string, service *service.SubforumRoleService) *SubforumRoleHandler {
	return &SubforumRoleHandler{
		SubforumRoleService: service,
		Endpoint:            endpoint,
	}
}

func (srh *SubforumRoleHandler) Route(w http.ResponseWriter, r *http.Request) {
	var funcToCall func(http.ResponseWriter, *http.Request)

	switch r.Method {
	case http.MethodGet:
		funcToCall = srh.GetSubforumRoleById
	case http.MethodPost:
		funcToCall = srh.CreateSubforumRole
	case http.MethodPut:
		funcToCall = srh.UpdateSRole
	case http.MethodDelete:
		funcToCall = srh.DeleteSRole
	default:
		http.Error(w, "Route Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	funcToCall(w, r)

}

func (srh *SubforumRoleHandler) CreateSubforumRole(w http.ResponseWriter, r *http.Request) {

	var subforumRole model.SubforumRole

	if err := json.NewDecoder(r.Body).Decode(&subforumRole); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	created, err := srh.SubforumRoleService.Create(&subforumRole)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func (srh *SubforumRoleHandler) GetSubforumRoleById(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len(srh.Endpoint):]

	SubforumRole, err := srh.SubforumRoleService.GetBySubforumId(uuid.FromStringOrNil(id))
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

	json.NewEncoder(w).Encode(SubforumRole)
}

func (srh *SubforumRoleHandler) UpdateSRole(w http.ResponseWriter, r *http.Request) {

	var subforumRole model.SubforumRole
	if err := json.NewDecoder(r.Body).Decode(&subforumRole); err != nil {
		http.Error(w, "Bad Body", http.StatusBadRequest)
		return
	}

	idFromURL := r.URL.Path[len(srh.Endpoint):]
	subforumRole.ID = uuid.FromStringOrNil(idFromURL)

	updated, err := srh.SubforumRoleService.UpdateSubforumRole(&subforumRole)
	if err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}
	json.NewEncoder(w).Encode(updated)
}

func (srh *SubforumRoleHandler) DeleteSRole(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Path[len(srh.Endpoint):]
	if err := srh.SubforumRoleService.DeleteSubforumRoleById(uuid.FromStringOrNil(id)); err != nil {
		httpErr := err.(*http_errors.HttpError)
		http.Error(w, httpErr.Error(), httpErr.Code)
		return
	}

}
