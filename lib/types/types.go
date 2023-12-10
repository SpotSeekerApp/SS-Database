package types

type UserRequest struct {
	UserId   int    `json:"user_id" firebase:"userID"`
	UserName string `json:"user_name" firebase:"userName"`
	Email    string `json:"email" firebase:"email"`
}

type FavoritePlaceRequest struct {
	UserId    string `json:"user_id"`
	PlaceId   string `json:"place_id"`
	PlaceName string
}

type FeedbackRequest struct {
	UserId  string            `json:"user_id"`
	PlaceId string            `json:"place_id"`
	Rating  map[string]string `json:"rating"`
}

type PlaceRequest struct {
	Id          string `json:"place_id"`
	Name        string `json:"place_name"`
	Location    string `json:"location"`
	Link2Photo  string `json:"photo_link"`
	PhoneNumber string `json:"phone_number"`
}

type ReviewRequest struct {
	PlaceId  string `json:"place_id"`
	Comment  string `json:"comment"`
	Rating   string `json:"rating"`
	ReviewId string `json:"review_id"`
}
