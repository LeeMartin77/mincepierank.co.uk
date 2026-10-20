package storage

import "testing"

func TestUserPieCanTransitionTo(t *testing.T) {
	tests := []struct {
		from   UserPieStatus
		to     UserPieStatus
		allow  bool
	}{
		// pending → everything
		{UserPieStatusPending, UserPieStatusApproved, true},
		{UserPieStatusPending, UserPieStatusPromoted, true},
		{UserPieStatusPending, UserPieStatusDuplicate, true},
		{UserPieStatusPending, UserPieStatusRejected, true},
		{UserPieStatusPending, UserPieStatusPending, false},

		// approved → promoted / duplicate / rejected
		{UserPieStatusApproved, UserPieStatusPromoted, true},
		{UserPieStatusApproved, UserPieStatusDuplicate, true},
		{UserPieStatusApproved, UserPieStatusRejected, true},
		{UserPieStatusApproved, UserPieStatusPending, false},
		{UserPieStatusApproved, UserPieStatusApproved, false},

		// terminal states — no transitions out
		{UserPieStatusPromoted, UserPieStatusRejected, false},
		{UserPieStatusPromoted, UserPieStatusPending, false},
		{UserPieStatusPromoted, UserPieStatusPromoted, false},
		{UserPieStatusDuplicate, UserPieStatusPending, false},
		{UserPieStatusDuplicate, UserPieStatusRejected, false},
		{UserPieStatusDuplicate, UserPieStatusDuplicate, false},
		{UserPieStatusRejected, UserPieStatusPending, false},
		{UserPieStatusRejected, UserPieStatusApproved, false},
		{UserPieStatusRejected, UserPieStatusRejected, false},
	}

	for _, tc := range tests {
		t.Run(string(tc.from)+"→"+string(tc.to), func(t *testing.T) {
			got := tc.from.CanTransitionTo(tc.to)
			if got != tc.allow {
				t.Errorf("CanTransitionTo(%s) = %v, want %v", tc.to, got, tc.allow)
			}
		})
	}
}

func TestUserPieIsPublic(t *testing.T) {
	tests := []struct {
		s   UserPieStatus
		pub bool
	}{
		{UserPieStatusPending, false},
		{UserPieStatusApproved, true},
		{UserPieStatusPromoted, true},
		{UserPieStatusDuplicate, false},
		{UserPieStatusRejected, false},
	}
	for _, tc := range tests {
		t.Run(string(tc.s), func(t *testing.T) {
			if got := tc.s.IsPublic(); got != tc.pub {
				t.Errorf("IsPublic(%s) = %v, want %v", tc.s, got, tc.pub)
			}
		})
	}
}

func TestUserPieIsOwnerEditable(t *testing.T) {
	tests := []struct {
		s          UserPieStatus
		editable   bool
	}{
		{UserPieStatusPending, true},
		{UserPieStatusApproved, false},
		{UserPieStatusPromoted, false},
		{UserPieStatusDuplicate, false},
		{UserPieStatusRejected, false},
	}
	for _, tc := range tests {
		t.Run(string(tc.s), func(t *testing.T) {
			if got := tc.s.IsOwnerEditable(); got != tc.editable {
				t.Errorf("IsOwnerEditable(%s) = %v, want %v", tc.s, got, tc.editable)
			}
		})
	}
}

func TestUserPieIsRankable(t *testing.T) {
	tests := []struct {
		s         UserPieStatus
		rankable  bool
	}{
		{UserPieStatusPending, false},
		{UserPieStatusApproved, true},
		{UserPieStatusPromoted, false},
		{UserPieStatusDuplicate, false},
		{UserPieStatusRejected, false},
	}
	for _, tc := range tests {
		t.Run(string(tc.s), func(t *testing.T) {
			if got := tc.s.IsRankable(); got != tc.rankable {
				t.Errorf("IsRankable(%s) = %v, want %v", tc.s, got, tc.rankable)
			}
		})
	}
}

func TestUserPieIsTerminal(t *testing.T) {
	tests := []struct {
		s        UserPieStatus
		terminal bool
	}{
		{UserPieStatusPending, false},
		{UserPieStatusApproved, false},
		{UserPieStatusPromoted, true},
		{UserPieStatusDuplicate, true},
		{UserPieStatusRejected, true},
	}
	for _, tc := range tests {
		t.Run(string(tc.s), func(t *testing.T) {
			if got := tc.s.IsTerminal(); got != tc.terminal {
				t.Errorf("IsTerminal(%s) = %v, want %v", tc.s, got, tc.terminal)
			}
		})
	}
}
