package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	desvio, _ := strconv.ParseFloat(os.Args[1], 64)

	relogioLocal := time.Now().Add(time.Duration(desvio) * time.Second)

	conn, _ := net.Dial("tcp", "localhost:5000")
	fmt.Printf("Conectado. Hora Inicial do Cliente: %s\n", relogioLocal.Format("15:04:05"))

	go func() {
		for {
			time.Sleep(1 * time.Second)
			relogioLocal = relogioLocal.Add(1 * time.Second)
		}
	}()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		if strings.TrimSpace(msg) == "GET_TIME" {
			diff := relogioLocal.Sub(time.Now()).Seconds()
			fmt.Fprintf(conn, "%.2f\n", diff)

			ajusteStr, _ := reader.ReadString('\n')
			ajuste, _ := strconv.ParseFloat(strings.TrimSpace(ajusteStr), 64)

			antes := relogioLocal.Format("15:04:05")
			relogioLocal = relogioLocal.Add(time.Duration(ajuste * float64(time.Second)))

			fmt.Printf("[%s] -> Ajuste de %.2fs recebido. Nova hora: %s\n",
				antes, ajuste, relogioLocal.Format("15:04:05"))
		}
	}
}
