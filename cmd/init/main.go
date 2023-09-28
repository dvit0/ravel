package main

import initPkg "github.com/valyentdev/ravel/internal/init"

func main() {
	// log.Println("started init 1")
	// log.Println("Initialisation...")

	// config, err := os.ReadFile("/valyent/run.json")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// var initConfig initPkg.InitConfig

	// err = json.Unmarshal(config, &initConfig)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// toExec := append(initConfig.ImageConfig.Entrypoint, initConfig.ImageConfig.Cmd...)

	// log.Println(append([]string{"Running"}, toExec...))

	// initPkg.PopulateProcessEnv(append(initConfig.ImageConfig.Env, initConfig.ExtraEnv...))

	// cmd := exec.Cmd{
	// 	Path: toExec[0],
	// 	Args: toExec[1:],
	// 	Env:  append(initConfig.ExtraEnv, initConfig.ImageConfig.Env...),
	// 	Dir:  initConfig.ImageConfig.WorkingDir,
	// }

	// cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	// if err := cmd.Run(); err != nil {
	// 	log.Fatal(err)
	// }

	initPkg.Init()

}
