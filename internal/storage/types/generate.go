package types

//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.Maker maker
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.MakerPieYearly maker_pie_yearly
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.MakerPieRankingYearly maker_pie_ranking_yearly
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.UserPieYearly user_pie_yearly
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.UserPieRankingYearly user_pie_ranking_yearly
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.Config config
//go:generate go run ../../../cmd/storage-generator/main.go github.com/leemartin77/mincepierank.co.uk/internal/storage/types.Admins admins
