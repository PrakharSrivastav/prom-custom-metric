package main

import (
	"bytes"
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
	"os"
	"os/exec"
	"sync"
	"time"
)

func main() {

	logger, err := setupLogger(logfilePath)
	if err != nil {
		panic(err)
	}
	logger.Info("START")
	defer logger.Info("END")

	c := new(Config)
	yamlFile, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("yaml.config.error", zap.Error(err))
		return
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		logger.Error("yaml.parse.error", zap.Error(err))
		return
	}

	if len(c.Items) <= 0 {
		logger.Error("config.not.provided")
		return
	}
	var finalStatus []string

	wg := sync.WaitGroup{}
	// run all the sftp connections in parallel
	for _, item := range c.Items {
		wg.Add(1)
		item := item
		go func() {
			defer wg.Done()
			l := logger.With(zap.String("name", item.Name))
			ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*15)
			defer cancelFunc()
			_, _, err := cleanArgsAndExecute(ctx, item.Args)
			if err != nil {
				l.Error("error.parsing.command", zap.Error(err))
				finalStatus = append(finalStatus, fmt.Sprintf(template, item.Name, 0))
				return
			}
			finalStatus = append(finalStatus, fmt.Sprintf(template, item.Name, 1))
			l.Info("complete")
			return
		}()
	}

	wg.Wait()

	f, err := os.OpenFile(promFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Error("error.opening.local.prom.file", zap.Error(err))
		return
	}

	_, _ = f.WriteString("# TYPE sftp_outbound_status gauge")
	for i := range finalStatus {
		_, _ = f.WriteString("\n")
		_, _ = f.WriteString(finalStatus[i])
	}
	_, _ = f.WriteString("\n")
	_ = f.Sync()
	_ = f.Close()

}

func setupLogger(filename string) (*zap.Logger, error) {
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.ISO8601TimeEncoder
	fileEncoder := zapcore.NewConsoleEncoder(config)
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	writer := zapcore.AddSync(logFile)
	defaultLogLevel := zapcore.DebugLevel
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, writer, defaultLogLevel),
	)
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	return logger, nil
}

func cleanArgsAndExecute(ctx context.Context, command string) (string, string, error) {
	if ctx.Err() != nil {
		return "", "", ctx.Err()
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.CommandContext(ctx, "bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return "", "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", "", err
	}

	return stdout.String(), stderr.String(), nil
}
