package k8s

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"

	"github.com/mskKandula/oes/code-executor/internal/model"
)

// ExecResult holds the raw output from a pod exec call.
type ExecResult struct {
	Stdout     string
	Stderr     string
	ExitCode   int
	DurationMs int64
}

// Execute runs the code from job inside the named pod using the Kubernetes
// pods/exec subresource (equivalent to `kubectl exec <pod> -- <cmd>`).
//
// Language → command mapping:
//
//	python  →  python3 -c <code>
//	nodejs  →  node    -e <code>
//	go      →  go run /dev/stdin  (code piped via stdin)
//
// The execution is bounded by job.TimeoutMs via context deadline.
func Execute(
	ctx context.Context,
	cfg *rest.Config,
	clientset *kubernetes.Clientset,
	namespace string,
	podName string,
	job model.CodeJob,
) (*ExecResult, error) {

	cmd, useStdin := buildCommand(job.Language, job.Code)

	execOpts := &corev1.PodExecOptions{
		Command: cmd,
		Stdout:  true,
		Stderr:  true,
		Stdin:   useStdin || job.Stdin != "",
		TTY:     false,
	}

	req := clientset.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(execOpts, scheme.ParameterCodec)

	executor, err := remotecommand.NewSPDYExecutor(cfg, "POST", req.URL())
	if err != nil {
		return nil, fmt.Errorf("failed to create SPDY executor: %w", err)
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	// Build stdin reader: for Go use code as stdin; for others use job.Stdin if provided
	var stdinReader *strings.Reader
	if useStdin {
		stdinReader = strings.NewReader(job.Code)
	} else if job.Stdin != "" {
		stdinReader = strings.NewReader(job.Stdin)
	}

	// Apply timeout from job
	execCtx, cancel := context.WithTimeout(ctx, time.Duration(job.TimeoutMs)*time.Millisecond)
	defer cancel()

	start := time.Now()

	streamOpts := remotecommand.StreamOptions{
		Stdout: &stdoutBuf,
		Stderr: &stderrBuf,
	}
	if stdinReader != nil {
		streamOpts.Stdin = stdinReader
	}

	err = executor.StreamWithContext(execCtx, streamOpts)

	durationMs := time.Since(start).Milliseconds()

	exitCode := 0
	if err != nil {
		// Check for context timeout
		if execCtx.Err() == context.DeadlineExceeded {
			return &ExecResult{
				Stdout:     stdoutBuf.String(),
				Stderr:     "execution timed out after " + fmt.Sprintf("%dms", job.TimeoutMs),
				ExitCode:   124, // standard timeout exit code (same as GNU timeout command)
				DurationMs: durationMs,
			}, nil
		}
		// Extract non-zero exit code from exec error
		exitCode = extractExitCode(err)
		if exitCode == 0 {
			// Unknown error — not an exit code error
			return nil, fmt.Errorf("exec stream error: %w", err)
		}
	}

	return &ExecResult{
		Stdout:     stdoutBuf.String(),
		Stderr:     stderrBuf.String(),
		ExitCode:   exitCode,
		DurationMs: durationMs,
	}, nil
}

// buildCommand returns the exec command and whether code should be passed via stdin.
//
//	python  → ["python3", "-c", code],  stdin=false
//	nodejs  → ["node",    "-e", code],  stdin=false
//	go      → ["go", "run", "/dev/stdin"], stdin=true  (code piped via stdin)
func buildCommand(language, code string) ([]string, bool) {
	switch language {
	case "python":
		return []string{"python3", "-c", code}, false
	case "nodejs":
		return []string{"node", "-e", code}, false
	case "go":
		// go run /dev/stdin reads the source from stdin — no temp file needed
		return []string{"go", "run", "/dev/stdin"}, true
	default:
		// Fallback: try to run as a shell command
		return []string{"sh", "-c", code}, false
	}
}

// extractExitCode attempts to extract the process exit code from a Kubernetes
// exec error. Returns 1 if the error is not an exit-code error.
func extractExitCode(err error) int {
	if err == nil {
		return 0
	}
	// Kubernetes wraps non-zero exit codes in a CodeExitError / execError
	// The error message format is: "command terminated with exit code N"
	msg := err.Error()
	if strings.Contains(msg, "exit code") {
		var code int
		if _, scanErr := fmt.Sscanf(msg, "command terminated with exit code %d", &code); scanErr == nil {
			return code
		}
	}
	return 1
}
