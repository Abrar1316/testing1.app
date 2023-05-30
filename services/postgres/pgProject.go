package postgres

import (
	"log"
	"time"

	"github.com/TFTPL/AWS-Cost-Calculator/services/types"
	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
)

// CreateProject creates a new project
func CreateProject(name, description string) (*types.DbProject, error) {
	project := &types.DbProject{
		Name:        name,
		Description: description,
	}
	err := Insert(project)

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:CreateProject failed - could not insert row")
	}

	return project, nil
}

// ReadProject reads a project
func ReadProject(id int64) (*types.DbProject, error) {
	project := &types.DbProject{
		ID: id,
	}

	err := getDB().Model(project).Where("id = ?", id).Select()

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:ReadProject failed - could not query row")
	}

	return project, nil
}

// UpdateProject updates a project
func UpdateProject(id int64, name, description string) (*types.DbProject, error) {

	project := &types.DbProject{
		ID:          id,
		Name:        name,
		Description: description,
		UpdatedAt:   time.Now(),
	}

	_, err := GetDB().Model(project).
		Set("name = ?name").
		Set("description = ?description").
		Set("updated_at=?updated_at").
		Where("id = ?id").
		Update()

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:UpdateProject failed - could not update row")
	}

	return project, nil
}

// DeleteProject deletes a project
func DeleteProject(id int64) (*types.DbProject, error) {
	project := &types.DbProject{
		ID: id,
	}

	err := Delete(project)

	if err != nil {
		return nil, errors.Wrap(err, "pgProject:DeleteProject failed - could not delete row")
	}

	return project, nil
}

func GetBulkProjectsByIds(ids []int64) ([]types.DbProject, error) {
	res := make([]types.DbProject, 0)

	query := GetDB().
		Model(&types.DbProject{}).
		Where("id = any(?)", pg.Array(ids))

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}

func GetProjectsByIsActive(isActive bool) ([]*types.DbProject, error) {
	res := make([]*types.DbProject, 0)

	query := GetDB().
		Model(&types.DbProject{}).
		Where("is_active = ?", isActive)

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

// func update project by pin
func UpdateProjectbyPin(id int64, isPinned bool) error {

	stmt := `UPDATE projects SET is_pinned = ? WHERE id = ?`
	_, err := GetDB().Exec(stmt, isPinned, id)
	if err != nil {
		log.Fatalf("Unable to execute query: %v\n", err)
	}

	return nil
}

// get all pinned projects
func GetAllPinnedProjects(ids []int64) ([]types.PinnedProjectResponse, error) {
	var res []types.PinnedProjectResponse

	err := GetDB().
		Model(&types.DbProject{}).
		Where("is_pinned = ?", true).
		Where("is_active = ?", true).
		WhereIn("id IN (?)", ids).
		Select(&res)

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return res, nil
}

// get all pinned projects by userid unused function
func GetActivePinnedProjectsByUserId(uid int64) ([]*types.DbProject, error) {
	pinnedProjects := make([]*types.DbProject, 0)

	err := GetDB().
		Model(&types.DbUserProject{}, &types.DbProject{}).
		Column("p.*").
		Where("p.is_pinned = ?", true).
		Where("p.is_active = ?", true).
		Where("user_projects.user_id = ?", uid).
		Join("JOIN projects p").JoinOn("user_projects.project_id = p.id").
		Select(&pinnedProjects)

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return pinnedProjects, nil
}

// func update project by project Name
func UpdateProjectbyProjectNmae(pName string, isPinned bool) error {

	stmt := `UPDATE projects SET is_pinned = ? WHERE name = ?`
	_, err := GetDB().Exec(stmt, isPinned, pName)
	if err != nil {
		log.Fatalf("Unable to execute query: %v\n", err)
	}

	return nil
}

func GetAllActiveProjectsByProjectIds(pIds []int64) ([]*types.DbProject, error) {
	res := make([]*types.DbProject, 0)

	query := GetDB().
		Model(&types.DbProject{}).
		Where("is_active =?", true).
		Where("id = any(?)", pg.Array(pIds))

	err := query.Select(&res)
	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return res, nil
}

// func update project by active status
func UpdateProjectbyActiveStatus(id int64, isActive bool) error {

	stmt := `UPDATE projects SET is_active = ? WHERE id = ?`
	_, err := GetDB().Exec(stmt, isActive, id)
	if err != nil {
		return err
	}

	return nil
}

// get all project details by userid unused function
func GetProjectDetailsByUserId(uid int64) ([]types.DbProject, error) {
	projects := make([]types.DbProject, 0)

	err := GetDB().
		Model(&types.DbUserProject{}, &types.DbProject{}).
		Column("p.*").
		Where("user_projects.user_id = ?", uid).
		Join("JOIN projects p").JoinOn("user_projects.project_id = p.id").
		Select(&projects)

	if err != nil {
		if errors.Is(err, pg.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return projects, nil
}
