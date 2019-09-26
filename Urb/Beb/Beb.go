/*
Módulo Best Effort Broadcast
*/

package Beb

import (
        "fmt"
        "net"
        "./PP2PLink"
        )

/*Para utilização de um objeto do módulo, espera-se que o usuário chame o
  contrutor abaixo, pois o struct não é exportado.*/
type Beb struct {
	Ind   chan PP2PLink.PP2PLink_Ind_Message
	Req   chan string
    ListenPort string
	P2P   PP2PLink.PP2PLink
	Nodes []string
}

/*Construtor. Retorna um ponteiro para um objeto Beb. É a única maneira possível
  do usuário obter um objeto Beb.*/
func NewBeb(appPort string, addresses []string) *Beb{
    var module Beb
    module.Ind = make(chan PP2PLink.PP2PLink_Ind_Message)
    module.Req = make(chan string)
    module.ListenPort = appPort
    module.Nodes = make([]string, len(addresses))
    copy(module.Nodes, addresses)
    
    /*Constrói um objeto PP2PLink que será usado para comunicação
      com os demais processos da rede. Nota: este programa não
      faz uso da função Init fornecida pelo módulo PP2PLink.
      Perceba-se que o canal de indicação deste objeto PP2PLink
      é o mesmo do objeto Beb que está sendo criado por esta função.*/
	module.P2P = PP2PLink.PP2PLink{
                 Ind: module.Ind,
                 Req: make(chan PP2PLink.PP2PLink_Req_Message),
                 Run: true,
                 Cache: make(map[string]net.Conn)}
	
    return &module
}

//Inicia o funcionamento de um módulo Beb
func (module *Beb) Start() {
    fmt.Println("Beb start")
    /*Inicia o funcionamento do Perfect P2P Link da camada inferior
      que efetivamente fará as comunicações com os demais processos
      via tcp*/
    module.P2P.Start(module.ListenPort)
	
    /*Cria um processo que aguarda requisições de broadcast da aplicação.
      Quando uma requisição é recebida, envia a mensagem via P2P para
      TODOS os processos que fazem parte da rede, incluindo o processo
      que fez a requisição*/
    go func() {
		for {
			message := <-module.Req
			for i:=0; i<len(module.Nodes); i++ {
			    msg := PP2PLink.PP2PLink_Req_Message {
			           To: module.Nodes[i],
			           Message: message}
			    module.P2P.Req <- msg
			}
		}
	}()
    
    /*Perceba-se que o módulo não inicia um processo em particular para
      indicações (delivery). Uma indicação ocorre pelo canal Ind. Este canal
      é iniciado acima no construtor NewBeb e é exatamente o mesmo canal usado
      para criar o Perfect P2P Link da camada inferior. Sendo assim, quando o
      P2P Link faz a indicação de uma mensagem recebida por tcp, o canal usado
      faz uma ligação direta com a processo que está utilizando o broadcast,
      ainda que este processo não tome conhecimento da camada inferior à
      camada Beb.*/
}
