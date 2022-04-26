package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk-go"
	"github.com/shopspring/decimal"
)

type Config struct {
	AccessToken string
	ClientID    string
	Scope       string
	Verbose     bool

	AssetID            string
	OpponentID         string
	StartTime          time.Time
	EndTime            time.Time
	SnapshotOutputPath string
}

func main() {
	cfg := initConfig()
	snapshots := getSnapshots(context.Background(), cfg)
	if len(snapshots) == 0 {
		log.Println("no records")
		os.Exit(0)
	}

	var total decimal.Decimal
	ids := make([]string, len(snapshots))
	for i, s := range snapshots {
		total = total.Add(s.Amount)
		ids[i] = s.SnapshotID
		if cfg.Verbose {
			fmt.Printf("%s -> (amount: %s, created_at: %s)\n", s.SnapshotID, s.Amount.String(), s.CreatedAt.Format(time.RFC3339))
		}
	}
	fmt.Println()

	fmt.Printf("ids: (%s)\n\n", strings.Join(ids, ", "))
	fmt.Printf("count: %d, total: %s\n", len(snapshots), total)
	if cfg.SnapshotOutputPath != "" {
		err := ioutil.WriteFile(cfg.SnapshotOutputPath, []byte(strings.Join(ids, ",")), 0644)
		fatalIfErr(err)
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
	var startTime, endTime string

	flag.StringVar(&cfg.AccessToken, "token", "", "Access token")
	flag.StringVar(&cfg.ClientID, "client", "", "Mixin client id")
	flag.StringVar(&cfg.Scope, "scope", "", "Mixin oauth scope (except SNAPSHOTS:READ)")
	flag.BoolVar(&cfg.Verbose, "verbose", false, "Verbose log")

	flag.StringVar(&cfg.AssetID, "asset", "", "Asset id")
	flag.StringVar(&cfg.OpponentID, "opponent", "", "Opponent id")
	flag.StringVar(&startTime, "start", "", "Start time, RFC3339 format")
	flag.StringVar(&endTime, "end", "", "End time, RFC3339 format")
	flag.StringVar(&cfg.SnapshotOutputPath, "output", "", "Snapshots id output path")

	flag.Parse()

	if cfg.AssetID == "" {
		log.Fatalln("asset not set")
	}

	if cfg.OpponentID == "" {
		log.Fatalln("opponent not set")
	}

	if startTime != "" {
		cfg.StartTime = mustParseTime(startTime)
	}
	if endTime != "" {
		cfg.EndTime = mustParseTime(endTime)
	}

	if cfg.AccessToken == "" {
		if cfg.ClientID == "" {
			log.Fatalln("you must set one of token and client")
		}

		scope := cfg.Scope
		if scope != "" {
			scope += "+"
		}
		scope += "SNAPSHOTS:READ"
		err := openBrowser(fmt.Sprintf("https://www.mixin.one/oauth/authorize?client_id=%s&scope=%s", cfg.ClientID, scope))
		fatalIfErr(err)

		fmt.Println("Enter Your Token: ")
		_, err = fmt.Scan(&cfg.AccessToken)
		fatalIfErr(err)
		fmt.Println()
	}

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
