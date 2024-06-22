package data

import "database/sql"

type ChildModel struct {
	DB *sql.DB
}

type Child struct {
	Name    string
	Dob     string //yymmdd
	UserId  string
	ChildId string
}

// Add a placeholder method for inserting a new record in the Children table.
func (m ChildModel) Insert(Child *Child) error {
	return nil
}

// Add a placeholder method for fetching a specific record from the Children table.
func (m ChildModel) Get(child_id string) (*Child, error) {
	return nil, nil
}

// Add a placeholder method for updating a specific record in the Children table.
func (m ChildModel) Update(Child *Child) error {
	return nil
}

// Add a placeholder method for deleting a specific record from the Childs table.
func (m ChildModel) Delete(child_id string) error {
	return nil
}
