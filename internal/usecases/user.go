package usecases

type User struct {
	UserId  int64            `json:"user_id"`
	Name    string           `json:"name"`
	Balance float64          `json:"balance"`
	Loyalty int64            `json:"loyalty"`
	Orders  map[string]Order `json:"orders"`
}
