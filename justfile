fermyon-spin:
  clear
  @cd demo && go run . --init-fermyon-spin-project

wws:
  clear
  @cd demo && go run . --run-wws

wws-with-remote:
  clear
  @cd demo && go run . --run-wws-with-remote

build-image:
  clear
  @cd demo && go run . --build-container-image

run-image:
  clear
  @cd demo && go run . --run-container-image