package service

import "kiddy-line-processor/internal/repo"

type KiddyLineServiceDeps struct {
	Sports map[string]repo.LineStorage
}

// type KiddyLineService struct {
// 	deps KiddyLineServiceDeps
// }
