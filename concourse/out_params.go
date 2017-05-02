package concourse

type OutParams struct {
	ErrandName  string `json:"name,omitempty"`
	KeepAlive   bool   `json:"keep_alive,omitempty"`
	WhenChanged bool   `json:"when_changed,omitempty"`
}
