package repositories

import (
	"fmt"

	"github.com/go-xorm/builder"

	"flower/common"
	"flower/entities"
	"flower/model"
	"github.com/pkg/errors"
)

type FuncRoleDal struct {
	*DB
}

func GetFuncRoleDal() *FuncRoleDal {
	return &FuncRoleDal{
		DB: GetDeepDB(),
	}
}

func (dal *FuncRoleDal) Query(query *model.FuncRoleQuery) ([]*entities.FuncRole, error) {
	tableName := entities.TableNameFuncRole
	results := make([]*entities.FuncRole, 0)
	err := dal.buildQuery(query).
		Sort(tableName, query).
		Page(query).
		Select(fmt.Sprintf("*, (SELECT COUNT(*) FROM %v WHERE %v.func_role_id = %v.func_role_id) AS user_count",
			entities.TableNameAccount,
			entities.TableNameAccount,
			entities.TableNameFuncRole)).
		Find(&results)
	return results, errors.WithStack(err)
}

func (dal *FuncRoleDal) QueryAndCount(query *model.FuncRoleQuery) ([]*entities.FuncRole, int, error) {
	results, err := dal.Query(query)
	if err != nil {
		return nil, 0, errors.WithStack(err)
	}
	c, err := dal.buildQuery(query).Count(new(entities.FuncRole))
	return results, int(c), errors.WithStack(err)
}

func (dal *FuncRoleDal) GetByID(id int64) (*entities.FuncRole, error) {
	entity := new(entities.FuncRole)
	_, err := WhereNoDelete(dal.Engine.Id(id)).Get(entity)
	if err != nil {
		return nil, err
	}
	return entity, nil
}

func (dal *FuncRoleDal) Insert(entity *entities.FuncRole, session *Session) error {
	if session == nil {
		session = dal.GetNewSession()
	}
	_, err := session.Insert(entity)
	return errors.WithStack(err)
}

func (dal *FuncRoleDal) Upate(entity *entities.FuncRole, session *Session) error {
	if session == nil {
		session = dal.GetNewSession()
	}
	_, err := session.ID(entity.FuncRoleId).Update(entity)
	return errors.WithStack(err)
}

func (dal *FuncRoleDal) SoftDelete(id int64, session *Session) error {
	e := &entities.FuncRole{
		FuncRoleId: id,
		Status:     int(common.TableStatus_Delete),
	}
	return errors.WithStack(dal.Upate(e, session))
}

func (dal *FuncRoleDal) GetUserIds(session *Session, ids ...int64) ([]int64, error) {
	tableName := entities.TableNameAccount
	results := make([]int64, 0)
	err := session.
		Where(builder.In("func_role_id", ids)).
		Select("account.user_id").
		Table(tableName).
		Find(&results)
	return results, errors.WithStack(err)
}

func (dal *FuncRoleDal) buildQuery(query *model.FuncRoleQuery) *Session {
	tableName := entities.TableNameFuncRole
	return dal.GetNewSession().
		WhereStatusNotDeleted(tableName).
		WhereTsMatch(tableName, query)
}
