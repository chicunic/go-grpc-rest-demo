package service

func paginate[T any](items []T, page, pageSize int32) ([]T, int32, int32, int32) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	total := int32(len(items))
	start := (page - 1) * pageSize

	if start >= total {
		return []T{}, total, page, pageSize
	}

	end := min(start+pageSize, total)

	return items[start:end], total, page, pageSize
}
