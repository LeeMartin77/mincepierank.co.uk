package storage

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/sqlcgen"
	"github.com/leemartin77/mincepierank.co.uk/internal/storage/types"
)

// UserPieWithOwner is a user pie that includes the owner's userid
// (only populated for admin viewers).
type UserPieWithOwner struct {
	Pie types.UserPieYearly
	// OwnerID is the submitter's email. Only set for admin viewers.
	OwnerID string
}

// AdminUserPieRow is returned by GetUserPieAdminQueue.
type AdminUserPieRow struct {
	Year           int32
	ID             string
	OwnerUserid    string
	Maker          string
	Location       string
	DisplayName    string
	Fresh          bool
	Labels         []string
	ImageFile      string
	WebLink        string
	PackCount      int32
	PackPriceInPence int32
	Clean          bool
	Status         string
	ModerationNote string
	DuplicateOf    string
	PromotedMakerPieO string
	SubmittedBy    string
	SubmittedAt    *time.Time
	ModeratedBy    string
	ModeratedAt    *time.Time
}

// CreateUserPie creates a new user pie submission and returns it.
func (o *OperationWrapper) CreateUserPie(c context.Context, pie types.UserPieYearly) (*types.UserPieYearly, error) {
	q := o.querier

	result, err := q.CreateUserPie(c, sqlcgen.CreateUserPieParams{
		Year:             pie.Year,
		ID:               pie.Id,
		OwnerUserid:      pie.OwnerUserid,
		Maker:            pie.Maker,
		Location:         pie.Location,
		Displayname:      pie.DisplayName,
		Fresh:            pie.Fresh,
		Labels:           pie.Labels,
		ImageFile:        pie.ImageFile,
		WebLink:          pie.WebLink,
		PackCount:        pie.PackCount,
		PackPriceInPence: pie.PackPriceInPence,
		Clean:            pie.Clean,
		SubmittedBy:      pie.SubmittedBy,
		DuplicateOf:      pie.DuplicateOf,
	})
	if err != nil {
		return nil, err
	}
	return fromSQLCUserPie(result), nil
}

// GetUserPie retrieves a user pie by year and id.
func (o *OperationWrapper) GetUserPie(c context.Context, year int32, id string) (*UserPieWithOwner, error) {
	q := o.querier

	result, err := q.GetUserPie(c, sqlcgen.GetUserPieParams{Year: year, ID: id})
	if err != nil {
		return nil, err
	}
	return &UserPieWithOwner{Pie: *fromSQLCUserPie(result)}, nil
}

// GetUserPieByOwner retrieves a user pie owned by a specific owner.
func (o *OperationWrapper) GetUserPieByOwner(c context.Context, year int32, id, ownerID string) (*types.UserPieYearly, error) {
	q := o.querier

	result, err := q.GetUserPieByOwner(c, sqlcgen.GetUserPieByOwnerParams{
		Year:        year,
		ID:          id,
		OwnerUserid: ownerID,
	})
	if err != nil {
		return nil, err
	}
	return fromSQLCUserPie(result), nil
}

// UpdateUserPie updates a user pie's fields.
func (o *OperationWrapper) UpdateUserPie(c context.Context, pie types.UserPieYearly) error {
	return o.querier.UpdateUserPie(c, sqlcgen.UpdateUserPieParams{
		Maker:            pie.Maker,
		Location:         pie.Location,
		Displayname:      pie.DisplayName,
		Fresh:            pie.Fresh,
		Labels:           pie.Labels,
		ImageFile:        pie.ImageFile,
		WebLink:          pie.WebLink,
		PackCount:        pie.PackCount,
		PackPriceInPence: pie.PackPriceInPence,
		Year:             pie.Year,
		ID:               pie.Id,
	})
}

// UpdateUserPieImage updates only the image file URL.
func (o *OperationWrapper) UpdateUserPieImage(c context.Context, year int32, id, imageFile string) error {
	return o.querier.UpdateUserPieImage(c, sqlcgen.UpdateUserPieImageParams{
		Year:      year,
		ID:        id,
		ImageFile: imageFile,
	})
}

// DeleteUserPie deletes a user pie by year and id.
func (o *OperationWrapper) DeleteUserPie(c context.Context, year int32, id string) error {
	return o.querier.DeleteUserPie(c, sqlcgen.DeleteUserPieParams{Year: year, ID: id})
}

// DeleteUserPieByOwner deletes a user pie only if it belongs to the owner.
func (o *OperationWrapper) DeleteUserPieByOwner(c context.Context, year int32, id, ownerID string) error {
	return o.querier.DeleteUserPieByOwner(c, sqlcgen.DeleteUserPieByOwnerParams{
		Year:        year,
		ID:          id,
		OwnerUserid: ownerID,
	})
}

// GetUserPieYearIDs returns all pie IDs for a given year.
func (o *OperationWrapper) GetUserPieYearIDs(c context.Context, year int32) ([]string, error) {
	return o.querier.GetUserPieYearIDs(c, year)
}

// GetUserPieCountByYear returns the total user pie count for a year.
func (o *OperationWrapper) GetUserPieCountByYear(c context.Context, year int32) (int64, error) {
	return o.querier.GetUserPieCountByYear(c, year)
}

// GetUserPieCountByOwnerAndYear returns the count of pies owned by a user in a year.
func (o *OperationWrapper) GetUserPieCountByOwnerAndYear(c context.Context, year int32, ownerID string) (int64, error) {
	return o.querier.GetUserPieCountByOwnerAndYear(c, sqlcgen.GetUserPieCountByOwnerAndYearParams{
		Year:        year,
		OwnerUserid: ownerID,
	})
}

// SetUserPieStatus atomically transitions a user pie's status if the current
// status permits the transition. Returns ErrInvalidTransition when the
// transition is not allowed.
func (o *OperationWrapper) SetUserPieStatus(c context.Context, year int32, id string, fromStatus, toStatus UserPieStatus) error {
	if !fromStatus.CanTransitionTo(toStatus) {
		return ErrInvalidTransition
	}

	query := `
		UPDATE user_pie_yearly
		SET status = $4,
		    moderated_by = COALESCE(moderated_by, ''),
		    moderated_at = NOW()
		WHERE year = $1::int
		  AND id = $2::text
		  AND status = $3::text
	`
	tx, err := o.db.Begin(c)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(c)
	}()

	result, err := tx.Exec(c, query, year, id, fromStatus, toStatus)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return ErrInvalidTransition
	}

	return tx.Commit(c)
}

// GetUserPieAdminQueue returns user pies for the admin moderation queue,
// limited and offset by limit/idx.
func (o *OperationWrapper) GetUserPieAdminQueue(c context.Context, year int32, limit int, idx int) ([]AdminUserPieRow, error) {
	query := `
		SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
		       image_file, web_link, pack_count, pack_price_in_pence, clean,
		       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
		       submitted_by, submitted_at, moderated_by, moderated_at
		FROM user_pie_yearly
		WHERE year = $1::int
		ORDER BY submitted_at DESC, id ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := o.db.Query(c, query, year, limit, idx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []AdminUserPieRow
	for rows.Next() {
		var row AdminUserPieRow
		err := rows.Scan(
			&row.Year, &row.ID, &row.OwnerUserid, &row.Maker, &row.Location,
			&row.DisplayName, &row.Fresh, &row.Labels, &row.ImageFile,
			&row.WebLink, &row.PackCount, &row.PackPriceInPence, &row.Clean,
			&row.Status, &row.ModerationNote, &row.DuplicateOf,
			&row.PromotedMakerPieO, &row.SubmittedBy, &row.SubmittedAt,
			&row.ModeratedBy, &row.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// GetUserPieAdminQueueCount returns the total count of user pies for a year.
func (o *OperationWrapper) GetUserPieAdminQueueCount(c context.Context, year int32) (int64, error) {
	var count int64
	err := o.db.QueryRow(c, "SELECT COUNT(*) FROM user_pie_yearly WHERE year = $1::int", year).Scan(&count)
	return count, err
}

// GetUserPieDuplicateCandidates returns user pie candidates that share a maker
// with the given pie (for dedupe review).
func (o *OperationWrapper) GetUserPieDuplicateCandidates(c context.Context, year int32, pieID, maker string) ([]types.UserPieYearly, error) {
	query := `
		SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
		       image_file, web_link, pack_count, pack_price_in_pence, clean,
		       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
		       submitted_by, submitted_at, moderated_by, moderated_at
		FROM user_pie_yearly
		WHERE year = $1::int
		  AND maker = $2::text
		  AND id != $3::text
		  AND status IN ('pending', 'approved')
		ORDER BY submitted_at ASC
	`
	rows, err := o.db.Query(c, query, year, maker, pieID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.UserPieYearly
	for rows.Next() {
		var pie types.UserPieYearly
		err := rows.Scan(
			&pie.Year, &pie.Id, &pie.OwnerUserid, &pie.Maker, &pie.Location,
			&pie.DisplayName, &pie.Fresh, &pie.Labels, &pie.ImageFile,
			&pie.WebLink, &pie.PackCount, &pie.PackPriceInPence, &pie.Clean,
			&pie.Status, &pie.ModerationNote, &pie.DuplicateOf,
			&pie.PromotedMakerPieO, &pie.SubmittedBy, &pie.SubmittedAt,
			&pie.ModeratedBy, &pie.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// UpsertUserPieRanking upserts a ranking for a user pie.
func (o *OperationWrapper) UpsertUserPieRanking(c context.Context, ranking types.UserPieRankingYearly) error {
	return o.querier.UpsertUserPieRanking(c, sqlcgen.UpsertUserPieRankingParams{
		Year:    ranking.Year,
		Pieid:   ranking.PieId,
		Userid:  ranking.UserId,
		Pastry:  int32(ranking.Pastry),
		Filling: int32(ranking.Filling),
		Topping: int32(ranking.Topping),
		Looks:   int32(ranking.Looks),
		Value:   int32(ranking.Value),
		Notes:   ranking.Notes,
	})
}

// GetUserPieRanking retrieves a user's ranking for a user pie.
func (o *OperationWrapper) GetUserPieRanking(c context.Context, year int32, pieID, userID string) (*types.UserPieRankingYearly, error) {
	r, err := o.querier.GetUserPieRanking(c, sqlcgen.GetUserPieRankingParams{
		Year:   year,
		Pieid:  pieID,
		Userid: userID,
	})
	if err != nil {
		return nil, err
	}
	return &types.UserPieRankingYearly{
		Year:    r.Year,
		PieId:   r.Pieid,
		UserId:  r.Userid,
		Pastry:  int8(r.Pastry.Int32),
		Filling: int8(r.Filling.Int32),
		Topping: int8(r.Topping.Int32),
		Looks:   int8(r.Looks.Int32),
		Value:   int8(r.Value.Int32),
		Notes:   r.Notes.String,
		LastUpdated: r.LastUpdated.Time,
	}, nil
}

// DeleteUserPieRankings deletes all rankings for a user pie.
func (o *OperationWrapper) DeleteUserPieRankings(c context.Context, year int32, pieID string) error {
	return o.querier.DeleteUserPieRankings(c, sqlcgen.DeleteUserPieRankingsParams{
		Year:  year,
		Pieid: pieID,
	})
}

// DeleteUserPieRankingsByUser deletes all rankings by a user for a year.
func (o *OperationWrapper) DeleteUserPieRankingsByUser(c context.Context, year int32, userID string) error {
	return o.querier.DeleteUserPieRankingsByUser(c, sqlcgen.DeleteUserPieRankingsByUserParams{
		Year:   year,
		Userid: userID,
	})
}

// GetUserPieRankingSummary returns the average ranking scores for a user pie.
func (o *OperationWrapper) GetUserPieRankingSummary(c context.Context, year int32, pieID string) (*types.RankingSummary, error) {
	query := `
		SELECT COALESCE(AVG(pastry), 0),
		       COALESCE(AVG(filling), 0),
		       COALESCE(AVG(topping), 0),
		       COALESCE(AVG(looks), 0),
		       COALESCE(AVG(value), 0),
		       COALESCE(AVG((pastry + filling + topping + looks + value)::float / 5.0), 0),
		       COUNT(*)
		FROM user_pie_ranking_yearly
		WHERE year = $1::int AND pieid = $2::text
	`
	var s types.RankingSummary
	err := o.db.QueryRow(c, query, year, pieID).Scan(
		&s.Pastry, &s.Filling, &s.Topping, &s.Looks, &s.Value,
		&s.Average, &s.Count,
	)
	return &s, err
}

// GetAllApprovedUserPies returns paginated approved user pies for a year.
func (o *OperationWrapper) GetAllApprovedUserPies(c context.Context, year int32, limit int, idx int) ([]types.UserPieYearly, error) {
	query := `
		SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
		       image_file, web_link, pack_count, pack_price_in_pence, clean,
		       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
		       submitted_by, submitted_at, moderated_by, moderated_at
		FROM user_pie_yearly
		WHERE year = $1::int AND status IN ('approved', 'promoted')
		ORDER BY submitted_at DESC, id ASC
		LIMIT $2 OFFSET $3
	`
	rows, err := o.db.Query(c, query, year, limit, idx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.UserPieYearly
	for rows.Next() {
		var pie types.UserPieYearly
		err := rows.Scan(
			&pie.Year, &pie.Id, &pie.OwnerUserid, &pie.Maker, &pie.Location,
			&pie.DisplayName, &pie.Fresh, &pie.Labels, &pie.ImageFile,
			&pie.WebLink, &pie.PackCount, &pie.PackPriceInPence, &pie.Clean,
			&pie.Status, &pie.ModerationNote, &pie.DuplicateOf,
			&pie.PromotedMakerPieO, &pie.SubmittedBy, &pie.SubmittedAt,
			&pie.ModeratedBy, &pie.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// GetAllApprovedUserPiesCount returns the total count of approved user pies for a year.
func (o *OperationWrapper) GetAllApprovedUserPiesCount(c context.Context, year int32) (int64, error) {
	var count int64
	err := o.db.QueryRow(c,
		"SELECT COUNT(*) FROM user_pie_yearly WHERE year = $1::int AND status IN ('approved', 'promoted')",
		year).Scan(&count)
	return count, err
}

// GetOwnerUserPies returns all pies owned by a user for a year.
func (o *OperationWrapper) GetOwnerUserPies(c context.Context, year int32, ownerID string) ([]types.UserPieYearly, error) {
	query := `
		SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
		       image_file, web_link, pack_count, pack_price_in_pence, clean,
		       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
		       submitted_by, submitted_at, moderated_by, moderated_at
		FROM user_pie_yearly
		WHERE year = $1::int AND owner_userid = $2::text
		ORDER BY submitted_at DESC, id ASC
	`
	rows, err := o.db.Query(c, query, year, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.UserPieYearly
	for rows.Next() {
		var pie types.UserPieYearly
		err := rows.Scan(
			&pie.Year, &pie.Id, &pie.OwnerUserid, &pie.Maker, &pie.Location,
			&pie.DisplayName, &pie.Fresh, &pie.Labels, &pie.ImageFile,
			&pie.WebLink, &pie.PackCount, &pie.PackPriceInPence, &pie.Clean,
			&pie.Status, &pie.ModerationNote, &pie.DuplicateOf,
			&pie.PromotedMakerPieO, &pie.SubmittedBy, &pie.SubmittedAt,
			&pie.ModeratedBy, &pie.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// fromSQLCUserPie converts a sqlc-generated UserPieYearly to the types package version.
func fromSQLCUserPie(sqlc sqlcgen.UserPieYearly) *types.UserPieYearly {
	return &types.UserPieYearly{
		Year:             sqlc.Year,
		Id:               sqlc.ID,
		OwnerUserid:      sqlc.OwnerUserid,
		Maker:            sqlc.Maker.String,
		Location:         sqlc.Location.String,
		DisplayName:      sqlc.Displayname.String,
		Fresh:            sqlc.Fresh.Bool,
		Labels:           sqlc.Labels,
		ImageFile:        sqlc.ImageFile.String,
		WebLink:          sqlc.WebLink.String,
		PackCount:        sqlc.PackCount.Int32,
		PackPriceInPence: sqlc.PackPriceInPence.Int32,
		Clean:            sqlc.Clean.Bool,
		Status:           sqlc.Status,
		ModerationNote:   sqlc.ModerationNote.String,
		DuplicateOf:      sqlc.DuplicateOf.String,
		PromotedMakerPieO: sqlc.PromotedMakerPieOid.String,
		SubmittedBy:      sqlc.SubmittedBy.String,
		SubmittedAt:      sqlc.SubmittedAt.Time,
		ModeratedBy:      sqlc.ModeratedBy.String,
		ModeratedAt:      sqlc.ModeratedAt.Time,
	}
}

// GetUserPieByOid retrieves a user pie by its OID.
func (o *OperationWrapper) GetUserPieByOid(c context.Context, oid uuid.UUID) (*UserPieWithOwner, error) {
	q := o.querier

	var pgOid pgtype.UUID
	copy(pgOid.Bytes[:], oid[:])

	result, err := q.GetUserPieByOid(c, pgOid)
	if err != nil {
		return nil, err
	}
	return &UserPieWithOwner{Pie: *fromSQLCUserPie(result)}, nil
}

// DeleteUserPieByOid deletes a user pie by its OID.
func (o *OperationWrapper) DeleteUserPieByOid(c context.Context, oid uuid.UUID) error {
	var pgOid pgtype.UUID
	copy(pgOid.Bytes[:], oid[:])

	return o.querier.DeleteUserPieByOid(c, pgOid)
}

// UpdateUserPieStatusWithNote atomically updates status, moderation note, and timestamps.
func (o *OperationWrapper) UpdateUserPieStatusWithNote(c context.Context, oid uuid.UUID, status, moderatedBy, moderationNote string) error {
	var pgOid pgtype.UUID
	copy(pgOid.Bytes[:], oid[:])

	return o.querier.UpdateUserPieStatusWithNote(c, sqlcgen.UpdateUserPieStatusWithNoteParams{
		Oid:            pgOid,
		Status:         status,
		ModeratedBy:    moderatedBy,
		ModerationNote: moderationNote,
	})
}

// UpdateUserPieDuplicateOf marks a submission as a duplicate.
func (o *OperationWrapper) UpdateUserPieDuplicateOf(c context.Context, oid uuid.UUID, duplicateOf, moderatedBy string) error {
	var pgOid pgtype.UUID
	copy(pgOid.Bytes[:], oid[:])

	return o.querier.UpdateUserPieDuplicateOf(c, sqlcgen.UpdateUserPieDuplicateOfParams{
		Oid:           pgOid,
		DuplicateOf:   duplicateOf,
		ModeratedBy:   moderatedBy,
		Status:        string(UserPieStatusDuplicate),
	})
}

// UpdateUserPiePromotedMakerPieOid updates the promoted brand pie OID and status.
func (o *OperationWrapper) UpdateUserPiePromotedMakerPieOid(c context.Context, oid uuid.UUID, promotedOid uuid.UUID, moderatedBy string) error {
	var pgOid pgtype.UUID
	copy(pgOid.Bytes[:], oid[:])

	var pgPromotedOid pgtype.UUID
	copy(pgPromotedOid.Bytes[:], promotedOid[:])

	return o.querier.UpdateUserPiePromotedMakerPieOid(c, sqlcgen.UpdateUserPiePromotedMakerPieOidParams{
		Oid:                   pgOid,
		PromotedMakerPieOid:   pgPromotedOid,
		ModeratedBy:           moderatedBy,
		Status:                string(UserPieStatusPromoted),
	})
}

// GetUserPieAdminQueueCountByStatus returns the count of user pies for a year filtered by status.
func (o *OperationWrapper) GetUserPieAdminQueueCountByStatus(c context.Context, year int32, status string) (int64, error) {
	return o.querier.GetUserPieAdminQueueCountByStatus(c, sqlcgen.GetUserPieAdminQueueCountByStatusParams{
		Year:   year,
		Status: status,
	})
}

// GetUserPieDuplicateGroups returns user pies that look like duplicates of the given pie.
func (o *OperationWrapper) GetUserPieDuplicateGroups(c context.Context, year int32, excludeID, maker, displayName, webLink string) ([]types.UserPieYearly, error) {
	query := `
		SELECT year, id, owner_userid, maker, location, displayname, fresh, labels,
		       image_file, web_link, pack_count, pack_price_in_pence, clean,
		       status, moderation_note, duplicate_of, promoted_maker_pie_oid,
		       submitted_by, submitted_at, moderated_by, moderated_at
		FROM user_pie_yearly
		WHERE year = $1::int
		  AND status IN ('pending', 'approved')
		  AND id != $2::text
		  AND (
		    (maker IS NOT NULL AND maker = $3::text
		     AND lower(regexp_replace(displayname, '[^a-z0-9]+', '', 'g')) =
		         lower(regexp_replace($4::text, '[^a-z0-9]+', '', 'g')))
		    OR
		    (web_link IS NOT NULL AND web_link = $5::text AND web_link != '')
		  )
		ORDER BY submitted_at ASC
	`
	rows, err := o.db.Query(c, query, year, excludeID, maker, displayName, webLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.UserPieYearly
	for rows.Next() {
		var pie types.UserPieYearly
		err := rows.Scan(
			&pie.Year, &pie.Id, &pie.OwnerUserid, &pie.Maker, &pie.Location,
			&pie.DisplayName, &pie.Fresh, &pie.Labels, &pie.ImageFile,
			&pie.WebLink, &pie.PackCount, &pie.PackPriceInPence, &pie.Clean,
			&pie.Status, &pie.ModerationNote, &pie.DuplicateOf,
			&pie.PromotedMakerPieO, &pie.SubmittedBy, &pie.SubmittedAt,
			&pie.ModeratedBy, &pie.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// GetUserPieDuplicateGroupsAll returns all user pies that have at least one
// potential duplicate (same maker + normalised displayname, or same web_link).
func (o *OperationWrapper) GetUserPieDuplicateGroupsAll(c context.Context, year int32) ([]types.UserPieYearly, error) {
	query := `
		WITH normalised AS (
			SELECT *,
				lower(regexp_replace(displayname, '[^a-z0-9]+', '', 'g')) AS norm_name
			FROM user_pie_yearly
			WHERE year = $1::int
			  AND status IN ('pending', 'approved')
		)
		SELECT DISTINCT ON (a.maker, a.norm_name, a.displayname)
			a.year, a.id, a.owner_userid, a.maker, a.location, a.displayname,
			a.fresh, a.labels, a.image_file, a.web_link, a.pack_count,
			a.pack_price_in_pence, a.clean, a.status, a.moderation_note,
			a.duplicate_of, a.promoted_maker_pie_oid,
			a.submitted_by, a.submitted_at, a.moderated_by, a.moderated_at
		FROM normalised a
		JOIN normalised b
		  ON a.maker = b.maker
		  AND a.norm_name = b.norm_name
		  AND a.id != b.id
		ORDER BY a.maker, a.norm_name, a.submitted_at ASC
	`
	rows, err := o.db.Query(c, query, year)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.UserPieYearly
	for rows.Next() {
		var pie types.UserPieYearly
		err := rows.Scan(
			&pie.Year, &pie.Id, &pie.OwnerUserid, &pie.Maker, &pie.Location,
			&pie.DisplayName, &pie.Fresh, &pie.Labels, &pie.ImageFile,
			&pie.WebLink, &pie.PackCount, &pie.PackPriceInPence, &pie.Clean,
			&pie.Status, &pie.ModerationNote, &pie.DuplicateOf,
			&pie.PromotedMakerPieO, &pie.SubmittedBy, &pie.SubmittedAt,
			&pie.ModeratedBy, &pie.ModeratedAt,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// GetExistingBrandPiesForYear returns existing brand pies that match the given maker or web link.
func (o *OperationWrapper) GetExistingBrandPiesForYear(c context.Context, year int32, makerid, webLink string) ([]types.MakerPieYearly, error) {
	query := `
		SELECT year, makerid, id, displayname, fresh, image_file, web_link,
		       pack_count, pack_price_in_pence, oid
		FROM maker_pie_yearly
		WHERE year = $1::int
		  AND (makerid = $2::text OR web_link = $3::text)
		ORDER BY displayname
	`
	rows, err := o.db.Query(c, query, year, makerid, webLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pies []types.MakerPieYearly
	for rows.Next() {
		var pie types.MakerPieYearly
		err := rows.Scan(
			&pie.Year, &pie.MakerId, &pie.Id, &pie.DisplayName, &pie.Fresh,
			&pie.ImageFile, &pie.WebLink,
			&pie.PackCount, &pie.PackPriceInPence, &pie.OId,
		)
		if err != nil {
			return nil, err
		}
		pies = append(pies, pie)
	}
	return pies, rows.Err()
}

// ErrInvalidTransition is returned when a status change is not allowed.
var ErrInvalidTransition = errors.New("invalid status transition")
