package main

import (
	stdContext "context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/kataras/iris"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	godemo "github.com/lsagiroglu/godemo"
)

var GoVersion string

func init() {

	PrgPath, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	GoVersion = runtime.Version()

	//	http://patorjk.com/software/taag/#p=display&f=Standard&t=GoDemo

	fmt.Println("   ____       ____                       ")
	fmt.Println("  / ___| ___ |  _ \\  ___ _ __ ___   ___  ")
	fmt.Println(" | |  _ / _ \\| | | |/ _ \\ '_ ` _ \\ / _ \\ ")
	fmt.Println(" | |_| | (_) | |_| |  __/ | | | | | (_) |")
	fmt.Printf("  \\____|\\___/|____/ \\___|_| |_| |_|\\___/   %s\n", godemo.Version)

	fmt.Println()
	fmt.Printf("Go   : %s\n", GoVersion)
	fmt.Printf("Iris : %s\n", iris.Version)
	fmt.Printf("Path : %s\n", PrgPath)
	fmt.Println()
}

func main() {

	app := iris.New()
	// Optionally, add two built'n handlers
	// that can recover from any http-relative panics
	// and log the requests to the terminal.
	app.Use(recover.New())
	app.Use(logger.New())

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		app.Logger().Warn("DOMAIN environment variable was not set")
	}

	email := os.Getenv("EMAIL")
	if domain == "" {
		app.Logger().Warn("EMAIL environment variable was not set")
	}

	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch,
			// kill -SIGINT XXXX or Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill  is equivalent with the syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
			syscall.SIGQUIT,
		)
		select {

		case msg := <-ch:
			fmt.Printf("shutdown: %s", msg)

			timeout := 5 * time.Second
			ctx, cancel := stdContext.WithTimeout(stdContext.Background(), timeout)
			defer cancel()
			app.Shutdown(ctx)
		}
	}()

	// NOTE: This will not work on domains like this,
	// use real whitelisted domain(or domains split by whitespaces)
	// and a non-public e-mail instead.

	if domain == "" {
		app.Any("/", godemo.InfoPage)
		app.Run(iris.Addr(":80"), iris.WithoutVersionChecker, iris.WithoutServerError(iris.ErrServerClosed))
	} else {
		app.Any("/", godemo.HomePage)
		app.Any("/ping", godemo.PingPage)

		app.Run(iris.AutoTLS(":443", domain, email), iris.WithoutVersionChecker, iris.WithoutServerError(iris.ErrServerClosed))
	}

}
