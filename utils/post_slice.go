package utils

import "git.01.alem.school/qjawko/forum/model"

func PostsFilter(posts []model.Post, predicate func(model.Post) bool) []model.Post {
	var res []model.Post
	for _, post := range posts {
		if predicate(post) {
			res = append(res, post)
		}
	}

	return res
}
