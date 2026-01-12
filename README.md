# Razpravljalnica

Projekt predstavlja forum, kjer lahko uporabniki:
- Ustvarjajo uporabniške račune
- Ustvarjajo teme za diskusijo
- Objavljajo sporočila v temah
- Všečkajo sporočila drugih uporabnikov
- Naročajo se na teme za prejemanje obvestil v realnem času

### Komponente

| Komponenta | Opis |
|------------|------|
| `server.go` | gRPC strežnik z implementacijo MessageBoard|
| `consumer.go` | odjemalec zgrajen s tview knjižnico |
| `razpravljalnica.proto` | Protocol Buffers definicija API-ja |
| `test_podatki.go` | Skripta za nalaganje testnih podatkov |

## Zagon

### 1. Zagon strežnika 

```bash
cd server
go run server.go 
```

### 2. Zagon odjemalca

```bash
cd consumer
go run consumer.go
```

### 3. Nalaganje testnih podatkov (opcijsko)

```bash
cd testi
go run test_podatki.go
```
