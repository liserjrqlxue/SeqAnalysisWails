package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"SeqAnalysis/pkg/seqAnalysis"

	"github.com/liserjrqlxue/goUtil/osUtil"
	"github.com/liserjrqlxue/goUtil/simpleUtil"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// global variables
var (
	LogLevel slog.LevelVar
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
	logFile := osUtil.Create("SeqAnalysisWails.log")
	mw := io.MultiWriter(logFile, a)
	LogLevel.Set(slog.LevelInfo)
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				mw,
				&slog.HandlerOptions{
					Level: &LogLevel,
				},
			),
		),
	)
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
	slog.Info("select file", "path", path, "err", err)
	if err != nil {
		slog.Info("cancel?")
		return
	}
	return
}

func (a *App) Write(p []byte) (n int, err error) {
	// fmt.Fprintf(os.Stderr, "%s", p)
	// 将输出发送到前端
	runtime.EventsEmit(a.ctx, "stderr-output", string(p))
	// log.Print(string(p))
	return len(p), nil
}

func (a *App) RunSeqAnalysis(input, workDir, outputPrefix string, long, plot, lessMem, debug bool, lineLimit, short int) (result string, err error) {
	defer func() {
		if e := recover(); e != nil {
			slog.Error("RunSeqAnalysis", "err", e)
			if err == nil {
				err = fmt.Errorf("%v", e)
			}
		}
	}()

	slog.Info("RunSeqAnalysis", "input", input, "workDir", workDir, "outputPrefix", outputPrefix, "long", long, "plot", plot, "lessMem", lessMem, "debug", debug, "lineLimit", lineLimit)
	if debug {
		slog.Info("Set log level to debug")
		LogLevel.Set(slog.LevelDebug)
		slog.Info("Set log level to debug Done")
		slog.Debug("Set log level to debug Done")
		slog.Debug("debug", "input", input, "workDir", workDir, "outputPrefix", outputPrefix, "long", long, "plot", plot, "lessMem", lessMem, "debug", debug, "lineLimit", lineLimit)
	} else {
		LogLevel.Set(slog.LevelInfo)
	}

	seqAnalysis.Short = short
	var batch = seqAnalysis.Batch{
		OutputPrefix: outputPrefix,
		BasePrefix:   filepath.Base(outputPrefix),

		Long:      long,
		LessMem:   lessMem,
		Plot:      plot,
		LineLimit: lineLimit,
		Zip:       true,
		UseRC:     true,
		// UseKmer:   true,

		Sheets:           make(map[string]string),
		SeqInfoMap:       make(map[string]*seqAnalysis.SeqInfo),
		ParallelStatsMap: make(map[string]*seqAnalysis.ParallelTest),
	}
	slog.Info("ENV", "cwd", simpleUtil.HandleError(os.Getwd()))
	slog.Info("ENV", "PATH", os.Getenv("PATH"))
	slog.Info("ENV", "os.Environ()", strings.Join(os.Environ(), "BREAK"))
	err = batch.BatchRun(input, workDir, exPath, etcEMFS, 0)
	slog.Info("ENV", "cwd", simpleUtil.HandleError(os.Getwd()))

	return
}
