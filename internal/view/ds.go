package view

import (
	"github.com/derailed/k9s/internal/client"
	"github.com/derailed/k9s/internal/dao"
	"github.com/derailed/k9s/internal/render"
	"github.com/derailed/k9s/internal/ui"
)

// DaemonSet represents a daemon set custom viewer.
type DaemonSet struct {
	ResourceViewer
}

// NewDaemonSet returns a new viewer.
func NewDaemonSet(gvr client.GVR) ResourceViewer {
	d := DaemonSet{
		ResourceViewer: NewPortForwardExtender(
			NewRestartExtender(
				NewLogsExtender(NewBrowser(gvr), nil),
			),
		),
	}
	d.SetBindKeysFn(d.bindKeys)
	d.GetTable().SetEnterFn(d.showPods)
	d.GetTable().SetColorerFn(render.DaemonSet{}.ColorerFunc())

	return &d
}

func (d *DaemonSet) bindKeys(aa ui.KeyActions) {
	aa.Add(ui.KeyActions{
		ui.KeyShiftD: ui.NewKeyAction("Sort Desired", d.GetTable().SortColCmd(1, true), false),
		ui.KeyShiftC: ui.NewKeyAction("Sort Current", d.GetTable().SortColCmd(2, true), false),
		ui.KeyShiftR: ui.NewKeyAction("Sort Ready", d.GetTable().SortColCmd(3, true), false),
		ui.KeyShiftU: ui.NewKeyAction("Sort UpToDate", d.GetTable().SortColCmd(4, true), false),
		ui.KeyShiftL: ui.NewKeyAction("Sort Available", d.GetTable().SortColCmd(5, true), false),
	})
}

func (d *DaemonSet) showPods(app *App, model ui.Tabular, _, path string) {
	var res dao.DaemonSet
	res.Init(app.factory, client.NewGVR(d.GVR()))

	ds, err := res.GetInstance(path)
	if err != nil {
		d.App().Flash().Err(err)
	}

	showPodsFromSelector(app, path, ds.Spec.Selector)
}
