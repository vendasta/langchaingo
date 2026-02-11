module github.com/vendasta/langchaingo/examples/openai-vision-example

go 1.23

toolchain go1.24.0

require github.com/vendasta/langchaingo v0.1.12-pre.6

replace github.com/vendasta/langchaingo => ../../

require (
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pkoukk/tiktoken-go v0.1.8 // indirect
)
