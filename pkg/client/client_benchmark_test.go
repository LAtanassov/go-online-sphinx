package client_test

import (
	"context"
	"crypto/sha256"
	"hash"
	"net/http"
	"os"
	"testing"

	"github.com/LAtanassov/go-online-sphinx/pkg/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	docker "github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func BenchmarkClient_Register_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Register(b, 2048, sha256.New)
}
func BenchmarkClient_Register_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Register(b, 2048, sha256.New)
}

func BenchmarkClient_Register_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Register(b, 3072, sha256.New)
}

func benchmarkClient_Register(b *testing.B, bits int, hash func() hash.Hash) {
	baseURL, err := startBench(b, bits, hash)
	if err != nil {
		_ = stopBench(b)
	}
	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		user, err := client.NewUser("new-user", bits)
		if err != nil {
			b.Errorf("NewUser() error = %v", err)
		}
		err = clt.Register(user)
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.StopTimer()
	_ = stopBench(b)
}

func BenchmarkClient_Login_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Login(b, 1024, sha256.New)
}

func BenchmarkClient_Login_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Login(b, 2048, sha256.New)
}

func BenchmarkClient_Login_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Login(b, 3072, sha256.New)
}

func benchmarkClient_Login(b *testing.B, bits int, hash func() hash.Hash) {
	baseURL, err := startBench(b, bits, hash)
	if err != nil {
		_ = stopBench(b)
	}
	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	user, err := client.NewUser("new-user", bits)
	if err != nil {
		b.Errorf("NewUser() error = %v", err)
	}
	err = clt.Register(user)
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		clt.Login("new-user", "password")
		if err != nil {
			b.Errorf("Register() error = %v", err)
		}
	}

	b.StopTimer()
	_ = stopBench(b)
}

func BenchmarkClient_Add_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Add(b, 1024, sha256.New)
}

func BenchmarkClient_Add_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Add(b, 2048, sha256.New)
}

func BenchmarkClient_Add_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Add(b, 3072, sha256.New)
}

func benchmarkClient_Add(b *testing.B, bits int, hash func() hash.Hash) {
	baseURL, err := startBench(b, bits, hash)
	if err != nil {
		_ = stopBench(b)
	}
	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	user, err := client.NewUser("new-user", bits)
	if err != nil {
		b.Errorf("NewUser() error = %v", err)
	}
	err = clt.Register(user)
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	clt.Login("new-user", "password")
	if err != nil {
		b.Errorf("Login() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		err = clt.Add(string(n))
		if err != nil {
			b.Errorf("Add() error = %v", err)
		}
	}

	b.StopTimer()
	_ = stopBench(b)
}

func BenchmarkClient_Get_SHA256_1024Bits(b *testing.B) {
	benchmarkClient_Get(b, 1024, sha256.New)
}

func BenchmarkClient_Get_SHA256_2048Bits(b *testing.B) {
	benchmarkClient_Get(b, 2048, sha256.New)
}

func BenchmarkClient_Get_SHA256_3072Bits(b *testing.B) {
	benchmarkClient_Get(b, 3072, sha256.New)
}

var pwd string

func benchmarkClient_Get(b *testing.B, bits int, hash func() hash.Hash) {
	baseURL, err := startBench(b, bits, hash)
	if err != nil {
		_ = stopBench(b)
	}
	clt := client.New(&http.Client{},
		client.NewConfiguration(baseURL, bits, hash),
		client.NewInMemoryUserRepository())

	user, err := client.NewUser("new-user", bits)
	if err != nil {
		b.Errorf("NewUser() error = %v", err)
	}
	err = clt.Register(user)
	if err != nil {
		b.Errorf("Register() error = %v", err)
	}

	err = clt.Login("new-user", "password")
	if err != nil {
		b.Errorf("Login() error = %v", err)
	}

	err = clt.Add("new-domain")
	if err != nil {
		b.Errorf("Add() error = %v", err)
	}

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		gpwd, err := clt.Get("new-domain")
		if err != nil {
			b.Errorf("Add() error = %v", err)
		}
		pwd = gpwd
	}

	b.StopTimer()
	_ = stopBench(b)
}

// === docker container utils ===
func startBench(b *testing.B, bits int, hashFn func() hash.Hash) (string, error) {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return "", err
	}

	imageName := "latanassov/ossrv:0.1.0"

	_, err = cli.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return "", err
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
		b.Fatal(err)
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		b.Fatal(err)
	}

	os.Setenv("OSSRV_DOCKER_ID", resp.ID)

	return "http://localhost:9090", nil
}

func stopBench(b *testing.B) error {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		return err
	}

	containerID := os.Getenv("OSSRV_DOCKER_ID")

	if err := cli.ContainerStop(ctx, containerID, nil); err != nil {
		return err
	}
	return nil
}
