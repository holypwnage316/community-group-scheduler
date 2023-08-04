package group

import "fmt"

type CommunityGroupMembers struct {
	Members []Member `json:"groupmembers"`
}

type Member struct {
	FamilyUnit       string   `json:"familyunit"`
	FirstName        string   `json:"firstname"`
	Gender           string   `json:"gender"`
	AgeGroup         string   `json:"agegroup"`
	DatesUnavailable []string `json:"datesunavailable"`
	TimesScheduled   int
}

func (g *CommunityGroupMembers) FindSecondPairMember(firstPairMember Member) (Member, error) {
	isFirstMemberTeenager := isTeenager(firstPairMember)
	maximumScheduled := getLargestScheduledNumber(g)

	for i, secondMemberCandidate := range g.Members {
		// Check to make sure not scheduling the same person twice
		if isThisMe(firstPairMember, secondMemberCandidate) {
			continue
		}

		// Check to make sure if first member is a teenager, we don't schedule another teenager at the same time
		if isFirstMemberTeenager && isTeenager(secondMemberCandidate) {
			continue
		}

		// Check to make sure they are not related
		if areMembersRelated(firstPairMember, secondMemberCandidate) {
			continue
		}

		// Check if one member is a teenager, do not schedule opposite gender adult
		if childAdultGenderCheck(firstPairMember, secondMemberCandidate) {
			continue
		}

		// Check to see they haven't been scheduled more than anyone else
		if fullyScheduled(secondMemberCandidate, maximumScheduled) {
			continue
		}

		g.Members[i].TimesScheduled = g.Members[i].TimesScheduled + 1

		return secondMemberCandidate, nil
	}

	return Member{}, fmt.Errorf("no member found")
}

func getLargestScheduledNumber(groupMembers *CommunityGroupMembers) int {
	tempCounter := 0
	largestScheduledNumber := 0

	for _, member := range groupMembers.Members {
		if member.TimesScheduled > tempCounter {
			tempCounter = member.TimesScheduled
			largestScheduledNumber = tempCounter
		}
	}

	return largestScheduledNumber
}

func isThisMe(firstMember Member, secondMemberCandidate Member) bool {
	return (firstMember.FamilyUnit == secondMemberCandidate.FamilyUnit) && (firstMember.FirstName == secondMemberCandidate.FirstName)
}

func isTeenager(member Member) bool {
	return member.AgeGroup == "teenager"
}

func areMembersRelated(firstMember Member, secondMemberCandidate Member) bool {
	return firstMember.FamilyUnit == secondMemberCandidate.FamilyUnit
}

func childAdultGenderCheck(firstMember Member, secondMemberCandidate Member) bool {
	return (isTeenager(firstMember) || isTeenager(secondMemberCandidate)) && (firstMember.Gender != secondMemberCandidate.Gender)
}

func fullyScheduled(secondMemberCandidate Member, maximumScheduled int) bool {
	return secondMemberCandidate.TimesScheduled >= maximumScheduled
}

func ResetTimesScheduled(groupMembers *CommunityGroupMembers) {
	for _, member := range groupMembers.Members {
		member.TimesScheduled = 0
	}
}
