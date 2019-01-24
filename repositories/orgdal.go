package repositories

import (
	"fmt"

	"github.com/go-xorm/xorm"
	"github.com/panzhenyu12/flower/common"
	"github.com/panzhenyu12/flower/entities"
	"github.com/pkg/errors"
)

//explain  select ext_data->'home'->>'you' as ayName  from org_structure where ext_data->'home'->>'you'='None';
//curd
type DeepOrgDal struct {
	*DB
}

func GetDeepOrgDal() *DeepOrgDal {
	return &DeepOrgDal{
		DB: GetDeepDB(),
	}
}

func (dal *DeepOrgDal) GetOrgByID(id int64) (*entities.OrgStructure, error) {
	entity := new(entities.OrgStructure)
	bl, err := WhereNoDelete(dal.Engine.Id(id)).Get(entity)
	if bl == false {
		return nil, errors.New("not found")
	}
	if err != nil {
		return nil, err
	}
	return entity, nil
}
func (dal *DeepOrgDal) GetDownOrg(id int64) ([]*entities.OrgStructure, error) {
	sqlon := "c.org_id = k.superior_org_id"
	sql := fmt.Sprintf(dal.getSearchBase(), sqlon)
	orgs := make([]*entities.OrgStructure, 0)
	err := dal.Engine.Sql(sql, id).Find(&orgs)
	if err != nil {
		return nil, err
	}
	return orgs, nil
}
func (dal *DeepOrgDal) GetUpOrg(id int64) {
	// sqlon := "c.superior_org_id = k.org_id"
	// sql := fmt.Sprintf(dal.getSearchBase(), sqlon)
}
func (dal *DeepOrgDal) AddOne(entity *entities.OrgStructure, session *xorm.Session) error {
	if session == nil {
		session = dal.Engine.NewSession()
		defer session.Close()
		err := session.Begin()
		if err != nil {
			return err
		}
		_, err = session.Insert(entity)
		if err != nil {
			session.Rollback()
			return err
		}
		err = session.Commit()
		if err != nil {
			return err
		}
	} else {
		_, err := session.Insert(entity)
		if err != nil {
			return err
		}
	}
	return nil
}
func (dal *DeepOrgDal) DeleteByID(id int64, session *xorm.Session) error {
	entity := new(entities.OrgStructure)
	if session == nil {
		_, err := dal.Engine.Id(id).Delete(entity)
		if err != nil {
			return err
		}
	} else {
		_, err := session.Id(id).Delete(entity)
		if err != nil {
			return err
		}
	}
	return nil
}
func (dal *DeepOrgDal) UpDate(entity *entities.OrgStructure, session *xorm.Session) error {
	if session == nil {
		_, err := dal.Engine.Id(entity.OrgId).Update(entity)
		if err != nil {
			return err
		}
	} else {
		_, err := session.Id(entity.OrgId).Update(entity)
		if err != nil {
			return err
		}
	}
	return nil
}
func (dal *DeepOrgDal) getSearchBase() string {
	return `with RECURSIVE cte as 
	(select a.* from org_structure a where org_id=? 
	union all
	select k.*  from org_structure k inner join cte c on %v  where k.org_id !=k.superior_org_id and k.status!=4
	)select * from cte order by cte.org_level asc,cte.ts asc;`
}

func (dal *DeepOrgDal) Delete(session *Session, ids ...int64) error {
	if session == nil {
		session = dal.GetNewSession()
	}
	entity := &entities.OrgStructure{
		Status: int32(common.TableStatus_Delete),
	}
	_, err := session.In("org_id", ids).Cols("status").Update(entity)
	return errors.WithStack(err)
}

func (dal *DeepOrgDal) HasChild(parentId int64, session *Session) (bool, error) {
	if session == nil {
		session = dal.GetNewSession()
	}
	tableName := entities.TableNameOrg
	var id int64
	got, err := session.
		WhereStatusNotDeleted(tableName).
		Where("superior_org_id = ?", parentId).
		Table(tableName).
		Select("org_id").
		Get(&id)
	return got, errors.WithStack(err)
}

func (dal *DeepOrgDal) HasAccount(id int64, session *Session) (bool, error) {
	if session == nil {
		session = dal.GetNewSession()
	}
	tableName := entities.TableNameAccount
	var accountId int64
	got, err := session.
		WhereStatusNotDeleted(tableName).
		Where("org_id = ?", id).
		Table(tableName).
		Select("user_id").
		Get(&accountId)
	return got, errors.WithStack(err)
}
