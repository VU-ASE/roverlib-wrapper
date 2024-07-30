package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	pb_core_messages "github.com/VU-ASE/rovercom/packages/go/core"
	roverlib "github.com/VU-ASE/roverlib/src"

	"github.com/rs/zerolog/log"
)

func run(
	service roverlib.ResolvedService,
	sysMan roverlib.CoreInfo,
	initialTuning *pb_core_messages.TuningState) error {

	// The last argument should always be the path to the service or command to run
	// if the number of arguments is odd, then the last argument is an option
	cmdIndex := len(os.Args) - 1
	if len(os.Args) < 2 || len(os.Args)%2 == 1 || os.Args[cmdIndex] == "" || strings.HasPrefix(os.Args[cmdIndex], "-") {
		fmt.Printf("\nNo service or command to run specified\n\n\tUsage: %s <service or command> [args...]\n\n", os.Args[0])

		return fmt.Errorf("no service or command to run specified")
	}

	serviceToRun := os.Args[cmdIndex]
	log.Info().Msgf("Running service '%s' with command '%s'", service.Name, serviceToRun)

	// Prepare the environment
	cmd := exec.Command("bash", "-c", serviceToRun)

	// According to the docs, these values should be set
	cmd.Env = os.Environ()

	// Service name
	nameEnv := fmt.Sprintf("ASE_SW_ServiceName=%s", service.Name)
	cmd.Env = append(cmd.Env, nameEnv)
	// Service PID
	pidEnv := fmt.Sprintf("ASE_SW_ServicePID=%d", os.Getpid())
	cmd.Env = append(cmd.Env, pidEnv)

	// Outputs
	for _, output := range service.Outputs {
		outputEnv := fmt.Sprintf("ASE_SW_Output_%s=%s", output.Name, output.Address)
		cmd.Env = append(cmd.Env, outputEnv)
	}

	// Tuning parameters
	for _, param := range initialTuning.DynamicParameters {
		if param.GetInt() != nil {
			typeEnv := fmt.Sprintf("ASE_SW_TuningParameterType_%s=Int", param.GetInt().Key)
			cmd.Env = append(cmd.Env, typeEnv)
			valueEnv := fmt.Sprintf("ASE_SW_TuningParameterValue_%s=%d", param.GetInt().Key, param.GetInt().Value)
			cmd.Env = append(cmd.Env, valueEnv)
		} else if param.GetFloat() != nil {
			typeEnv := fmt.Sprintf("ASE_SW_TuningParameterType_%s=Float", param.GetFloat().Key)
			cmd.Env = append(cmd.Env, typeEnv)
			valueEnv := fmt.Sprintf("ASE_SW_TuningParameterValue_%s=%f", param.GetFloat().Key, param.GetFloat().Value)
			cmd.Env = append(cmd.Env, valueEnv)
		} else if param.GetString_() != nil {
			typeEnv := fmt.Sprintf("ASE_SW_TuningParameterType_%s=String", param.GetString_().Key)
			cmd.Env = append(cmd.Env, typeEnv)
			valueEnv := fmt.Sprintf("ASE_SW_TuningParameterValue_%s=%s", param.GetString_().Key, param.GetString_().Value)
			cmd.Env = append(cmd.Env, valueEnv)
		}
	}

	// Resolved dependencies
	for _, dep := range service.Dependencies {
		depEnv := fmt.Sprintf("ASE_SW_Dependency_%s_%s=%s", dep.ServiceName, dep.OutputName, dep.Address)
		cmd.Env = append(cmd.Env, depEnv)
	}

	// Run the command and log the output
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	err = cmd.Start()
	if err != nil {
		return err
	}

	fmt.Printf("\nThe output of your program will be logged here\n💡 If you do not see output here, make sure to flush stdout and stderr in your program\n\n")

	go logOutput(stdout)
	go logOutput(stderr)
	return cmd.Wait()
}

func logOutput(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
}

func onTuningState(newtuning *pb_core_messages.TuningState) {
	log.Warn().Msg("Received OTA tuning update. This will be ignored since programs launched through mod-ServiceWrapper do not support OTA tuning updates.")
}

func onTerminate(signal os.Signal) {
	log.Warn().Msgf("Terminating due to signal: %s", signal)
}

func main() {
	roverlib.Run(run, onTuningState, onTerminate, false)
}
