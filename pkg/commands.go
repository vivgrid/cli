package pkg

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type command struct {
	tid     string
	apiAddr string
	secret  string
	tool    string
	envs    []string
}

func parseEnvs(values []string) (map[string]string, error) {
	if len(values) == 0 {
		return map[string]string{}, nil
	}
	envs := make(map[string]string, len(values))
	for _, value := range values {
		key, val, ok := strings.Cut(value, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("invalid env %q, expected key=value", value)
		}
		envs[key] = val
	}
	return envs, nil
}

func Execute(rootCmd *cobra.Command, configFile string, tid string, defaultAPIAddr string) error {
	c := &command{tid: tid}

	rootCmd.PersistentFlags().StringVar(&c.apiAddr, "api", defaultAPIAddr, "REST API endpoint")
	rootCmd.PersistentFlags().StringVar(&c.secret, "secret", "", "app secret")
	rootCmd.PersistentFlags().StringVar(&c.tool, "tool", "my_first_llm_tool", "serverless LLM tool name")

	uploadCmd := c.addUploadCmd(rootCmd)
	removeCmd := c.addRemoveCmd(rootCmd)
	createCmd := c.addCreateCmd(rootCmd)

	c.addVersionCmd(rootCmd)
	c.addStatusCmd(rootCmd)
	c.addLogsCmd(rootCmd)
	c.addDeployCmd(rootCmd, uploadCmd, removeCmd, createCmd)
	c.addDocCmd(rootCmd)

	rootCmd.AddGroup(&cobra.Group{ID: groupIDGeneral, Title: colorBlue + "General" + colorReset})
	rootCmd.AddGroup(&cobra.Group{ID: groupIDDeployment, Title: colorBlue + "Manage serverless deployment" + colorReset})
	rootCmd.AddGroup(&cobra.Group{ID: groupIDMonitoring, Title: colorBlue + "Observability" + colorReset})

	if configFile != "" {
		v := viper.GetViper()
		v.SetConfigFile(configFile)
		if err := v.ReadInConfig(); err != nil {
			return err
		}
		if v.IsSet("api") {
			c.apiAddr = v.GetString("api")
		}
		if v.IsSet("secret") {
			c.secret = v.GetString("secret")
		}
		if v.IsSet("tool") {
			c.tool = v.GetString("tool")
		}
	}

	return rootCmd.Execute()
}

func (c *command) apiURL(parts ...string) string {
	return c.apiAddr + "/api/tool/" + c.tool + "/" + strings.Join(parts, "/")
}

func (c *command) do(req *http.Request) ([]byte, error) {
	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return nil, errors.New(strings.TrimSpace(string(body)))
	}
	return body, nil
}

func (c *command) doStream(req *http.Request) error {
	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return errors.New(strings.TrimSpace(string(body)))
	}

	s := bufio.NewScanner(resp.Body)
	for s.Scan() {
		line := s.Text()
		if strings.HasPrefix(line, ":") {
			continue
		}
		if strings.HasPrefix(line, "event: end") {
			return nil
		}
		if strings.HasPrefix(line, "data:") {
			fmt.Println(strings.TrimSpace(strings.TrimPrefix(line, "data:")))
			continue
		}
		if line != "" {
			fmt.Println(line)
		}
	}
	if err := s.Err(); err != nil && !errors.Is(err, io.EOF) {
		return err
	}
	return nil
}

func (c *command) addDocCmd(rootCmd *cobra.Command) {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.DisableAutoGenTag = true
	rootCmd.AddCommand(&cobra.Command{Use: "doc", Short: "Generate documentation for the CLI commands", Hidden: true, RunE: func(_ *cobra.Command, _ []string) error { return GenDoc(rootCmd) }})
}

func (c *command) addVersionCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{Use: "version", Short: "Show version", Args: cobra.ExactArgs(0), Run: func(_ *cobra.Command, _ []string) { fmt.Println("version:", CliVersion) }})
}

func (c *command) addUploadCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{Use: "upload src_file[.go|.zip|dir]", Short: "Upload the source code and compile", Args: cobra.ExactArgs(1), RunE: func(_ *cobra.Command, args []string) error {
		src := args[0]
		info, err := os.Stat(src)
		if err != nil {
			return err
		}
		var data []byte
		if info.IsDir() {
			f, err := os.CreateTemp("", "app-*.zip")
			if err != nil {
				return err
			}
			zipPath := f.Name()
			defer os.Remove(zipPath)
			defer f.Close()
			if err := ZipWithExclusions(src, zipPath); err != nil {
				return err
			}
			data, err = os.ReadFile(zipPath)
			if err != nil {
				return err
			}
		} else {
			switch path.Ext(src) {
			case ".zip":
				data, err = os.ReadFile(src)
				if err != nil {
					return err
				}
			case ".go":
				buf := new(bytes.Buffer)
				writer := zip.NewWriter(buf)
				f, err := writer.Create("app.go")
				if err != nil {
					return err
				}
				content, err := os.ReadFile(src)
				if err != nil {
					return err
				}
				if _, err = f.Write(content); err != nil {
					return err
				}
				if err := writer.Close(); err != nil {
					return err
				}
				data = buf.Bytes()
			default:
				return errors.New("unsupported src file type")
			}
		}
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		if err := writer.WriteField("language", "auto"); err != nil {
			return err
		}
		part, err := writer.CreateFormFile("zip_file", "app.zip")
		if err != nil {
			return err
		}
		if _, err := part.Write(data); err != nil {
			return err
		}
		if err := writer.Close(); err != nil {
			return err
		}
		req, err := http.NewRequest(http.MethodPost, c.apiURL("build"), body)
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		req.Header.Set("Accept", "text/event-stream")
		return c.doStream(req)
	}}
	rootCmd.AddCommand(cmd)
	return cmd
}

func (c *command) addCreateCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{Use: "create", Short: "Create serverless deployment and start it", Args: cobra.ExactArgs(0), RunE: func(_ *cobra.Command, _ []string) error {
		envs, err := parseEnvs(c.envs)
		if err != nil {
			return err
		}
		payload, err := json.Marshal(map[string]any{"envs": envs})
		if err != nil {
			return err
		}
		req, err := http.NewRequest(http.MethodPost, c.apiURL("create"), bytes.NewReader(payload))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		res, err := c.do(req)
		if err != nil {
			return err
		}
		fmt.Print(string(res))
		return nil
	}}
	rootCmd.AddCommand(cmd)
	cmd.Flags().StringArrayVar(&c.envs, "env", nil, "Set environment variable as key=value")
	return cmd
}

func (c *command) addRemoveCmd(rootCmd *cobra.Command) *cobra.Command {
	cmd := &cobra.Command{Use: "remove", Short: "Delete current serverless deployment", Args: cobra.ExactArgs(0), RunE: func(_ *cobra.Command, _ []string) error {
		req, err := http.NewRequest(http.MethodDelete, c.apiURL("remove"), nil)
		if err != nil {
			return err
		}
		res, err := c.do(req)
		if err != nil {
			return err
		}
		fmt.Print(string(res))
		return nil
	}, GroupID: groupIDDeployment}
	rootCmd.AddCommand(cmd)
	return cmd
}

func (c *command) addStatusCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{Use: "status", Short: "Show serverless status", Args: cobra.ExactArgs(0), RunE: func(_ *cobra.Command, _ []string) error {
		req, err := http.NewRequest(http.MethodGet, c.apiURL("status"), nil)
		if err != nil {
			return err
		}
		return c.doStream(req)
	}, GroupID: groupIDMonitoring})
}

func (c *command) addLogsCmd(rootCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{Use: "logs", Short: "Observe serverless logs in real-time", Args: cobra.ExactArgs(0), RunE: func(_ *cobra.Command, _ []string) error {
		req, err := http.NewRequest(http.MethodGet, c.apiURL("logs"), nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "text/event-stream")
		return c.doStream(req)
	}, GroupID: groupIDMonitoring})
}

func (c *command) addDeployCmd(rootCmd *cobra.Command, uploadCmd *cobra.Command, removeCmd *cobra.Command, createCmd *cobra.Command) {
	rootCmd.AddCommand(&cobra.Command{Use: "deploy src_file[.go|.zip|dir]", Short: "Deploy your serverless, this is an alias of chaining commands (upload -> remove -> create)", Args: cobra.ExactArgs(1), Run: func(cmd *cobra.Command, args []string) {
		uploadCmd.RunE(uploadCmd, args)
		removeCmd.RunE(removeCmd, nil)
		createCmd.RunE(createCmd, nil)
		fmt.Println("Successfully!")
	}, GroupID: groupIDGeneral})
	rootCmd.Flags().StringArrayVar(&c.envs, "env", nil, "Set environment variables as key=value")
}

const (
	groupIDDeployment = "deployment"
	groupIDMonitoring = "monitoring"
	groupIDGeneral    = "general"

	colorReset = "\033[0m"
	colorBlue  = "\033[34m"
)
