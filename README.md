# Exchange Track

Serviço escrito em **Go** para monitoramento automatizado de cotações de moedas estrangeiras em relação ao Real Brasileiro (BRL), com coleta histórica e atualização periódica em tempo real.

---

## Sobre o projeto

O Exchange Track realiza a coleta e armazenamento de taxas de câmbio utilizando a API da AwesomeAPI.

O serviço opera com dois modos de monitoramento:

- **Histórico** — coleta as cotações diárias dos últimos 90 dias para cada moeda configurada.
- **Tempo real** — coleta uma nova cotação a cada 5 minutos.

Todas as cotações são registradas tendo o BRL como moeda base.

A aplicação é executada via Docker e utiliza banco de dados relacional ([MariaDB](https://mariadb.org/)) para persistência dos dados.

---

## Executando localmente

### Pré-requisitos

- Docker
- Docker Compose
- Token de acesso da https://awesomeapi.com.br/

---

### 1. Clonar o repositório

```bash
git clone https://github.com/cardosodg/exchange-track-go.git
cd exchange-track-go
```
### 2. Configuração das variáveis de ambiente

Crie o arquivo `.env` a partir do modelo:
```bash
cp .env.example .env
```

Edite o arquivo .env com os valores adequados:
```env
# Configurações do banco de dados
DB_USER=seu_usuario
DB_PASS=sua_senha
DB_HOST=mariadb
DB_PORT=3306
DB_NAME=exchange_db

# Moedas monitoradas
EXCHANGE_RT=USD-BRL,EUR-BRL
EXCHANGE_HIST=USD-BRL,EUR-BRL

# Token da API AwesomeAPI
EXCHANGE_KEY=sua_api_key
```

#### Variáveis de Ambiente

| Variável        | Tipo   | Obrigatória | Exemplo        | Descrição |
|---------------|--------|------------|----------------|------------|
| `DB_USER`     | string | Sim        | root           | Usuário utilizado para autenticação no banco de dados |
| `DB_PASS`     | string | Sim        | secret         | Senha do usuário do banco de dados |
| `DB_HOST`     | string | Sim        | mariadb      | Endereço do servidor do banco de dados |
| `DB_PORT`     | int    | Sim        | 3306           | Porta de conexão do banco de dados |
| `DB_NAME`     | string | Sim        | exchange_db    | Nome do banco de dados utilizado pela aplicação |
| `EXCHANGE_RT` | string | Sim        | USD-BRL,EUR-BRL        | Lista de moedas monitoradas em tempo real (separadas por vírgula) |
| `EXCHANGE_HIST` | string | Sim     | USD-BRL,EUR-BRL        | Lista de moedas utilizadas para coleta histórica (separadas por vírgula) |
| `EXCHANGE_KEY` | string | Sim       | abc123token    | Token de autenticação da API AwesomeAPI |


As moedas devem ser informadas conforme o padrão suportado pela AwesomeAPI. Consulte a documentação oficial em https://docs.awesomeapi.com.br/api-de-moedas para detalhes sobre os códigos aceitos.


### 3. Subir os containers

Com as variáveis de ambiente configuradas, inicie os serviços com:

```bash
docker compose up -d
```

Após a inicialização, o serviço começará automaticamente:

- A coleta histórica (últimos 90 dias)
- A coleta em tempo real (a cada 5 minutos)

Para acompanhar os logs do serviço:

```bash
docker compose logs -f exchangetrack
```

Para acessar o banco de dados manualmente:

```bash
docker exec -it mariadb mysql -u <DB_USER> -p
```

## Visualização dos dados

Para criação de dashboards e visualização das métricas coletadas, pode-se utilizar o [Grafana](https://grafana.com/).

O Grafana pode ser conectado diretamente ao MariaDB. basta adicionar uma nova *Data Source* do tipo **MySQL** e configurar a conexão apontando para o container do MariaDB.

### Exemplo de serviço Grafana no Docker Compose

Você pode adicionar o seguinte serviço ao seu `docker-compose.yml`:

```yaml
  grafana:
    image: grafana/grafana:latest
    container_name: grafana
    ports:
      - "3000:3000"
    restart: unless-stopped
    depends_on:
      - mariadb
```

Após subir o container:

```bash
docker compose up -d grafana
```

O Grafana estará disponível em:

```
http://localhost:3000
```

Credenciais padrão:

- **Usuário:** admin  
- **Senha:** admin  

Lembre-se de trocar as credenciais padrão do Grafana.
