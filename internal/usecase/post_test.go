package usecase

//func TestCreatePost(t *testing.T) {
//	t.Run("Should create a post with validated data", func(t *testing.T) {
//		post := entity.Post{
//			Title:   "Title 1",
//			Content: "Content 1",
//		}
//
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.CreatePost(&post)
//		if err != nil {
//			t.Errorf("CreatePost should not return an error for a valid post data. Post: %v. Error: %v", post, err)
//		} else if post.Id != usecase.NEW_POST_ID {
//			t.Errorf("CreatePost should set an id for user. User id: %v. User id expected: %v",
//				post.Id,
//				usecase.NEW_POST_ID,
//			)
//		}
//	})
//
//	t.Run("Should not create a post with non-validated data", func(t *testing.T) {
//		scenarios := []entity.Post{
//			{
//				Title: "Title a",
//			},
//			{
//				Content: "Content b",
//			},
//		}
//
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		for _, scenario := range scenarios {
//			err := postUseCase.CreatePost(&scenario)
//			var epv *errorType.ErrorPostValidation
//			if !errors.As(err, &epv) {
//				t.Errorf("CreatePost should return an ErrorPostValidation error for a non-valid post data. Post: %v Returned: %v. Error expected: %T",
//					scenario,
//					err,
//					epv.Error(),
//				)
//			} else if scenario.Id != 0 {
//				t.Errorf("CreatePost should set 0 on post id for non-valid post data. Got: %v.", scenario.Id)
//			}
//		}
//	})
//}
//
//func TestUpdatePost(t *testing.T) {
//	t.Run("Should update post with valid id", func(t *testing.T) {
//		post := entity.Post{Title: "Title test", Content: "Content test"}
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.Update(usecase.MockPosts[0].AuthorId, usecase.MockPosts[0].Id, post)
//		if err != nil {
//			t.Errorf("Update should not return an error with valid data. Data sended: %v. Post updated: %v Error: %v",
//				post,
//				usecase.MockPosts[0],
//				err,
//			)
//		} else if post.Title != usecase.MockPosts[0].Title || post.Content != usecase.MockPosts[0].Content {
//			t.Errorf("Update should update title and content of mock post 0. Data sended: %v. User updated: %v",
//				post,
//				usecase.MockUsers[0],
//			)
//		}
//	})
//
//	t.Run("Should not update user with non-valid post id", func(t *testing.T) {
//		post := entity.Post{Title: "Title test", Content: "Content test"}
//		originalPosts := usecase.MockPosts
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.Update(usecase.MockUsers[0].Id, 0, post)
//		if !errors.Is(err, sql.ErrNoRows) {
//			t.Errorf("Update should return sql.ErrNoRows error with non-valid post id. Data sended: %v. User updated: %v Error: %v",
//				post,
//				usecase.MockPosts[0],
//				err,
//			)
//		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
//			t.Errorf("Update should not update none of posts with non-valid id. Data sended: %v. Posts: %v", post, usecase.MockPosts)
//		}
//	})
//
//	t.Run("Should not update post with non-valid author id", func(t *testing.T) {
//		post := entity.Post{Title: "Title test", Content: "Content test"}
//		originalPosts := usecase.MockPosts
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.Update("wrong-author-id", usecase.MockPosts[0].Id, post)
//		if !errors.Is(err, ErrAccessDenied) {
//			t.Errorf("Update should return ErrAccessDenied error with non-valid author id. Data sended: %v. Post updated: %v Error: %v",
//				post,
//				usecase.MockPosts[0],
//				err,
//			)
//		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
//			t.Errorf("Update should not update none of posts with non-valid author id. Data sended: %v. Posts: %v",
//				post,
//				usecase.MockPosts,
//			)
//		}
//	})
//
//	t.Run("Should not update post with non-valid data", func(t *testing.T) {
//		scenarios := []entity.Post{
//			{
//				Title: "Title a",
//			},
//			{
//				Content: "Content a",
//			},
//		}
//
//		originalPosts := usecase.MockPosts
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		for _, scenario := range scenarios {
//			err := postUseCase.Update(usecase.MockUsers[0].Id, usecase.MockPosts[0].Id, scenario)
//			var euv *errorType.ErrorPostValidation
//			if !errors.As(err, &euv) {
//				t.Errorf("Update should return an ErrorPostValidation error for a non-valid post data. Post: %v Error returned: %v. Error expected: %T",
//					scenario,
//					err,
//					euv.Error(),
//				)
//			} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
//				t.Errorf("Update should not update none of posts with non-valid data. Data sended: %v. Original posts: %v", scenario, originalPosts)
//			}
//		}
//	})
//}
//
//func TestDeletePost(t *testing.T) {
//	t.Run("Should not delete post with non-valid author id", func(t *testing.T) {
//		originalPosts := usecase.MockPosts
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.Delete(usecase.MockPosts[0].Id, "wrong-author-id")
//		if !errors.Is(err, ErrAccessDenied) {
//			t.Errorf("Delete should return ErrAccessDenied error with non-valid author id. Post to delete: %v Error: %v",
//				usecase.MockPosts[0],
//				err,
//			)
//		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
//			t.Errorf("Delete should not delete with non-valid author id. Posts: %v", usecase.MockPosts)
//		}
//	})
//
//	t.Run("Should not delete post with non-valid id", func(t *testing.T) {
//		originalPosts := usecase.MockPosts
//		postUseCase := NewPostUseCase(usecase.NewMockPostRepository())
//		err := postUseCase.Delete(0, usecase.MockPosts[0].AuthorId)
//		if !errors.Is(err, sql.ErrNoRows) {
//			t.Errorf("Delete should return sql.ErrNoRows error with non-valid id. Error: %v", err)
//		} else if originalPosts[0] != usecase.MockPosts[0] || originalPosts[1] != usecase.MockPosts[1] {
//			t.Errorf("Delete should not delete with non-valid id. Posts: %v", usecase.MockPosts)
//		}
//	})
//}
