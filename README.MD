```txt
PPGCC
Disciplina: Sistemas Distribuídos
Aluno: Matheus Lázaro Honório da Silva
```

# Distributed Consensus

Este projeto implementa **dois** algoritmos clássicos de consenso distribuído, pedagógica e didaticamente, em Java:

1. **Flooding (FloodSet) (Consenso baseado em inundações)**  
   Algoritmo síncrono determinístico que tolera até *f* falhas por meio de *f + 1* rodadas de flooding (algoritmo de inundação - distribuir informação para todos nós). Baseado em Nancy A. Lynch, _Distributed Algorithms_, Cap. 6.

2. **Ben-Or (Consenso Assíncrono Probabilístico)**  
   Algoritmo assíncrono, randomizado, que garante terminação quase certa (probabilidade 1) tolerando falhas de processo. Baseado em Lynch, Cap. 21.

---

## Organização

- **`src/main/java/utils`**  
  Classes de suporte à *simulação de rede* e modelo de nó/mensagem:
  - `NetworkSimulator`: simula entrega de mensagens (com opção de perda/atraso).
  - `Node`: abstração de processo participante.
  - `Message`: estrutura básica de mensagem.

- **`src/main/java/flooding/FloodSet.java`**  
  Implementa FloodSet: em cada rodada, cada nó “floods” seu valor conhecido para todos os vizinhos; ao final de *f + 1* rodadas, todos chegam ao mesmo valor.

- **`src/main/java/benor/BenOrConsensus.java`**  
  Implementa Ben-Or: cada nó vota em seu valor, troca votos, e usa coin-flip (aleatório) em caso de empates; repete até consenso.

---

## Como executar

1. Compilar:
```bash
cd distributedConsensus
javac -d out \
  src/main/java/utils/*.java \
  src/main/java/flooding/*.java \
  src/main/java/benor/*.java



```

2. Executar FloodSet:
```bash
java -cp out flooding.FloodSet
```

3. Executar Ben-Or:
```bash
java -cp out benor.BenOrConsensus
```


Ambos simulam um pequeno grupo de N nós, com até f falhas.