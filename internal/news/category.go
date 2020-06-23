package news

type Category struct {
	ID   string
	Name string
}

var (
	HumanRightsAndSociety   = Category{"5951db87507c6a0d00ab063c", "human_rights_and_society"}
	EnvironmentAndEducation = Category{"5951db9b507c6a0d00ab063d", "environment_and_education"}
	PoliticsAndEconomy      = Category{"5951dbc2507c6a0d00ab0640", "politics_and_economy"}
	CultureAndArt           = Category{"57175d923970a5e46ff854db", "culture_and_art"}
	International           = Category{"5743d35a940ee41000e81f4a", "international"}
	LivingAndMedicalCare    = Category{"59783ad89092de0d00b41691", "living_and_medical_care"}
	Review                  = Category{"573177cb8c0c261000b3f6d2", "review"}
	Photography             = Category{"574d028748fa171000c45d48", "photography"}
)
