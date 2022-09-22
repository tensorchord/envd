package envd

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/tensorchord/envd-server/client"
	"github.com/tensorchord/envd/pkg/lang/ir"
	"github.com/tensorchord/envd/pkg/types"
)

type envdServerEngine struct {
	*client.Client
}

func (e *envdServerEngine) ListImage(ctx context.Context) ([]types.EnvdImage, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) ListImageDependency(ctx context.Context, image string) (*types.Dependency, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) GetImage(ctx context.Context, image string) (dockertypes.ImageSummary, error) {
	return dockertypes.ImageSummary{}, errors.New("not implemented")
}

func (e *envdServerEngine) GetInfo(ctx context.Context) (*types.EnvdInfo, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) GPUEnabled(ctx context.Context) (bool, error) {
	return false, errors.New("not implemented")
}

func (e *envdServerEngine) PauseEnvironment(ctx context.Context, env string) (string, error) {
	return "", errors.New("not implemented")
}

func (e *envdServerEngine) ResumeEnvironment(ctx context.Context, env string) (string, error) {
	return "", errors.New("not implemented")
}

func (e *envdServerEngine) ListEnvironment(ctx context.Context) ([]types.EnvdEnvironment, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) ListEnvDependency(ctx context.Context, env string) (*types.Dependency, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) ListEnvPortBinding(ctx context.Context, env string) ([]types.PortBinding, error) {
	return nil, errors.New("not implemented")
}

func (e *envdServerEngine) CleanEnvdIfExists(ctx context.Context, name string, force bool) error {
	return errors.New("not implemented")
}

// StartEnvd creates the container for the given tag and container name.
func (e *envdServerEngine) StartEnvd(ctx context.Context, tag, name, buildContext string,
	gpuEnabled bool, numGPUs int, sshPort int, g ir.Graph, timeout time.Duration,
	mountOptionsStr []string) (string, string, error) {
	return "", "", errors.New("not implemented")
}

func (e *envdServerEngine) IsRunning(ctx context.Context, name string) (bool, error) {
	return false, errors.New("not implemented")
}

func (e *envdServerEngine) Exists(ctx context.Context, name string) (bool, error) {
	return false, errors.New("not implemented")
}

func (e *envdServerEngine) WaitUntilRunning(ctx context.Context, name string, timeout time.Duration) error {
	return errors.New("not implemented")
}
