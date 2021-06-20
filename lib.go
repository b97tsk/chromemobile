package chromemobile

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/b97tsk/chrome"
)

type LogWriter interface {
	WriteLog(s string)
}

type ChromeService struct {
	manager     *chrome.Manager
	filesDir    string
	workingPath string
}

func NewChromeService(filesDir string, w LogWriter) *ChromeService {
	chrome := &ChromeService{
		manager:  newManager(),
		filesDir: filepath.Clean(filesDir),
	}

	if w != nil {
		chrome.manager.SetLogOutput(writerFunc(
			func(p []byte) (n int, err error) {
				w.WriteLog(string(p))
				return len(p), nil
			},
		))
	}

	_ = chrome.loadWorking()

	return chrome
}

func (chrome *ChromeService) Shutdown() {
	chrome.manager.Shutdown()
}

func (chrome *ChromeService) IsWorking() bool {
	return chrome.workingPath != ""
}

func (chrome *ChromeService) LoadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return chrome.Load(file)
}

func (chrome *ChromeService) LoadURL(url string) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	return chrome.Load(resp.Body)
}

func (chrome *ChromeService) Load(r io.Reader) (err error) {
	dir, err := os.MkdirTemp(chrome.filesDir, workingDirPattern)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			os.RemoveAll(dir)
		} else {
			chrome.setWorking(dir)
		}
	}()

	file, err := os.OpenFile(
		filepath.Join(dir, configFileName),
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0o600,
	)
	if err != nil {
		return err
	}

	if _, err := io.Copy(file, r); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return chrome.loadDir(dir)
}

func (chrome *ChromeService) loadDir(dir string) error {
	if err := os.Chdir(dir); err != nil {
		return err
	}

	return chrome.manager.LoadFile(configFileName)
}

func (chrome *ChromeService) setWorking(path string) {
	filename := filepath.Join(chrome.filesDir, workingFileName)
	_ = os.WriteFile(filename, []byte(path), 0o600)

	if chrome.workingPath != "" {
		os.RemoveAll(chrome.workingPath)
	}

	chrome.workingPath = path
}

func (chrome *ChromeService) loadWorking() error {
	filename := filepath.Join(chrome.filesDir, workingFileName)

	content, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	dir := string(content)

	if err := chrome.loadDir(dir); err != nil {
		return err
	}

	chrome.workingPath = dir

	return nil
}

type writerFunc func(p []byte) (n int, err error)

func (f writerFunc) Write(p []byte) (n int, err error) {
	return f(p)
}

const (
	configFileName    = "chrome.config"
	workingDirPattern = "tmp"
	workingFileName   = "working"
)
