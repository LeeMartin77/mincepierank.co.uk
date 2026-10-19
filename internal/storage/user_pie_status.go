package storage

// UserPieStatus represents the moderation status of a user-submitted pie.
type UserPieStatus string

const (
	UserPieStatusPending    UserPieStatus = "pending"     // submitted, awaiting review
	UserPieStatusApproved   UserPieStatus = "approved"    // published as a user pie
	UserPieStatusPromoted   UserPieStatus = "promoted"    // became a brand pie
	UserPieStatusDuplicate  UserPieStatus = "duplicate"   // collapsed into another submission
	UserPieStatusRejected   UserPieStatus = "rejected"    // removed by admin
)

// The transition table. Anything not listed here is refused.
//
// promoted, duplicate and rejected are terminal: to remove a submission
// that has reached one of them, the admin deletes it.
var userPieTransitions = map[UserPieStatus][]UserPieStatus{
	UserPieStatusPending: {
		UserPieStatusApproved, UserPieStatusPromoted,
		UserPieStatusDuplicate, UserPieStatusRejected,
	},
	UserPieStatusApproved: {
		UserPieStatusPromoted, UserPieStatusDuplicate, UserPieStatusRejected,
	},
	UserPieStatusPromoted:  {}, // terminal - removal is a delete
	UserPieStatusDuplicate: {}, // terminal - removal is a delete
	UserPieStatusRejected:  {}, // terminal - a fresh submission is required
}

// CanTransitionTo reports whether a submission in status s may move to next.
func (s UserPieStatus) CanTransitionTo(next UserPieStatus) bool {
	for _, t := range userPieTransitions[s] {
		if t == next {
			return true
		}
	}
	return false
}

// IsPublic reports whether the status is visible on public pages.
func (s UserPieStatus) IsPublic() bool {
	return s == UserPieStatusApproved || s == UserPieStatusPromoted
}

// IsOwnerEditable reports whether the owner may edit or withdraw the submission.
func (s UserPieStatus) IsOwnerEditable() bool {
	return s == UserPieStatusPending
}

// IsRankable reports whether the status accepts user rankings.
func (s UserPieStatus) IsRankable() bool {
	return s == UserPieStatusApproved
}

// IsTerminal reports whether the status is settled. Terminal statuses have no
// outgoing transitions: to remove a submission that has reached one, the admin
// deletes it.
func (s UserPieStatus) IsTerminal() bool {
	switch s {
	case UserPieStatusPromoted, UserPieStatusDuplicate, UserPieStatusRejected:
		return true
	}
	return false
}
