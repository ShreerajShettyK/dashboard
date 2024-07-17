package models

type EC2Instance struct {
	InstanceId       string  `json:"InstanceId"`
	LastActivity     string  `json:"LastActivity"`
	LastActivityDays int     `json:"LastActivityDays"`
	Cost             float64 `json:"Cost"`
	Region           string  `json:"Region"`
	InstanceType     string  `json:"InstanceType"`
}

type AWSEfficiencyViewData struct {
	Instances       []string
	InstanceDetails *EC2Instance
}
