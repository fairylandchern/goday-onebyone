package main

import (
	"log"
	"os/exec"
	"time"
)

func init() {
	log.Println("plugin init function called")
}

type BadNastyDoctor string

func (g BadNastyDoctor) HealthCheck() error {
	bs, err := exec.Command("bash", "-c", "ls").CombinedOutput()
	if err != nil {
		return err

	}
	log.Println("now is", g)
	log.Println("shell has executed ->>>>>", string(bs))
	return nil
}

//go build -buildmode=plugin -o=plugin_doctor.so plugin_bad_docter.go

// exported as symbol named "Doctor"
var Doctor = BadNastyDoctor(time.Now().Format(time.RFC3339))
