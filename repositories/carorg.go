package repositories

type OrgDal struct {
	*DB
}

func GetOrgDal() *OrgDal {
	return &OrgDal{
		DB: GetBiDB(),
	}
}
