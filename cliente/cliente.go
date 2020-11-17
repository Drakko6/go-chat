package main

import (
	"bufio"
	"os"	
	"fmt"
	"net"
	"bytes"
	"io"
	"strconv"
	

)


const bufferTam = 256
const finLinea = 10

var nick string
var input *bufio.Reader




func cerrarCliente(conn net.Conn, ){

	terminar := (append([]byte("1"), []byte(nick)...))

	// terminar := "1"
	conn.Write(append([]byte(terminar), finLinea) )
	conn.Close()

	
}


func pedirChat(conn net.Conn ){

	mostrar:= "3"

	conn.Write(append([]byte(mostrar), finLinea) )

	
}

func main(){

	input = bufio.NewReader(os.Stdin)

	for nick ==""{
		fmt.Println("Escribe tu nick: ")
		buf, _, _ := input.ReadLine()
		nick = string(buf)
	}



	var conn net.Conn
	var err error

	for {

		conn, err = net.Dial("tcp", ":9999")
		if err == nil {
			break
		}

	}

    defer conn.Close()
	

	enviarNick(conn)

	go recibirMensajes(conn)

	op := 5

	for op != 0 {
		
		fmt.Println("")
		fmt.Println("	CHAT EN GO")
		fmt.Println("1.- Enviar mensaje")
		fmt.Println("2.- Enviar archivo")
		fmt.Println("3.- Mostrar chat")
		fmt.Println("0.- Salir")

		//fmt.Scan(&op)
		fmt.Println("Opcion:")
		buf, _, _ := input.ReadLine()
		op, _ = strconv.Atoi(string(buf))
		

		switch {
		case op == 1:
			fmt.Println("Opción 1")
			fmt.Printf("Mensaje: ")
			handleConnection(conn)
			

		case op == 2:
			fmt.Println("Opción 2")
			// fmt.Printf("Nombre del archivo: ")
			// nombre, _, _ := input.ReadLine()


			// archivo, err := os.Create(string(nombre)+".txt")
			// if err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			fmt.Printf("Introduce la ubicación del archivo: ")
			
			handleArchivo(conn)




			// for _, c:= range ascendente {
			// 	archivo.WriteString(c + "\n")
			// }

			// defer archivo.Close()
					
		
		
		case op == 3:
			fmt.Println("Opción 3")
			pedirChat(conn)
 

		case op == 0:
			fmt.Println("Saliendo...")
			cerrarCliente(conn)			
			

		default:
			fmt.Println("Opción inválida")

		}


	}






}

func enviarNick(conn net.Conn){

	conn.Write(append([]byte(nick), finLinea))

}



func handleArchivo(conn net.Conn){

	buf, _, _ := input.ReadLine()
	if len(buf) > 0 {
			conn.Write(append([]byte(nick + "-> Archivo:"), append(buf, finLinea)...))
	}

}

func handleConnection(conn net.Conn){

		buf, _, _ := input.ReadLine()
		if len(buf) > 0 {
				conn.Write(append([]byte(nick + "->"), append(buf, finLinea)...))
		}

}

func recibirMensajes(c net.Conn) {

	var datos[]byte
	buffer := make([]byte, bufferTam)

	for {

		for {
			n, err := c.Read(buffer)
			if err != nil {
				if err == io.EOF{
					fmt.Println(err)
				    break

				}
				
			}

			buffer = bytes.Trim(buffer[:n], "\x00")

			datos = append(datos, buffer...)
			if datos[len(datos)-1] == finLinea {
				break
			}

		}

		fmt.Printf("%s",datos)

		datos = make([]byte, 0)

	}

}