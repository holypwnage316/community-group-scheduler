package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"communitygroupscheduler/group"
)

const (
	FileName  = "groupInfo.json"
	StartDate = "09/01/2023"
	EndDate   = "12/01/2023"
)

type ScheduleItem struct {
	Date string
	Pair string
}

func main() {
	communityGroupMembers, err := parseCommunityGroupMembers(FileName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Scheduling %d community group members to lead the kids...\n\n", len(communityGroupMembers.Members))
	fmt.Print(createSchedule(StartDate, EndDate, communityGroupMembers))
}

// parseCommunityGroupMembers opens a JSON file and parses the group members out of it
func parseCommunityGroupMembers(fileName string) (group.CommunityGroupMembers, error) {
	var groupMembers group.CommunityGroupMembers

	// Open the JSON file
	jsonFile, err := os.Open(FileName)
	if err != nil {
		return group.CommunityGroupMembers{}, err
	}

	// Extract and marshal the group members out of the file
	byteValue, _ := io.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &groupMembers)

	return groupMembers, nil
}

func createSchedule(startDate string, endDate string, groupMembers group.CommunityGroupMembers) (string, error) {
	var pair []group.Member
	var scheduleItems []ScheduleItem

	// Convert startDate to a time type and intialize our currentDate to track through the algorithm
	startDateTime, _ := time.Parse("01/02/2006", startDate)
	endDateTime, _ := time.Parse("01/02/2006", endDate)
	currentDate := startDateTime

	// Seed the first schedule item
	scheduleItem := ScheduleItem{
		Date: startDate,
	}

	// Loop until we have met the end date
	for !currentDate.After(endDateTime) {
		// Shuffle the list of group members
		rand.Shuffle(len(groupMembers.Members), func(i, j int) {
			groupMembers.Members[i], groupMembers.Members[j] = groupMembers.Members[j], groupMembers.Members[i]
		})

		for i, member := range groupMembers.Members {
			// If there isn't a current pair schedule this member
			if len(pair) < 1 {
				groupMembers.Members[i].TimesScheduled += 1
				pair = append(pair, member)
				continue
			}

			// Find a second pair member thru the algorithm
			secondPairMember, err := groupMembers.FindSecondPairMember(pair[0])
			if err != nil {
				fmt.Print("Errored\n")
				// for _, member := range groupMembers.Members {
				// 	fmt.Print(member.FirstName, " scheduled: ", member.TimesScheduled, "\n")
				// }
				// Reset everyone 'times scheduled' back to zero and try again
				group.ResetTimesScheduled(&groupMembers)

				secondPairMember, err = groupMembers.FindSecondPairMember(pair[0])
				if err != nil {
					return "", err
				}
			}

			pair = append(pair, secondPairMember)

			// Save the schedule item into the schedule
			scheduleItem.Pair = pair[0].FirstName + ", " + pair[1].FirstName
			scheduleItems = append(scheduleItems, scheduleItem)

			// Move to the next date 7 days out
			currentDate = currentDate.AddDate(0, 0, 7)

			// Reinitialize the pair
			pair = make([]group.Member, 0)

			// Initialize the next schedule item
			scheduleItem.Date = usaDateFormat(currentDate)
			scheduleItem.Pair = ""

			break
		}
	}

	for _, member := range groupMembers.Members {
		fmt.Print(member.FirstName, " scheduled: ", member.TimesScheduled, "\n")
	}

	return createPrintableSchedule(scheduleItems), nil
}

// usaDateFormat returns a date type into a standard USA date format MM/DD/YYYY
func usaDateFormat(date time.Time) string {
	return date.Format("01/02/2006")
}

// createPrintableSchedule creates a printable schedule like: date - name, name
func createPrintableSchedule(scheduleItems []ScheduleItem) string {
	var schedule string

	for _, item := range scheduleItems {
		schedule = schedule + item.Date + " - " + item.Pair + "\n"
	}

	return schedule
}
