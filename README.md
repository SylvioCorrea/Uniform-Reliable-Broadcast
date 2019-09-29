# Uniform-Reliable-Broadcast
Implementação e teste de um módulo Uniform Reliable Broadcast para sistemas distribuídos.

Para testar:

```>go run App2.go <n>```

onde 3<=n<=10 é o número de aplicações simuladas distribuídas em uma rede. IP e porta destas aplicações devem estar definidas no arquivo network.txt no formato ```número_de_ip:porta```.

O arquivo já contém 3 nodos predefinidos que usam endereço de loopback mais portas próprias. Um teste pode ser imediatamente executado com ```>go run App2.go 3```.

Cada aplicação simulada faz uso de um módulo próprio Uniform Reliable Broadcast que implementa especificação usual com entrega de mensagens quando atingido um quórum entre as aplicações. Cada aplicação manda 150 mensagens via URB e recebe as mensagens das demais aplicações da rede.

A aplicação 0 falha (crash failure) após gerar 100 mensagens. Espera-se que todas as mensagens entregues (delivered) por 0 tenham sido entregues pelas demais aplicações, e que todas as alpicações corretas tenham entregue exatamente as mesmas mensagens, não importando a ordem. Mensagens entregues são escritas por todas as aplicações em arquivos txt. O teste de entregas é feito já na main, ao término das aplicações.
