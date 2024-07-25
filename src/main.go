package main

import (
	pb_core_messages "github.com/VU-ASE/rovercom/packages/go/core"
	servicerunner "github.com/VU-ASE/roverlib/src"

	"os"
	"github.com/rs/zerolog/log"
)

func run(
	service servicerunner.ResolvedService,
	sysMan servicerunner.SystemManagerInfo,
	initialTuning *pb_core_messages.TuningState) error {

	log.Info().Str("Planet", "Earth").Msg("Hello world")

	//TODO: Implement the service logic here. Likely this will involve creating a pub/sub and some main logic.
	//      The de facto standard is to have some read (zmq/IO), some handling logic (may be several items),
	//      and some write (zmq/IO). The go routines typically communicate via channels.

	return nil
}

func onTuningState(newtuning *pb_core_messages.TuningState) {
	log.Info().Str("Value", newtuning.String()).Msg("Received tuning state from system manager")
	//TODO: Update this service based on the new tuning state
}



func onTerminate(signal os.Signal) {
	log.Info().Msg("onTerminate callback!")
}


func main() {
	servicerunner.Run(run, onTuningState, onTerminate, false)
}
