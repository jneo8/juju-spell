package view

import "github.com/jneo8/jujuspell/internal/model"

type Viewer interface {
	model.Component
}

type TableViewer interface {
	Viewer
}

type ResourceReviewer interface {
	TableViewer
}
