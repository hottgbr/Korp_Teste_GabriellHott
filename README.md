# Desafio Korp - Sistema de Emissão de Notas Fiscais

Projeto desenvolvido como desafio técnico para implementação de um sistema de emissão de notas fiscais, utilizando Angular no frontend e uma arquitetura baseada em microserviços no backend.

A aplicação permite o cadastro de produtos, controle de estoque, criação de notas fiscais com múltiplos itens, fechamento da nota com baixa automática de estoque e impressão após o processamento.

---

## Funcionalidades

### Produtos

- Cadastro de produtos
- Código único por produto
- Descrição do produto
- Saldo inicial em estoque
- Listagem de produtos
- Indicação visual de disponibilidade de estoque

### Notas fiscais

- Criação de notas fiscais com múltiplos produtos
- Definição da quantidade por item
- Numeração sequencial das notas
- Status `OPEN` e `CLOSED`
- Consulta de notas emitidas
- Visualização detalhada da nota
- Exibição do código e descrição dos produtos
- Validação de estoque antes do fechamento
- Baixa automática do estoque
- Impressão da nota fiscal
- Reimpressão de notas fechadas

### Tratamento de falhas

A aplicação também trata cenários em que a comunicação entre os microserviços não pode ser concluída.

Por exemplo, caso o Stock Service esteja indisponível durante o fechamento de uma nota:

- o Billing Service retorna uma resposta apropriada;
- o frontend apresenta feedback ao usuário;
- a nota permanece aberta;
- a impressão não é realizada.

---

## Arquitetura

O backend foi dividido em dois microserviços independentes.

```text
                         ┌──────────────────┐
                         │     Angular      │
                         │     Frontend     │
                         └────────┬─────────┘
                                  │
                  ┌───────────────┴───────────────┐
                  │                               │
                  ▼                               ▼
        ┌─────────────────┐             ┌──────────────────┐
        │  Stock Service  │             │ Billing Service  │
        │     :8081       │◄────────────│      :8082       │
        └────────┬────────┘    HTTP     └────────┬─────────┘
                 │                                │
                 ▼                                ▼
        ┌─────────────────┐             ┌──────────────────┐
        │   stock_db      │             │   billing_db     │
        │ PostgreSQL      │             │ PostgreSQL       │
        └─────────────────┘             └──────────────────┘
