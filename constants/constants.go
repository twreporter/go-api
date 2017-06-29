package constants

const (
	MembershipController = "membership_controller"
	NewsController       = "news_controller"
	FacebookController   = "facebook_controller"
	GoogleController     = "google_controller"

	RegistrationTable = "registrations"

	DefaultOrderBy        = "updated_at desc"
	DefaultLimit      int = 0
	DefaultOffset     int = 0
	DefaultActiveCode int = 2

	NewsLetter = "news_letter"

	Activate = "activate"

	// index page sections //
	LastestSection     = "latest"
	EditorPicksSection = "editor_picks"
	LatestTopicSection = "latest_topic"
	ReviewsSection     = "reviews"
	CategoriesSection  = "categories_posts"
	TopicsSection      = "topics"
	PhotoSection       = "photos"
	InfographicSection = "infographics"
)
