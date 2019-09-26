/*
  Construido como parte da disciplina de Sistemas Distribuidos
  Semestre 2018/2  -  PUCRS - Escola Politecnica
  Estudantes:  Andre Antonitsch e Rafael Copstein
  Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)
  Algoritmo baseado no livro:
  Introduction to Reliable and Secure Distributed Programming
  Christian Cachin, Rachid Gerraoui, Luis Rodrigues

  Melhorado como parte da disciplina de Sistemas Distribuídos
  Reaproveita conexões TCP já abertas
  Semestre 2019/1
  Vinicius Sesti e Gabriel Waengertner
  
  
  =============Modificações para entrega 2019/2================
  
  Correções adicionais de checagem de erro de conexão feitas
  durante a aula.
  
  As funções Start e Send foram modificadas para receber um
  PONTEIRO PP2PLink.
*/

package PP2PLink

import (
        "fmt"
        "net"
        )

func errorCheck(e error) {
    if e != nil {
        fmt.Println(e.Error())
    }
}

type PP2PLink_Req_Message struct {
	To      string
	Message string
}

type PP2PLink_Ind_Message struct {
	From    string
	Message string
}

type PP2PLink struct {
	Ind   chan PP2PLink_Ind_Message
	Req   chan PP2PLink_Req_Message
	Run   bool
	Cache map[string]net.Conn
}

func (module PP2PLink) Init(address string) {

	fmt.Println("Init PP2PLink!")
	if module.Run {
		return
	}

	module.Cache = make(map[string]net.Conn)
	module.Run = true
	module.Start(address)
}

func (module *PP2PLink) Start(address string) {
    
    fmt.Println("PP2PLink start")
    
    go func() {

		listen, err2 := net.Listen("tcp4", address)
        errorCheck(err2)
        
        readLimit := 12 //máximo número de bytes a serem recebidos por leitura tcp
        
		for {

			// aceita repetidamente tentativas novas de conexao
			conn, err3 := listen.Accept()
            errorCheck(err3)

			go func() {

				// quando aceita, repetidamente recebe mensagens na conexao TCP (sem fechar)
				// e passa para cima
				for {
					buf := make([]byte, readLimit)

					len, err := conn.Read(buf)
					if err != nil {
						errorCheck(err)
                        continue
					}

					content := make([]byte, len)
					copy(content, buf)

					msg := PP2PLink_Ind_Message{
						From:    conn.RemoteAddr().String(),
						Message: string(content)}

					module.Ind <- msg
				}
			}()
		}
	}()

	go func() {
		for {
			message := <-module.Req
			module.Send(message)
		}
	}()

}

func (module *PP2PLink) Send(message PP2PLink_Req_Message) {
    
	var conn net.Conn
	var ok bool
	var err error

	// ja existe uma conexao aberta para aquele destinatario?
	if conn, ok = module.Cache[message.To]; ok {
		//fmt.Printf("Reusing connection %v to %v\n", conn.LocalAddr(), message.To)
	} else { // se nao tiver, abre e guarda na cache
		conn, err = net.Dial("tcp", message.To)
		if err != nil {
			fmt.Println(err)
			return
		}
		module.Cache[message.To] = conn
	}
    
	fmt.Fprintf(conn, message.Message)
}
