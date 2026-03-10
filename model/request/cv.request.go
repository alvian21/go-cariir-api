package request

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobSearchRequest struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PreferredLocation string             `bson:"preferredLocation" json:"preferredLocation"`
	MaxJobs           string             `bson:"maxJobs" json:"maxJobs"`
	UserEmail         string             `bson:"userEmail" json:"userEmail"`
	FirecrawlBase     string             `bson:"firecrawlBase" json:"firecrawlBase"`
	CrawlRunDate      string             `bson:"crawl_run_date" json:"crawl_run_date"`
	CVData            CVData             `bson:"cvData" json:"cvData"`
	SearchKeyword     string             `bson:"searchKeyword" json:"searchKeyword"`
	SearchLocation    string             `bson:"searchLocation" json:"searchLocation"`
	ListingUrl        string             `bson:"listingUrl" json:"listingUrl"`
}

type CVData struct {
	Name            string   `bson:"name" json:"name"`
	Skills          []string `bson:"skills" json:"skills"`
	ExperienceYears int      `bson:"experience_years" json:"experience_years"`
	EducationLevel  string   `bson:"education_level" json:"education_level"`
	JobTitles       []string `bson:"job_titles" json:"job_titles"`
	Industries      []string `bson:"industries" json:"industries"`
	Languages       []string `bson:"languages" json:"languages"`
	Summary         string   `bson:"summary" json:"summary"`
	SearchKeywords  []string `bson:"search_keywords" json:"search_keywords"`
}
