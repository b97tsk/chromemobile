package chromemobile

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/b97tsk/chrome/service"
)

type ChromeService struct {
	manager     *service.Manager
	filesDir    string
	workingPath string
}

func NewChromeService(filesDir string) *ChromeService {
	chrome := &ChromeService{
		manager:  newManager(),
		filesDir: filepath.Clean(filesDir),
	}
	chrome.loadWorking()
	return chrome
}

func (chrome *ChromeService) Shutdown() {
	chrome.manager.Shutdown()
}

func (chrome *ChromeService) IsWorking() bool {
	return chrome.workingPath != ""
}

func (chrome *ChromeService) Load(filename string) error {
	fd, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fd.Close()
	return chrome.load(fd)
}

func (chrome *ChromeService) LoadURL(url string) (err error) {
	fd, err := ioutil.TempFile(chrome.filesDir, "~")
	if err != nil {
		return
	}
	defer os.Remove(fd.Name())
	defer fd.Close()
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	_, err = io.Copy(fd, resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}
	return chrome.load(fd)
}

func (chrome *ChromeService) load(fd *os.File) (err error) {
	tmpDir, err := ioutil.TempDir(chrome.filesDir, "tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpDir)
		} else {
			chrome.setWorking(tmpDir)
		}
	}()

	filesize, _ := fd.Seek(0, io.SeekEnd)
	_, _ = fd.Seek(0, io.SeekStart)
	if zr, err := zip.NewReader(fd, filesize); err == nil {
		return chrome.loadZipFile(zr, tmpDir)
	}

	_, _ = fd.Seek(0, io.SeekStart)
	data, _ := ioutil.ReadAll(fd)
	filename := filepath.Join(tmpDir, "chrome.yaml")
	if err := ioutil.WriteFile(filename, data, 0666); err != nil {
		return err
	}
	chrome.loadYAML(filename)
	return nil
}

func (chrome *ChromeService) loadZipFile(zr *zip.Reader, tmpDir string) error {
	var ok bool
	for _, file := range zr.File {
		if file.Name == "chrome.yaml" {
			ok = true
			break
		}
	}
	if !ok {
		return errors.New("chrome.yaml not found in archive")
	}
	for _, file := range zr.File {
		if strings.HasSuffix(file.Name, "/") {
			continue // Skip directories.
		}
		rc, err := file.Open()
		if err != nil {
			return err
		}
		filename := filepath.Join(tmpDir, file.Name)
		if !strings.HasPrefix(filename, tmpDir+string(os.PathSeparator)) {
			continue // Relative path goes up too far.
		}
		_ = os.MkdirAll(filepath.Dir(filename), 0666)
		fd, err := os.Create(filename)
		if err == nil {
			_, err = io.Copy(fd, rc)
			if cerr := fd.Close(); err == nil {
				err = cerr
			}
		}
		rc.Close()
		if err != nil {
			return err
		}
	}
	filename := filepath.Join(tmpDir, "chrome.yaml")
	chrome.loadYAML(filename)
	return nil
}

func (chrome *ChromeService) loadYAML(filename string) {
	dir, base := filepath.Dir(filename), filepath.Base(filename)
	_ = os.Chdir(dir)
	chrome.manager.LoadFile(base)
}

func (chrome *ChromeService) setWorking(path string) {
	filename := filepath.Join(chrome.filesDir, "working")
	_ = ioutil.WriteFile(filename, []byte(path), 0666)
	if chrome.workingPath != "" {
		os.RemoveAll(chrome.workingPath)
	}
	chrome.workingPath = path
}

func (chrome *ChromeService) loadWorking() {
	filename := filepath.Join(chrome.filesDir, "working")
	data, err := ioutil.ReadFile(filename)
	if err == nil {
		path := string(data)
		filename := filepath.Join(path, "chrome.yaml")
		if _, err := os.Stat(filename); err == nil {
			chrome.loadYAML(filename)
			chrome.workingPath = path
		}
	}
}
