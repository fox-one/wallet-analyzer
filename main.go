package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/manifoldco/promptui"
	"github.com/shopspring/decimal"
)

type Config struct {
	AccessToken  string
	ClientID     string
	ClientSecret string

	AssetID        string
	OpponentID     string
	StartTime      time.Time
	EndTime        time.Time
	OutputPath     string
	FormatTemplate *template.Template
}

func main() {
	cfg := initConfig()
	snapshots := getSnapshots(context.Background(), cfg)
	if len(snapshots) == 0 {
		log.Println("no records")
		os.Exit(0)
	}

	writer := new(bytes.Buffer)

	ids := make([]string, len(snapshots))
	asm := map[string][]*mixin.Snapshot{}
	for i, s := range snapshots {
		asm[s.AssetID] = append(asm[s.AssetID], s)
		ids[i] = fmt.Sprintf("'%s'", s.SnapshotID)
		tpl := new(bytes.Buffer)
		fatalIfErr(cfg.FormatTemplate.Execute(tpl, s))
		fmt.Fprintf(writer, "%s\n", tpl.String())
	}

	fmt.Fprintf(writer, "\nids: (%s)\n\n", strings.Join(ids, ", "))

	for asset, ss := range asm {
		count := len(ss)
		var total decimal.Decimal
		for _, s := range ss {
			total = total.Add(s.Amount)
		}

		fmt.Fprintf(writer, "asset: %s -> (count: %d, total: %s)\n", asset, count, total)
	}

	if cfg.OutputPath != "" {
		err := ioutil.WriteFile(cfg.OutputPath, writer.Bytes(), 0644)
		fatalIfErr(err)
	} else {
		io.Copy(os.Stdout, writer)
	}
}

func getSnapshots(ctx context.Context, cfg *Config) []*mixin.Snapshot {
	var result []*mixin.Snapshot
	limit := 500

	offset := cfg.StartTime
	for {
		ss, err := mixin.ReadSnapshots(ctx, cfg.AccessToken, cfg.AssetID, offset, "ASC", limit)
		fatalIfErr(err)
		for _, s := range ss {
			if !cfg.EndTime.IsZero() && s.CreatedAt.After(cfg.EndTime) {
				break
			}

			if s.OpponentID == cfg.OpponentID {
				result = append(result, s)
			}
		}

		if len(ss) < limit {
			break
		}
		offset = ss[len(ss)-1].CreatedAt
	}

	return result
}

func initConfig() *Config {
	cfg := &Config{}
	var startTime, endTime, format string

	flag.StringVar(&cfg.AccessToken, "token", "", "Access token")
	flag.StringVar(&cfg.ClientID, "client", "", "Mixin client id")
	flag.StringVar(&cfg.ClientSecret, "secret", "", "Mixin client secret")

	flag.StringVar(&cfg.AssetID, "asset", "", "Asset id")
	flag.StringVar(&cfg.OpponentID, "opponent", "", "Opponent id")
	flag.StringVar(&startTime, "start", "", "Start time, RFC3339 format")
	flag.StringVar(&endTime, "end", "", "End time, RFC3339 format")
	flag.StringVar(&cfg.OutputPath, "output", "", "Output file path")
	flag.StringVar(&format, "format", "id: {{ .SnapshotID }} -> (asset: {{ .AssetID }}, amount: {{ .Amount }})", "Snapshot format")

	flag.Parse()

	if cfg.AccessToken == "" && (cfg.ClientID == "" || cfg.ClientSecret == "") {
		log.Fatalln("you must set one of token and (client+secret)")
	}

	if startTime != "" {
		cfg.StartTime = mustParseTime(startTime)
	}
	if endTime != "" {
		cfg.EndTime = mustParseTime(endTime)
	}

	if cfg.AccessToken == "" {
		scope := "SNAPSHOTS:READ+PROFILE:READ"
		err := openBrowser(fmt.Sprintf("https://mixin-oauth.fox.one?client_id=%s&scope=%s", cfg.ClientID, scope))
		fatalIfErr(err)

		prompt := promptui.Prompt{
			Label: "OAuth Code",
		}
		code, err := prompt.Run()
		fatalIfErr(err)

		cfg.AccessToken, _, err = mixin.AuthorizeToken(context.Background(), cfg.ClientID, cfg.ClientSecret, code, "")
		fatalIfErr(err)
		fmt.Printf("\ntoken: %s\n\n", cfg.AccessToken)
	}

	cfg.FormatTemplate = template.Must(template.New("").Parse(format))

	return cfg
}

func mustParseTime(v string) time.Time {
	t, err := time.Parse(time.RFC3339, v)
	fatalIfErr(err)
	return t
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
