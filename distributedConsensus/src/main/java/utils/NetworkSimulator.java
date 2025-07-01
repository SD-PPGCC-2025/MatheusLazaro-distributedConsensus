package utils;

import java.util.*;
import java.util.concurrent.*;

/**
 * Simula uma rede ponto-a-ponto entre nós:
 *  - cada nó registra sua fila de entrada
 *  - mensagens são entregues com atraso aleatório (até MAX_DELAY ms)
 *  - possibilidade de perda controlada por LOSS_PROB (hoje zero)
 *  - fornece receive() bloqueante e poll() não-bloqueante
 */
public class NetworkSimulator {
    private static final Random RAND = new Random();
    private static final double LOSS_PROB = 0.0;    // prob. de perda (0 = sem perda)
    private static final int MAX_DELAY = 50;        // atraso máximo em ms

    // Mapeia nodeId PARA fila de mensagens desse nó
    private final Map<Integer, BlockingQueue<Message>> inboxes =
        new ConcurrentHashMap<>();

    /**
     * Deve ser chamado por cada nó no início para criar sua fila de entrada. 
     */
    public void registerNode(int nodeId) {
        inboxes.put(nodeId, new LinkedBlockingQueue<>());
    }

    /**
     * Envia uma mensagem a destId:
     *  - descarta com probabilidade LOSS_PROB
     *  - agenda entrega após atraso aleatório [0, MAX_DELAY)
     */
    public void send(int destId, Message msg) {
        if (RAND.nextDouble() < LOSS_PROB) return;  // simula perda
        int delay = RAND.nextInt(MAX_DELAY);
        Executors.newSingleThreadScheduledExecutor()
                 .schedule(() -> {
                     BlockingQueue<Message> q = inboxes.get(destId);
                     if (q != null) q.offer(msg);
                 }, delay, TimeUnit.MILLISECONDS);
    }

    /**
     * Bloqueante: aguarda até que haja ao menos uma mensagem na fila.
     */
    public Message receive(int nodeId) throws InterruptedException {
        return inboxes.get(nodeId).take();
    }

    /**
     * Não-bloqueante: retorna a próxima mensagem ou null se não houver.
     */
    public Message poll(int nodeId) {
        return inboxes.get(nodeId).poll();
    }
}
