package utils

import "git.01.alem.school/qjawko/forum/model"

func CommentsFilter(comments []model.Comment, predicate func(model.Comment) bool) []model.Comment {
	var res []model.Comment
	for _, comment := range comments {
		if predicate(comment) {
			res = append(res, comment)
		}
	}

	return res
}
