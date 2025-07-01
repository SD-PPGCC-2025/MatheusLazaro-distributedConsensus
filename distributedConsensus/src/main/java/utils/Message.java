package utils;

// Mensagem simples trocada entre nós participantes de um algoritmo de consenso:
//   senderId: ID do nó que enviou
//   round:    rodada (ou fase) do protocolo
//   payload:  conteúdo genérico (por exemplo, voto ou valor conhecido)
public class Message {
    // identificador do nó que enviou esta mensagem
    public final int senderId;

    // número da rodada (ou fase) em que a mensagem foi enviada
    public final int round;

    // conteúdo transportado, que pode ser qualquer objeto relevante ao protocolo
    public final Object payload;

    // constrói uma nova mensagem imutável com os campos especificados
    public Message(int senderId, int round, Object payload) {
        this.senderId = senderId;  // armazena o ID do remetente
        this.round    = round;     // armazena o número da rodada
        this.payload  = payload;   // armazena o conteúdo da mensagem
    }
}
