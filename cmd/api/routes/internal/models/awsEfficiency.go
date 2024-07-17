package models

type EC2Instance struct {
	InstanceId       string  `json:"InstanceId"`
	LastActivity     string  `json:"LastActivity"`
	LastActivityDays int     `json:"LastActivityDays"`
	Cost             float64 `json:"Cost"`
}

type AWSEfficiencyViewData struct {
	Instances       []string
	InstanceDetails *EC2Instance
}
