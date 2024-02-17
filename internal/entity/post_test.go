package entity

import (
	"testing"
	"time"
)

const AUTHOR_ID = "eedf21bf-dde8-4c85-b50b-89a1cba87c2e"

type PostScenarios struct {
	post     Post
	expected Post
}

func TestPostPrepare(t *testing.T) {
	t.Run("Should format and validate data with valid post", func(t *testing.T) {
		createdAt := time.Now()
		scenarios := []PostScenarios{
			{
				Post{1, "test title", "test content", AUTHOR_ID, "nick", 0, createdAt},
				Post{1, "test title", "test content", AUTHOR_ID, "nick", 0, createdAt},
			},
			{
				Post{1, "   test title   ", "  test content   ", AUTHOR_ID, "nick", 0, createdAt},
				Post{1, "test title", "test content", AUTHOR_ID, "nick", 0, createdAt},
			},
			{
				Post{1, " test  title ", " test  content ", AUTHOR_ID, "nick", 0, createdAt},
				Post{1, "test  title", "test  content", AUTHOR_ID, "nick", 0, createdAt},
			},
		}

		for _, scenario := range scenarios {
			err := scenario.post.Prepare()

			if err != nil {
				t.Errorf("Post prepare should not return an error for a valid post: %v. Scenario: %v", err, scenario.post)
			}
			if scenario.post == scenario.expected {
				t.Errorf(
					"Post prepare should correctly format the post data. Post: %v. Post expected: %v",
					scenario.post,
					scenario.expected,
				)
			}
		}
	})

	t.Run("Should return error if title is empty", func(t *testing.T) {
		createdAt := time.Now()
		post := Post{1, "", "content", AUTHOR_ID, "nick", 0, createdAt}
		err := post.Prepare()

		if err.Error() != "title is required" {
			t.Errorf("Post prepare should return a 'title is required' error if title is empty. Error: %v", err)
		}
	})

	t.Run("Should return error if content is empty", func(t *testing.T) {
		createdAt := time.Now()
		post := Post{1, "title", "", AUTHOR_ID, "nick", 0, createdAt}
		err := post.Prepare()

		if err.Error() != "content is required" {
			t.Errorf("Post prepare should return a 'content is required' error if content is empty. Error: %v", err)
		}
	})
}
