package polls

type poll struct {
	ID      string
	Title   string
	Options []string
	Status  PollStatus
	Outcome OutcomeStatus
}

type Poll interface {
	GetID() string
	GetTitle() string
	GetOptions() []string
	GetStatus() PollStatus
	GetOutcome() OutcomeStatus
}

func (p *poll) GetID() string                    { return p.ID }
func (p *poll) GetTitle() string                 { return p.Title }
func (p *poll) SetTitle(title string)            { p.Title = title }
func (p *poll) GetOptions() []string             { return p.Options }
func (p *poll) SetOptions(options []string)      { p.Options = options }
func (p *poll) GetStatus() PollStatus            { return p.Status }
func (p *poll) SetStatus(status PollStatus)      { p.Status = status }
func (p *poll) GetOutcome() OutcomeStatus        { return p.Outcome }
func (p *poll) SetOutcome(outcome OutcomeStatus) { p.Outcome = outcome }

type PollStatus int

const (
	Open PollStatus = iota
	Closed
)

type OutcomeStatus int

const (
	Option1 OutcomeStatus = iota
	Option2
	Pending
)
