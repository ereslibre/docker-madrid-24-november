package main

import (
	"os"
	"os/exec"

	demo "github.com/saschagrunert/demo"
)

func main() {
	d := demo.New()

	d.Add(runStrace(), "strace", "Strace example")
	d.Add(initFermyonSpinProject(), "init-fermyon-spin-project", "Creates an empty Fermyon Spin project")
	d.Add(runWws(), "run-wws", "Run Wasm Workers Server (wws)")
	d.Add(runWwsWithRemoteApp(), "run-wws-with-remote", "Run Wasm Workers Server with remote app (wws)")
	d.Add(buildContainerImage(), "build-container-image", "Builds a WebAssembly container image")
	d.Add(runContainerImage(), "run-container-image", "Runs a WebAssembly container image")
	d.Add(runPHPScript(), "run-php-script", "Runs a PHP script")
	d.Add(runSpinKube(), "run-spinkube", "Runs SpinKube")

	d.Run()
}

func runStrace() *demo.Run {
	r := demo.NewRun(
		"Strace HTTP request example",
	)

	r.Step(demo.S(
		"Perform HTTP request with summarized output",
	), demo.S(
		"strace -c --follow-forks -- ",
		"curl -s https://www.google.com 1> /dev/null",
	))

	r.Step(demo.S(
		"Perform HTTP request with regular output",
	), demo.S(
		"strace --follow-forks -- ",
		"curl -s https://www.google.com 1> /dev/null",
	))

	return r
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
		"List wws endpoints",
	), demo.S(
		"tree ../wws-root",
	))

	r.Step(demo.S(
		"Show function contents",
	), demo.S(
		"bat ../wws-root/endpoint-1.js",
	))

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
		"Show function contents",
	), demo.S(
		"curl https://raw.githubusercontent.com/webassemblylabs/wasm-workers-server/refs/heads/main/examples/ruby-basic/index.rb |",
		"bat --language=ruby",
	))

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
		"Build program",
	), demo.S(
		"clang --target=wasm32-wasi ../container-image/main.c -o ../container-image/hello-docker.wasm",
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
		"Pull container image",
	), demo.S(
		"docker pull --platform=wasi/wasm docker.io/ereslibre/wasm-example:0.0.1",
	))

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

func runPHPScript() *demo.Run {
	r := demo.NewRun(
		"Run a PHP script",
	)

	r.Step(demo.S(
		"Inspect PHP interpreter wasm32-wasi module",
	), demo.S(
		"file php-cgi-8.2.6.wasm",
	))

	r.Step(demo.S(
		"Inspect PHP interpreter wasm32-wasi module",
	), demo.S(
		"du -hs php-cgi-8.2.6.wasm",
	))

	r.Step(demo.S(
		"Inspect PHP interpreter wasm32-wasi module metadata",
	), demo.S(
		"wasm-tools metadata show php-cgi-8.2.6.wasm",
	))

	r.Step(demo.S(
		"Inspect PHP script",
	), demo.S(
		"bat script.php",
	))

	r.Step(demo.S(
		"Run PHP script",
	), demo.S(
		"wasmtime run --dir .::/demo",
		"php-cgi-8.2.6.wasm -- /demo/script.php",
	))

	return r
}

func runSpinKube() *demo.Run {
	r := demo.NewRun(
		"Run SpinKube",
	)

	r.Step(demo.S(
		"Inspect manifest",
	), demo.S(
		"bat spinkube-app.yaml",
	))

	r.Step(demo.S(
		"Deploy SpinKube app",
	), demo.S(
		"kubectl apply -f spinkube-app.yaml",
	))

	r.Step(demo.S(
		"Wait for SpinKube app to be ready",
	), demo.S(
		"kubectl wait --for=jsonpath='{.status.readyReplicas}'=1",
		"spinapp.core.spinoperator.dev/simple-spinapp",
	))

	r.Step(demo.S(
		"Forward port",
	), demo.S(
		"kubectl port-forward services/simple-spinapp 8090:80 &",
	))

	r.Step(demo.S(
		"Make a request (Rust endpoint)",
	), demo.S(
		"curl -vvv http://localhost:8090/hello",
	))

	r.Step(demo.S(
		"Make a request (Go endpoint)",
	), demo.S(
		"curl -vvv http://localhost:8090/go-hello",
	))

	r.Setup(func() error {
		cleanupSpinKube()
		return nil
	})

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

func cleanupSpinKube() error {
	cmd := exec.Command("pkill", "kubectl")
	cmd.Run()
	cmd = exec.Command("kubectl", "delete", "-A", "--all", "spinapp.core.spinoperator.dev")
	cmd.Run()
	return nil
}
