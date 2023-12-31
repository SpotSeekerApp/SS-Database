package types

type UserRequest struct {
	UserId         string            `json:"user_id" firebase:"userID"`
	UserName       string            `json:"user_name" firebase:"userName"`
	Email          string            `json:"email" firebase:"email"`
	FavoritePlaces map[string]string `json:"favorite_places" firebase:"favoritePlaces"`
	PlaceId        string            `json:"place_id" firebase:"placeId"`
	UserType       string            `json:"user_type" firebase:"userType"`
	PlaceName      string
}

type FeedbackRequest struct {
	FeedbackId int
	UserId     int               `json:"user_id"`
	PlaceId    string            `json:"place_id"`
	Rating     map[string]string `json:"rating"`
}

type PlaceRequest struct {
	UserId       string             `json:"user_id"`
	PlaceId      string             `json:"place_id"`
	PlaceName    string             `json:"place_name"`
	MainCategory string             `json:"main_category"`
	Link         string             `json:"link"`
	Tags         map[string]float32 `json:"tags"`
	Review       string             `json:"review"`
}

type ReviewRequest struct {
	ReviewId     string  `json:"review_id"`
	ReviewerName string  `json:"reviewer_name"`
	Rating       float32 `json:"rating"`
	Comment      string  `json:"comment"`
	Date         string  `json:"date"`
}
