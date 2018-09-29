package globals

const (
	// controller name
	MembershipController = "membership_controller"
	NewsController       = "news_controller"

	RegistrationTable = "registrations"

	DefaultOrderBy        = "updated_at desc"
	DefaultLimit      int = 10
	DefaultOffset     int = 0
	DefaultActiveCode int = 2

	NewsLetter = "news_letter"

	Activate = "activate"

	// index page sections //
	LastestSection     = "latest_section"
	EditorPicksSection = "editor_picks_section"
	LatestTopicSection = "latest_topic_section"
	ReviewsSection     = "reviews_section"
	CategoriesSection  = "categories_posts_section"
	TopicsSection      = "topics_section"
	PhotoSection       = "photos_section"
	InfographicSection = "infographics_section"

	// index page categories
	HumanRightsAndSociety   = "human_rights_society"
	EnvironmentAndEducation = "environment_education"
	PoliticsAndEconomy      = "politics_economy"
	CultureAndArt           = "culture_art"
	International           = "international"
	LivingAndMedicalCare    = "living_medical_care"

	// table name
	TableUsersBookmarks = "users_bookmarks"
	TableBookmarks      = "bookmarks"

	// oauth type
	GoogleOAuth   = "Google"
	FacebookOAuth = "Facebook"

	// donation
	PeriodicDonationType = "periodic_donation"
	PrimeDonaitionType   = "prime"
	TokenDonationType    = "token"
)
