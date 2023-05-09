package rotation

var AllowedSteps = []string{"createSecret", "setSecret", "testSecret", "finishSecret"}

type StagingLabels struct {
	Current  string
	Pending  string
	Previous string
}

type Steps struct {
	Create string
	Set    string
	Test   string
	Finish string
}

func GetStagingLabels() StagingLabels {
	return StagingLabels{
		Current:  "AWSCURRENT",
		Pending:  "AWSPENDING",
		Previous: "AWSPREVIOUS",
	}
}

func GetSteps() Steps {
	return Steps{
		Create: "createSecret",
		Set:    "setSecret",
		Test:   "testSecret",
		Finish: "finishSecret",
	}
}
