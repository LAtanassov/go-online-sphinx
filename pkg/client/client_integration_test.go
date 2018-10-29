// +build integration

package client

import (
	"context"
	"crypto/sha256"
	"math/big"
	"net/http"
	"net/http/cookiejar"
	"os"
	"reflect"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/contract"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func TestITClient_Register(t *testing.T) {

	baseURL := "http://localhost:8080"

	t.Run("should register a new user ID", func(t *testing.T) {
		clt := New(&http.Client{}, Configuration{baseURL: baseURL, registerPath: "/v1/register"}, NewInMemoryUserRepository())
		user, err := NewUser("new-user", 8)
		if err != nil {
			t.Errorf("NewUser() error = %v", err)
		}
		err = clt.Register(user)
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
	})

	t.Run("should not be able to register with an existing user ID", func(t *testing.T) {
		clt := New(&http.Client{}, Configuration{baseURL: baseURL, registerPath: "/v1/register"}, NewInMemoryUserRepository())
		// given
		user, _ := NewUser("another-new-user", 8)
		err := clt.Register(user)
		if err != nil {
			t.Errorf("Register() error = %v", err)
		}
		// when
		err = clt.Register(user)
		if err == nil {
			t.Errorf("Register() error = %v wantErr = %v", err, contract.ErrRegistrationFailed)
		}
	})

}

func TestITClient_Login(t *testing.T) {

	var two = big.NewInt(2)

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	baseURL := "http://localhost:8080"
	bits := 8

	max := new(big.Int)
	max.Exp(two, big.NewInt(int64(bits)), nil)

	clt := New(client, Configuration{
		hash:          sha256.New,
		bits:          bits,
		baseURL:       baseURL,
		registerPath:  "/v1/register",
		expkPath:      "/v1/login/expk",
		challengePath: "/v1/login/challenge",
		metadataPath:  "/v1/metadata",
		addPath:       "/v1/add",
		getPath:       "/v1/get",
	}, NewInMemoryUserRepository())

	user, _ := NewUser("username", bits)
	err := clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should login with a valid password", func(t *testing.T) {
		err := clt.Login(user.username, "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})

	t.Run("should recv. common error if wrong password", func(t *testing.T) {
		err := clt.Login(user.username, "wrong-password")
		if err == nil {
			t.Errorf("Login() error = %v wantErr = %v", err, contract.ErrAuthenticationFailed)
		}
		err = clt.Challenge()
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		clt.Logout()
	})
}

func TestITClient_GetMetadata(t *testing.T) {

	baseURL := "http://localhost:8080"
	bits := 8

	cookieJar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: cookieJar,
	}

	clt := New(client, Configuration{
		hash:          sha256.New,
		bits:          bits,
		baseURL:       baseURL,
		registerPath:  "/v1/register",
		expkPath:      "/v1/login/expk",
		challengePath: "/v1/login/challenge",
		metadataPath:  "/v1/metadata",
		addPath:       "/v1/add",
		getPath:       "/v1/get",
	}, NewInMemoryUserRepository())

	user, _ := NewUser("metadata-user", bits)
	err := clt.Register(user)
	if err != nil {
		t.Errorf("Register() error = %v", err)
	}

	t.Run("should have no domains", func(t *testing.T) {
		err := clt.Login(user.username, "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}

		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if len(domains) != 0 {
			t.Errorf("domains = %v wantDomains = %v", domains, []string{})
		}
	})

	t.Run("should have google.com domain", func(t *testing.T) {
		// given
		wantDomains := []string{"google.com"}
		err := clt.Login(user.username, "password")
		if err != nil {
			t.Errorf("Login() error = %v", err)
		}
		err = clt.Add(wantDomains[0])
		if err != nil {
			t.Errorf("Add() error = %v", err)
		}

		// when
		domains, err := clt.GetMetadata()
		if err != nil {
			t.Errorf("GetMetadata() error = %v", err)
		}

		if !reflect.DeepEqual(domains, wantDomains) {
			t.Errorf("domains = %v wantDomains = %v", domains, wantDomains)
		}
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
