package flooding;

import utils.*;
import java.util.*;

/**
 * FloodSet (síncrono, determinístico):
 *  - tolera até f crash‐failures em f+1 rodadas
 *  - cada nó inicia com knownValue = seu próprio ID
 *  - em cada rodada:
 *      1) Broadcast do knownValue atual
 *      2) recepção de N−1 valores de nós
 *      3) atualização para o mínimo recebido
 *  - ao fim de todas as rodadas, cada nó decide seu knownValue
 */
public class FloodSet {
    public static void main(String[] args) throws Exception {
        final int N      = 3;       // total de nós na rede
        final int f      = 2;       // tolerância a f falhas crash
        final int rounds = f + 1;   // número de rodadas = f+1
        NetworkSimulator net = new NetworkSimulator();
        List<Thread> threads = new ArrayList<>();

        // criação e início de uma thread para cada nó
        for (int i = 1; i <= N; i++) {
            // inicializa cada nó com ID e valor inicial = ID
            FloodNode node = new FloodNode(i, N, rounds, net, i);
            Thread t = new Thread(node, "Node-" + i);
            threads.add(t);
            t.start();  // dispara FloodNode.run()
        }

        // aguarda todas as threads terminarem antes de encerrar
        for (Thread t : threads) {
            t.join();
        }
    }
}

/**
 * Representa um nó do protocolo FloodSet.
 */
class FloodNode extends Node {
    private int knownValue;    // valor corrente a ser propagado
    private final int rounds;  // total de rodadas a executar

    /**
     * @param id        identificador do nó
     * @param N         número total de nós na rede
     * @param rounds    número de rodadas (f+1)
     * @param net       simulador de rede compartilhado
     * @param initial   valor inicial de knownValue (aqui, igual ao id)
     */
    public FloodNode(int id, int N, int rounds,
                     NetworkSimulator net, int initial) {
        super(id, N, net);
        this.knownValue = initial;  // mantém o valor inicial
        this.rounds     = rounds;
        System.out.printf("Node %d -> knownValue inicial = %d%n", id, initial);
    }

    @Override
    public void run() {
        try {
            // executa exatamente 'rounds' rodadas
            for (int r = 1; r <= rounds; r++) {
                // --- 1) Broadcast do knownValue atual ---
                System.out.printf("Node %d [r=%d] -> broadcast(%d)%n",
                                  id, r, knownValue);
                broadcast(r, knownValue);

                // --- 2) Recepção de N-1 valores de outros nós ---
                for (int i = 1; i < N; i++) {
                    Message m = receive();              // bloqueia até chegar mensagem (SINCRONICIDADE)
                    int v = (Integer) m.payload;
                    System.out.printf("Node %d [r=%d] <- de Node%d: %d%n",
                                      id, r, m.senderId, v);

                    // --- 3) Atualização para o mínimo recebido ---
                    if (v < knownValue) {
                        knownValue = v;                // adota o menor valor
                        System.out.printf("Node %d [r=%d]    -> knownValue = %d%n",
                                          id, r, knownValue);
                    }
                }

                // log após término da rodada
                System.out.printf("Node %d [r=%d] fim rodada, knownValue = %d%n",
                                  id, r, knownValue);
            }

            // --- decisão final: knownValue após todas as rodadas ---
            System.out.printf("Node %d -> DECIDE %d%n", id, knownValue);

        } catch (InterruptedException ex) {
            Thread.currentThread().interrupt();
        }
    }
}
