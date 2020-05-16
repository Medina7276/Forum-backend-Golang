package operations

import (
	"encoding/json"
	"net/http"

	"git.01.alem.school/qjawko/forum/model"
	"git.01.alem.school/qjawko/forum/service"
	"git.01.alem.school/qjawko/forum/utils"
)

type LikeOperations struct {
	likeService *service.LikeService
}

func NewLikeOperations(likeService *service.LikeService) *LikeOperations {
	return &LikeOperations{
		likeService: likeService,
	}
}

func (lo *LikeOperations) Rate(w http.ResponseWriter, r *http.Request) {
	var newrate model.Like
	if err := json.NewDecoder(r.Body).Decode(&newrate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := utils.GetUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if user.ID != newrate.UserID {
		http.Error(w, "UserId and rate.UserId are not equal", http.StatusBadRequest)
		return
	}

	var res *model.Like
	rate, err := lo.likeService.GetLikeByID(newrate.ID)
	if err != nil { //returns err when no rates found so we create it
		res, err = lo.likeService.CreateLike(&newrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		if rate.IsUpVote == newrate.IsUpVote {
			return
		}

		res, err = lo.likeService.UpdateLike(&newrate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	json.NewEncoder(w).Encode(res)
}
