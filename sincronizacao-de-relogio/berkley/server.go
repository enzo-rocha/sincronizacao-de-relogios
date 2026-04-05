package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Cliente struct {
	conn net.Conn
	diff float64
}

func main() {
	ln, _ := net.Listen("tcp", ":5000")
	fmt.Println("Servidor Berkeley iniciado na porta 5000...")

	var clientes []Cliente

	go func() {
		for {
			conn, _ := ln.Accept()
			clientes = append(clientes, Cliente{conn: conn})
			fmt.Printf("Novo cliente conectado: %s\n", conn.RemoteAddr())
		}
	}()

	for {
		time.Sleep(10 * time.Second)
		if len(clientes) == 0 {
			continue
		}

		fmt.Println("\n--- Iniciando Sincronização ---")

		somaDiffs := 0.0
		for i := range clientes {
			fmt.Fprintf(clientes[i].conn, "GET_TIME\n")
			msg, _ := bufio.NewReader(clientes[i].conn).ReadString('\n')
			val, _ := strconv.ParseFloat(strings.TrimSpace(msg), 64)
			clientes[i].diff = val
			somaDiffs += val
		}

		media := somaDiffs / float64(len(clientes)+1)
		fmt.Printf("Média calculada: %.2f\n", media)

		for i := range clientes {
			ajuste := media - clientes[i].diff
			fmt.Fprintf(clientes[i].conn, "%.2f\n", ajuste)
		}
		fmt.Println("Ajustes enviados aos clientes.")
	}
}
