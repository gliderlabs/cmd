package sentry

const (
	EventNewIssue = "Sentry:NewIssue"
)

type IssueEvent struct {
	Project     string `json:"project"`
	ProjectName string `json:"project_name"`
	Culprit     string `json:"culprit"`
	Level       string `json:"level"`
	URL         string `json:"url"`
	Logger      string `json:"logger"`
	Message     string `json:"message"`
	ID          string `json:"id"`
	Event       struct {
		Received                float64 `json:"received"`
		SentryInterfacesMessage struct {
			Message string `json:"message"`
		} `json:"sentry.interfaces.Message"`
		Version string `json:"version"`
		Extra   struct {
			RuntimeNumGoroutine int    `json:"runtime.NumGoroutine"`
			RuntimeVersion      string `json:"runtime.Version"`
			RuntimeNumCPU       int    `json:"runtime.NumCPU"`
			RuntimeGOMAXPROCS   int    `json:"runtime.GOMAXPROCS"`
		} `json:"extra"`
		EventID     string        `json:"event_id"`
		Fingerprint []string      `json:"fingerprint"`
		ID          int64         `json:"id"`
		Errors      []interface{} `json:"errors"`
		RefVersion  int           `json:"_ref_version"`
		Ref         int           `json:"_ref"`
		Metadata    struct {
			Title string `json:"title"`
		} `json:"metadata"`
		Type                       string `json:"type"`
		SentryInterfacesStacktrace struct {
			Frames []struct {
				Function string `json:"function"`
				AbsPath  string `json:"abs_path"`
				Module   string `json:"module"`
				Filename string `json:"filename"`
				Lineno   int    `json:"lineno"`
				InApp    bool   `json:"in_app"`
			} `json:"frames"`
			HasSystemFrames bool        `json:"has_system_frames"`
			Registers       interface{} `json:"registers"`
			FramesOmitted   interface{} `json:"frames_omitted"`
		} `json:"sentry.interfaces.Stacktrace"`
		Tags [][]string `json:"tags"`
		Sdk  struct {
			ClientIP string `json:"client_ip"`
			Version  string `json:"version"`
			Name     string `json:"name"`
		} `json:"sdk"`
	} `json:"event"`
}

func (e IssueEvent) EventName() string {
	return EventNewIssue
}
