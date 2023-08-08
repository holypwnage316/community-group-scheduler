package group

import (
	"fmt"
	"time"
)

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

func (g *CommunityGroupMembers) FindPairMember(firstPairMember Member, currentDate time.Time) (Member, error) {
	isFirstMemberTeenager := isTeenager(firstPairMember)
	maximumScheduled := getLargestScheduledNumber(g)
	timesToSearch := 0

SearchAgain:
	for i, memberCandidate := range g.Members {
		// Check to make sure not scheduling the same person twice
		if isThisMe(firstPairMember, memberCandidate) {
			continue
		}

		// Check to make sure if first member is a teenager, we don't schedule another teenager at the same time
		if isFirstMemberTeenager && isTeenager(memberCandidate) {
			continue
		}

		// Check to make sure they are not related
		if areMembersRelated(firstPairMember, memberCandidate) {
			continue
		}

		// Check if one member is a teenager, do not schedule opposite gender adult
		if childAdultGenderCheck(firstPairMember, memberCandidate) {
			continue
		}

		// Check if they are unavailable for the current date
		if unavailableForDate(memberCandidate, currentDate) {
			continue
		}

		// Check if the same gender
		if isSameGender(firstPairMember, memberCandidate) {
			continue
		}

		// Check to see they haven't been scheduled more than anyone else
		if fullyScheduled(memberCandidate, maximumScheduled) {
			continue
		}

		g.Members[i].TimesScheduled = g.Members[i].TimesScheduled + 1
		fmt.Println("Scheduling ", g.Members[i].FirstName)

		return memberCandidate, nil
	}

	if timesToSearch < 5 {
		timesToSearch += 1
		ResetTimesScheduled(g)
		goto SearchAgain
	}

	return Member{}, fmt.Errorf("failed to find pair member")
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

// isThisMe returns true if the member being paired is the same person as the first pair member
func isThisMe(firstMember Member, secondMember Member) bool {
	return (firstMember.FamilyUnit == secondMember.FamilyUnit) && (firstMember.FirstName == secondMember.FirstName)
}

// isTeenager returns true if the member is a teenager
func isTeenager(member Member) bool {
	return member.AgeGroup == "teenager"
}

// areMembersRelated returns true if the two members are related
func areMembersRelated(firstMember Member, secondMember Member) bool {
	return firstMember.FamilyUnit == secondMember.FamilyUnit
}

// childAdultGenderCheck returns true if one member is a teenager and the member they are being paired up to is not the same gender
func childAdultGenderCheck(firstMember Member, secondMember Member) bool {
	return (isTeenager(firstMember) || isTeenager(secondMember)) && (firstMember.Gender != secondMember.Gender)
}

// isSameGender returns true if the genders of the two members are not the same
func isSameGender(firstMember Member, secondMember Member) bool {
	return firstMember.Gender != "" && secondMember.Gender != "" && firstMember.Gender != secondMember.Gender
}

// fullyScheduled returns true if this candidate has already been scheduled and we haven't hit the maxium scheduled threshold
func fullyScheduled(pairCandidate Member, maximumScheduled int) bool {
	return pairCandidate.TimesScheduled >= maximumScheduled && maximumScheduled != 0
}

// ResetTimesScheduled resets everyone's time scheduled back to 0 to reuse them in scheduling
func ResetTimesScheduled(groupMembers *CommunityGroupMembers) {
	for i := range groupMembers.Members {
		groupMembers.Members[i].TimesScheduled = 0
	}
}

// unavailableForDate returns true if the member is unavailable for the date being paired with
func unavailableForDate(member Member, currentDate time.Time) bool {
	for i := range member.DatesUnavailable {
		if UsaDateFormat(currentDate) == (member.DatesUnavailable[i]) {
			return true
		}
	}

	return false
}

// UsaDateFormat returns a date type into a standard USA date format MM/DD/YYYY
func UsaDateFormat(date time.Time) string {
	return date.Format("01/02/2006")
}
