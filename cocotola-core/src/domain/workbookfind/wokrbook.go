package workbookfind

type Parameter struct {
	PageNo   int
	PageSize int
}

type WorkbookModel struct {
	ID   int
	Name string
}

type Result struct {
	TotalCount int
	Results    []*WorkbookModel
}
