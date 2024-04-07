package main

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os/exec"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
// NewApp 创建一个新的 App 应用程序
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
// startup 在应用程序启动时调用
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	// 在这里执行初始化设置
	a.ctx = ctx
}

// domReady is called after the front-end dom has been loaded
// domReady 在前端Dom加载完毕后调用
func (a *App) domReady(ctx context.Context) {
	// Add your action here
	// 在这里添加你的操作
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
// beforeClose在单击窗口关闭按钮或调用runtime.Quit即将退出应用程序时被调用.
// 返回 true 将导致应用程序继续，false 将继续正常关闭。
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	return false
}

// shutdown is called at application termination
// 在应用程序终止时被调用
func (a *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	// 在此处做一些资源释放的操作
}

func (a *App) SelectFile(title string) (path string, err error) {
	path, err = runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: title,
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
		ShowHiddenFiles: true,
	})
	slog.Info("[%s]:[%v]\n", path, err)
	if err != nil {
		slog.Info("cancel?")
		return
	}
	return
}

func (a *App) Write(p []byte) (n int, err error) {
	// 将输出发送到前端
	runtime.EventsEmit(a.ctx, "stderr", string(p))
	return len(p), nil
}

func (a *App) RunSeqAnalysis(args []string) (result string, err error) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("RunSeqAnalysis", "err", e)
			if err == nil {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	// SeqAnalysis -i baseName
	cmd := exec.Command("SeqAnalysis.exe", args...)
	slog.Info("run SeqAnalysis", "cmd", cmd.String())
	// // 使用自定义Writer捕获Stderr
	// cmd.Stderr = &CustomWriter{ctx: a.ctx}

	// err = cmd.Run()
	// if err != nil {
	// 	slog.Error("SeqAnalysis.exe", "err", err)
	// }
	// return

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		runtime.LogError(a.ctx, "创建 Stderr 管道失败: "+err.Error())
		return
	}

	if err = cmd.Start(); err != nil {
		runtime.LogError(a.ctx, "启动命令失败: "+err.Error())
		return
	}

	scanner := bufio.NewScanner(stderrPipe)
	for scanner.Scan() {
		msg := scanner.Text()
		slog.Info("SeqAnalysis.exe", "msg", msg)
		runtime.EventsEmit(a.ctx, "stderr-output", msg)
	}

	if err = cmd.Wait(); err != nil {
		runtime.LogError(a.ctx, "命令执行出错: "+err.Error())
	}
	return
}

// type CustomWriter struct {
// 	ctx context.Context
// }

// func (cw *CustomWriter) Write(p []byte) (n int, err error) {
// 	// 将输出发送到前端
// 	runtime.EventsEmit(cw.ctx, "stderr", string(p))
// 	return len(p), nil
// }
