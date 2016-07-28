// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"fmt"

	"github.com/aporeto-inc/elemental"
)

// ListIdentity represents the Identity of the object
var ListIdentity = elemental.Identity{
	Name:     "list",
	Category: "lists",
}

// ListsList represents a list of Lists
type ListsList []*List

// List represents the model of a list
type List struct {
	ID          string `json:"ID,omitempty"`
	ParentID    string `json:"parentID,omitempty"`
	ParentType  string `json:"parentType,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
}

// NewList returns a new *List
func NewList() *List {

	return &List{}
}

// Identity returns the Identity of the object.
func (o *List) Identity() elemental.Identity {

	return ListIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *List) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *List) SetIdentifier(ID string) {

	o.ID = ID
}

// Validate valides the current information stored into the structure.
func (o *List) Validate() elemental.Errors {

	errors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// TaskIdentity represents the Identity of the object
var TaskIdentity = elemental.Identity{
	Name:     "task",
	Category: "tasks",
}

// TasksList represents a list of Tasks
type TasksList []*Task

// Task represents the model of a task
type Task struct {
	ID          string `json:"ID,omitempty"`
	ParentID    string `json:"parentID,omitempty"`
	ParentType  string `json:"parentType,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Description string `json:"description,omitempty"`
	Name        string `json:"name,omitempty"`
	Status      string `json:"status,omitempty"`
}

// NewTask returns a new *Task
func NewTask() *Task {

	return &Task{
		Status: "TODO",
	}
}

// Identity returns the Identity of the object.
func (o *Task) Identity() elemental.Identity {

	return TaskIdentity
}

// Identifier returns the value of the object's unique identifier.
func (o *Task) Identifier() string {

	return o.ID
}

// SetIdentifier sets the value of the object's unique identifier.
func (o *Task) SetIdentifier(ID string) {

	o.ID = ID
}

// Validate valides the current information stored into the structure.
func (o *Task) Validate() elemental.Errors {

	errors := elemental.Errors{}

	if err := elemental.ValidateRequiredString("name", o.Name); err != nil {
		errors = append(errors, err)
	}

	if err := elemental.ValidateStringInList("status", o.Status, []string{"DONE", "PROGRESS", "TODO"}); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// UnmarshalableListIdentity is mock
var UnmarshalableListIdentity = elemental.Identity{Name: "list", Category: "lists"}

// UnmarshalableList is mock.
type UnmarshalableList struct {
	List
}

// NewUnmarshalableList returns a new UnmarshalableListIdentity
func NewUnmarshalableList() *UnmarshalableList {
	return &UnmarshalableList{List: List{}}
}

// Identity is mock.
func (o *UnmarshalableList) Identity() elemental.Identity { return UnmarshalableListIdentity }

// UnmarshalJSON is mock.
func (o *UnmarshalableList) UnmarshalJSON([]byte) error {
	return fmt.Errorf("error unmarshalling")
}

// MarshalJSON is mock.
func (o *UnmarshalableList) MarshalJSON() ([]byte, error) {
	return nil, fmt.Errorf("error marshalling")
}

// Validate is mock.
func (o *UnmarshalableList) Validate() elemental.Errors { return nil }