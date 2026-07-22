package project

import (
	"errors"
	"testing"
)

func TestServiceCreateSuccess(t *testing.T) {
	repository := NewMemoryRepository()
	Service := NewService(repository)

	request := CreateRequest{
		Name:        "Task",
		Description: "A application like jira",
	}

	createdProject, err := Service.Create(request)
	if err != nil {
		t.Fatalf("Create() returned an unexpected error : %v", err)
	}

	if createdProject.ID != 1 {
		t.Errorf("expected project id 1, got %d", createdProject.ID)
	}
	if createdProject.Name != "Task" {
		t.Errorf(
			"expected project name %q, got %q",
			"Task",
			createdProject.Name,
		)
	}

	if createdProject.Description != request.Description {
		t.Errorf(
			"expected description %q, got %q",
			request.Description,
			createdProject.Description,
		)
	}

	if createdProject.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	if createdProject.UpdatedAt.IsZero() {
		t.Error("expected UpdatedAt to be set")
	}

	if !createdProject.CreatedAt.Equal(createdProject.UpdatedAt) {
		t.Errorf(
			"expected timestamps to be equal, CreatedAt=%v UpdatedAt=%v",
			createdProject.CreatedAt,
			createdProject.UpdatedAt,
		)
	}
}

func TestServiceCreateNameRequired(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)

	request := CreateRequest{
		Name:        "   ",
		Description: "Invalid project",
	}

	createdProject, err := service.Create(request)

	if !errors.Is(err, ErrNameRequired) {
		t.Fatalf(
			"expected ErrNameRequired, got %v",
			err,
		)
	}

	if createdProject != (Project{}) {
		t.Errorf(
			"expected an empty project, got %+v",
			createdProject,
		)
	}
}
