package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v3"
	"bufio"
)

type Step struct {
	Type    string   `yaml:"type"`
	Var     string   `yaml:"var,omitempty"`
	Message string   `yaml:"message,omitempty"`
	Default string   `yaml:"default,omitempty"`
	Options []string `yaml:"options,omitempty"`
	When    string   `yaml:"when,omitempty"`
	Command string   `yaml:"command,omitempty"`
	Args    []string `yaml:"args,omitempty"`
	Path    string   `yaml:"path,omitempty"`
	Content string   `yaml:"content,omitempty"`
	Text    string   `yaml:"text,omitempty"`
}

type ToolchainConfig struct {
	Name    string            `yaml:"name"`
	Steps   []Step            `yaml:"steps"`
	EnvVars map[string]string `yaml:"env_vars"`
}

// LoadToolchainConfig loads a toolchain YAML config from file
func LoadToolchainConfig(path string) (*ToolchainConfig, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var cfg ToolchainConfig
	dec := yaml.NewDecoder(f)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// RunToolchainSteps executes the steps in the toolchain config
func RunToolchainSteps(cfg *ToolchainConfig, vars map[string]string) error {
	reader := bufio.NewReader(os.Stdin)
	for _, step := range cfg.Steps {
		if step.When != "" && !evalCondition(step.When, vars) {
			continue
		}
		switch step.Type {
		case "prompt":
			msg := substitute(step.Message, vars)
			def := substitute(step.Default, vars)
			fmt.Printf("%s [%s]: ", msg, def)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)
			if input == "" {
				input = def
			}
			vars[step.Var] = input
		case "select":
			msg := substitute(step.Message, vars)
			fmt.Println(msg)
			for i, opt := range step.Options {
				fmt.Printf("  [%d] %s\n", i+1, opt)
			}
			fmt.Printf("Enter number [1]: ")
			var sel int
			fmt.Scanf("%d\n", &sel)
			if sel < 1 || sel > len(step.Options) {
				sel = 1
			}
			vars[step.Var] = step.Options[sel-1]
		case "run":
			cmdStr := substitute(step.Command, vars)
			args := make([]string, len(step.Args))
			for i, a := range step.Args {
				args[i] = substitute(a, vars)
			}
			fmt.Printf("Running: %s %s\n", cmdStr, strings.Join(args, " "))
			cmd := exec.Command(cmdStr, args...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("command failed: %v", err)
			}
		case "file":
			path := substitute(step.Path, vars)
			content := substitute(step.Content, vars)
			os.WriteFile(path, []byte(content), 0644)
		case "message":
			fmt.Println(substitute(step.Text, vars))
		}
	}
	return nil
}

// substitute replaces {var} in s with values from vars
func substitute(s string, vars map[string]string) string {
	for k, v := range vars {
		s = strings.ReplaceAll(s, "{"+k+"}", v)
	}
	return s
}

// evalCondition supports simple equality checks: "{var}" == "value"
func evalCondition(cond string, vars map[string]string) bool {
	cond = substitute(cond, vars)
	parts := strings.Split(cond, "==")
	if len(parts) != 2 {
		return false
	}
	return strings.TrimSpace(parts[0]) == strings.TrimSpace(strings.Trim(parts[1], " \""))
}
