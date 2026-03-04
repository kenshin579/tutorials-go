package backend

import (
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) CreateMenu() *menu.Menu {
	appMenu := menu.NewMenu()

	fileMenu := appMenu.AddSubmenu("파일")
	fileMenu.AddText("불러오기...", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
		a.ImportTodos()
		runtime.EventsEmit(a.ctx, "todos:reload")
	})
	fileMenu.AddText("내보내기...", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
		a.ExportTodos()
	})
	fileMenu.AddSeparator()
	fileMenu.AddText("종료", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})

	return appMenu
}
