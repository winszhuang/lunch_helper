package constant

type Directive string

const (
	Search              Directive = "/s"
	SearchLocation      Directive = "/sl"
	SearchText          Directive = "/st"
	SearchRadius        Directive = "/sr"
	SearchAI            Directive = "/sai"
	FavoriteRestaurants Directive = "/fr"
	FavoriteFoods       Directive = "/ff"
	PickRestaurant      Directive = "/pr"
	NotificationSetting Directive = "ns"
	UserOption          Directive = "/uo"
	SearchOption        Directive = "/so"
)

func IsDirective(text string) bool {
	switch text {
	case
		string(Search),
		string(SearchLocation),
		string(SearchText),
		string(SearchRadius),
		string(SearchAI),
		string(FavoriteRestaurants),
		string(FavoriteFoods),
		string(PickRestaurant),
		string(NotificationSetting),
		string(UserOption),
		string(SearchOption):
		return true
	}
	return false
}
