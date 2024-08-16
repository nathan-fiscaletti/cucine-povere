package datatypes

type Site struct {
	Logo           string   `yaml:"logo"`
	Title          string   `yaml:"title"`
	TagLine        string   `yaml:"tagline"`
	DateFormat     string   `yaml:"date_format"`
	SharePlatforms []string `yaml:"share_platforms"`
	FontAwesomeKit string   `yaml:"font_awesome_kit"`
	CopyrightYear  int
}
