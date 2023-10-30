package globals

const (
	LocalhostPort = "8080"

	// environment
	DevelopmentEnvironment = "development"
	StagingEnvironment     = "staging"
	ProductionEnvironment  = "production"

	// client URLs
	MainSiteOrigin            = "https://www.twreporter.org"
	MainSiteDevOrigin         = "http://localhost:3000"
	MainSiteStagingOrigin     = "https://staging.twreporter.org"
	SupportSiteOrigin         = "https://support.twreporter.org"
	SupportSiteDevOrigin      = "http://localhost:3000"
	SupportSiteStagingOrigin  = "https://staging-support.twreporter.org"
	AccountsSiteOrigin        = "https://accounts.twreporter.org"
	AccountsSiteDevOrigin     = "http://localhost:3000"
	AccountsSiteStagingOrigin = "https://staging-accounts.twreporter.org"

	// route path
	SendOtpRoutePath             = "mail/send_otp"
	SendActivationRoutePath      = "mail/send_activation"
	SendAuthenticationRoutePath  = "mail/send_authentication"
	SendSuccessDonationRoutePath = "mail/send_success_donation"
	SendRoleExplorerRoutePath    = "mail/send_role_explorer"
	SendRoleActiontakerRoutePath = "mail/send_role_actiontaker"
	SendRoleTrailblazerRoutePath = "mail/send_role_trailblazer"
	SendRoleDowngradeRoutePath   = "mail/send_role_downgrade"

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

	// Deprecated once v1 news endpoints are removed
	// index page sections //
	LatestSection      = "latest_section"
	EditorPicksSection = "editor_picks_section"
	LatestTopicSection = "latest_topic_section"
	ReviewsSection     = "reviews_section"
	CategoriesSection  = "categories_posts_section"
	TopicsSection      = "topics_section"
	PhotoSection       = "photos_section"
	InfographicSection = "infographics_section"

	// Deprecated once v1 news endpoints are removed
	// index page categories
	HumanRightsAndSociety   = "human_rights_and_society"
	EnvironmentAndEducation = "environment_and_education"
	PoliticsAndEconomy      = "politics_and_economy"
	CultureAndArt           = "culture_and_art"
	International           = "international"
	LivingAndMedicalCare    = "living_and_medical_care"

	// table name
	TableUsersBookmarks            = "users_bookmarks"
	TableBookmarks                 = "bookmarks"
	TablePayByPrimeDonations       = "pay_by_prime_donations"
	TablePayByCardTokenDonations   = "pay_by_card_token_donations"
	TablePayByOtherMethodDonations = "pay_by_other_method_donations"

	// oauth type
	GoogleOAuth   = "Google"
	FacebookOAuth = "Facebook"

	// donation
	PeriodicDonationType = "periodic_donation"
	PrimeDonationType    = "prime"
	TokenDonationType    = "token"
	OthersDonationType   = "others"

	// userType
	UserType = "user"

	// jwt prefix
	MailServiceJWTPrefix = "mail-service-jwt-"

	// custom context key
	AuthUserIDProperty = "auth-user-id"
)
