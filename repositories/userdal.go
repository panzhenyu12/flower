package repositories

import (
	"strings"
	"sync"

	"flower/entities"
	"flower/model"
	"github.com/pkg/errors"
)

type AccountDal struct {
	*DB
}

func GetAccountDal() *AccountDal {
	return &AccountDal{
		DB: GetDeepDB(),
	}
}

func (dal *AccountDal) GetValidByUserName(username string) (*entities.Account, error) {
	e := &entities.Account{}
	s := dal.GetNewSession()
	exist, err := dal.buildValidQuery(s).
		Where("user_name = ?", username).Get(e)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !exist {
		return nil, nil
	}
	return e, nil
}

func (dal *AccountDal) GetValidExtendByUserName(username string) (*entities.AccountExtend, error) {
	e := &entities.AccountExtend{}
	s := dal.GetNewSession()
	exist, err := dal.buildValidQuery(s).
		InnerJoinById(entities.TableNameAccount, entities.TableNameFuncRole, "func_role_id").
		Where("user_name = ?", username).
		Get(e)
	if err != nil {
		return nil, errors.WithStack(err)
	} else if !exist {
		return nil, nil
	}
	return e, nil
}

func (dal *AccountDal) GetByOrgIDs(ids []int64) ([]*entities.AccountExtend, error) {
	tableName := entities.TableNameAccount
	accounts := make([]*entities.AccountExtend, 0)
	err := dal.GetNewSession().
		InnerJoinById(tableName, entities.TableNameFuncRole, "func_role_id").
		WhereStatusNotDeleted(tableName).
		Find(&accounts)
	return accounts, errors.WithStack(err)
}

func (dal *AccountDal) Query(query model.BaseQuery, ids ...int64) ([]*entities.AccountExtend, error) {
	accounts := make([]*entities.AccountExtend, 0)
	session := dal.GetNewSession()
	if ids == nil || len(ids) == 0 {
		return accounts, nil
	}
	tableName := entities.TableNameAccount
	orderbytb := ""
	if strings.ToLower(query.OrderBy) != "ts" {
		orderbytb = entities.TableNameOrg
		query.OrderBy = "org_level"
	} else {
		orderbytb = tableName
	}
	//session.In()
	session.InnerJoinById(tableName, entities.TableNameFuncRole, "func_role_id").
		InnerJoinById(tableName, entities.TableNameOrg, "org_id").
		WhereTsMatch(tableName, query).
		WhereStatusNotDeleted(tableName).
		WhereCond(buildWhereInCond(tableName, "org_id", ids)).
		Sort(orderbytb, query).OrderBy(orderbytb, "ts", true).Page(query)
	err := session.Find(&accounts)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return accounts, nil
}
func (dal *AccountDal) Count(query model.BaseQuery, ids ...int64) (int64, error) {
	session := dal.GetNewSession()
	if ids == nil || len(ids) == 0 {
		return 0, nil
	}
	tableName := entities.TableNameAccount
	session.InnerJoinById(tableName, entities.TableNameFuncRole, "func_role_id").
		InnerJoinById(tableName, entities.TableNameOrg, "org_id").
		WhereTsMatch(tableName, query).
		WhereStatusNotDeleted(tableName).
		WhereCond(buildWhereInCond(tableName, "org_id", ids))
	count, err := session.Count(new(entities.Account))
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return count, nil
}
func (dal *AccountDal) QueryAndCount(query model.BaseQuery, ids []int64) (int64, []*entities.AccountExtend, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var count int64
	data := make([]*entities.AccountExtend, 0)
	var err1, err2 error
	go func() {
		defer wg.Done()
		count, err1 = dal.Count(query, ids...)
	}()
	go func() {
		defer wg.Done()
		data, err2 = dal.Query(query, ids...)
	}()
	wg.Wait()
	if err1 != nil {
		return 0, nil, errors.WithStack(err1)
	}
	if err2 != nil {
		return 0, nil, errors.WithStack(err2)
	}
	return count, data, nil
}
func (dal *AccountDal) Insert(entity *entities.Account) error {
	_, err := dal.GetNewSession().Insert(entity)
	return errors.WithStack(err)
}

func (dal *AccountDal) Update(entity *entities.Account) error {
	s := dal.GetNewSession()
	s.MustCols("is_valid", "real_name", "comment")
	return errors.WithStack(dal.UpdateCore(entity, s))
}

func (dal *AccountDal) UpdateCore(entity *entities.Account, session *Session) error {
	if session == nil {
		session = dal.GetNewSession()
	}
	_, err := session.
		ID(entity.UserId).
		Update(entity)
	return errors.WithStack(err)
}

func (dal *AccountDal) UpdateSecurityToken(userName, securityToken string) error {
	e := &entities.Account{
		UserName:      userName,
		SecurityToken: securityToken,
	}
	_, err := dal.GetNewSession().
		Where("user_name = ?", userName).
		Update(e)
	return errors.WithStack(err)
}

func (dal *AccountDal) Delete(id int64) error {
	_, err := dal.GetNewSession().ID(id).Delete(&entities.Account{})
	return errors.WithStack(err)
}

func (dal *AccountDal) BatchUpdateSecurityToken(m map[int64]string, session *Session) error {
	fields := []string{"user_id", "security_token"}
	values := make([][]interface{}, len(m))
	i := 0
	for id, token := range m {
		values[i] = []interface{}{id, token}
		i++
	}
	return errors.WithStack(session.BatchUpdateFields(entities.TableNameAccount, fields, values))
}

func (dal *AccountDal) buildValidQuery(s *Session) *Session {
	s.WhereStatusNotDeleted(entities.TableNameAccount).
		Where("is_valid = ?", true)
	return s
}
