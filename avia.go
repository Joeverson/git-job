package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Project is struct ri data default about project
var (
	TaskName      = ""
	IDTask        = ""
	IDProject, _  = strconv.Atoi(os.Getenv("PROJECT_ID"))
	branchUpLevel = os.Getenv("BRANCH_UPLEVEL")
	baseAPI       = os.Getenv("SERVER_OP")
	token         = os.Getenv("TOKEN")
)

func main() {
	args := os.Args

	if len(args[1]) == 0 {
		fmt.Println("You need give some parameter, about more know: 'avia --help'")
	} else {
		switch command := args[1]; command {
		case "--help":
			help()
		case "--list":
			help()
		case "--stop":
			timer(false)
		default:
			_, err := strconv.ParseInt(command, 10, 0)

			// caso não seja um inteiro(code open project tasks) ele fala comando invalido
			if err == nil {
				IDTask = args[1]
				GetNameTaskOpenProject()

				if TaskName != "" {
					if len(args) > 2 {
						branchUpLevel = args[2]
					}

					DoGitFlow()
					// timer(true)
				} else {
					fmt.Println("Fail connect the server Open Project, see the settings")
				}
			} else {
				fmt.Println("Invalid command, about more know: 'avia --help'")
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
	fmt.Println("\n @@ Script para facilitar a criação e atuallização de novas Branchs @@")
	fmt.Println("\n Exemplo:")
	fmt.Println("\n- Criando uma nova branch atualizada com as ultimas coisas da branch de desenvolvimento")
	fmt.Println("\n\t$ avia [task_id] [branch-master-development]")
	fmt.Println("\n- Listar todas as tasks vinculadas a você")
	fmt.Println("\n\t$ avia --list")
	fmt.Println("\n- Abrindo o Help")
	fmt.Println("\n\t$ avia --help")
	fmt.Println("\n- Finalizando o timer da task")
	fmt.Println("\n\t$ avia --stop")
}

// # ------------------------------
// #
// # TIMER
// #
// # ------------------------------

// # DISPLAY TIMER

// func screenTimer() {
//   CallClear()
//   string display = `
//                       / \\   \\    / |   / \\
//                      / - \\   \\  /  |  / - \\
//                     /     \\   \\/   | /     \\
//                 ----------------------------------

//                           [ Work timer ]
//                           $TASK_NAME        "
//   printf "                            02dh:02dm:02ds         " $h $m $s
//   echo -e "\n
//                 ----------------------------------"
//   echo ":: 'd' for done task, 'c' for continue and 'p' para pause ::"
//   `

//   fmt.Println(display)
// }

// func CallClear() {
//     value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
//     if ok { //if we defined a clear func for that platform:
//         value()  //we execute it
//     } else { //unsupported platform
//         panic("Your platform is unsupported! I can't clear terminal screen :(")
//     }
// }

// func clock() {
//   _screenTimer
//   sleep 1
//   s=$((s+1))
//   [ $s -eq 60 ] && m=$((m+1)) && s=00
//   [ $m -eq 60 ] && h=$((h+1)) && m=00
// }

// function _pausar() {
//   while :
//   do
//       _screenTimer
//       sleep 1
//       read tecla
//       [ "$tecla" = "c" ] && clear && break
//   done
// }

func timer(stopClock bool) {
	// Poe o terminal em modo especial de interpretacao de caracteres

	timer1 := time.NewTimer(2 * time.Second)

	<-timer1.C
	fmt.Println("Timer 1 fired")

	timer2 := time.NewTimer(time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 fired")
	}()
	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}

	time.Sleep(2 * time.Second)
}

//   while :
//   do
//       [ "$tecla" = "d" ] && _finish && break && clear
//       [ "$tecla" = "p" ] && _pausar
//       _clock
//       read tecla
//   done

//   # Restaura o modo padrao
//   stty sane

//   exit 0
// }

// # ------------------------------
// #
// # Github
// #
// # ------------------------------

// DoGitFlow fazendo o fluxo do git para alterar a branch baixar as coisas e atualizar
func DoGitFlow() {
	cmd := exec.Command("bash", "bash.sh", branchUpLevel, TaskName)
	defer cmd.Wait()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	// fmt.Println(string(cmd))
	//   $(git checkout $BRANCH_UPLEVEL && git pull origin $BRANCH_UPLEVEL &&  git checkout -b $task_id"_"$TASK_NAME)
}

// # ------------------------------
// #
// # Get data in open project
// #
// # ------------------------------

// GetNameTaskOpenProject conectando a api do open project e pegando o nome da task
func GetNameTaskOpenProject() {
	result := GetJSON("/api/v3/work_packages/" + IDTask)
	TaskName = IDTask + "-" + strings.ReplaceAll(result["subject"].(string), " ", "-")
}

// # ------------------------------
// #
// # Utils
// #
// # ------------------------------

// prepare login auth
func basicAuth() string {
	auth := "apikey" + ":" + token
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

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
		fmt.Println(resp.Status)
		os.Exit(1)
	}

	jsonErr := json.Unmarshal(respBody, &result)

	if jsonErr != nil {
		log.Fatal(jsonErr)
	}

	return
}
