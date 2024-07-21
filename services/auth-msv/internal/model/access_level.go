package model

type AccessLevel int

const (
	Undefined AccessLevel = iota
	Admin
	WarehouseWorker
	Courier
	Buyer
)

func (lvl AccessLevel) IsValid() bool {
	return lvl >= Admin && lvl <= Buyer
}

func (lvl AccessLevel) String() string {
	if !lvl.IsValid() {
		return ""
	}

	return []string{"Undefined", "Admin", "WarehouseWorker", "Courier", "Buyer"}[lvl]
}
