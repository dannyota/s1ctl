package cli

import (
	"danny.vn/s1/graphql"
	"danny.vn/s1/mgmt"
)

func fetchAllREST[T any](resource string, fn func(cursor string) ([]T, *mgmt.Pagination, error)) ([]T, int, error) {
	var all []T
	var cursor string
	var total int
	for {
		items, pag, err := fn(cursor)
		if err != nil {
			clearProgress()
			return nil, 0, err
		}
		all = append(all, items...)
		if pag != nil {
			total = pag.TotalItems
		}
		printProgress(resource, len(all), total)
		if pag == nil || pag.NextCursor == "" {
			break
		}
		cursor = pag.NextCursor
	}
	clearProgress()
	return all, total, nil
}

func fetchAllGQL[T any](resource string, fn func(after string) (*graphql.Connection[T], error)) ([]T, int64, error) {
	var all []T
	var after string
	var total int64
	for {
		conn, err := fn(after)
		if err != nil {
			clearProgress()
			return nil, 0, err
		}
		total = conn.TotalCount
		for _, e := range conn.Edges {
			all = append(all, e.Node)
		}
		printProgress(resource, len(all), int(total))
		if !conn.PageInfo.HasNextPage {
			break
		}
		after = conn.PageInfo.EndCursor
	}
	clearProgress()
	return all, total, nil
}
