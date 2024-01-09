package docker

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"os"
	"regexp"
	"strings"
)

type Client struct {
	Cli  *client.Client
	host string
	re   *regexp.Regexp //	匹配加载镜像时返回的镜像ID,new 初始化一次
}

func NewClient(host, ca, cert, key string) *Client {

	opts := []client.Opt{
		client.WithHost(host),
		client.WithAPIVersionNegotiation(),
	}
	if ca != "" && cert != "" && key != "" {
		opts = append(opts, client.WithTLSClientConfig(ca, cert, key))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		log.Errorf("[docker] create docker client failed: %v", err)
		panic("")
	}
	log.Info("[docker] create docker client host: %v", err, host)

	pattern := `sha256:[a-f0-9]+`
	return &Client{
		Cli:  cli,
		host: host,
		re:   regexp.MustCompile(pattern),
	}
}

// LoadImage
//
//	imagePath: /opt/raw.tar
//	return: sha256
func (c *Client) LoadImage(ctx context.Context, imagePath string) (string, error) {
	// 打开要导入的镜像文件
	inputFile, err := os.Open(imagePath)
	if err != nil {
		return "", fmt.Errorf("镜像文件打开失败：%v", err)
	}
	defer inputFile.Close()

	// 发起导入镜像的请求
	ctx = context.Background()
	resp, err := c.Cli.ImageLoad(ctx, inputFile, true)
	if err != nil {
		return "", fmt.Errorf("镜像导入失败：%v", err)
	}
	defer resp.Body.Close()

	// 读取并输出导入过程的结果
	var buffer strings.Builder
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {

		return "", fmt.Errorf("输出结果失败：%v", err)
	}

	match := c.re.FindString(buffer.String())
	if match != "" {
		// 去除前缀 "sha256:"
		return match[7:], nil
	} else {
		return "", fmt.Errorf("sha256不存在")
	}
}

// ImageTag imageID => image@latest
func (c *Client) ImageTag(ctx context.Context, imageID, imageTag string) error {
	ctx = context.Background()
	err := c.Cli.ImageTag(ctx, imageID, imageTag)
	if err != nil {
		return fmt.Errorf("镜像重命名失败：%v", err)
	}
	return nil
}

func (c *Client) ContainerStop(ctx context.Context, containerID string, timeout int) error {
	ctx = context.Background()
	// 设置终止容器的超时时间（可选，如果不需要超时可以忽略此步骤）
	// 使用cli.ContainerStop方法终止容器
	var err error
	if timeout == 0 {
		err = c.Cli.ContainerStop(context.Background(), containerID, container.StopOptions{})
	} else {
		err = c.Cli.ContainerStop(context.Background(), containerID, container.StopOptions{
			Timeout: &timeout,
		})
	}
	if err != nil {
		fmt.Println("终止容器出错：", err)
		return err
	}
	fmt.Println("容器已成功终止！")
	return nil
}
func (c *Client) ContainerRemove(ctx context.Context, containerIDorName string, force bool) error {
	ctx = context.Background()
	// 删除容器
	var err error
	if force {
		err = c.Cli.ContainerRemove(context.Background(), containerIDorName, types.ContainerRemoveOptions{
			RemoveVolumes: true, // 如果容器使用了卷，同时删除关联的卷
			Force:         true, // 强制删除，即使容器在运行时也删除
		})
	} else {
		err = c.Cli.ContainerRemove(context.Background(), containerIDorName, types.ContainerRemoveOptions{})
	}

	if err != nil {
		fmt.Println("容器删除失败:", err)
		return err
	}

	fmt.Println("容器删除成功")
	return nil
}

// ContainerCreate
// 1. 配置要启动的容器
//
//	containerConfig := &container.Config{
//		Image: "docker@latest",                    // 指定要使用的镜像
//		Cmd:   []string{"echo", "Hello, Docker!"}, // 指定容器启动时要执行的命令
//		Env:   []string{"LOCAL_USER_ID=" + "id -u $USER"},	//
//	}
//
// 2. 配置容器自动重启 网络 磁盘挂载等
//
//	hostConfig := &container.HostConfig{
//		RestartPolicy: container.RestartPolicy{
//			Name: "always", // 设置重启策略为"always"，容器将总是自动重启
//			// 可选的重启策略：
//			// - "no"：无重启策略
//			// - "always"：容器总是自动重启
//			// - "on-failure"：容器在非零退出状态时重启（默认最多重启3次）
//			// - "unless-stopped"：除非手动停止，否则容器总是自动重启
//		},
//		NetworkMode: "host", // 设置主机网络模式，
//		// 可选的策略:
//		// - "none"
//		// - "default"
//		// - "host"：与主机共享网络 无需配置端口映射
//		Privileged: true,
//		CapAdd:     []string{"SYS_ADMIN", "IPC_LOCK"},
//		//PortBindings: nat.PortMap{
//		//	"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}}, // 将容器的80端口映射到宿主机的8080端口
//		//},
//		Mounts: []mount.Mount{
//			{
//				Type:     mount.TypeBind,
//				Source:   "/opt", // 宿主机上要挂载的文件夹路径
//				Target:   "/opt", // 容器内挂载的路径
//				ReadOnly: false,  // 是否只读
//			},
//		},
//	}
func (c *Client) ContainerCreate(ctx context.Context, containerName string, containerConfig *container.Config, hostConfig *container.HostConfig) error {
	//// 配置要启动的容器
	//containerConfig := &container.Config{
	//	Image: "docker@latest",                    // 指定要使用的镜像
	//	Cmd:   []string{"echo", "Hello, Docker!"}, // 指定容器启动时要执行的命令
	//	Env:   []string{"LOCAL_USER_ID=" + "id -u $USER"},
	//}
	//
	//// 配置容器自动重启
	//hostConfig := &container.HostConfig{
	//	RestartPolicy: container.RestartPolicy{
	//		Name: "always", // 设置重启策略为"always"，容器将总是自动重启
	//		// 可选的重启策略：
	//		// - "no"：无重启策略
	//		// - "always"：容器总是自动重启
	//		// - "on-failure"：容器在非零退出状态时重启（默认最多重启3次）
	//		// - "unless-stopped"：除非手动停止，否则容器总是自动重启
	//	},
	//	NetworkMode: "host", // 设置主机网络模式，
	//	// 可选的策略:
	//	// - "none"
	//	// - "default"
	//	// - "host"：与主机共享网络 无需配置端口映射
	//	Privileged: true,
	//	CapAdd:     []string{"SYS_ADMIN", "IPC_LOCK"},
	//	//PortBindings: nat.PortMap{
	//	//	"80/tcp": []nat.PortBinding{{HostIP: "0.0.0.0", HostPort: "8080"}}, // 将容器的80端口映射到宿主机的8080端口
	//	//},
	//	Mounts: []mount.Mount{
	//		{
	//			Type:     mount.TypeBind,
	//			Source:   "/opt", // 宿主机上要挂载的文件夹路径
	//			Target:   "/opt", // 容器内挂载的路径
	//			ReadOnly: false,  // 是否只读
	//		},
	//	},
	//}

	// 设置容器名称
	// 创建并启动容器
	ctx = context.Background()
	resp, err := c.Cli.ContainerCreate(context.Background(), containerConfig, hostConfig, nil, nil, containerName)
	if err != nil {
		fmt.Println("容器创建失败:", err)
		return err
	}

	// 启动容器
	if err := c.Cli.ContainerStart(context.Background(), resp.ID, types.ContainerStartOptions{}); err != nil {
		fmt.Println("容器启动失败:", err)
		return err
	}

	fmt.Println("容器ID:", resp.ID)

	return nil
}
