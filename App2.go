/*
Programa de teste de uma implementação do módulo Uniform Reliable Broadcast.
O programa roda uma quantidade de aplicações passada por parâmetro no terminal.
Cada aplicação é representada por uma go routine Application definida abaixo.
Cada aplicação recebe um número de identidade 0 a n. A aplicação 0 falha depois de enviar 100 mensagens.
*/

package main

import (
        "fmt"
        "io/ioutil"
        "os"
        "time"
        "strings"
        "bytes"
        "strconv"
        "./Urb"
        )

func cancer() {
    os.Exit(1)
}

func errorCheck(e error) {
    if e != nil {
        panic(e)
    }
}


func main() {
    if len(os.Args)<2 {
        fmt.Println("Usage:\n>go run App2.go n")
        fmt.Println("where n is the number of nodes in the network.")
        fmt.Println("Do not exceed the number of nodes defined in file network.txt.")
        os.Exit(1)
    }
    
    /*Executa uma quantidade de funções Application igual ao parâmetro
      passado pela linha de comando*/
    nodes, err := strconv.Atoi(os.Args[1])
    errorCheck(err)
    outputFiles := make([]string, 0, nodes)
    signal := make(chan int, nodes)
    for i:=0; i<nodes; i++ {
        idString:= strconv.Itoa(i)
        outputFiles = append(outputFiles, "App"+idString+"delivers.txt")
        go Application(i, outputFiles[i],  signal)
    }
    
    /*Aguarda a finalização de todas as funções*/
    for i:=0; i<nodes; i++ {
        <-signal
    }
    fmt.Println("All nodes done.")
    CheckResults(nodes, outputFiles)
}


//Função que representa o funcionamento de uma aplicação
func Application(nodeId int, outputFile string, signal chan int) {
    
    //Le a topologia de rede em um arquivo de texto
    fileData, err := ioutil.ReadFile("network.txt")
    errorCheck(err)
    fileString := string(fileData)
    
    nodeIdString := strconv.Itoa(nodeId)
    fmt.Println("Node id:", nodeIdString)
    
    /*Substitui caracteres de nova linha do windows se necessario*/
    fileString = strings.Replace(fileString, "\r\n", "\n", -1)
    /*Cada linha contem um endereço de ip seguido de uma porta*/
    network := strings.Split(fileString, "\n")
    /*Obtem a porta a ser escutada nesta execucao*/
    listenPort := ":"+strings.Split(network[nodeId], ":")[1]
    
    urb := Urb.NewUrb(nodeIdString, listenPort, network)
    urb.Start()
    messageN := 0
    quit := make(chan int)
    quit2 := make(chan int)
    var delivered bytes.Buffer
    
    //Sending thread
    go func() {
        header := nodeIdString + ":"
        time.Sleep(3 * time.Second)
        fmt.Println("Start sending")
        for messageN < 150 {
            if messageN == 100 && nodeId == 0 {
                urb.NotFail = false
            }
            urb.Req <- fmt.Sprintf("%s%8d", header, messageN)
            messageN++
            //time.Sleep(200 * time.Millisecond)
        }
        fmt.Println("Stop sending")
        time.Sleep(5 * time.Second)
        fmt.Println("signalling")
        quit <- 0
    }()
    
    //Receiving thread
    go func() {
        forloop: //label do for
        for {
            select {
            case <- quit :
                fmt.Println("Stop receiving")
                break forloop //se não usa label interrompe apenas o select
            case  message := <- urb.Ind:
                //fmt.Println("\n", message, "\n")
                delivered.WriteString(message+"\n")
            }
        }
        err := ioutil.WriteFile(outputFile, delivered.Bytes(), 0644)
        errorCheck(err)
        fmt.Println("Closing")
        quit2 <- 0
    }()
    
    <- quit2
    
    signal <- 0 //Signal main the function is done
}

func CheckResults(nodes int, files []string) {
    delivers := make([][]string, 0, nodes)
    for i:=0; i<nodes; i++ {
        fileData, err := ioutil.ReadFile(files[i])
        errorCheck(err)
        delivers = append(delivers, strings.Split(string(fileData), "\n"))
    }
    
    /*Laço de teste geral: verifica as entregas de todos os processos corretos.
      Processos corretos devem ter entregado exatamente as mesmas mensagens.
      O único processo incorreto é o 0, logo, começa o teste na lista de entregas
      do processo 1.*/
    for i:=1; i<nodes; i++ {
        
        for _, msg := range(delivers[i]){
            
            for j:=1; j<nodes; j++ {
                
                if j!=i {
                    found := Contains(delivers[j], msg)
                    if found==false {
                        fmt.Printf("Message \"%s\" from application %d was not delivered by application %d.\n", msg, i, j)
                    }
                }
            }
        }
    }
    
    /*Laço de teste do processo falho 0. Todas as mensagens presentes em 0 devem
      estar presentes nos demais processos.*/
    for _, msg := range(delivers[0]){
            
        for j:=1; j<nodes; j++ {
            
            found := Contains(delivers[j], msg)
            if found==false {
                fmt.Printf("Message \"%s\" from application %d was not delivered by application %d.\n", msg, 0, j)
            }
        }
            
    }
    
    fmt.Println("Results checked")
}

func Contains(strSlice []string, str1 string) bool {
    for _, str2 := range(strSlice){
        if str1 == str2 {
            return true
        }
    }
    return false
}













