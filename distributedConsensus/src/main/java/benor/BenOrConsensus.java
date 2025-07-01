package benor;

import utils.*;
import java.util.*;

/**
 * Ben-Or (assíncrono probabilístico):
 * - tolera até f crash failures se N ≥ 3f + 1
 * - cada nó inicia com um estimate ∈ {0,1}
 * - **assíncrono**: não há rounds lock-step; cada nó avança independentemente
 * - em cada rodada:
 * 1) broadcast do estimate atual (envio **assíncrono** via NetworkSimulator)
 * 2) coleta N−f votos de nós (espera bloqueante, mas sem barreira global)
 * 3) se todos os N−f votos forem iguais, decide esse valor e encerra
 * 4) caso contrário, realiza um coin-flip para novo estimate e repete
 */
public class BenOrConsensus {
    public static void main(String[] args) throws Exception {
        final int N = 50; // total de nós (é preciso N ≥ 3f + 1)
        final int f = 15; // tolera até f crash failures
        NetworkSimulator net = new NetworkSimulator();
        List<Thread> threads = new ArrayList<>();

        // cria e inicia uma thread independente para cada nó
        for (int i = 1; i <= N; i++) {
            int init = new Random().nextBoolean() ? 1 : 0;
            BenOrNode node = new BenOrNode(i, N, f, net, init);
            Thread t = new Thread(node, "Node-" + i);
            threads.add(t);
            t.start(); // dispara a execução assíncrona de run()
        }

        // aguarda todas as threads terminarem
        for (Thread t : threads) {
            t.join();
        }
    }
}

/**
 * Cada instância representa um nó no protocolo Ben-Or.
 * - A rede é assíncrona: NetworkSimulator introduz atrasos e não há
 * sincronização de clock.
 * - O nó processa mensagens de decisão via poll() antes de iniciar cada rodada.
 */
class BenOrNode extends Node {
    private int estimate; // valor corrente (0 ou 1)
    private final int N, f; // total de nós e falhas toleradas
    private final Random rand; // gerador de coin-flip independente
    private boolean decided = false;

    public BenOrNode(int id, int N, int f, NetworkSimulator net, int init) {
        super(id, N, net);
        this.estimate = init;
        this.N = N;
        this.f = f;
        this.rand = new Random(id * System.nanoTime()); // Long
        System.out.printf("Node %d -> estimate inicial = %d%n", id, init);
    }

    @Override
    public void run() {
        try {
            int round = 0;
            while (!decided) {
                // (1) Processa **assincronamente** mensagens de decisão que possam já ter
                // chegado
                Message dm;
                while ((dm = net.poll(id)) != null) { // (ASSINCRONICIDADE) Não existe essa barreira global de rounds.
                                                      // Cada nó faz poll() para ver se já chegou uma decisão.
                    if (dm.round == -1) {
                        finish((Integer) dm.payload); // encerra sem aguardar nós
                        return;
                    }
                }

                // (2) Broadcast assíncrono do estimate atual
                System.out.printf("Node %d [r=%d] -> broadcast estimate = %d%n",
                        id, round, estimate);
                broadcast(round, estimate); // net.send agenda entrega com atraso

                // (3) Recebe N−f mensagens de nós de forma bloqueante,
                // mas cada nó o faz no seu tempo, sem sincronizar rounds globalmente
                Map<Integer, Integer> votes = new HashMap<>();
                while (votes.size() < N - f) {
                    Message m = receive(); // bloqueia até chegar mensagem
                    if (m.round == -1) { // decisão antecipada
                        finish((Integer) m.payload);
                        return;
                    }
                    if (m.round != round) {
                        // descarta votes de rounds diferentes, mantendo assíncrono
                        continue;
                    }
                    votes.put(m.senderId, (Integer) m.payload);
                    System.out.printf("Node %d [r=%d] <- voto de Node%d = %d%n",
                            id, round, m.senderId, m.payload);
                }

                // (4) Verifica unanimidade local (sem barreira global)
                // votes.values() retorna uma coleção de Integers com todos os votos recebidos
                // .stream() cria um fluxo sobre essa coleção
                // .filter(v -> v == 0) mantém apenas os votos iguais a 0
                // .count() conta quantos elementos passaram pelo filtro
                long c0 = votes.values().stream().filter(v -> v == 0).count();

                // Mesma lógica para contar votos iguais a 1
                long c1 = votes.values().stream().filter(v -> v == 1).count();
                System.out.printf("Node %d [r=%d] votos: 0 -> %d, 1 -> %d%n",
                        id, round, c0, c1);

                if (c0 == N - f || c1 == N - f) {
                    // unanimidade: anuncia decisão de forma assíncrona
                    int decision = (c1 == N - f) ? 1 : 0;
                    System.out.printf("Node %d [r=%d] -> unanimidade, DECIDE %d%n",
                            id, round, decision);
                    broadcast(-1, decision); // round = -1 sinaliza finish a todos
                    finish(decision);
                    return;
                }

                // (5) Sem unanimidade: coin-flip e próximo round de forma independente
                int previous = estimate;
                estimate = rand.nextBoolean() ? 1 : 0;
                System.out.printf("Node %d [r=%d] -> coinflip %d -> %d%n",
                        id, round, previous, estimate);

                round++;
            }
        } catch (InterruptedException e) {
            Thread.currentThread().interrupt();
        }
    }

    /**
     * Finaliza o nó após decisão:
     * - Espera confirmações de finish (round = -1) de todos os N−1 nós
     * - Cada receive() é bloqueante, mas cada nó faz isso independentemente
     */
    private void finish(int decision) throws InterruptedException {
        if (!decided) {
            System.out.printf("Node %d -> FINALIZA decide %d%n", id, decision);
            decided = true;
        }
        int confirmations = 1; // conta a confirmação local
        while (confirmations < N) {
            Message m = receive(); // bloqueia até próxima mensagem
            if (m.round == -1) {
                confirmations++;
                System.out.printf("Node %d recebeu confirmação de Node%d%n",
                        id, m.senderId);
            }
        }
    }
}