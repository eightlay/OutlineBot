package bot

const (
	stateIDLE uint = iota
	statePAY
	stateHELP
)

type user struct {
	OutlineID   string `json:"olID"`
	Admin       bool   `json:"admin"`
	PaymentCode string `json:"pc"`
	State       uint   `json:"state"`
}
