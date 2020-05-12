package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Project is struct ri data default about project
var (
	TaskName      = getDataCache("TaskName")
	IDTask        = getDataCache("IDTask")
	IDProject, _  = strconv.Atoi(os.Getenv("PROJECT_ID"))
	branchUpLevel = os.Getenv("BRANCH_UPLEVEL")
	baseAPI       = os.Getenv("SERVER_OP")
	token         = os.Getenv("TOKEN")
	initTime, _   = time.Parse("20060102150405", getDataCache("initTime"))
	endTime, _    = time.Parse("20060102150405", getDataCache("endTime"))
)

func main() {
	args := os.Args

	if len(args[1]) == 0 {
		fmt.Println("You need give some parameter, about more know: 'avia --help'")
	} else {
		switch command := args[1]; command {
		case "--help", "h":
			help()
		case "--list", "l":
			fmt.Println("TODO - this feature")
		case "--resume", "r":
			timer(true)
		case "--pause", "p":
			SaveTimeTask()
		case "--stop", "s":
			timer(false)
		case "--time", "t":
			seeTimer()
		case "--projects", "pp":
			seeProjects()
		case "--clear", "c":
			clearDataCache()
		default:
			_, err := strconv.ParseInt(command, 10, 0)

			// test if is a number
			if err == nil {
				IDTask = args[1]
				setDataCache("IDTask", args[1])

				GetNameTaskOpenProject() // get name task in open project

				if TaskName != "" {
					if len(args) > 2 {
						branchUpLevel = args[2]
					}

					DoGitFlow()
					banner()    // show banner when start the work
					timer(true) // start time
				} else {
					fmt.Println("Fail connect the server Open Project, see the settings")
				}
			} else {
				fmt.Println("Invalid command, about more: 'avia --help'")
			}

		}
	}
}

/**
----------------
	HELP
----------------
**/

// pequeno help de comandos
func help() {
	fmt.Println("@@ Script para facilitar a criação e atuallização de novas Branchs @@")
	fmt.Println("Exemplo:")
	fmt.Println("- Criando uma nova branch atualizada com as ultimas coisas da branch de desenvolvimento")
	fmt.Println("\t$ avia <task_id> [branch-master-development]")
	fmt.Println("- Listar todas as tasks vinculadas a você")
	fmt.Println("\t$ avia --list ou l")
	fmt.Println("- Abrindo o Help")
	fmt.Println("\t$ avia --help ou h")
	fmt.Println("- Pausando o timer da task")
	fmt.Println("\t$ avia --pause ou p")
	fmt.Println("- Retomando a atividade depois de pausar")
	fmt.Println("\t$ avia --resume ou r")
	fmt.Println("- Finalizando o timer da task")
	fmt.Println("\t$ avia --stop ou s")
	fmt.Println("- Visualizando os projetos e Ids de que você faz parte")
	fmt.Println("\t$ avia --projects ou pp")
	fmt.Println("- Limpando a memoria do app")
	fmt.Println("\t$ avia --clear ou c")
}

// # ------------------------------
// #
// # UI
// #
// # ------------------------------
func banner() {
	fmt.Printf(`
                      / \\   \\    / |   / \\
                     / - \\   \\  /  |  / - \\
                    /     \\   \\/   | /     \\
                ----------------------------------

                          [ Work timer ]
            %s        
                        
                ----------------------------------
  	:: avia --help to more info::
  `, TaskName)
}

// # ------------------------------
// #
// # Timer control
// #
// # ------------------------------
func timer(stopClock bool) {
	if stopClock {
		initTime = time.Now()
		setDataCache("initTime", time.Now().Format("20060102150405"))
	} else {
		SaveTimeTask()
	}
}

// CalcTimer calculando o tempo decorrido
func calcTimer() (hs float64, ms float64, ss float64) {
	t1 := initTime
	t2, _ := time.Parse("20060102150405", time.Now().Format("20060102150405"))
	var mf, sf float64

	hs = t1.Sub(t2).Hours() * -1
	ms = mf * 60

	ms, sf = math.Modf(ms)
	ss = sf * 60

	return
}

// seeTimer to see timer decorred
func seeTimer() {
	hs, ms, ss := calcTimer()
	fmt.Printf("Current time is: %.2fh %.2fm %.2fs", hs, ms, ss)
}

// # ------------------------------
// #
// # Github
// #
// # ------------------------------

// DoGitFlow fazendo o fluxo do git para alterar a branch baixar as coisas e atualizar
func DoGitFlow() {
	vars := make(map[string]interface{})
	vars["BranchUpLevel"] = branchUpLevel
	vars["TaskName"] = TaskName

	cmd := exec.Command("bash", "-c", processString("git checkout {{.BranchUpLevel}} && git pull origin {{.BranchUpLevel}} &&  git checkout -b {{.TaskName}}", vars))
	cmd.Env = strings.Split(os.Getenv("PATH"), ":")
	defer cmd.Wait()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatalf("Failed with %s\n", err)
	}
}

// # ------------------------------
// #
// # Get data in open project
// #
// # ------------------------------

// GetNameTaskOpenProject conectando a api do open project e pegando o nome da task
func GetNameTaskOpenProject() {
	result := GetJSON("/work_packages/" + IDTask)
	TaskName = strings.ToLower(IDTask + "_" + strings.ReplaceAll(result["subject"].(string), " ", "-"))
	setDataCache("TaskName", TaskName)
}

// see my projects
func seeProjects() {
	result := GetJSON("/projects")
	v := result["_embedded"].(map[string]interface{})
	for _, data := range v["elements"].([]interface{}) {
		datac := data.(map[string]interface{})
		fmt.Println("ID: ", datac["id"], " Project Name: ", datac["name"])
	}
}

// SaveTimeTask save the time task
func SaveTimeTask() {
	hs, ms, _ := calcTimer()
	conversionTimer := hs + ms/60

	jsonString := `
	{
		"_links":
		{
			"workPackage":
			{
				"href": "/api/v3/work_packages/{{.IDTask}}"
			}
		},
		"hours": "PT{{.ConversionTimer}}H",
		"spentOn": "{{.Time}}"
	}`

	vars := make(map[string]interface{})
	vars["IDTask"] = IDTask
	vars["ConversionTimer"] = float64(int(conversionTimer*100)) / 100
	vars["Time"] = time.Now().Format("2006-01-02")

	// process a template string
	tempatestring := processString(jsonString, vars)
	PostJSON("/time_entries", tempatestring)

	// clear cache start and end time
	setDataCache("initTime", "nil")
	setDataCache("endTime", "nil")
}

// # ------------------------------
// #
// # Utils
// #
// # ------------------------------
func processString(str string, vars interface{}) string {
	tmpl, err := template.New("tmpl").Parse(str)

	if err != nil {
		panic(err)
	}
	return process(tmpl, vars)
}

// process applies the data structure 'vars' onto an already
// parsed template 't', and returns the resulting string.
func process(t *template.Template, vars interface{}) string {
	var tmplBytes bytes.Buffer

	err := t.Execute(&tmplBytes, vars)
	if err != nil {
		panic(err)
	}
	return tmplBytes.String()
}

// prepare login auth
func basicAuth() string {
	auth := "apikey" + ":" + token
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// set data cache save in file temp
func setDataCache(key string, data string) {
	dat, _ := ioutil.ReadFile("/tmp/avia")
	newdata := key + ":" + data
	var errWriter error
	var exists bool = false
	arr := strings.Split(string(dat), ":")

	for i, field := range arr {
		if field == key {
			arr[i+1] = data
			exists = true
		}
	}

	if exists {
		errWriter = ioutil.WriteFile("/tmp/avia", []byte(strings.Join(arr, ":")), 0644)
	} else {
		errWriter = ioutil.WriteFile("/tmp/avia", []byte(strings.Join(arr, ":")+":"+newdata), 0644)
	}

	if errWriter != nil {
		panic(errWriter)
	}
}

// clear storage
func clearDataCache() {
	errWriter := ioutil.WriteFile("/tmp/avia", []byte(""), 0644)

	if errWriter != nil {
		panic(errWriter)
	}
}

// get data cache save in file temp
func getDataCache(key string) string {
	dat, err := ioutil.ReadFile("/tmp/avia")
	if err != nil {
		return ""
	}

	arr := strings.Split(string(dat), ":")
	for i, field := range arr {
		if field == key {
			return arr[i+1]
		}
	}

	return ""
}

// clear the terminal
func clearTerminal() {
	os.Stdout.WriteString("\x1b[3;J\x1b[H\x1b[2J")
}

// # ------------------------------
// #
// # Requests
// #
// # ------------------------------

// GetJSON pega os json e reoassa para quem pedir
func GetJSON(url string) (result map[string]interface{}) {
	// json data
	client := &http.Client{}

	req, _ := http.NewRequest("GET", baseAPI+url, nil)
	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("error", resp.Status)
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(respBody, &result)

	if jsonErr != nil {
		fmt.Println(string(respBody))
		log.Fatal(jsonErr)
	}

	return
}

// PostJSON pega os json e reoassa para quem pedir
func PostJSON(url string, jsonStr string) (result map[string]interface{}) {
	// json data
	client := &http.Client{}

	req, _ := http.NewRequest("POST", baseAPI+url, bytes.NewBuffer([]byte(jsonStr)))
	req.Header.Add("Authorization", basicAuth())
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusMethodNotAllowed || resp.StatusCode == http.StatusInternalServerError {
		fmt.Println("Error aqui ", resp.Status)
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(respBody, &result)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return
}
