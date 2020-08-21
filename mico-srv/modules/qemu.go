package modules

import (
	"fmt"
	"net"
	"time"

	"github.com/ThomasRooney/gexpect"
	"github.com/digitalocean/go-qemu/hypervisor"
	"github.com/jstang9527/opentest/ops"
)

var (
	// conn  net.Conn               //old
	hv    *hypervisor.Hypervisor //管理驱动
	mnet  string
	maddr string
	mtime time.Duration
)

// VM 虚拟机 ...
type VM struct {
	Domain string
}

// Init 创建qemu客户端  v2
func Init(network, address string, timeout time.Duration) (err error) {
	newConn := func() (net.Conn, error) {
		return net.DialTimeout(network, address, timeout)
	}
	driver := hypervisor.NewRPCDriver(newConn)
	hv = hypervisor.New(driver)
	ops.Console("connect to hypervisor success")
	return
}

// NewVM 结构体对象
func NewVM(domain string) *VM {
	return &VM{
		Domain: domain,
	}
}

// GetDomainList 返回VM对象列表(通过qemu模块实现)
func GetDomainList() (VMObj []*VM, err error) {
	domains, err := hv.Domains()
	if err != nil {
		return
	}
	//创建对象
	for _, dom := range domains {
		VMObj = append(VMObj, NewVM(dom.Name))
	}
	return
}

// IsRunning 查询主机是否运行(通过qemu模块实现)
func (v *VM) IsRunning() (status bool, err error) {
	dom, err := hv.Domain(v.Domain)
	if err != nil {
		return
	}
	s, err := dom.Status()
	if err != nil {
		status = false
		return // err =所需操作无效，域没有在运行
	}
	if s.String() == "StatusRunning" {
		status = true
		return
	}
	return false, fmt.Errorf("Unknown status, info: %v", s.String())
}

// Recover 恢复快照(通过交互程序expect实现)
func (v *VM) Recover() error {
	dom, err := hv.Domain(v.Domain)
	if err != nil {
		if dom == nil { // domain不存在
			return fmt.Errorf("domain [%v] not found", v.Domain)
		}
		s, _ := dom.Status()
		status := s.String()
		return fmt.Errorf("domain [%v] unknown errors, info: %v", v.Domain, status)
	}
	s, _ := dom.Status()
	status := s.String()
	//域被关闭，可以启快照
	if status == "StatusDebug" {
		bash, err := gexpect.Spawn("sh")
		if err != nil {
			return err
		}
		defer bash.Close()

		bash.Expect("sh-4.2#")
		bash.SendLine("virsh snapshot-revert " + v.Domain + " --snapshotname init")
		bash.Expect("sh-4.2#")
		ops.Console("sh-4.2# execute resume finished. [%v]", v.Domain)
	}
	return err
}

// Start 开机(通过交互程序expect实现)
func (v *VM) Start() error {
	dom, err := hv.Domain(v.Domain)
	if err != nil {
		if dom == nil { // domain不存在
			return fmt.Errorf("domain [%v] not found", v.Domain)
		}
		// domain存在,但是交互出现问题
		s, _ := dom.Status()
		status := s.String()
		return fmt.Errorf("domain [%v] unknown errors, info: %v", v.Domain, status)
	}
	s, _ := dom.Status()
	status := s.String()
	//域被关闭，可以开机
	if status == "StatusDebug" {
		bash, err := gexpect.Spawn("sh")
		if err != nil {
			return err
		}
		defer bash.Close()

		bash.Expect("sh-4.2#")
		bash.SendLine("virsh start " + v.Domain)
		bash.Expect("sh-4.2#")
		ops.Console("sh-4.2# execute start finished. [%v]", v.Domain)
	}
	return err
}

// ShutDown 关闭虚拟机(通过qemu模块实现)
func (v *VM) ShutDown() (err error) {
	vm, err := hv.Domain(v.Domain)
	if err != nil {
		return
	}
	err = vm.SystemPowerdown()
	return
}
