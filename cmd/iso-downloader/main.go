package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	pb "github.com/schollz/progressbar/v3"
)

type Distro struct {
	BaseURL          string   `json:"base_url"`
	URLTemplate      string   `json:"url_template"`
	FilenameTemplate string   `json:"filename_template"`
	ListURLTemplate  string   `json:"list_url_template"`
	FilenameGlob     string   `json:"filename_glob"`
	Versions         []string `json:"versions"`
}

type Family struct {
	Distros map[string]Distro `json:"distros"`
}

type Config struct {
	Families map[string]Family `json:"families"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	familyKey, distroKey, versions, err := promptSelection(cfg)
	if err != nil {
		return err
	}
	urls, err := resolveURLs(cfg, familyKey, distroKey, versions)
	if err != nil {
		return err
	}
	out, err := promptOutputDir(distroKey)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(out, 0o755); err != nil {
		return err
	}
	for _, u := range urls {
		if err := download(u.URL, filepath.Join(out, u.Filename)); err != nil {
			return err
		}
	}
	fmt.Println("Saved to:", out)
	return nil
}

func promptOutputDir(distro string) (string, error) {
	def := filepath.Join("downloads", distro)
	// Try to use yazi as a directory chooser if available
	if _, err := exec.LookPath("yazi"); err == nil {
		cmd := exec.Command("yazi", "--chooser-dir")
		cmd.Stdin = os.Stdin
		out, err := cmd.Output()
		if err == nil {
			chosen := strings.TrimSpace(string(out))
			if chosen != "" {
				return chosen, nil
			}
		}
	}
	// fallback simple input
	var ans string
	if err := survey.AskOne(&survey.Input{Message: "Save directory:", Default: def}, &ans); err != nil {
		return "", err
	}
	ans = strings.TrimSpace(ans)
	if ans == "" {
		ans = def
	}
	return ans, nil
}

func loadConfig() (Config, error) {
	var cfg Config
	path := os.Getenv("ISO_DOWNLOADER_CONFIG")
	if path == "" {
		// Prefer JSON in Go data path
		primary := filepath.Join("data", "distros.json")
		if _, err := os.Stat(primary); err == nil {
			path = primary
		} else {
			return cfg, fmt.Errorf("config not found: %s", primary)
		}
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	if err := json.Unmarshal(b, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

func promptSelection(cfg Config) (string, string, []string, error) {
	if len(cfg.Families) == 0 {
		return "", "", nil, errors.New("no families configured")
	}
	familyKeys := make([]string, 0, len(cfg.Families))
	for k := range cfg.Families {
		familyKeys = append(familyKeys, k)
	}
	sort.Strings(familyKeys)
	var family string
	if err := survey.AskOne(&survey.Select{Message: "Select a family:", Options: familyKeys}, &family); err != nil {
		return "", "", nil, err
	}
	distros := cfg.Families[family].Distros
	if len(distros) == 0 {
		return "", "", nil, fmt.Errorf("no distros under family %s", family)
	}
	distroKeys := make([]string, 0, len(distros))
	for k := range distros {
		distroKeys = append(distroKeys, k)
	}
	sort.Strings(distroKeys)
	var distro string
	if err := survey.AskOne(&survey.Select{Message: "Select a distribution:", Options: distroKeys}, &distro); err != nil {
		return "", "", nil, err
	}
	versions := distros[distro].Versions
	opts := append([]string{"All versions"}, versions...)
	var chosen []string
	if err := survey.AskOne(&survey.MultiSelect{Message: "Select versions (space to toggle, enter to confirm):", Options: opts}, &chosen); err != nil {
		return "", "", nil, err
	}
	if len(chosen) == 0 || contains(chosen, "All versions") {
		return family, distro, versions, nil
	}
	return family, distro, chosen, nil
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

type URLItem struct {
	URL      string
	Filename string
	Version  string
}

func resolveURLs(cfg Config, familyKey string, distroKey string, versions []string) ([]URLItem, error) {
	fam, ok := cfg.Families[familyKey]
	if !ok {
		return nil, fmt.Errorf("family not found: %s", familyKey)
	}
	d, ok := fam.Distros[distroKey]
	if !ok {
		return nil, fmt.Errorf("distro not found: %s", distroKey)
	}
	var items []URLItem
	for _, v := range versions {
		if d.ListURLTemplate != "" && d.FilenameGlob != "" {
			listURL := strings.NewReplacer("{base_url}", d.BaseURL, "{version}", v, "{distro}", distroKey).Replace(d.ListURLTemplate)
			files, err := listDir(listURL)
			if err != nil {
				return nil, err
			}
			re := globToRegexp(strings.NewReplacer("{version}", v, "{distro}", distroKey).Replace(d.FilenameGlob))
			match := ""
			for _, f := range files {
				if re.MatchString(f) {
					match = f
					break
				}
			}
			if match == "" {
				return nil, fmt.Errorf("no ISO matching %q at %s", d.FilenameGlob, listURL)
			}
			items = append(items, URLItem{URL: strings.TrimRight(listURL, "/") + "/" + match, Filename: match, Version: v})
		} else if d.URLTemplate != "" && d.FilenameTemplate != "" {
			url := strings.NewReplacer("{base_url}", d.BaseURL, "{version}", v, "{distro}", distroKey).Replace(d.URLTemplate)
			filename := strings.NewReplacer("{version}", v, "{distro}", distroKey).Replace(d.FilenameTemplate)
			items = append(items, URLItem{URL: url, Filename: filename, Version: v})
		} else {
			return nil, errors.New("invalid config for distro: missing templates")
		}
	}
	return items, nil
}

func listDir(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	// naive href extraction
	re := regexp.MustCompile(`href=["']([^"']+)["']`)
	matches := re.FindAllStringSubmatch(string(b), -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		if len(m) > 1 && strings.HasSuffix(strings.ToLower(m[1]), ".iso") {
			out = append(out, m[1])
		}
	}
	return out, nil
}

func globToRegexp(glob string) *regexp.Regexp {
	// very small glob: * -> .*
	re := regexp.QuoteMeta(glob)
	re = strings.ReplaceAll(re, "\\*", ".*")
	return regexp.MustCompile("^" + re + "$")
}

func download(url, dest string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()
	// We won't use ctx directly with http.Get, but keep structure for future
	_ = ctx
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed: %s", resp.Status)
	}

	// Determine total size
	total := resp.ContentLength
	var bar *pb.ProgressBar
	if total > 0 {
		bar = pb.NewOptions64(total,
			pb.OptionSetDescription(filepath.Base(dest)),
			pb.OptionSetWidth(30),
			pb.OptionShowBytes(true),
			pb.OptionSetPredictTime(true),
			pb.OptionThrottle(65*time.Millisecond),
			pb.OptionShowCount(),
		)
	} else {
		bar = pb.NewOptions(-1,
			pb.OptionSetDescription(filepath.Base(dest)),
			pb.OptionSetWidth(30),
		)
	}
	defer bar.Finish()

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	buf := make([]byte, 256*1024)
	for {
		n, rerr := resp.Body.Read(buf)
		if n > 0 {
			if _, werr := f.Write(buf[:n]); werr != nil {
				return werr
			}
			bar.Add(n)
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return rerr
		}
	}
	return nil
}
