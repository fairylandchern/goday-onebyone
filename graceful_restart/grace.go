package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

// graceful restart app using flag to identify
var (
	graceful bool
	ln       net.Listener
	srv      *http.Server
)

func init() {
	flag.BoolVar(&graceful, "upgrade", false, "need graceful restart")
}

//first step construct the infrastructure
func main() {
	flag.Parse()
	srv = &http.Server{Addr: ":8080"}
	var err error
	if graceful {
		fd := os.NewFile(3, "") //附加的文件描述符，界定为3
		ln, err = net.FileListener(fd)
		if err != nil {
			log.Fatal("err open file listener:", err)
			return
		}
		log.Println("extends from parent")
		fd.Close()
	} else {
		ln, err = net.Listen("tcp", srv.Addr)
		if err != nil {
			log.Fatal("err when listen tcp:", err)
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		//log.Println("process request now")
		// print pid,ppid
		writer.Write([]byte(fmt.Sprintf("hello,pid:%v,ppid:%v", syscall.Getpid(), syscall.Getppid())))
	})

	mux.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte(fmt.Sprintf("imporved in next version,say hello,world,pid:%v,ppid:%v",syscall.Getpid(),syscall.Getppid())))
	})
	srv.Handler = mux

	go func() {
		// can mov process to goroutine
		err = srv.Serve(ln)
		if err != nil {
			log.Fatal("err when version1 work:", err)
		}
	}()

	// add side process
	setupsignal()
	log.Println("over")
}

func setupsignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR2, syscall.SIGINT, syscall.SIGTERM)
	sig := <-ch
	log.Println("get signal now:", sig)
	fmt.Println("get signal now:", sig)
	switch sig {
	case syscall.SIGUSR2:
		log.Println("fork process now:", os.Getppid(), " pid:", os.Getpid())
		fmt.Println("fork process now:", os.Getppid(), " pid:", os.Getpid())
		err := forkprocess()
		err = srv.Shutdown(context.Background())
		if err != nil {
			log.Println("fork err when shutdown srv:", err)
		}
		log.Println("fork process now:", os.Getppid(), " pid:", os.Getpid())
		fmt.Println("fork process now:", os.Getppid(), " pid:", os.Getpid())
	case syscall.SIGINT, syscall.SIGTERM:
		signal.Stop(ch)
		close(ch)
		err := srv.Shutdown(context.Background())
		if err != nil {
			log.Println("shutdown err:", err)
		}
		log.Println("shutdown now,pid:", os.Getpid(), " ppid:", os.Getppid())
		fmt.Println("shutdown now,pid:", os.Getpid(), " ppid:", os.Getppid())
	}
}

func forkprocess() error {
	flags := []string{"-upgrade"}
	cmd := exec.Command(os.Args[0], flags...)
	//cmd.Stderr=os.Stderr
	//cmd.Stdout=os.Stdout
	l, _ := ln.(*net.TCPListener)
	lfd, err := l.File()
	if err != nil {
		log.Fatal("fork err get lfd listener:", err)
		return err
	}
	cmd.ExtraFiles = []*os.File{lfd}
	return cmd.Start()
}
