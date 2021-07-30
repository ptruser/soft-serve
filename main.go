package main

import (
	"log"
	"smoothie/server"
	bm "smoothie/server/middleware/bubbletea"
	gm "smoothie/server/middleware/git"
	lm "smoothie/server/middleware/logging"
	"smoothie/tui"

	"github.com/meowgorithm/babyenv"
)

type Config struct {
	Port         int    `env:"SMOOTHIE_PORT" default:"23231"`
	KeyPath      string `env:"SMOOTHIE_KEY_PATH" default:".ssh/smoothie_server_ed25519"`
	RepoAuthPath string `env:"SMOOTHIE_KEY_PATH" default:".ssh/smoothie_git_authorized_keys"`
	RepoPath     string `env:"SMOOTHIE_REPO_PATH" default:".repos"`
}

func main() {
	var cfg Config
	err := babyenv.Parse(&cfg)
	if err != nil {
		log.Fatalln(err)
	}
	s, err := server.NewServer(
		cfg.Port,
		cfg.KeyPath,
		bm.Middleware(tui.SessionHandler(cfg.RepoPath)),
		gm.Middleware(cfg.RepoPath, cfg.RepoAuthPath),
		lm.Middleware(),
	)
	if err != nil {
		log.Fatalln(err)
	}
	err = s.Start()
	if err != nil {
		log.Fatalln(err)
	}
}