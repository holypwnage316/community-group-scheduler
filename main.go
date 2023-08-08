package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"

	"communitygroupscheduler/group"
)

const (
	FileName  = "groupInfo.json"
	StartDate = "08/14/2023"
	EndDate   = "12/05/2023"
)

type ScheduleItem struct {
	Date string
	Pair string
}

func main() {
	communityGroupMembers, err := parseCommunityGroupMembers(FileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println("Scheduling " + strconv.Itoa(len(communityGroupMembers.Members)) + " community group members to lead the kids...")

	schedule, err := createSchedule(StartDate, EndDate, communityGroupMembers)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("\n" + schedule)
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

// createSchedule is the main function creating the leading the kids schedule
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

		if len(pair) < 1 {
			pairMember, err := groupMembers.FindPairMember(group.Member{}, currentDate)
			if err != nil {
				fmt.Println("Failed to find first member of a pair")
			}
			pair = append(pair, pairMember)
			continue
		}

		// Find a second pair member thru the algorithm
		secondPairMember, err := groupMembers.FindPairMember(pair[0], currentDate)
		if err != nil {
			fmt.Println("Failed to find second member of a pair")
		}

		pair = append(pair, secondPairMember)

		// Save the schedule item into the schedule
		scheduleItem.Pair = pair[0].FirstName + ", " + pair[1].FirstName
		scheduleItems = append(scheduleItems, scheduleItem)

		// Move to the next date 14 days out
		currentDate = currentDate.AddDate(0, 0, 14)

		// Reinitialize the pair
		pair = make([]group.Member, 0)

		// Initialize the next schedule item
		scheduleItem.Date = group.UsaDateFormat(currentDate)
		scheduleItem.Pair = ""
	}

	return createPrintableSchedule(scheduleItems), nil
}

// createPrintableSchedule creates a printable schedule like: date - name, name
func createPrintableSchedule(scheduleItems []ScheduleItem) string {
	var schedule string

	for _, item := range scheduleItems {
		schedule = schedule + item.Date + " - " + item.Pair + "\n"
	}

	return schedule
}
