package model

type OrgQuery struct {
	// OrgId            int64 `form:"-"`
	RecursiveDepth   int64 `form:"recursivedepth"`
	IncludeFaceRepos bool  `form:"facerepos"`
	IncludeFaceRules bool  `form:"facerules"`
	IncludeStations  bool  `form:"stations"`
	IncludeUsers     bool  `form:"users"`
}
