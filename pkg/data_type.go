package pkg

import "time"

type ReqMsgUpload struct {
	ZipData []byte `json:"zip_data"`
}
type ResMsgUpload struct {
	Log string `json:"log"`
}

type ReqMsgCreate struct {
	Envs map[string]string `json:"envs"`
}
type ResMsgCreate struct{}

type ReqMsgRemove struct{}
type ResMsgRemove struct{}

type ReqMsgStatus struct{}
type ResMsgStatus struct {
	Status    string    `json:"status"`
	StartedAt time.Time `json:"started_at"`
}

type ReqMsgLogs = struct{}
type ResMsgLogs = struct {
	Log string `json:"log"`
}
