package juju

type JujuError struct {
	Msg string
}

func (j *JujuError) Error() string {
	return j.Msg
}
