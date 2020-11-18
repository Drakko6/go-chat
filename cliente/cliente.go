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
					

			var nombre string
			fmt.Printf("Introduce el nombre del archivo: ")
			fmt.Scanln(&nombre)
			// nombre, _, _ := input.ReadLine()


			archivo, err := os.Create(string(nombre)+".txt")
			if err != nil {
				fmt.Println(err)
			 	return
			}

			var nombreBytes []byte
			nombreBytes = append(nombreBytes, []byte(nombre)...)

			var n uint
			fmt.Printf("¿Cuántas líneas tendrá el archivo? ")
			fmt.Scanln(&n)

			var sliceBytes []byte
			var archivoMandar []byte
			archivoMandar = append(archivoMandar, []byte("2")...)
			archivoMandar = append(archivoMandar, append(nombreBytes, finLinea)...)


			var i uint =1
			var s []byte
			for i<= n{
				fmt.Printf("Línea %d :", i)
				s, _, _ = input.ReadLine()
				sliceBytes = append(sliceBytes, append(s, finLinea)...)
				archivoMandar = append(archivoMandar, append(s, finLinea)...)
				i ++
			}


			archivo.Write(sliceBytes)

			archivo.Close()
		
		
			handleArchivo(conn, archivoMandar, nombreBytes)
		
		
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



func handleArchivo(conn net.Conn, archivo []byte, nombre []byte){

	//mandar archivo
	conn.Write(append(archivo, finLinea))
	// conn.Write(archivo)

	//mandar nombre
	nombre = (append(nombre, []byte(".txt")...))

	conn.Write(append([]byte(nick + "-> Archivo:"), append([]byte(nombre), finLinea)...))
	

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

		if(string(datos[0])=="2"){
			//es un archivo, no imprimir...Sino crear el archivo

			var nombre []byte

			var primerIndice int 
			for i:=1; i<len(datos); i++{
				if datos[i] ==finLinea{
					primerIndice = i+1
					break
				}
				nombre= append(nombre, datos[i])

			}


			var sliceBytes []byte
			for j:=primerIndice; j<len(datos); j++{
				
				sliceBytes= append(sliceBytes, datos[j])

			}


			archivo, err := os.Create(string(nombre)+".txt")
			if err != nil {
			 	fmt.Println(err)
			 	
			 }

			 archivo.Write(sliceBytes)

			 archivo.Close()


		}else{
			
			fmt.Printf("%s",datos)


		}


		datos = make([]byte, 0)

	}

}