package types

type UserRequest struct {
	Id      string `json:"user_id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Born    string `json:"born"`
	Email   string `json:"email"`
}

type FavoritePlaceRequest struct {
	UserId    string `json:"user_id"`
	PlaceId   string `json:"place_id"`
	PlaceName string
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
