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

	logo, err := GetFullSizeLogoPNG()
	if err != nil {
		panic(fmt.Sprintf("error reading logo: %v", err))
	}
	v.logoImg = canvas.NewImageFromImage(logo)
	v.logoImg.FillMode = canvas.ImageFillOriginal

	return &v
}

func (v *HomeView) Foreground() fyne.CanvasObject {
	v.mu.Lock()
	defer v.mu.Unlock()
	if v.foreground() {
		todayBtn := NewWhiteTextButton("Today's List", func() {
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
				v.app.RenderCreateListView()
				return
			}
			v.app.RenderTaskListView(*latestList)
		})
		todayBtn.Importance = widget.MediumImportance

		createListBtn := NewWhiteTextButton("Create List", v.app.RenderCreateListView)

		return container.NewCenter(

			container.NewVBox(
				v.logoImg,
				todayBtn,
				createListBtn,
			),
		)
	}
	return nil
}

func (v *HomeView) Background() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.background()
}
