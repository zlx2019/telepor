package main

// @Title       main.go
// @Author      Zero.
// @Create      2024-08-09 15:20

func main() {
	server := MixedServer{
		Addr: "0.0.0.0:9946",
	}
	if err := server.Startup(); err != nil {
		Logger.PanicSf("Startup server failed: %s", err)
	}
}
