package ws

// SubscriptionMsg is used to create a WS subscription at CCData
type SubscriptionMsg struct {
	Action       Action           `json:"action"`
	Type         SubscriptionType `json:"type"`
	Groups       []Group          `json:"groups"`
	Subcriptions []Subscription   `json:"subscriptions,omitempty"`
}

// Action contains the subscription request action
type Action string

// SubscriptionType is used to specify the type of subscription
type SubscriptionType string

// Group is used to specify a required group in the response
type Group string

const (
	AddSubscription    Action           = "SUB_ADD"
	RemoveSubscription Action           = "SUB_REMOVE"
	Value              Group            = "VALUE"
	CurrentHour        Group            = "CURRENT_HOUR"
	CADLITick          SubscriptionType = "1101"
)

// Subscription contains a specific market & instrument
type Subscription struct {
	Market     string `json:"market"`
	Instrument string `json:"instrument"`
}

// DefaultCADLIGroups contains two groups used in CADLI subscriptions
var DefaultCADLIGroups []Group = []Group{
	Value,
	CurrentHour,
}
