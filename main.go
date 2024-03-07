package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

// outside perspective
func main() {
	reader := bufio.NewReader(os.Stdin)
	clientset := getClientset()

	for {
		fmt.Println("Select an option (1-5) or 0 to exit:")
		fmt.Println("1. Pod")
		fmt.Println("2. Deployment")
		fmt.Println("3. Service")
		fmt.Println("4. Job")
		fmt.Println("5. Cronjob")
		fmt.Println("0. Exit")

		input, _ := reader.ReadString('\n')
		input = input[:len(input)-1]

		option, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		switch option {
		case 0:
			fmt.Println("Exiting...")
			return
		case 1:
			createPod(clientset)
		case 2:
			createDeployment(clientset)
		case 3:
			createService(clientset)
		case 4:
			createJob(clientset)
		case 5:
			createCronJob(clientset)
		default:
			fmt.Println("Invalid option. Please select a number between 0 and 5.")
		}
	}
}
