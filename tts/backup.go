package tts

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"regexp"
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

	err = downloadModels(s, dirs["Models"])

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

	dirs["Models"] = dirs["Root"] + "/" + "Models"

	err = os.Mkdir(dirs["Models"], 0777)
	if err != nil {
		return map[string]string{}, err
	}

	return dirs, nil
}

func downloadImages(s SaveFile, rootDir string) error {

	urls := map[string]bool{}

	if s.TableURL != "" {
		urls[s.TableURL] = true
	}

	if s.SkyURL != "" {
		urls[s.SkyURL] = true
	}

	for _, o := range s.ObjectStates {
		if o.CustomImage.ImageURL != "" {
			urls[o.CustomImage.ImageURL] = true
		}
		if o.CustomMesh.DiffuseURL != "" {
			urls[o.CustomMesh.DiffuseURL] = true
		}
		if o.CustomMesh.NormalURL != "" {
			urls[o.CustomMesh.NormalURL] = true
		}
		for _, d := range o.CustomDeck {
			if d.FaceURL != "" {
				urls[d.FaceURL] = true
			}
			if d.BackURL != "" {
				urls[d.BackURL] = true
			}
		}
	}

	downloadURLs(urls, rootDir)

	return nil
}

func downloadModels(s SaveFile, rootDir string) error {

	urls := map[string]bool{}

	for _, o := range s.ObjectStates {
		if o.CustomMesh.MeshURL != "" {
			urls[o.CustomMesh.MeshURL] = true
		}
		if o.CustomMesh.ColliderURL != "" {
			urls[o.CustomMesh.ColliderURL] = true
		}
	}

	downloadURLs(urls, rootDir)

	return nil

}

func downloadURLs(urls map[string]bool, rootDir string) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return err
	}

	for k := range urls {
		filename := reg.ReplaceAllString(k, "")

		err := downloadFile(k, rootDir+"/"+filename)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url string, filename string) error {

	out, err := os.Create(filename)
	if err != nil {
		return err
	}

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	out.Seek(0, 0)

	fileType, err := getFileContentType(out)
	if err != nil {
		return err
	}

	out.Close()

	if fileType == "image/png" {
		err := os.Rename(filename, filename+".png")
		if err != nil {
			return err
		}
	} else if fileType == "image/jpeg" {
		err := os.Rename(filename, filename+".jpg")
		if err != nil {
			return err
		}
	} else if fileType == "text/plain; charset=utf-8" {
		err := os.Rename(filename, filename+".obj")
		if err != nil {
			return nil
		}
	}

	return nil
}

func getFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	return contentType, nil
}
