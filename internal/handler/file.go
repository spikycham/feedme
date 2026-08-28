package handler

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	"github.com/spikycham/feedme/internal/constant"
	"github.com/spikycham/feedme/pkg/network"
)

type FileHandler struct{}

func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

type ResponseUpload struct {
	URL string `json:"url"`
}

func (f *FileHandler) Upload(w http.ResponseWriter, r *http.Request) error {
	file, header, err := r.FormFile("file")
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	hash := md5.Sum(data)
	hashStr := hex.EncodeToString(hash[:])

	// Validate the file extension.
	li := strings.LastIndex(header.Filename, ".")
	if li == -1 {
		network.Error(w, http.StatusBadRequest)
		return constant.InvalidFileExt
	}
	ext := strings.ToLower(header.Filename[li+1:])
	if !slices.Contains([]string{"jpg", "jpeg", "png", "webp"}, ext) {
		network.Error(w, http.StatusBadRequest)
		fmt.Print("EXT: ", ext)
		return constant.InvalidFileExt
	}
	filename := hashStr + "." + ext

	entries, err := os.ReadDir("upload/")
	if err != nil {
		network.Error(w, http.StatusForbidden)
		return err
	}

	prefix := os.Getenv("ASSET_URL_PREFIX")
	resp := ResponseUpload{URL: prefix + filename}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == filename {
			network.Write(w, &resp)
			return nil
		}
	}

	// TODO: bind the uploads/ to ~/assets/feedme.
	dst := "upload/" + filename
	if err := os.WriteFile(dst, data, 0o644); err != nil {
		network.Error(w, http.StatusInternalServerError)
		return err
	}

	network.Write(w, &resp)
	return nil
}
