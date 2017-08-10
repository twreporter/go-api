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
)
