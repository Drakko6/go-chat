package main

import(

	"net"
	"fmt"
	"bytes"
	"io"
	"strconv"
	"bufio"
	"os"


	
)




func buscarConexion(conn net.Conn)int{
	var indice int
	for i:=0; i<len(clientes); i++{
		if clientes[i] == conn {
		indice = i
		return indice
		}
	}
	return -1
}


const bufferTam = 256
const finLinea = 10

var clientes[] net.Conn

var chatCompleto []byte = []byte("CHAT COMPLETO\n")


func handleClient (c net.Conn){

	defer c.Close()

	

	var datos[]byte
	buffer := make([]byte, bufferTam)
	cont := 0

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
				cont++
				break
			}

		}

		if cont == 1{
			fmt.Printf("Se conectó %s", datos)
	
			datos = make([]byte, 0)


		}else{

			if string(datos[0]) == "1"{


				//agregar nick para decir quien se desconectó
				var nick []byte
				for j:= 1; j<len(datos)-1; j++{

					nick =append( nick, datos[j])

				}
				fmt.Println("Se desconectó", string(nick) )


				nick = append(nick, []byte(" se desconectó\n")...)
				
				
				//buscar conexion y borrar de clientes
				i:= buscarConexion(c)
				clientes = append(clientes[:i], clientes[i+1:]...)

				//enviar mensaje a los demás
				for h:=0; h<len(clientes); h++{
		
					if clientes[h] != c {
						clientes[h].Write(nick)

					}

			
				}

				chatCompleto = append(chatCompleto, nick...)

				break


			}else if string(datos[0]) =="3"{

			 enviarChatCompleto(c)
			 datos = make([]byte, 0)



			}else{
				enviarOtrosClientes(c, datos)
				chatCompleto = append(chatCompleto, datos...)
				datos = make([]byte, 0)

			}

		}



		

	}
	
}





func enviarChatCompleto(c net.Conn){

	c.Write(chatCompleto)

}

func enviarOtrosClientes (emisor net.Conn, datos []byte) {

	for i:=0; i<len(clientes); i++{
		
			if clientes[i] != emisor{
				clientes[i].Write(datos)

			}

		
	}

}

func menu(s net.Listener) {

	input := bufio.NewReader(os.Stdin)

	op := 5


	for op != 0 {

		
		fmt.Println("")
		fmt.Println("	CHAT EN GO -- SERVIDOR")
		fmt.Println("1.- Mostrar los mensajes/nombres de archivos")
		fmt.Println("2.- Respaldar mensajes en archivo")
		fmt.Println("0.- Terminar servidor")

		//fmt.Scan(&op)
		fmt.Println("Opcion:")
		buf, _, _ := input.ReadLine()
		op, _ = strconv.Atoi(string(buf))
		

		switch {
		case op == 1:
			fmt.Println("Opción 1")		
			fmt.Printf("%s", chatCompleto)	

		case op == 2:
			fmt.Println("Opción 2")
			// input := bufio.NewReader(os.Stdin)
			// fmt.Printf("Introduce el nombre del archivo: ")
			// nombre, _, _ := input.ReadLine()


			archivo, err := os.Create("respaldo.txt")
			if err != nil {
			 	fmt.Println(err)
			 	
			}

			archivo.WriteString(string(chatCompleto))
			archivo.Close()
			fmt.Println("Archivo guardado como respaldo.txt")

		case op == 0:
			fmt.Println("Saliendo...")
			s.Close()
			
		default:
			fmt.Println("Opción inválida")

		}


	}

}


func main () {


	
	clientes = make([]net.Conn, 0)
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}


	for {
		c, _:= s.Accept()
		 if err != nil{
			fmt.Println(err)
			continue
			
		 }

		
		clientes = append(clientes, c)


		go handleClient(c)	

		go menu(s)


	}


}