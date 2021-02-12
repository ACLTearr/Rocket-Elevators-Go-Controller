package main

import (
	"fmt"
	"math"
	"sort"
)

var columnID int = 1
var floorRequestButtonID int = 1
var callButtonID = 1
var buttonFloor = 1

//Battery class definition
type Battery struct {
	ID                        int
	amountOfColumns           int
	status                    string
	amountOfFloors            int
	amountOfBasements         int
	amountOfElevatorPerColumn int
	servedFloors              []int
	columnsList               []Column
	floorRequestButtonsList   []FloorRequestButton
}

//Method to create basement column
func (battery *Battery) makeBasementColumn(amountOfBasements int, amountOfElevatorPerColumn int) {
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
	column := Column{
		ID:                columnID,
		status:            "online",
		amountOfElevators: amountOfElevatorPerColumn,
		servedFloorsList:  servedFloors,
		isBasement:        true,
		elevatorsList:     []Elevator{},
		callButtonsList:   []CallButton{},
	}
	column.makeElevators(servedFloors, column.amountOfElevators)
	column.makeCallButtons(column.servedFloorsList, battery.amountOfBasements, column.isBasement)
	battery.columnsList = append(battery.columnsList, column)
	columnID++
}

//Method to create columns
func (battery *Battery) makeColumns(amountOfColumns int, amountOfFloors int, amountOfElevatorPerColumn int) {
	amountOfFloorsPerColumn := int(math.Ceil(float64(amountOfFloors / amountOfColumns)))
	floor := 1
	for i := 0; i < amountOfColumns; i++ {
		var servedFloors []int
		for x := 0; x < amountOfFloorsPerColumn; x++ {
			if i == 0 { //For first above ground column
				servedFloors = append(servedFloors, floor)
				floor++
			} else {
				if x == 0 { //Adding lobby to other columns
					servedFloors = append(servedFloors, 1)
				}
				servedFloors = append(servedFloors, floor)
				floor++
			}
		}
		column := Column{
			ID:                columnID,
			status:            "online",
			amountOfElevators: amountOfElevatorPerColumn,
			servedFloorsList:  servedFloors,
			isBasement:        false,
			elevatorsList:     []Elevator{},
			callButtonsList:   []CallButton{},
		}
		column.makeElevators(servedFloors, column.amountOfElevators)
		column.makeCallButtons(column.servedFloorsList, battery.amountOfBasements, column.isBasement)
		battery.columnsList = append(battery.columnsList, column)
		columnID++
	}
}

//Method to create basement floor request buttons
func (battery *Battery) makeBasementFloorRequestButtons(amountOfBasements int) {
	buttonFloor := -1
	for i := 0; i < amountOfBasements; i++ {
		floorRequestButton := FloorRequestButton{
			ID:     floorRequestButtonID,
			status: "off",
			floor:  buttonFloor,
		}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		buttonFloor--
		floorRequestButtonID++
	}
}

//Method to create buttons to request a floor
func (battery *Battery) makeFloorRequestButtons(amountOfFloors int) {
	buttonFloor := 1
	for i := 0; i < amountOfFloors; i++ {
		floorRequestButton := FloorRequestButton{
			ID:     floorRequestButtonID,
			status: "off",
			floor:  buttonFloor,
		}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		buttonFloor++
		floorRequestButtonID++
	}
}

//Method to find the appropriate elevator within the appropriate column to serve user
func (battery *Battery) assignElevator(requestedFloor int, direction string) {
	fmt.Println("A request for an elevator is made from the lobby for floor", requestedFloor, "going", direction, ".")
	var column Column = battery.findBestColumn(requestedFloor)
	fmt.Println("Column", column.ID, "is the column that can handle this request.")
	var elevator Elevator = column.findBestElevator(1, direction)
	stopFloor := elevator.floorRequestList[0]
	fmt.Println("Elevator", elevator.ID, "is the best elevator, so it is sent.")
	if elevator.status == "moving" {
		elevator.moveElevator(stopFloor)
	}
	elevator.floorRequestList = append(elevator.floorRequestList, requestedFloor)
	elevator.sortFloorlist()
	fmt.Println("Elevator is moving.")
	elevator.moveElevator(stopFloor)
	fmt.Println("Elevator is", elevator.status, ".")
	elevator.doorController()
	if len(elevator.floorRequestList) == 0 {
		elevator.direction = ""
		elevator.status = "idle"
	}
	fmt.Println("Elevator is", elevator.status, ".")
}

//Method to find appropriate column to serve user
func (battery *Battery) findBestColumn(requestedFloor int) Column {
	var bestColumn Column
	for _, column := range battery.columnsList {
		foundColumn := contains(requestedFloor, column.servedFloorsList)
		if foundColumn {
			bestColumn = column
		}
	}
	return bestColumn
}

//method to check if users floor is in columns served floors
func contains(a int, list []int) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
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

//Creating elevator method
func (column *Column) makeElevators(servedFloorsList []int, amountOfElevators int) {
	elevatorID := 1
	for i := 0; i < amountOfElevators; i++ {
		elevator := Elevator{
			ID:               elevatorID,
			status:           "idle",
			servedFloorsList: servedFloorsList,
			currentFloor:     1,
			direction:        "",
			door:             Door{ID: elevatorID, status: "closed"},
			floorRequestList: []int{},
		}
		column.elevatorsList = append(column.elevatorsList, elevator)
		elevatorID++
	}
}

//Method to create call buttons
func (column *Column) makeCallButtons(floorsServed []int, amountOfBasements int, isBasement bool) {
	if isBasement {
		buttonFloor := -1
		for i := 0; i < amountOfBasements; i++ {
			callButton := CallButton{
				ID:        callButtonID,
				status:    "off",
				floor:     buttonFloor,
				direction: "up",
			}
			column.callButtonsList = append(column.callButtonsList, callButton)
			buttonFloor--
			callButtonID++
		}
	} else {
		for _, floor := range floorsServed {
			callButton := CallButton{
				ID:        callButtonID,
				status:    "off",
				floor:     floor,
				direction: "down",
			}
			column.callButtonsList = append(column.callButtonsList, callButton)
			buttonFloor++
			callButtonID++
		}
	}
}

//When a user calls an elevator form a floor, not the lobby
func (column *Column) requestElevator(userFloor int, direction string) {
	fmt.Println("A request for an elevator is made from", userFloor, "going", direction, "to the lobby.")
	var elevator Elevator = column.findBestElevator(userFloor, direction)
	fmt.Println("Elevator", elevator.ID, "is the best elevator, so it is sent.")
	elevator.floorRequestList = append(elevator.floorRequestList, 1)
	elevator.sortFloorlist()
	fmt.Println("Elevator is moving.")
	elevator.moveElevator(userFloor)
	fmt.Println("Elevator is", elevator.status, ".")
	elevator.doorController()
	if len(elevator.floorRequestList) == 0 {
		elevator.direction = ""
		elevator.status = "idle"
	}
	fmt.Println("Elevator is", elevator.status, ".")
}

//Find best elevator to send
func (column *Column) findBestElevator(floor int, direction string) Elevator {
	requestedFloor := floor
	requestedDirection := direction
	bestElevatorInfo := BestElevatorInfo{
		bestElevator: Elevator{},
		bestScore:    6,
		referenceGap: 10000000,
	}

	if requestedFloor == 1 {
		for _, elevator := range column.elevatorsList {
			//Elevator is at lobby with some requests, and about to leave but has not yet
			if 1 == elevator.currentFloor && elevator.status == "stopped" {
				bestElevatorInfo = column.checkBestElevator(1, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is at lobby with no requests
			} else if 1 == elevator.currentFloor && elevator.status == "idle" {
				bestElevatorInfo = column.checkBestElevator(2, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is lower than user and moving up. Shows user is requesting to go to basement, and elevator is moving to them.
			} else if 1 > elevator.currentFloor && elevator.direction == "up" {
				bestElevatorInfo = column.checkBestElevator(3, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is higher than user and moving down. Shows user is requesting to go to a floor, and elevator is moving to them.
			} else if 1 < elevator.currentFloor && elevator.direction == "down" {
				bestElevatorInfo = column.checkBestElevator(3, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is not at lobby floor, but has no requests
			} else if elevator.status == "idle" {
				bestElevatorInfo = column.checkBestElevator(4, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is last resort
			} else {
				bestElevatorInfo = column.checkBestElevator(5, elevator, bestElevatorInfo, requestedFloor)
			}
		}
	} else {
		for _, elevator := range column.elevatorsList {
			//Elevator is at floor going to lobby
			if requestedFloor == elevator.currentFloor && elevator.status == "stopped" && requestedDirection == elevator.direction {
				bestElevatorInfo = column.checkBestElevator(1, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is lower than user and moving through them to destination
			} else if requestedFloor > elevator.currentFloor && elevator.direction == "up" && requestedDirection == "up" {
				bestElevatorInfo = column.checkBestElevator(2, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is higher than user and moving through them to destination
			} else if requestedFloor < elevator.currentFloor && elevator.direction == "down" && requestedDirection == "down" {
				bestElevatorInfo = column.checkBestElevator(2, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is idle
			} else if elevator.status == "idle" {
				bestElevatorInfo = column.checkBestElevator(3, elevator, bestElevatorInfo, requestedFloor)
				//Elevator is last resort
			} else {
				bestElevatorInfo = column.checkBestElevator(4, elevator, bestElevatorInfo, requestedFloor)
			}
		}
	}
	return bestElevatorInfo.bestElevator
}

func (column *Column) checkBestElevator(scoreToCheck int, newElevator Elevator, bestElevatorInfo BestElevatorInfo, floor int) BestElevatorInfo {
	//If elevators situation is more favourable, set to best elevator
	if scoreToCheck < bestElevatorInfo.bestScore {
		bestElevatorInfo.bestScore = scoreToCheck
		bestElevatorInfo.bestElevator = newElevator
		bestElevatorInfo.referenceGap = int(math.Abs(float64(newElevator.currentFloor - floor)))
		//If elevators are in a similar situation, set the closest one to the best elevator
	} else if bestElevatorInfo.bestScore == scoreToCheck {
		gap := int(math.Abs(float64(newElevator.currentFloor - floor)))
		if bestElevatorInfo.referenceGap > gap {
			bestElevatorInfo.bestScore = scoreToCheck
			bestElevatorInfo.bestElevator = newElevator
			bestElevatorInfo.referenceGap = gap
		}
	}
	return bestElevatorInfo
}

//Elevator class definition
type Elevator struct {
	ID               int
	status           string
	servedFloorsList []int
	currentFloor     int
	direction        string
	door             Door
	floorRequestList []int
}

func (elevator *Elevator) moveElevator(stopFloor int) {
	for len(elevator.floorRequestList) != 0 {
		destination := elevator.floorRequestList[0]
		elevator.status = "moving"
		if elevator.currentFloor < destination {
			elevator.direction = "up"
			for elevator.currentFloor < destination {
				if elevator.currentFloor == stopFloor {
					elevator.status = "stopped"
					elevator.doorController()
					elevator.currentFloor++
				} else {
					elevator.currentFloor++
				}
				if elevator.currentFloor == 0 {
					//Do nothing, so that moving from basement to/from 1 doesnt show 0
				} else {
					fmt.Println("Elevator is at floor:", elevator.currentFloor)
				}
			}
		} else if elevator.currentFloor > destination {
			elevator.direction = "down"
			for elevator.currentFloor > destination {
				if elevator.currentFloor == stopFloor {
					elevator.status = "stopped"
					elevator.doorController()
					elevator.currentFloor--
				} else {
					elevator.currentFloor--
				}
				if elevator.currentFloor == 0 {
					//Do nothing, so that moving from basement to/from 1 doesnt show 0
				} else {
					fmt.Println("Elevator is at floor:", elevator.currentFloor)
				}
			}
		}
		elevator.status = "stopped"
		elevator.floorRequestList = elevator.floorRequestList[1:]
	}
}

func (elevator *Elevator) sortFloorlist() {
	if elevator.direction == "up" {
		sort.Slice(elevator.floorRequestList, func(a, b int) bool { return elevator.floorRequestList[a] < elevator.floorRequestList[b] })
	} else {
		sort.Slice(elevator.floorRequestList, func(a, b int) bool { return elevator.floorRequestList[a] > elevator.floorRequestList[b] })
	}
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

//Door operation controller
func (elevator *Elevator) doorController() {
	overweight := false
	obstruction := false
	elevator.door.status = "opened"
	fmt.Println("Elevator doors are", elevator.door.status)
	fmt.Println("Waiting for occupant(s) to transition")
	//Wait 5 seconds
	if !overweight {
		elevator.door.status = "closing"
		fmt.Println("Elevator doors are", elevator.door.status)
		if !obstruction {
			elevator.door.status = "closed"
			fmt.Println("Elevator doors are", elevator.door.status)
		} else {
			//Wait for obstruction to clear
			obstruction = false
			elevator.doorController()
		}
	} else {
		for overweight {
			//Ring alarm and wait until not overweight
			overweight = false
		}
		elevator.doorController()
	}
}

//Defining scenario 1
func (battery *Battery) scenario1() {
	battery.columnsList[1].elevatorsList[0].currentFloor = 20
	battery.columnsList[1].elevatorsList[0].direction = "down"
	battery.columnsList[1].elevatorsList[0].status = "moving"
	battery.columnsList[1].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 5)

	battery.columnsList[1].elevatorsList[1].currentFloor = 3
	battery.columnsList[1].elevatorsList[1].direction = "up"
	battery.columnsList[1].elevatorsList[1].status = "moving"
	battery.columnsList[1].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 15)

	battery.columnsList[1].elevatorsList[2].currentFloor = 13
	battery.columnsList[1].elevatorsList[2].direction = "down"
	battery.columnsList[1].elevatorsList[2].status = "moving"
	battery.columnsList[1].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[1].elevatorsList[3].currentFloor = 15
	battery.columnsList[1].elevatorsList[3].direction = "down"
	battery.columnsList[1].elevatorsList[3].status = "moving"
	battery.columnsList[1].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 2)

	battery.columnsList[1].elevatorsList[4].currentFloor = 6
	battery.columnsList[1].elevatorsList[4].direction = "down"
	battery.columnsList[1].elevatorsList[4].status = "moving"
	battery.columnsList[1].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(20, "up")
}

//Defining scenario 2
func (battery *Battery) scenario2() {
	battery.columnsList[2].elevatorsList[0].currentFloor = 1
	battery.columnsList[2].elevatorsList[0].direction = "up"
	battery.columnsList[2].elevatorsList[0].status = "stopped"
	battery.columnsList[2].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 21)

	battery.columnsList[2].elevatorsList[1].currentFloor = 23
	battery.columnsList[2].elevatorsList[1].direction = "up"
	battery.columnsList[2].elevatorsList[1].status = "moving"
	battery.columnsList[2].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 28)

	battery.columnsList[2].elevatorsList[2].currentFloor = 33
	battery.columnsList[2].elevatorsList[2].direction = "down"
	battery.columnsList[2].elevatorsList[2].status = "moving"
	battery.columnsList[2].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[2].elevatorsList[3].currentFloor = 40
	battery.columnsList[2].elevatorsList[3].direction = "down"
	battery.columnsList[2].elevatorsList[3].status = "moving"
	battery.columnsList[2].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 24)

	battery.columnsList[2].elevatorsList[4].currentFloor = 39
	battery.columnsList[2].elevatorsList[4].direction = "down"
	battery.columnsList[2].elevatorsList[4].status = "moving"
	battery.columnsList[2].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(36, "up")
}

//Defining scenario 3
func (battery *Battery) scenario3() {
	battery.columnsList[3].elevatorsList[0].currentFloor = 58
	battery.columnsList[3].elevatorsList[0].direction = "down"
	battery.columnsList[3].elevatorsList[0].status = "moving"
	battery.columnsList[3].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[3].elevatorsList[1].currentFloor = 50
	battery.columnsList[3].elevatorsList[1].direction = "up"
	battery.columnsList[3].elevatorsList[1].status = "moving"
	battery.columnsList[3].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 60)

	battery.columnsList[3].elevatorsList[2].currentFloor = 46
	battery.columnsList[3].elevatorsList[2].direction = "up"
	battery.columnsList[3].elevatorsList[2].status = "moving"
	battery.columnsList[3].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 58)

	battery.columnsList[3].elevatorsList[3].currentFloor = 1
	battery.columnsList[3].elevatorsList[3].direction = "up"
	battery.columnsList[3].elevatorsList[3].status = "moving"
	battery.columnsList[3].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 54)

	battery.columnsList[3].elevatorsList[4].currentFloor = 60
	battery.columnsList[3].elevatorsList[4].direction = "down"
	battery.columnsList[3].elevatorsList[4].status = "moving"
	battery.columnsList[3].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[3].requestElevator(54, "down")
}

//Defining scenario 4
func (battery *Battery) scenario4() {
	battery.columnsList[0].elevatorsList[0].currentFloor = -4

	battery.columnsList[0].elevatorsList[1].currentFloor = 1

	battery.columnsList[0].elevatorsList[2].currentFloor = -3
	battery.columnsList[0].elevatorsList[2].direction = "down"
	battery.columnsList[0].elevatorsList[2].status = "moving"
	battery.columnsList[0].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, -5)

	battery.columnsList[0].elevatorsList[3].currentFloor = -6
	battery.columnsList[0].elevatorsList[3].direction = "up"
	battery.columnsList[0].elevatorsList[3].status = "moving"
	battery.columnsList[0].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[0].elevatorsList[4].currentFloor = -1
	battery.columnsList[0].elevatorsList[4].direction = "down"
	battery.columnsList[0].elevatorsList[4].status = "moving"
	battery.columnsList[0].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, -6)

	battery.columnsList[0].requestElevator(-3, "up")
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

	if battery.amountOfBasements > 0 {
		battery.makeBasementColumn(battery.amountOfBasements, battery.amountOfElevatorPerColumn)
		battery.makeBasementFloorRequestButtons(battery.amountOfBasements)
		battery.amountOfColumns--
	}

	battery.makeColumns(battery.amountOfColumns, battery.amountOfFloors, battery.amountOfElevatorPerColumn)
	battery.makeFloorRequestButtons(battery.amountOfFloors)

	//Uncomment to run scenario 1
	// battery.scenario1()

	//Uncomment to run scenario 2
	// battery.scenario2()

	//Uncomment to run scenario 3
	// battery.scenario3()

	//Uncomment to run scenario 4
	// battery.scenario4()

}
