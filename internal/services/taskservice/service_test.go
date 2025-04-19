package taskservice_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/tehrelt/workmate-testovoe/internal/models"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice"
	"github.com/tehrelt/workmate-testovoe/internal/services/taskservice/mocks"
)

func TestCreateTask(t *testing.T) {
	type Case struct {
		name    string
		title   string
		wantErr bool
	}
	cases := []Case{
		Case{
			name:    "base",
			title:   "Test Task",
			wantErr: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ctx := context.Background()

			tasksaver := mocks.NewTaskSaver(t)
			taskprovider := mocks.NewTaskProvider(t)
			taskprocessor := mocks.NewTaskProcessor(t)

			input := &models.CreateTask{
				Title: c.title,
			}

			savedTask := &models.Task{
				Id:    uuid.New(),
				Title: c.title,
			}
			tasksaver.
				On("Save", mock.Anything, input).
				Return(
					savedTask,
					nil,
				)

			taskprocessor.
				On("Push", mock.Anything, savedTask).
				Return(nil)

			service := taskservice.New(tasksaver, taskprovider, taskprocessor)

			task, err := service.CreateTask(ctx, input)
			if err != nil && !c.wantErr {
				t.Errorf("unexpected error: %v", err)
			}

			assert.NotEqual(t, uuid.Nil, task.Id)
			assert.Equal(t, c.title, task.Title)
		})
	}

}
