package deviantart

type LicenseOptions struct {
	CreativeCommons bool `url:"creative_commons,omitempty"`
	Commercial      bool `url:"commercial,omitempty"`
	// Valid values: yes, no, share
	Modify string `url:"modify,omitempty"`
}
