package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/da4nik/feedback/internal/config"
	"github.com/da4nik/feedback/internal/log"
	"github.com/da4nik/feedback/internal/server"
	"github.com/da4nik/feedback/pkg/handlers"
	"github.com/da4nik/feedback/pkg/mandrill"
	"github.com/da4nik/feedback/pkg/text"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	conf := config.LoadConfig()

	log.InitLogger(log.LoggerOpts{})

	txt := text.NewText(conf.TwilioPhone,
		conf.TwilioSID,
		conf.TwilioKey)

	mandr, err := mandrill.NewMandrill(conf.MandrillKey)
	if err != nil {
		log.Errorf("Unable to create Mandrill email proxy: %s", err.Error())
		return
	}

	hndlr := handlers.NewHandlers(
		txt,
		mandr,
		conf.TargetEmail,
		[]string{
			"captureproof.com",
			"cp2.div-art.com.ua",
		})

	httpServer := server.NewServer(conf.Port, hndlr)
	httpServer.Start()

	log.Infof("Starting on port %d", conf.Port)
	log.Infof("Crtl-C to interrupt.")
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT)
	<-sigs

	log.Infof("Exiting...")
}
