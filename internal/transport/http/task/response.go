package task

import (
	"github.com/Nikkoz/task-service/internal/domain/task"
)

type (
	Response struct {
		ID uint64 `json:"id" binding:"required"`

		Short
	}

	Short struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		DueDate     *string `json:"due_date"`
	}

	List struct {
		//Total uint64 `json:"total" default:"0"`
		Limit uint64 `json:"limit" default:"10"`
		Page  uint64 `json:"page" default:"0"`

		Data []*Response `json:"data"`
	}
)

func toResponse(task task.Task) *Response {
	short := Short{
		Title:       task.Title.String(),
		Description: task.Description.String(),
	}

	if task.DueDate != nil {
		dueDate := task.DueDate.String()
		short.DueDate = &dueDate
	}

	return &Response{
		ID:    task.ID,
		Short: short,
	}
}

func toListResponse(limit, page uint64, tasks []task.Task) List {
	list := List{
		//Total: count,
		Limit: limit,
		Page:  page,
		Data:  []*Response{},
	}

	for _, value := range tasks {
		list.Data = append(list.Data, toResponse(value))
	}

	return list
}
