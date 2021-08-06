package pkg

import (
	"context"
	"fmt"
	"io"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"os/exec"
	"path"
	"sync"
)

var client *kubernetes.Clientset

func LoadKubeClient() {
	home := homedir.HomeDir()
	kubeconfig := path.Join(home, ".kube/config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	// creates the client
	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
}

// KubeConfiguration //
type KubeConfiguration struct {
	Kube struct {
		Name       string
		Deployment string
		Namespace  string
		Container  string
		Port       string
	}
	Format FormatterConfiguration
}

type Kube struct {
	KubeConfiguration
	cmd *exec.Cmd
	out chan Message
}

func (r Kube) String() string {
	return fmt.Sprintf("%v: %v", r.KubeConfiguration.Kube.Name, r.KubeConfiguration.Kube.Deployment)
}

func NewKube(c KubeConfiguration) (*Kube, error) {
	out := make(chan Message)
	runner := Kube{KubeConfiguration: c, out: out}
	return &runner, nil
}

func (r *Kube) ID() string {
	return r.KubeConfiguration.Kube.Name
}

func (r *Kube) Produce(ctx context.Context) <-chan Message {
	return r.out
}

func (r *Kube) Start(ctx context.Context) {
	go r.Kube(ctx)
}

func (r *Kube) Kube(ctx context.Context) {
	// We need to get one pod
	pods, err := client.CoreV1().Pods(r.KubeConfiguration.Kube.Namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	if len(pods.Items) == 0 {
		panic("no pod found")
	}
	// TODO Fix
	pod := pods.Items[0]

	// Forward
	err = r.Forward(pod)

	// Log
	args := []string{"kubectl", "logs", pod.Name, "-n", r.KubeConfiguration.Kube.Namespace, "-f"}
	if r.KubeConfiguration.Kube.Container != "" {
		args = append(args, r.KubeConfiguration.Kube.Container)
	}
	r.cmd = exec.Command(args[0], args[1:]...)
	stdout, _ := r.cmd.StdoutPipe()
	stderr, _ := r.cmd.StderrPipe()
	err = r.cmd.Start()
	if err != nil {
		log.Fatalf("%v: cmd.Start() failed with '%s'\n", r, err)
	}
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, err = io.Copy(NewLineBreaker(r.out), stdout)
		wg.Done()
	}()
	_, _ = io.Copy(NewLineBreaker(r.out), stderr)

	wg.Wait()
	_ = r.cmd.Wait()
}

func (r *Kube) Restart(ctx context.Context) {
	if err := r.cmd.Process.Kill(); err != nil {
		log.Fatalf("failed to kill process %v: %v", r.KubeConfiguration.Kube.Name, err)
	}
	r.out <- Message{Content: "*** restarting ***"}
	go r.Start(ctx)
}

func (r *Kube) Forward(pod v1.Pod) error {
	args := []string{"kubectl", "port-forward", fmt.Sprintf("pod/%v", pod.Name), "-n", r.KubeConfiguration.Kube.Namespace, r.KubeConfiguration.Kube.Port}
	cmd := exec.Command(args[0], args[1:]...)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("%v: cmd.Start() failed with '%s'\n", r, err)
	}
	return nil
}
