package checker

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DependencyCheck representa uma verificação de dependência
type DependencyCheck struct {
	Name        string
	Check       func() error
	FixCommand  string
	Required    bool
	Description string
}

// ValidationResult contém o resultado da validação
type ValidationResult struct {
	Passed   bool
	Errors   []string
	Warnings []string
}

// RunAllChecks executa todas as verificações de dependências
func RunAllChecks() error {
	slog.Info("[CHECKER] Iniciando verificação de dependências...")

	checks := []DependencyCheck{
		{
			Name:        "Virtual Environment",
			Check:       checkVirtualEnv,
			FixCommand:  "python3.10 -m venv .venv",
			Required:    true,
			Description: "Ambiente virtual Python",
		},
		{
			Name:        "Python",
			Check:       checkPython,
			FixCommand:  "source .venv/bin/activate && python --version",
			Required:    true,
			Description: "Python do virtualenv",
		},
		{
			Name:        "FFmpeg",
			Check:       checkFFmpeg,
			FixCommand:  "sudo apt-get install ffmpeg (Linux) ou brew install ffmpeg (Mac)",
			Required:    true,
			Description: "FFmpeg para processamento de mídia",
		},
		{
			Name:        "yt-dlp",
			Check:       checkYtDlp,
			FixCommand:  "curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o ~/.local/bin/yt-dlp && chmod a+rx ~/.local/bin/yt-dlp",
			Required:    true,
			Description: "Downloader de vídeos do YouTube",
		},
		{
			Name:        "transformers",
			Check:       checkPythonPackage("transformers", "4.47.0"),
			FixCommand:  "source .venv/bin/activate && pip install transformers==4.47.0",
			Required:    true,
			Description: "Biblioteca HuggingFace Transformers",
		},
		{
			Name:        "sentencepiece",
			Check:       checkPythonPackage("sentencepiece", ""),
			FixCommand:  "source .venv/bin/activate && pip install sentencepiece==0.2.0",
			Required:    true,
			Description: "Tokenizador SentencePiece",
		},
		{
			Name:        "sacremoses",
			Check:       checkPythonPackage("sacremoses", ""),
			FixCommand:  "source .venv/bin/activate && pip install sacremoses==0.1.1",
			Required:    true,
			Description: "Ferramenta de tokenização Moses",
		},
		{
			Name:        "openai-whisper",
			Check:       checkPythonPackage("whisper", ""),
			FixCommand:  "source .venv/bin/activate && pip install openai-whisper==20240930",
			Required:    true,
			Description: "Modelo Whisper para transcrição",
		},
		{
			Name:        "coqui-tts",
			Check:       checkCoquiTTS,
			FixCommand:  "source .venv/bin/activate && pip install coqui-tts==0.26.0",
			Required:    true,
			Description: "TTS (Text-to-Speech)",
		},
		{
			Name:        "stable-ts",
			Check:       checkStableTs,
			FixCommand:  "source .venv/bin/activate && pip install stable-ts==2.17.5 --no-deps",
			Required:    true,
			Description: "Sincronização de legendas",
		},
		{
			Name:        "torchcodec (ausente)",
			Check:       checkTorchcodecAbsent,
			FixCommand:  "source .venv/bin/activate && pip uninstall torchcodec -y",
			Required:    true,
			Description: "Verificar que torchcodec NÃO está instalado",
		},
		{
			Name:        "Variáveis de Ambiente",
			Check:       checkEnvVars,
			FixCommand:  "cp .env.template .env && # editar com suas credenciais",
			Required:    true,
			Description: "Configuração de variáveis de ambiente",
		},
		{
			Name:        "PostgreSQL",
			Check:       checkPostgreSQL,
			FixCommand:  "docker compose up -d postgres",
			Required:    true,
			Description: "Banco de dados PostgreSQL",
		},
	}

	var hasErrors bool
	var hasWarnings bool

	for _, check := range checks {
		slog.Info(fmt.Sprintf("[CHECKER] Verificando: %s", check.Description))

		if err := check.Check(); err != nil {
			if check.Required {
				slog.Error(fmt.Sprintf("[CHECKER] ❌ ERRO: %s", check.Name))
				slog.Error(fmt.Sprintf("[CHECKER]    Motivo: %s", err.Error()))
				slog.Error(fmt.Sprintf("[CHECKER]    Correção: %s", check.FixCommand))
				hasErrors = true
			} else {
				slog.Warn(fmt.Sprintf("[CHECKER] ⚠️  AVISO: %s", check.Name))
				slog.Warn(fmt.Sprintf("[CHECKER]    Motivo: %s", err.Error()))
				slog.Warn(fmt.Sprintf("[CHECKER]    Correção: %s", check.FixCommand))
				hasWarnings = true
			}
		} else {
			slog.Info(fmt.Sprintf("[CHECKER] ✅ OK: %s", check.Name))
		}
	}

	if hasErrors {
		return fmt.Errorf("\nVERIFICACAO FALHOU\n\nExistem dependencias obrigatorias nao satisfeitas.\nCorrija os erros acima antes de iniciar o servidor.\n\nVeja o README.md para instrucoes detalhadas.")
	}

	if hasWarnings {
		slog.Warn("[CHECKER] Verificacao concluida com avisos, mas prosseguindo...")
	}

	slog.Info("[CHECKER] OK: Todas as dependencias estao satisfeitas!")
	return nil
}

// checkVirtualEnv verifica se o virtualenv existe
func checkVirtualEnv() error {
	venvPaths := []string{"./.venv", ".venv", "venv"}

	for _, path := range venvPaths {
		if info, err := os.Stat(path); err == nil && info.IsDir() {
			// Verificar se o Python existe dentro do venv
			pythonPath := filepath.Join(path, "bin", "python")
			if _, err := os.Stat(pythonPath); err == nil {
				return nil
			}
		}
	}

	return fmt.Errorf("virtualenv não encontrado em ./.venv")
}

// checkPython verifica se o Python do virtualenv funciona
func checkPython() error {
	pythonPath := "./.venv/bin/python"

	if _, err := os.Stat(pythonPath); os.IsNotExist(err) {
		return fmt.Errorf("Python não encontrado em %s", pythonPath)
	}

	cmd := exec.Command(pythonPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao executar Python: %v", err)
	}

	version := strings.TrimSpace(string(output))
	slog.Info(fmt.Sprintf("[CHECKER]    Versão: %s", version))

	return nil
}

// checkFFmpeg verifica se o FFmpeg está instalado
func checkFFmpeg() error {
	cmd := exec.Command("ffmpeg", "-version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("FFmpeg não encontrado no PATH")
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		slog.Info(fmt.Sprintf("[CHECKER]    Versão: %s", strings.TrimSpace(lines[0])))
	}

	return nil
}

// checkYtDlp verifica se o yt-dlp está instalado
func checkYtDlp() error {
	// Verificar no PATH
	cmd := exec.Command("which", "yt-dlp")
	if err := cmd.Run(); err != nil {
		// Verificar em ~/.local/bin
		localBin := filepath.Join(os.Getenv("HOME"), ".local", "bin", "yt-dlp")
		if _, err := os.Stat(localBin); err != nil {
			return fmt.Errorf("yt-dlp não encontrado")
		}
	}

	cmd = exec.Command("yt-dlp", "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("yt-dlp não executável: %v", err)
	}

	version := strings.TrimSpace(string(output))
	slog.Info(fmt.Sprintf("[CHECKER]    Versão: %s", version))

	return nil
}

// checkPythonPackage verifica se um pacote Python está instalado
func checkPythonPackage(packageName, expectedVersion string) func() error {
	return func() error {
		pythonPath := "./.venv/bin/python"

		// Verificar importação do pacote
		script := fmt.Sprintf("import %s; print(%s.__version__)", packageName, packageName)
		if packageName == "whisper" {
			script = "import whisper; print(whisper.__version__)"
		}

		cmd := exec.Command(pythonPath, "-c", script)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("pacote %s não encontrado", packageName)
		}

		version := strings.TrimSpace(string(output))
		slog.Info(fmt.Sprintf("[CHECKER]    Versão: %s", version))

		// Verificar versão esperada (se especificada)
		if expectedVersion != "" && !strings.HasPrefix(version, expectedVersion) {
			slog.Warn(fmt.Sprintf("[CHECKER]    ⚠️  Versão diferente da esperada: %s", expectedVersion))
		}

		return nil
	}
}

// checkCoquiTTS verifica se o Coqui TTS está instalado
func checkCoquiTTS() error {
	pythonPath := "./.venv/bin/python"

	// Verificar se o pacote TTS pode ser importado
	cmd := exec.Command(pythonPath, "-c", "import TTS; print(TTS.__version__)")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Coqui TTS não encontrado")
	}

	version := strings.TrimSpace(string(output))
	slog.Info(fmt.Sprintf("[CHECKER]    Versão: %s", version))

	// Verificar se o comando tts existe
	ttsPath := "./.venv/bin/tts"
	if _, err := os.Stat(ttsPath); err != nil {
		return fmt.Errorf("comando tts não encontrado em %s", ttsPath)
	}

	return nil
}

// checkStableTs verifica se o stable-ts está instalado
func checkStableTs() error {
	stableTsPath := "./.venv/bin/stable-ts"

	if _, err := os.Stat(stableTsPath); err != nil {
		return fmt.Errorf("stable-ts não encontrado em %s", stableTsPath)
	}

	// Testar execução
	cmd := exec.Command(stableTsPath, "--help")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil && !strings.Contains(stderr.String(), "usage") {
		return fmt.Errorf("stable-ts não executável")
	}

	return nil
}

// checkTorchcodecAbsent verifica que torchcodec NÃO está instalado (causa erros)
func checkTorchcodecAbsent() error {
	pythonPath := "./.venv/bin/python"

	cmd := exec.Command(pythonPath, "-c", "import torchcodec")
	err := cmd.Run()

	if err == nil {
		return fmt.Errorf("torchcodec está instalado, mas deve ser removido (causa incompatibilidades)")
	}

	// Erro esperado - pacote não encontrado é o que queremos
	return nil
}

// checkEnvVars verifica se as variáveis de ambiente essenciais estão configuradas
func checkEnvVars() error {
	requiredVars := []string{"DATABASE_URL"}
	optionalVars := []string{"CLIENT_ID", "CLIENT_SECRET", "REDIRECT_URI"}

	var missingRequired []string
	var missingOptional []string

	for _, v := range requiredVars {
		if os.Getenv(v) == "" {
			missingRequired = append(missingRequired, v)
		}
	}

	for _, v := range optionalVars {
		if os.Getenv(v) == "" {
			missingOptional = append(missingOptional, v)
		}
	}

	if len(missingRequired) > 0 {
		return fmt.Errorf("variáveis obrigatórias não definidas: %s", strings.Join(missingRequired, ", "))
	}

	if len(missingOptional) > 0 {
		slog.Warn(fmt.Sprintf("[CHECKER]    Variáveis opcionais não definidas: %s", strings.Join(missingOptional, ", ")))
	}

	return nil
}

// checkPostgreSQL verifica se o PostgreSQL está acessível
func checkPostgreSQL() error {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return fmt.Errorf("DATABASE_URL não definida")
	}

	// Verificar formato da URL
	if !strings.HasPrefix(dbURL, "postgresql://") && !strings.HasPrefix(dbURL, "postgres://") {
		return fmt.Errorf("DATABASE_URL deve começar com postgresql://")
	}

	// Extrair host e porta
	host := "localhost"
	port := "5432"

	if strings.Contains(dbURL, "@") {
		parts := strings.Split(dbURL, "@")
		if len(parts) == 2 {
			hostPart := strings.Split(parts[1], "/")[0]
			if strings.Contains(hostPart, ":") {
				hostPort := strings.Split(hostPart, ":")
				host = hostPort[0]
				port = hostPort[1]
			} else {
				host = hostPart
			}
		}
	}

	// Tentar conectar
	address := fmt.Sprintf("%s:%s", host, port)

	// Verificar se a porta está aberta
	conn, err := exec.Command("nc", "-z", host, port).CombinedOutput()
	if err != nil {
		// Tentar com timeout se nc não funcionar
		cmd := exec.Command("timeout", "2", "bash", "-c", fmt.Sprintf("echo > /dev/tcp/%s/%s", host, port))
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("não foi possível conectar ao PostgreSQL em %s: %v (verifique se o container está rodando: docker compose up -d postgres)", address, err)
		}
	}

	_ = conn
	slog.Info(fmt.Sprintf("[CHECKER]    Conectado: %s", address))

	return nil
}
