package main

import (
	"fmt"
)

//Battery class definition
type Battery struct {
	ID                        int
	amountOfColumns           int
	status                    string
	amountOfFloors            int
	amountOfBasements         int
	amountOfElevatorPerColumn int
	servedFloors              []int
	columnID                  int
	floorRequestButtonID      int
	columnsList               []Column
	floorRequestButtonsList   []FloorRequestButton
}

func (b *Battery) makeBasementColumn(amountOfBasements int, amountOfElevatorPerColumn int) {
	var servedFloors []int
	floor := -1
	for i := 0; i < (amountOfBasements + 1); i++ {
		if i == 0 {
			servedFloors = append(servedFloors, 1)
		} else {
			servedFloors = append(servedFloors, floor)
			floor--
		}
	}
	fmt.Printf("Floor: %v", servedFloors)
}

//Column class definition
type Column struct {
	ID                int
	status            string
	amountOfElevators int
	servedFloorsList  []int
	isBasement        bool
	elevatorsList     []Elevator
	callButtonsList   []CallButton
}

//Elevator class definition
type Elevator struct {
	ID               int
	status           string
	servedFloorsList []int
	currentFloor     int
	direction        string
	door             []Door
	floorRequestList []int
}

//BestElevatorInfo class definition
type BestElevatorInfo struct {
	bestElevator Elevator
	bestScore    int
	referenceGap int
}

//CallButton class definition
type CallButton struct {
	ID        int
	status    string
	floor     int
	direction string
}

//FloorRequestButton class definition
type FloorRequestButton struct {
	ID     int
	status string
	floor  int
}

//Door class definition
type Door struct {
	ID     int
	status string
}

func main() {
	battery := Battery{
		ID:                        1,
		amountOfColumns:           4,
		status:                    "online",
		amountOfFloors:            60,
		amountOfBasements:         6,
		amountOfElevatorPerColumn: 5,
		columnsList:               []Column{},
		floorRequestButtonsList:   []FloorRequestButton{},
	}

	fmt.Println("keeping fmt from throwing fucking error")

	battery.makeBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)

}
