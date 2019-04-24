package common

const (
	OutputTypeMySQL  = "mysql"
	OutputTypeCSV    = "csv"
	OutputTypeStdout = "stdout"
)

type MTS struct {
	ID     uint64
	Status TaskStatus
}
