package main

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var _ View = (*HomeView)(nil)

type HomeView struct {
	*baseView

	logoImg *canvas.Image
}

func NewHomeView(app *TaskApp) *HomeView {
	v := HomeView{
		baseView: newBaseView("home", app),
	}

	v.logoImg = GetAssetImageCanvas(GetConstrainedImage(AssetImageLogo, 400))

	return &v
}

func (v *HomeView) Title() []fyne.CanvasObject {
	return nil
}

func (v *HomeView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if !v.foreground() {
		return nil
	}

	todayBtn := widget.NewButton("Today's List", func() {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func() {
			<-v.deactivated
			cancel()
		}()
		latestList, err := FindOneModel[TaskList](ctx, v.app.DB(), WithSort("Date desc"))
		if err != nil {
			log.Error("Error finding latest task list", "err", err)
			panic(fmt.Sprintf("Error finding latest task list: %v", err))
		}
		if latestList == nil {
			v.app.RenderMutateTaskListView(nil)
			return
		}
		v.app.RenderListOfTasksView("Today's List", nil, todaysTasksModelQueryOpt())
	})
	todayBtn.Importance = widget.MediumImportance

	createListBtn := widget.NewButton("Create List", func() {
		v.app.RenderMutateTaskListView(nil)
	})

	return container.NewCenter(
		container.NewVBox(
			v.logoImg,
			todayBtn,
			createListBtn,
		),
	)
}

func (v *HomeView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
