package common

const (
	OutputTypeMySQL = "mysql"
	OutputTypeCSV   = "csv"
)

type MTS struct {
	ID     uint64
	Status TaskStatus
}
