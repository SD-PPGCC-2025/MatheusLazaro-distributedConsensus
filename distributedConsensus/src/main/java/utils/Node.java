package utils;

import java.util.*;

/**
 * Abstração de nó em algoritmo de consenso.
 * Cada nó roda em sua própria Thread.
 */
public abstract class Node implements Runnable {
    protected final int id;               // identificador único do nó
    protected final int N;                // total de nós na rede
    protected final NetworkSimulator net; // simulador de rede compartilhado

    public Node(int id, int N, NetworkSimulator net) {
        this.id  = id;
        this.N   = N;
        this.net = net;
        net.registerNode(id);             // registra a fila deste nó no simulador
    }

    /** broadcast: envia a todos, exceto a si mesmo */
    protected void broadcast(int round, Object payload) {
        Message m = new Message(id, round, payload);
        for (int dst = 1; dst <= N; dst++) {
            if (dst != id) 
                net.send(dst, m);         // envia a mensagem a dst
        }
    }

    /** Unicast simples */
    protected void send(int dst, int round, Object payload) {
        net.send(dst, new Message(id, round, payload));
    }

    /** Bloqueante: espera e retorna próxima mensagem desta fila */
    protected Message receive() throws InterruptedException {
        return net.receive(id);
    }

    @Override
    public abstract void run();          // definida pela implementação concreta
}
