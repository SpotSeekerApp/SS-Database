package types

var TAGLIST = []string{
	"Family-Friendly",
	"Cozy",
	"Live Music",
	"Romantic",
	"Quiet Atmosphere",
	"Good Hygiene",
	"Good Service",
	"Pet-Friendly",
	"Comfortable Seating",
	"Unique Decoration",
	"Quick Service",
	"Late-Night Spot",
	"Outdoor Seating",
	"Affordable",
	"Sports-Friendly",
	"Easy Location",
	"Vegan",
	"Study-Friendly",
	"Board Games",
	"Healthy Options",
}

type UserRequest struct {
	UserId         string             `json:"user_id" firebase:"userID"`
	UserName       string             `json:"user_name" firebase:"userName"`
	Email          string             `json:"email" firebase:"email"`
	FavoritePlaces map[string]string  `json:"favorite_places" firebase:"favoritePlaces"`
	PlaceId        string             `json:"place_id" firebase:"placeId"`
	UserType       string             `json:"user_type" firebase:"userType"`
	Tags           map[string]float64 `json:"tags" firebase:"tags"`
	PlaceName      string
}

type FeedbackRequest struct {
	FeedbackId string
	UserId     string            `json:"user_id"`
	PlaceId    string            `json:"place_id"`
	Rating     map[string]string `json:"rating"`
}

type FilterRequest struct {
	UserId string   `json:"user_id" firebase:"userId"`
	Tags   []string `json:"tags" firebase:"tags"`
}

type PlaceRequest struct {
	UserId       string             `json:"user_id" firebase:"userId"`
	PlaceId      string             `json:"place_id" firebase:"placeId"`
	PlaceName    string             `json:"place_name" firebase:"placeName"`
	MainCategory string             `json:"main_category" firebase:"mainCategory"`
	Link         string             `json:"link" firebase:"link'"`
	Tags         map[string]float64 `json:"tags"`
	FirstReview  string             `json:"first_review" firebase:"firstReview"`
	Images       []string           `json:"images" firebase:"images"`
}

type ReviewRequest struct {
	ReviewId     string `json:"review_id" firebase:"reviewId"`
	ReviewerName string `json:"reviewer_name" firebase:"reviewerName"`
	Rating       int    `json:"rating" firebase:"rating"`
	Comment      string `json:"comment" firebase:"comment"`
	Date         string `json:"date" firebase:"date"`
	UserId       string `json:"user_id" firebase:"userId"`
	PlaceId      string `json:"place_id" firebase:"placeId"`
}
