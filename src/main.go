package main

func main() {
    go wsHub.run()
    newDockerClient()
    newConsulRegistry()
    startWebserver()
}
