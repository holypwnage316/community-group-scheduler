package group

type CommunityGroupMembers struct {
	Members []Member `json:"groupmembers"`
}

type Member struct {
	FamilyUnit       string   `json:"familyunit"`
	FirstName        string   `json:"firstname"`
	Gender           string   `json:"gender"`
	AgeGroup         string   `json:"agegroup"`
	DatesUnavailable []string `json:"datesunavailable"`
}
