build-image:
  clear
  @cd demo && go run . --build-container-image

run-image:
  clear
  @cd demo && go run . --run-container-image

wws:
  clear
  @cd demo && go run . --run-wws

wws-with-remote:
  clear
  @cd demo && go run . --run-wws-with-remote

run-php-script:
  clear
  @cd demo && go run . --run-php-script

fermyon-spin:
  clear
  @cd demo && go run . --init-fermyon-spin-project

spinkube:
  clear
  @cd demo && go run . --run-spinkube

fmt:
  find . -name "*.nix" | xargs alejandra
  find demo -name "*.go" | xargs gofmt -w
