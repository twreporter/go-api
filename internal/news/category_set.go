package news

type CategorySet struct {
	Name string
	Key  string
}

var (
	World              = CategorySet{"world", "63206383207bf7c5f871622c"}
	Humanrights        = CategorySet{"humanrights", "63206383207bf7c5f8716234"}
	PoliticsAndSociety = CategorySet{"politics_and_society", "63206383207bf7c5f871623d"}
	Health             = CategorySet{"health", "63206383207bf7c5f8716245"}
	Environment        = CategorySet{"environment", "63206383207bf7c5f871624d"}
	Econ               = CategorySet{"econ", "63206383207bf7c5f8716254"}
	Culture            = CategorySet{"culture", "63206383207bf7c5f8716259"}
	Education          = CategorySet{"education", "63206383207bf7c5f8716260"}
	Podcast            = CategorySet{"podcast", "63206383207bf7c5f8716266"}
	Opinion            = CategorySet{"opinion", "63206383207bf7c5f8716269"}
)
