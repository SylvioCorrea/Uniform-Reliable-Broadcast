/*
Módulo Uniform Reliable Broadcast

Módulo prevê a possibilidade de falha a partir de determinado momento (variável NotFail).
Após a falha o módulo ignora pedidos de envio e entregas recebidas.
*/

package Urb

import (
        "./Beb"
        "./Beb/PP2PLink"
        "strings"
        "fmt"
        //"net"
        )

type Urb struct {
    UserID string
    NetworkSize int
    Quorum int
    
    Ind chan string
	Req chan string
	
    Ack map[string][]string
    
    beb *Beb.Beb
    
    NotFail bool
}

func NewUrb(userID string, appPort string, addresses []string) *Urb {
    var module Urb
    module.UserID = string(userID)
    module.NetworkSize = len(addresses)
    module.Quorum = module.NetworkSize/2 + 1
    module.beb = Beb.NewBeb(appPort, addresses)
    module.Ind = make(chan string)
    module.Req = make(chan string)
    
    module.Ack = make(map[string][]string)
    module.NotFail = true
    
    return &module
}

func (module *Urb) Start() {
    fmt.Println("Urb starting")
    module.beb.Start()
    
    var messageSend string
    var messageReceive PP2PLink.PP2PLink_Ind_Message
    
    go func() {
        for {
            select {
            case messageSend = <- module.Req :
                if module.NotFail {
                    module.send(messageSend)
                }
            case messageReceive = <- module.beb.Ind :
                if module.NotFail {
                    module.receive(messageReceive.Message)
                }
            }
        }
    }()
}

func (module *Urb) send(message string){
    module.Ack[message] = make([]string, 0, module.NetworkSize)
    message = module.UserID + ":" + message
    module.beb.Req <- message 
}

func (module *Urb) receive(message string) {
        
    msgSplit := strings.SplitN(message, ":", 2)
    /*msgSplit[0] é o remetente da mensagem
      msgSplit[1] é a mensagem na forma "criador:mensagem"*/
    //fmt.Println("Urb received:", message, "   Split", msgSplit[0], "|", msgSplit[1])
    
    if acks, ok := module.Ack[msgSplit[1]]; ok { //Mensagem já é de conhecimento deste módulo
        acks = append(acks, msgSplit[0]) //Adiciona ack
        //fmt.Println("appended:", acks)
        module.Ack[msgSplit[1]] = acks    
        if len(acks) == module.Quorum {
            //Quorum acaba de ser atingido. Deliver
            //fmt.Println("acks:", acks)
            module.Ind <- msgSplit[1]
        } else {
            /*else, quorum ainda não foi atingido ou ja foi atingido anteriormente*/
            //fmt.Println("receive if: done")
        }
        
    } else { //Mensagem vista pela primeira vez. Necessario broadcast de ack
        //Cria entrada no dicionário e adiciona o ack do remetente
        newAckList := make([]string, 0, module.NetworkSize)
        newAckList = append(newAckList, msgSplit[0])
        //fmt.Println("appended(new msg):", newAckList)
        module.Ack[msgSplit[1]] = newAckList
        //Broadcast o próprio ack para esta mensagem
        //fmt.Println("receive else: sending acks")
        module.beb.Req <- module.UserID + ":" + msgSplit[1]
    }
    //fmt.Println("receive else: done")

}


/*
App1: "1:string"
Urb1: "1:1:string"

Urb2: "1:1:string" --> "2:1:string"
*/
