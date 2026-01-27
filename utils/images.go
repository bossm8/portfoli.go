// Copyright (c) 2023, Boss Marco <bossm8@hotmail.com>
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its
//    contributors may be used to endorse or promote products derived from
//    this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package utils

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	appconfig "github.com/bossm8/portfoli.go/config"
)

var (
	imageCacheEnabled bool
	imageCacheDir     string
	imageCachePublic  string
)

// SetImageCacheConfig configures image caching and prepares the cache directory.
func SetImageCacheConfig(enabled bool, force bool, cacheDir string) error {
	imageCacheEnabled = enabled
	imageCacheDir = ""
	imageCachePublic = ""

	if !enabled {
		return nil
	}
	if cacheDir == "" {
		cacheDir = filepath.Join("img", "cache")
	}
	if filepath.IsAbs(cacheDir) {
		imageCacheEnabled = false
		return fmt.Errorf("cache dir must be relative to static dir")
	}

	cacheDir = filepath.Join(appconfig.StaticContentPath(), cacheDir)
	absCacheDir, err := filepath.Abs(cacheDir)
	if err != nil {
		imageCacheEnabled = false
		return fmt.Errorf("resolving image cache dir: %w", err)
	}

	rel, err := filepath.Rel(appconfig.StaticContentPath(), absCacheDir)
	if err != nil {
		imageCacheEnabled = false
		return fmt.Errorf("resolving image cache path: %w", err)
	}
	if rel == "." || rel == "" {
		imageCachePublic = "/static"
	} else if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		imageCacheEnabled = false
		return fmt.Errorf("cache dir must be inside static dir (%s)", appconfig.StaticContentPath())
	} else {
		imageCachePublic = filepath.ToSlash(filepath.Join("/static", rel))
	}

	imageCacheDir = absCacheDir
	if force {
		if err := os.RemoveAll(imageCacheDir); err != nil {
			return fmt.Errorf("clearing image cache: %w", err)
		}
	}
	if err := os.MkdirAll(imageCacheDir, 0775); err != nil {
		return fmt.Errorf("creating image cache dir: %w", err)
	}
	return nil
}

// MaybeCacheImage returns the cached local path for a remote image if enabled.
func MaybeCacheImage(image string) string {
	if !imageCacheEnabled {
		return image
	}
	u, err := url.Parse(image)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return image
	}
	if imageCacheDir == "" {
		return image
	}

	filename := cacheFileName(image, u)
	localPath := filepath.Join(imageCacheDir, filename)
	publicPath := filepath.ToSlash(filepath.Join(imageCachePublic, filename))

	if _, err := os.Stat(localPath); err == nil {
		return publicPath
	}
	if err := downloadImage(image, localPath); err != nil {
		log.Printf("[WARNING] Failed to cache image %s: %s\n", image, err)
		return image
	}
	return publicPath
}

// cacheFileName generates a deterministic filename based on the URL and its extension.
func cacheFileName(raw string, u *url.URL) string {
	hash := sha1.Sum([]byte(raw))
	ext := path.Ext(u.Path)
	if ext == "" || len(ext) > 8 {
		ext = ".img"
	}
	return hex.EncodeToString(hash[:]) + ext
}

// downloadImage downloads the image to dest using a temp file for atomic writes.
func downloadImage(src string, dest string) error {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	if err := os.MkdirAll(filepath.Dir(dest), 0775); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(dest), "img-*.tmp")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil && !os.IsNotExist(err) {
			log.Printf("[WARNING] Failed to cleanup temp file %s: %s\n", tmp.Name(), err)
		}
	}()
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), dest)
}
