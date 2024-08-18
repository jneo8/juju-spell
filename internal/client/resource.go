package client

type Resource struct {
	Controller string
	Model      string
	Name       string
}

func (r *Resource) ResourceName() string {
	s := r.Controller + ":" + r.Model + ":" + r.Name
	return s
}
