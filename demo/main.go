package main

import (
	"os"
	"os/exec"

	demo "github.com/saschagrunert/demo"
)

func main() {
	d := demo.New()

	d.Add(initFermyonSpinProject(), "init-fermyon-spin-project", "Creates an empty Fermyon Spin project")
	d.Add(runWws(), "run-wws", "Run Wasm Workers Server (wws)")
	d.Add(runWwsWithRemoteApp(), "run-wws-with-remote", "Run Wasm Workers Server with remote app (wws)")
	d.Add(buildContainerImage(), "build-container-image", "Builds a WebAssembly container image")
	d.Add(runContainerImage(), "run-container-image", "Runs a WebAssembly container image")

	d.Run()
}

func initFermyonSpinProject() *demo.Run {
	r := demo.NewRun(
		"Create an empty Fermyon Spin project",
	)

	r.Step(demo.S(
		"List available templates",
	), demo.S(
		"spin template list",
	))

	r.Step(demo.S(
		"Create Rust HTTP handler example",
	), demo.S(
		"spin new --template http-rust http-rust-example",
	))

	r.Step(demo.S(
		"Build the example",
	), demo.S(
		"cd http-rust-example && spin build",
	))

	r.Step(demo.S(
		"Serve the example",
	), demo.S(
		"cd http-rust-example && spin up &",
	))

	r.Step(demo.S(
		"Make a request",
	), demo.S(
		"curl -vvv http://localhost:3000",
	))

	r.Step(demo.S(
		"Check spin.toml",
	), demo.S(
		"bat http-rust-example/spin.toml",
	))

	r.Setup(func() error {
		stopSpin()
		cleanupFermyonSpinProject()
		return nil
	})

	r.Cleanup(func() error {
		stopSpin()
		cleanupFermyonSpinProject()
		return nil
	})

	return r
}

func runWws() *demo.Run {
	r := demo.NewRun(
		"Run Wasm Workers Server",
	)

	r.Step(demo.S(
		"Run wws",
	), demo.S(
		"wws ../wws-root &",
	))

	r.Step(demo.S(
		"Make a request",
	), demo.S(
		"curl -vvv http://localhost:8080/endpoint-1",
	))

	r.Setup(func() error {
		stopWws()
		return nil
	})

	r.Cleanup(func() error {
		stopWws()
		return nil
	})

	return r
}

func runWwsWithRemoteApp() *demo.Run {
	r := demo.NewRun(
		"Run Wasm Workers Server with a remote app",
	)

	r.Step(demo.S(
		"Run wws with a remote app",
	), demo.S(
		"wws https://github.com/webassemblylabs/wasm-workers-server.git -i",
		"--git-folder \"examples/ruby-basic\" &",
	))

	r.Step(demo.S(
		"Make a request",
	), demo.S(
		"curl -vvv http://localhost:8080/",
	))

	r.Setup(func() error {
		stopWws()
		return nil
	})

	r.Cleanup(func() error {
		stopWws()
		return nil
	})

	return r
}

func buildContainerImage() *demo.Run {
	r := demo.NewRun(
		"Build Container Image",
	)

	r.Step(demo.S(
		"Program to be built",
	), demo.S(
		"bat ../container-image/main.c",
	))

	r.Step(demo.S(
		"Dockerfile",
	), demo.S(
		"bat ../container-image/Dockerfile",
	))

	r.Step(demo.S(
		"Build container image",
	), demo.S(
		"cd ../container-image &&",
		"docker buildx build --platform wasi/wasm",
		"-t docker.io/ereslibre/wasm-example:0.0.1 .",
	))

	r.Step(demo.S(
		"Push container image",
	), demo.S(
		"docker push docker.io/ereslibre/wasm-example:0.0.1",
	))

	return r
}

func runContainerImage() *demo.Run {
	r := demo.NewRun(
		"Run Container Image",
	)

	r.Step(demo.S(
		"Inspect container image",
	), demo.S(
		"docker inspect docker.io/ereslibre/wasm-example:0.0.1 | jq '.[].Architecture'",
	))

	r.Step(demo.S(
		"Run container image",
	),
		demo.S(
			"docker run --platform=wasi/wasm --runtime=io.containerd.shim.wasmtime.v1",
			"docker.io/ereslibre/wasm-example:0.0.1",
		))

	return r
}

func stopWws() error {
	cmd := exec.Command("pkill", "wws")
	cmd.Run()
	return nil
}

func stopSpin() error {
	cmd := exec.Command("pkill", "spin")
	cmd.Run()
	return nil
}

func cleanupFermyonSpinProject() error {
	os.RemoveAll("http-rust-example")
	return nil
}
