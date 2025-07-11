package app

type ViewState int

const (
	MainView ViewState = iota
	PrintView
	ConfigView
	AccountView
	ThemeView
	SetupView
)

type ViewController struct {
	Current ViewState
}

func NewViewController() *ViewController {
	return &ViewController{Current: MainView}
}

func (vc *ViewController) Set(view ViewState) {
	vc.Current = view
}

func (vc *ViewController) Get() ViewState {
	return vc.Current
}
