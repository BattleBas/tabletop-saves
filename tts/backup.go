package tts

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

// Backup parses a Tabeltop Simulator JSON save file and
// saving all dependencies to a zip
func Backup(filename string) error {

	usr, err := user.Current()
	if err != nil {
		return err
	}
	dir := usr.HomeDir
	windowsPath := "\\Documents\\My Games\\Tabletop Simulator\\Mods\\Workshop\\"

	f, err := os.Open(dir + windowsPath + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var s SaveFile
	err = json.Unmarshal(b, &s)
	if err != nil {
		return err
	}

	dirs, err := createDirectories(s.SaveName)
	if err != nil {
		return nil
	}

	err = downloadImages(s, dirs["Images"])
	if err != nil {
		return err
	}

	fmt.Println(s.SaveName)
	fmt.Println(s.TableURL)

	return nil
}

func createDirectories(gameName string) (map[string]string, error) {

	dirs := map[string]string{}

	dirs["Root"] = strings.Replace(gameName, " ", "", -1)

	err := os.Mkdir(dirs["Root"], 0777)
	if err != nil {
		return map[string]string{}, err
	}

	dirs["Images"] = dirs["Root"] + "/" + "Image"

	err = os.Mkdir(dirs["Images"], 0777)
	if err != nil {
		return map[string]string{}, err
	}

	return dirs, nil
}

func downloadImages(s SaveFile, dir string) error {

	urls := []string{
		s.BackURL,
		s.DiffuseURL,
		s.FaceURL,
		s.ImageURL,
		s.NormalURL,
		s.SkyURL,
		s.TableURL,
	}

	for _, u := range urls {
		if u == "" {
			continue
		}
		p := path.Base(u)
		err := downloadFile(u, dir+"/"+p)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url string, file string) error {

	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
