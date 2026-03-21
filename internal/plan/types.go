package plan

type Status string

const (
	StatusInbox     Status = "inbox"
	StatusPlan      Status = "plan"
	StatusActive    Status = "active"
	StatusDone      Status = "done"
	StatusCancelled Status = "cancelled"
)

type Frontmatter struct {
	ID     string   `yaml:"id" json:"id"`
	Title  string   `yaml:"title" json:"title"`
	Status Status   `yaml:"status" json:"status"`
	Tags   []string `yaml:"tags" json:"tags"`
	Parent string   `yaml:"parent" json:"parent"`
}

type Document struct {
	Path        string      `json:"path"`
	Frontmatter Frontmatter `json:"frontmatter"`
	Body        string      `json:"body"`
}

type ValidationIssue struct {
	Path    string `json:"path"`
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (s Status) Valid() bool {
	switch s {
	case StatusInbox, StatusPlan, StatusActive, StatusDone, StatusCancelled:
		return true
	default:
		return false
	}
}
