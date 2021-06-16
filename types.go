package main

type responseBody struct {
	Articles []struct {
		AuthorID          int64         `json:"author_id"`
		Body              string        `json:"body"`
		CommentsDisabled  bool          `json:"comments_disabled"`
		CreatedAt         string        `json:"created_at"`
		Draft             bool          `json:"draft"`
		EditedAt          string        `json:"edited_at"`
		HTMLURL           string        `json:"html_url"`
		ID                int64         `json:"id"`
		LabelNames        []string      `json:"label_names"`
		Locale            string        `json:"locale"`
		Name              string        `json:"name"`
		Outdated          bool          `json:"outdated"`
		OutdatedLocales   []interface{} `json:"outdated_locales"`
		PermissionGroupID int64         `json:"permission_group_id"`
		Position          int64         `json:"position"`
		Promoted          bool          `json:"promoted"`
		SectionID         int64         `json:"section_id"`
		SourceLocale      string        `json:"source_locale"`
		Title             string        `json:"title"`
		UpdatedAt         string        `json:"updated_at"`
		URL               string        `json:"url"`
		UserSegmentID     int64         `json:"user_segment_id"`
		VoteCount         int64         `json:"vote_count"`
		VoteSum           int64         `json:"vote_sum"`
	} `json:"articles"`
	Count        int64       `json:"count"`
	NextPage     string      `json:"next_page"`
	Page         int64       `json:"page"`
	PageCount    int64       `json:"page_count"`
	PerPage      int64       `json:"per_page"`
	PreviousPage interface{} `json:"previous_page"`
	SortBy       string      `json:"sort_by"`
	SortOrder    string      `json:"sort_order"`
}
