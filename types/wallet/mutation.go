package types

import "time"

type (
	Mutation struct {
		ID          string    `db:"id"`
		ReferenceID string    `db:"reference_id"`
		CreatedAt   time.Time `db:"created_at"`
		CreatedBy   string    `db:"created_by"`
		Action      int       `db:"action"`
		Status      int       `db:"status"`
		Amount      float64   `db:"amount"`
	}
)

type MutationStatus int
type MutationAction int

const (
	MutationStatusSuccess MutationStatus = iota + 1
	MutationStatusFailed
)

const (
	MutationActionDeposit MutationAction = iota + 1
	MutationActionWithdraw
)

const (
	MutationStatusSuccessString string = "success"
	MutationStatusFailString    string = "fail"
)

var (
	MutationStatusMap = map[MutationStatus]string{
		MutationStatusSuccess: MutationStatusSuccessString,
		MutationStatusFailed:  MutationStatusFailString,
	}
)

func (m *Mutation) GetStatusString() string {
	return MutationStatusMap[MutationStatus(m.Status)]
}
