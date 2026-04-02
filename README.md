# Dublago - Sistema de Dublagem Automática de Vídeos

Sistema que automatiza o processo de dublagem de vídeos do YouTube usando transcrição (Whisper), tradução (MarianMT/Google) e TTS (Coqui TTS).

---

## 📋 Pré-requisitos

Antes de começar, certifique-se de ter instalado:

- **Go** (versão 1.23.2 ou superior) - [Download](https://go.dev/dl/)
- **Python** (versão 3.10) com pip
- **Docker** e **Docker Compose** - [Instruções](https://docs.docker.com/get-docker/)
- **FFmpeg** (será instalado via dependências Python)
- **yt-dlp** (necessário para download de vídeos do YouTube) - [Download](https://github.com/yt-dlp/yt-dlp)

---

## 🚀 Instalação

### 1. Clonar o Repositório

```bash
git clone <url-do-repositorio>
cd dublago
```

### 2. Configurar o Ambiente Python

Criar e ativar o virtual environment:

```bash
# Criar virtualenv (certifique-se de estar usando Python 3.10)
python3.10 -m venv .venv

# Ativar o ambiente
source .venv/bin/activate  # Linux/Mac
# OU
.venv\Scripts\activate    # Windows
```

### 3. Instalar Dependências Python

**⚠️ IMPORTANTE:** Antes de instalar os pacotes do `requirements.txt`, atualize o `setuptools` e `wheel`:

```bash
# ETAPA CRÍTICA - Atualizar setuptools primeiro
pip install --upgrade setuptools wheel

# Agora instalar as dependências
pip install -r requirements.txt
```

> **Nota sobre o erro `No module named 'pkg_resources'`:**
> Se encontrar esse erro ao instalar `openai-whisper`, significa que o `setuptools` está desatualizado. Execute o comando acima para atualizar antes de tentar novamente.

> **⚠️ Verificação importante:**
> Certifique-se de que TODOS os pacotes foram instalados corretamente. Se encontrar erros de `ModuleNotFoundError` ou `exec: "tts": executable file not found in $PATH` mais tarde, execute:
> ```bash
> # Reinstalar todas as dependências
> pip install -r requirements.txt --force-reinstall
> 
> # Ou instalar pacotes específicos que faltam
> pip install coqui-tts==0.26.0 transformers==4.39.3 sentencepiece==0.2.0
> ```

### 4. Instalar yt-dlp

O sistema usa `yt-dlp` para baixar vídeos do YouTube. Baixe o binário diretamente:

**Linux:**
```bash
# Baixar o binário
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp

# Tornar executável
chmod a+rx ~/.local/bin/yt-dlp

# Verificar instalação
yt-dlp --version
```

**macOS:**
```bash
# Baixar o binário
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp

# Tornar executável
chmod a+rx /usr/local/bin/yt-dlp

# Verificar instalação
yt-dlp --version
```

> **Nota:** Certifique-se de que o diretório `~/.local/bin` (Linux) ou `/usr/local/bin` (macOS) esteja no seu `$PATH`.

### 5. Configurar Variáveis de Ambiente

Copiar o template e configurar suas credenciais:

```bash
cp .env.template .env
```

Editar o arquivo `.env` com suas credenciais:

```env
CLIENT_ID=seu_client_id_google
CLIENT_SECRET=seu_client_secret_google
REDIRECT_URI=http://localhost:4000/api/auth/callback
TOKEN_URL=https://oauth2.googleapis.com/token
USER_INFO_URL=https://www.googleapis.com/oauth2/v2/userinfo
FRONTEND_REDIRECT_URL=http://localhost:3000
DATABASE_URL=postgresql://postgres:postgres@localhost:5000/dubmasterdb
```

---

## 🐳 Configuração do Banco de Dados (PostgreSQL)

### Opção A: Usar Docker Compose (Recomendado)

```bash
# Iniciar o PostgreSQL
docker compose up -d postgres

# Verificar se está rodando
docker compose ps

# Ver logs
docker compose logs -f postgres
```

O banco ficará disponível em:
- **Host:** localhost
- **Porta:** 5000
- **Database:** dubmasterdb
- **Usuário:** postgres
- **Senha:** postgres

### Opção B: PostgreSQL Local

Se preferir usar um PostgreSQL local, ajuste a `DATABASE_URL` no `.env`:

```env
DATABASE_URL=postgresql://seu_user:sua_senha@localhost:5432/dubmasterdb
```

---

## ▶️ Executando o Projeto

### Modo Desenvolvimento (Go)

```bash
# Certifique-se que o virtualenv Python está ativado
source .venv/bin/activate

# Executar a aplicação
go run ./cmd/main.go
```

A API estará disponível em: `http://localhost:4000`

### Usando Docker Compose (Completo)

```bash
# Iniciar todos os serviços (PostgreSQL + Aplicação)
docker compose up -d

# Ver logs
docker compose logs -f

# Parar
docker compose down
```

---

## 📚 API Endpoints

| Método | Endpoint | Descrição | Autenticação |
|--------|----------|-----------|--------------|
| GET | `/api/videos/:filename` | Baixar vídeo processado | ❌ Público |
| POST | `/api/videos` | Enviar vídeo para processamento | ✅ Requer login |
| GET | `/api/me` | Informações do usuário logado | ✅ Requer login |
| GET | `/api/auth/callback` | Callback OAuth Google | ❌ Público |

---

## 🔧 Pipeline de Processamento

O sistema processa vídeos em 6 etapas:

1. **Download** (`input/`) - Baixa vídeo do YouTube
2. **Transcrição** (`transcript/`) - Extrai áudio e transcreve com Whisper
3. **Tradução** (`translator/`) - Traduz texto com MarianMT ou Google Translate
4. **Síntese de Voz** (`tts/`) - Gera áudio em português com Coqui TTS
5. **Legendas** (`caption/`) - Gera e sincroniza legendas
6. **Finalização** (`media/`) - Combina vídeo, áudio e legendas com FFmpeg

---

## 🧪 Scripts Python Standalone

### Tradutor MarianMT

```bash
source .venv/bin/activate

# Traduzir texto simples
python scripts/translator_marianmt.py --text "Hello world" --target-lang pt_BR

# Salvar em arquivo
python scripts/translator_marianmt.py --text "Text to translate" --output resultado.txt

# Usar modelo específico
python scripts/translator_marianmt.py \
  --text "Hello" \
  --target-lang pt_BR \
  --model "Helsinki-NLP/opus-mt-en-ROMANCE"
```

---

## 🐛 Troubleshooting

### Erro: `No module named 'pkg_resources'`

**Causa:** `setuptools` desatualizado

**Solução:**
```bash
pip install --upgrade setuptools wheel
pip install openai-whisper
```

### Erro: `ModuleNotFoundError: No module named 'transformers'`

**Causa:** O pacote `transformers` não foi instalado ou está instalado no Python do sistema ao invés do virtualenv

**Solução:**
```bash
# Certifique-se de que o virtualenv está ativado
source .venv/bin/activate

# Instalar as dependências do tradutor
pip install transformers==4.47.0 sentencepiece==0.2.0 sacremoses==0.1.1

# Verificar instalação
python -c "from transformers import MarianMTModel; print('OK')"
```

### Erro: `AttributeError: module 'transformers.pytorch_utils' has no attribute 'isin_mps_friendly'`

**Causa:** Versão do `transformers` incompatível com o `coqui-tts`. O Coqui TTS 0.26.0 requer `transformers>=4.47.0`.

**Solução:**
```bash
# Atualizar transformers para versão compatível
pip install transformers==4.47.0
```

### Erro: `exec: "tts": executable file not found in $PATH`

**Causa:** O Coqui TTS não foi instalado ou está instalado no Python do sistema ao invés do virtualenv

**Solução:**
```bash
# Certifique-se de que o virtualenv está ativado
source .venv/bin/activate

# Instalar o Coqui TTS
pip install coqui-tts==0.26.0

# Verificar instalação
which tts
# Deve retornar: /home/seu-user/dublago/.venv/bin/tts
```

Se o erro persistir, o código Go já foi atualizado para buscar o `tts` no virtualenv automaticamente. Certifique-se apenas de que:
1. O virtualenv está criado em `./.venv`
2. O Coqui TTS está instalado nele

### Erro: `EOFError: EOF when reading a line` (TTS)

**Causa:** O Coqui TTS está pedindo para aceitar os Termos de Serviço (ToS) interativamente, mas não há terminal disponível quando rodando via código Go.

**Solução:** O código Go já foi atualizado para aceitar automaticamente os ToS (`cmd.Stdin = strings.NewReader("y\n")`). Se ainda encontrar este erro, certifique-se de estar usando a versão mais recente do código.

### Erro: `ImportError: TorchCodec is required for load_with_torchcodec`

**Causa:** O pacote `torchcodec` é necessário para carregar arquivos de áudio, mas não foi instalado automaticamente com o Coqui TTS.

**Solução:** NÃO instale o torchcodec! Em vez disso, use versões compatíveis do PyTorch e TorchAudio sem o backend torchcodec:

```bash
# Certifique-se de que o virtualenv está ativado
source .venv/bin/activate

# Remover torchcodec se instalado
pip uninstall torchcodec -y

# Instalar versões compatíveis do PyTorch e TorchAudio
pip install torch==2.1.0 torchaudio==2.1.0 --force-reinstall
```

### Erro: `RuntimeError: Could not load libtorchcodec` / `OSError: libnppicc.so.13`

**Causa:** O torchcodec foi instalado mas tem incompatibilidade com as bibliotecas CUDA/FFmpeg do sistema.

**Solução:** Remova o torchcodec e use versões mais antigas do PyTorch/TorchAudio:

```bash
# Remover torchcodec
pip uninstall torchcodec -y

# Reinstalar torch e torchaudio com versões compatíveis
pip install torch==2.1.0 torchaudio==2.1.0 --force-reinstall
```

### Erro: `panic: exec: "stable-ts": executable file not found in $PATH`

**Causa:** O stable-ts não foi instalado ou o código Go não está encontrando o binário no virtualenv.

**Solução:**
```bash
# Certifique-se de que o virtualenv está ativado
source .venv/bin/activate

# Instalar stable-ts sem dependências (para evitar conflitos)
pip install stable-ts==2.17.5 --no-deps

# Verificar instalação
which stable-ts
# Deve retornar: /home/seu-user/dublago/.venv/bin/stable-ts
```

### Erro: `panic: exec: "yt-dlp": executable file not found in $PATH`

**Causa:** O `yt-dlp` não está instalado ou não está no PATH

**Solução:**
```bash
# Baixar o binário
curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp

# Tornar executável
chmod a+rx ~/.local/bin/yt-dlp

# Verificar se está no PATH
which yt-dlp
yt-dlp --version
```

Se o diretório `~/.local/bin` não estiver no PATH:
```bash
# Adicione ao seu ~/.bashrc ou ~/.zshrc
export PATH="$HOME/.local/bin:$PATH"
```

### Erro: `connection refused` no banco de dados

**Causa:** PostgreSQL não está rodando

**Solução:**
```bash
docker compose up -d postgres
```

### Erro: `permission denied` no FFmpeg

**Causa:** FFmpeg não instalado ou sem permissão

**Solução:**
```bash
# Ubuntu/Debian
sudo apt-get install ffmpeg

# macOS
brew install ffmpeg
```

### Verificar instalação do Whisper

```bash
source .venv/bin/activate
python -c "import whisper; print('Whisper OK:', whisper.__version__)"
```

---

## 📁 Estrutura do Projeto

```
dublago/
├── cmd/
│   └── main.go              # Entry point da aplicação Go
├── config/
│   └── env.go               # Configuração de variáveis de ambiente
├── controller/
│   ├── video.go             # Controllers da API REST
│   └── user.go
├── database/
│   └── database.go          # Conexão com PostgreSQL
├── input/
│   ├── cli.go               # Parser de argumentos CLI
│   ├── downloader.go        # Download de vídeos
│   └── youtube.go           # Integração com YouTube
├── transcript/
│   ├── whisper.go           # Integração com Whisper
│   └── transcripter.go      # Interface de transcrição
├── translator/
│   ├── marianmt.go          # Integração MarianMT
│   ├── googletranslator.go  # Google Translate API
│   └── translator.go        # Interface de tradução
├── tts/
│   ├── coquitts.go          # Coqui TTS
│   └── tts.go               # Interface TTS
├── caption/
│   ├── subtitler.go         # Interface de legendas
│   └── stablets.go          # Stable-TS para sincronização
├── media/
│   ├── ffmpegWrapper.go     # Wrapper FFmpeg
│   └── mediaHandler.go      # Handler de mídia
├── pipeline/
│   └── pipeline.go          # Orquestração do processamento
├── scripts/
│   └── translator_marianmt.py  # Script standalone Python
├── requirements.txt         # Dependências Python
├── docker-compose.yaml      # Configuração Docker
├── .env.template            # Template de variáveis de ambiente
└── init.sql                 # Schema inicial do banco
```

---

## 📝 Notas Importantes

1. **Python 3.10 é obrigatório** para compatibilidade com as versões das dependências
2. **Sempre ative o virtualenv** antes de executar comandos Python
3. **Atualize setuptools** antes de instalar as dependências
4. O projeto usa **Go Fiber v3** (versão beta) para a API
5. A primeira execução do Whisper baixará modelos automaticamente (~1-2GB)

---

## 👥 Contato

Para dúvidas ou problemas, verifique:
- Logs da aplicação: `docker compose logs -f` ou saída do terminal
- Issues do projeto no GitHub
- Documentação das dependências: Whisper, Coqui TTS, Stable-TS
