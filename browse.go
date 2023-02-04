package deviantart

import (
	"time"

	"github.com/dghubble/sling"
)

// TODO: BROWSE
// 	/dailydeviations
// 	/deviantsyouwatch
// 	/hot
// 	/morelikethis
// 	/morelikethis/preview
// 	/newest
// 	/popular
// 	/posts/deviantsyouwatch
// 	/recommended
// 	/tags
// 	/tags/search
// 	/topic
// 	/topics
// 	/toptopics
// 	/undiscovered
// 	/user/journals

type browseService struct {
	sling *sling.Sling
}

func newBrowseService(sling *sling.Sling) *browseService {
	return &browseService{
		sling: sling.Path("browse/"),
	}
}

type DailyDeviationsParams struct {
	Date        time.Time `url:"date,omitempty" layout:"2006-01-02"`
	WithSession bool      `url:"with_session,omitempty"`
}

func (s *browseService) DailyDeviations(params *DailyDeviationsParams) {
	// TODO: do
}
