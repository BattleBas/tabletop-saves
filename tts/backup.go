package tts

import (
	"archive/zip"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
)

// Backup parses a Tabeltop Simulator JSON save file and
// saving all dependencies to a zip
func Backup(filename string) error {

	s, err := readSaveFile(filename)

	z, err := os.Create(strings.Replace(s.SaveName, " ", "", -1) + ".zip")
	if err != nil {
		return err
	}
	defer z.Close()

	zipWriter := zip.NewWriter(z)
	defer zipWriter.Close()

	err = addSaveFile(filename, zipWriter)
	if err != nil {
		return err
	}

	err = downloadImages(s, "Image", zipWriter)
	if err != nil {
		return err
	}

	err = downloadModels(s, "Models", zipWriter)
	if err != nil {
		return err
	}

	return nil
}

func getModDir() string {
	usr, err := user.Current()
	if err != nil {
		return ""
	}
	dir := usr.HomeDir
	windowsPath := "\\Documents\\My Games\\Tabletop Simulator\\Mods\\Workshop\\"

	return dir + windowsPath
}

func addSaveFile(filename string, zipWriter *zip.Writer) error {
	dir := getModDir()
	f, err := os.Open(dir + filename)
	if err != nil {
		return err
	}

	zipF, err := zipWriter.Create(filename)

	_, err = io.Copy(zipF, f)
	if err != nil {
		return err
	}

	return nil
}

func readSaveFile(filename string) (SaveFile, error) {

	dir := getModDir()

	f, err := os.Open(dir + filename)
	if err != nil {
		return SaveFile{}, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return SaveFile{}, err
	}

	var s SaveFile
	err = json.Unmarshal(b, &s)
	if err != nil {
		return SaveFile{}, err
	}

	return s, nil
}

func downloadImages(s SaveFile, rootDir string, zipWriter *zip.Writer) error {

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

	downloadURLs(urls, rootDir, zipWriter)

	return nil
}

func downloadModels(s SaveFile, rootDir string, zipWriter *zip.Writer) error {

	urls := map[string]bool{}

	for _, o := range s.ObjectStates {
		if o.CustomMesh.MeshURL != "" {
			urls[o.CustomMesh.MeshURL] = true
		}
		if o.CustomMesh.ColliderURL != "" {
			urls[o.CustomMesh.ColliderURL] = true
		}
	}

	downloadURLs(urls, rootDir, zipWriter)

	return nil

}

func downloadURLs(urls map[string]bool, rootDir string, zipWriter *zip.Writer) error {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		return err
	}

	for k := range urls {
		filename := reg.ReplaceAllString(k, "")

		err := downloadFile(k, path.Join(rootDir, filename), zipWriter)
		if err != nil {
			return err
		}
	}

	return nil
}

func downloadFile(url string, filename string, zipWriter *zip.Writer) error {

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	contentType := http.DetectContentType(b)

	if contentType == "image/png" {
		filename += ".png"
	} else if contentType == "image/jpeg" {
		filename += ".jpg"
	} else if contentType == "text/plain; charset=utf-8" {
		filename += ".obj"
	}

	f, err := zipWriter.Create(filename)
	if err != nil {
		return err
	}

	_, err = f.Write(b)
	if err != nil {
		return err
	}

	return nil
}
