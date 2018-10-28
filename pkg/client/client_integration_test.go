// +build integration

package client

import (
	"context"
	"crypto/sha256"
	"math/big"
	"net/http"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func TestITClient_Register(t *testing.T) {

	baseURL := "http://localhost:8080"

	t.Run("should register a new user ID", func(t *testing.T) {
		clt := New(&http.Client{}, Configuration{baseURL: baseURL, registerPath: "/v1/register"}, NewInMemoryUserRepository())
		err := clt.Register("new-user")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should be able to register an existing user ID", func(t *testing.T) {
		clt := New(&http.Client{}, Configuration{baseURL: baseURL, registerPath: "/v1/register"}, NewInMemoryUserRepository())
		err := clt.Register("another-new-user")
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}

		clt = New(&http.Client{}, Configuration{baseURL: baseURL, registerPath: "/v1/register"}, NewInMemoryUserRepository())
		err = clt.Register("another-new-user")
		if err == ErrRegistrationFailed {
			t.Errorf("Register() error = %v wantErr = %v", err, ErrRegistrationFailed)
		}
	})

}

func TestITClient_Login(t *testing.T) {

	// often used big.Int
	var two = big.NewInt(2)

	baseURL := "http://localhost:8080"
	bits := 8

	max := new(big.Int)
	max.Exp(two, big.NewInt(int64(bits)), nil)

	clt := New(&http.Client{}, Configuration{
		hash:         sha256.New,
		bits:         big.NewInt(int64(bits)),
		baseURL:      baseURL,
		registerPath: "/v1/register",
		expkPath:     "/v1/login/expk",
		verifyPath:   "/v1/login/verify",
	}, NewInMemoryUserRepository())

	err := clt.Register("user")
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should login with an valid password", func(t *testing.T) {
		err := clt.Login("user", "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})

	t.Run("should recv. common error if login failed because of wrong password", func(t *testing.T) {
		err := clt.Login("user", "wrong-password")
		if err == ErrAuthenticationFailed {
			t.Errorf("Login() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
		clt.Logout()
	})

	t.Run("should recv. error if configuration is invalid", func(t *testing.T) {
		err := clt.Login("wrong-user", "password")
		if err == ErrAuthenticationFailed {
			t.Errorf("Login() error = %v wantErr = %v", err, ErrAuthenticationFailed)
		}
		clt.Logout()
	})
}

func before(t *testing.T) string {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}

	imageName := "latanassov/ossrv:0.1.0"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		t.Fatal(err)
	}

	config := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"8080/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"8080/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "9090",
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, config, hostConfig, nil, "")
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		t.Fatal(err)
	}

	os.Setenv("OSSRV_DOCKER_ID", resp.ID)

	return "http://localhost:9090"
}

func after(t *testing.T) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}

	containerID := os.Getenv("OSSRV_DOCKER_ID")

	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		t.Fatal(err)
	}
}
